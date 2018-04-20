package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhengchun/objectid"
)

type fileType int

const (
	typeNcx fileType = iota
	typeToc
	typeOpf
	typeText
	typeImage
)

type file struct {
	Type     fileType
	Path     string
	Resource []*file
	Prop     map[string]interface{}
}

type section struct {
	Title string
	List  []*file
}

func (f *file) Id() string {
	return strconv.Itoa(int(fnvHash(f.Path)))
}

func (f *file) MediaType() (typ string) {
	typ = "application/xhtml+xml"

	switch {
	case f.Type == typeNcx:
		typ = "application/x-dtbncx+xml"
	case f.Type == typeImage:
		ext := filepath.Ext(f.Path)
		switch ext {
		case ".jpg", ".jpeg":
			typ = "image/jpeg"
		case ".bmp":
			typ = "image/bmp"
		case ".gif":
			typ = "image/gif"
		case ".png":
			typ = "image/png"
		}
	}
	return typ
}

var (
	imgRegex = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
)

func adjustImage(dir string, r io.ReadCloser, resizeImage func(image.Image) (image.Image, error)) (*file, error) {
	// resize image
	img, typ, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img, err = resizeImage(img)
	if err != nil {
		return nil, fmt.Errorf("resize image failed: %v", err)
	}

	ext := ".jpg"
	switch typ {
	case "gif":
		ext = ".gif"
	case "png":
		ext = ".png"
	}

	name := filepath.Join(dir, objectid.New().String()+ext)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch typ {
	case "gif":
		if err := gif.Encode(f, img, &gif.Options{NumColors: 256}); err != nil {
			return nil, err
		}
	case "png":
		if err := png.Encode(f, img); err != nil {
			return nil, err
		}
	default: // jpeg or other types
		if err := jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return nil, err
		}
	}
	return &file{
		Type: typeImage,
		Path: name,
	}, nil
}

func createDetailFile(client *http.Client, dir string, post *article, resizeImage func(image.Image) (image.Image, error)) (*file, error) {
	var (
		imgs []*file
		wg   sync.WaitGroup
	)
	u, _ := url.Parse(post.Url)
	// download all image that show in an article.
	for _, g := range imgRegex.FindAllStringSubmatch(post.Content, -1) {
		imgSrc := g[1]
		var imgURL string
		if strings.HasPrefix(imgSrc, "http://") || strings.HasPrefix(imgSrc, "https://") {
			imgURL = imgSrc
		} else {
			if u2, err := u.Parse(imgSrc); err == nil {
				imgURL = u2.String()
			}
		}
		if imgURL == "" {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			logrus.Debugf("downloading image %s", imgURL)
			resp, err := client.Get(imgURL)
			if err != nil {
				logrus.Warnf("image %s download failed.", err)
				return
			}
			defer resp.Body.Close()
			f, err := adjustImage(dir, resp.Body, resizeImage)
			if err != nil {
				if err == image.ErrFormat {
					// not supported image format.
					post.Content = strings.Replace(post.Content, imgSrc, "", -1)
				}
				logrus.Warnf("resize image failed, %v", err)
				return
			}
			post.Content = strings.Replace(post.Content, imgSrc, f.Path, -1)
			imgs = append(imgs, f)
		}()
	}
	wg.Wait()
	logrus.Debug("all image download completed")
	name := path.Join(dir, fmt.Sprintf("%d.html", fnvHash(post.Url)))
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := detailTmpl.Execute(f, post); err != nil {
		return nil, err
	}
	logrus.Debugf("%s generate successful", post.Url)
	return &file{
		Type:     typeText,
		Path:     name,
		Resource: imgs,
		Prop: map[string]interface{}{
			"title":       post.Title,
			"author":      post.Author,
			"description": post.Description,
		},
	}, nil
}

func createNcx(unid uuid.UUID, dir string, toc *file, list []*section) (*file, error) {
	type navPoint struct {
		Order int
		File  *file
		Child []*navPoint
	}

	n := 0
	createNav := func(f *file) (nav *navPoint) {
		nav = &navPoint{Order: n, File: f}
		n++
		return nav
	}

	rootNav := createNav(toc)
	for _, section := range list {
		if len(section.List) == 0 {
			continue
		}
		nav := createNav(&file{Type: section.List[0].Type, Path: section.List[0].Path, Prop: map[string]interface{}{"title": section.Title}})
		rootNav.Child = append(rootNav.Child, nav)
		for _, f := range section.List {
			nav.Child = append(nav.Child, createNav(f))
		}
	}

	name := filepath.Join(dir, fmt.Sprintf("%s.ncx", objectid.New()))
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := struct {
		UUID     uuid.UUID
		Title    string
		Author   string
		NavPoint []*navPoint
	}{
		UUID:     unid,
		Title:    time.Now().Format("2006-01-02"),
		Author:   "KindlePush",
		NavPoint: []*navPoint{rootNav},
	}
	if err := ncxTmpl.Execute(f, data); err != nil {
		return nil, err
	}

	return &file{
		Type: typeNcx,
		Path: name,
	}, nil
}

func createToc(dir string, list []*section) (*file, error) {
	name := filepath.Join(dir, fmt.Sprintf("%s.html", objectid.New()))
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := struct {
		Title string
		List  []*section
	}{
		Title: "Table of Contents",
		List:  list,
	}
	if err := tocTmpl.Execute(f, data); err != nil {
		return nil, err
	}
	return &file{
		Type: typeToc,
		Path: name,
		Prop: map[string]interface{}{"title": data.Title},
	}, nil
}

func createOpf(dir string, list []*section) (*file, error) {
	var (
		toc, ncx *file
		err      error
		now      = time.Now()
		unid     = uuid.New()

		metadata = map[string]interface{}{
			"title":       "Kindle Magazine",
			"language":    "en",
			"creator":     "KindlePush",
			"publisher":   "KindlePush",
			"description": "",
			"date":        now,
		}
	)

	if toc, err = createToc(dir, list); err != nil {
		return nil, fmt.Errorf("create file failed - %s", err)
	}
	if ncx, err = createNcx(unid, dir, toc, list); err != nil {
		return nil, fmt.Errorf("create file failed - %s", err)
	}

	name := filepath.Join(dir, fmt.Sprintf("%s.opf", objectid.New()))
	f, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("create file failed - %s", err)
	}
	defer f.Close()

	var (
		manifest = []*file{toc, ncx}
		spine    = []*file{toc}
	)
	for _, section := range list {
		spine = append(spine, section.List...)
		manifest = append(manifest, section.List...)
		for _, f := range section.List {
			if len(f.Resource) > 0 {
				manifest = append(manifest, f.Resource...)
			}
		}
	}
	data := struct {
		UUID     uuid.UUID
		Manifest []*file
		Spine    []*file
		Ncx      *file
		Toc      *file
		Metadata map[string]interface{}
	}{
		UUID:     unid,
		Manifest: manifest,
		Spine:    spine,
		Ncx:      ncx,
		Toc:      toc,
		Metadata: metadata,
	}
	if err := opfTmpl.Execute(f, data); err != nil {
		os.Remove(toc.Path)
		os.Remove(ncx.Path)
		return nil, err
	}

	return &file{
		Path:     name,
		Type:     typeOpf,
		Resource: []*file{toc, ncx},
	}, nil
}
