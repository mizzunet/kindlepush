package main

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
	"github.com/zhengchun/objectid"
	"github.com/zhengchun/syndfeed"
)

type server struct {
	conf   *appConfig
	once   sync.Once
	client *http.Client
}

// About An article information and body.
type article struct {
	Title       string
	Url         string
	Content     string
	Description string
	Author      string
	Category    string
	Published   time.Time
}

func (s *server) init() {
	var proxyFunc = http.ProxyFromEnvironment
	if s.conf.Proxy != "" {
		proxyURL, err := url.Parse(s.conf.Proxy)
		if err != nil {
			logrus.Warn("proxy address is invalid")
		} else {
			proxyFunc = http.ProxyURL(proxyURL)
		}
	}
	s.client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 proxyFunc,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func (s *server) run() {
	s.once.Do(s.init)
	// Make sure directory is exists.
	os.Mkdir(s.conf.CacheDir, os.ModePerm)
	if len(s.conf.Feeds) == 0 {
		logrus.Warn("no RSS feed in config.yaml")
		return
	}
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		list []*section
	)
	for _, feed := range s.conf.Feeds {
		wg.Add(1)
		feedURL := feed
		go func() {
			defer wg.Done()
			sec, err := s.fetchFeed(feedURL)
			if err != nil {
				logrus.Warn(err)
			}
			mu.Lock()
			list = append(list, sec)
			mu.Unlock()
		}()
	}
	wg.Wait()
	f, err := s.buildMobi(list)
	if err != nil {
		logrus.Warnf("mobi build failed, %v", err)
		return
	}
	//
	//defer os.Remove(f.Name())
	if f.Size() > int64(s.conf.MaxFileSize*1024*1024) {
		logrus.Warnf("mobi size exceeds the allowable limit %dM", s.conf.MaxFileSize)
		return
	}
	if err := s.sendMobi(filepath.Join(s.conf.CacheDir, f.Name())); err != nil {
		logrus.Warnf("send ebook failed: %v", err)
		return
	}
}

func (s *server) sendMobi(attachFile string) error {
	logrus.Info("sending mail to your kindle device, it will takes a little time")
	smtpConf := s.conf.Smtp
	senderAddr := smtpConf.SenderAddress
	if senderAddr == "" {
		senderAddr = smtpConf.Account
	}
	m := gomail.NewMessage()
	m.SetHeader("From", senderAddr)
	m.SetHeader("To", s.conf.KindleAddress)
	m.SetHeader("Subject", "mobi")
	m.SetBody("text/html", "")
	m.Attach(attachFile)

	var (
		host string = smtpConf.HostAndPort
		port int    = 25
	)
	if i := strings.Index(smtpConf.HostAndPort, ":"); i > 0 {
		host = host[:i]
		v, err := strconv.ParseInt(smtpConf.HostAndPort[i+1:], 10, 64)
		if err != nil {
			return errors.New("smtp.hostAndPort is invalid in config.yaml")
		}
		port = int(v)
	}
	d := gomail.NewDialer(host, port, smtpConf.Account, smtpConf.Password)
	d.SSL = smtpConf.SSL
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	logrus.Infof("%s has been send successfully", filepath.Base(attachFile))
	return nil
}

func (s *server) buildMobi(list []*section) (os.FileInfo, error) {
	logrus.Info("building mobi file")
	f, err := createOpf(s.conf.CacheDir, list)
	if err != nil {
		return nil, err
	}
	outname := filepath.Join(s.conf.CacheDir, fmt.Sprintf("%s.mobi", objectid.New()))
	cmd := exec.Command(s.conf.Kindlegen, f.Path, "-c1", "-o", filepath.Base(outname))
	cmd.Dir = s.conf.CacheDir
	cmd.Run()

	stat, err := os.Stat(outname)
	if err != nil {
		return nil, err
	}
	logrus.Infof("build successful, %s in %d bytes", stat.Name(), stat.Size())
	return stat, nil
}

func (s *server) fetchFeed(urlStr string) (*section, error) {
	logrus.Debugf("downloading feed %s", urlStr)
	res, err := s.client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("feed download failed: %v", err)
	}
	defer res.Body.Close()

	feed, err := syndfeed.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("feed parse failed: %v", err)
	}
	sec := &section{Title: feed.Title}
	for _, item := range feed.Items {
		post := &article{
			Title:       item.Title,
			Content:     item.Content,
			Description: item.Summary,
			Published:   item.PublishDate,
		}
		if post.Content == "" {
			post.Content = item.Summary
		}
		if len(item.Links) > 0 {
			post.Url = item.Links[0].URL
		}
		if len(item.Categories) > 0 {
			post.Category = item.Categories[0]
		}
		if len(item.Authors) > 0 {
			post.Author = item.Authors[0].Name
		}
		f, err := createDetailFile(s.client, s.conf.CacheDir, post, s.resizeImageHandler)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		sec.List = append(sec.List, f)
	}
	return sec, nil
}

func (s *server) resizeImageHandler(src image.Image) (image.Image, error) {
	if s.conf.ResizeImage == "" {
		return src, nil
	}
	var maxWidth, maxHeight uint
	a := strings.Split(s.conf.ResizeImage, "x")
	if len(a) != 2 {
		logrus.Warn("resizeImage argument is invalid")
		return src, nil
	}
	if v, err := strconv.ParseUint(a[0], 10, 64); err == nil {
		maxWidth = uint(v)
	}
	if v, err := strconv.ParseUint(a[1], 10, 64); err == nil {
		maxHeight = uint(v)
	}
	dstImg := resize.Thumbnail(maxWidth, maxHeight, src, resize.NearestNeighbor)
	return dstImg, nil
}

func newServer(conf *appConfig) *server {
	return &server{
		conf: conf,
	}
}
