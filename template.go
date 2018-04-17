package main

import (
	"text/template"
	"time"
	"unicode/utf8"
)

var (
	defaultFuncMap = template.FuncMap{
		"fdate": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"substr": func(s string, n int) string {
			if utf8.RuneCountInString(s) <= n {
				return s
			}
			return string([]rune(s)[:n]) + "..."
		},
	}

	detailTmpl, _ = template.New("text").Funcs(defaultFuncMap).Parse(`<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>{{.Title}}</title>
	<meta name="author" content="{{.Author}}" />
	<meta name="description" content="{{substr .Description 64}}" />
</head>
<body>
	<h3>{{.Title}}</h3>
	<hr/>
	<div>{{fdate .Published "Jan 02, 2006"}}{{if ne .Author ""}} by {{.Author}}{{end}}</div>
	<section id="content">{{.Content}}</section>
	<hr/>
	<footer>{{.Url}}</footer>
	</body>
</html>`)

	tocTmpl, _ = template.New("toc").Parse(`<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>
	<meta content="text/html; charset=utf-8" http-equiv="Content-Type"/>
	<title>{{.Title}}</title>
</head>
<body>
	<h1>{{.Title}}</h1>
	{{ range $_,$section:=.List }}
	<h4>{{ $section.Title }}</h4>
	<ul>
		{{ range $,$f:= $section.List}}<li><a href="{{ $f.Path }}">{{ index $f.Prop "title" }}</a></li>{{end}}
	</ul>
	{{ end }}
</body>
</html>`)

	ncxTmpl, _ = template.New("ncx").Parse(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx xmlns:mbp="http://mobipocket.com/ns/mbp" xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="en-US">
	<head>
		<meta name="dtb:uid" content="urn:uuid:{{.UUID}}" />
		<meta name="dtb:depth" content="2" />
		<meta name="dtb:totalPageCount" content="0" />
		<meta name="dtb:maxPageNumber" content="0" />
	</head>
	<docTitle><text>{{.Title}}</text></docTitle>
	<docAuthor><text>{{.Author}}</text></docAuthor>	
	<navMap>
	{{range $,$nav:=.NavPoint}}
		<navPoint playOrder="{{$nav.Order}}" class="periodical" id="periodical">
			<navLabel><text>{{index $nav.File.Prop "title"}}</text></navLabel>
			<content src="{{$nav.File.Path}}" />
			{{range $,$nav:=$nav.Child}}
			<navPoint playOrder="{{$nav.Order}}" class="section" id="section-{{$nav.Order}}">
				<navLabel><text>{{index $nav.File.Prop "title"}}</text></navLabel>
				<content src="{{$nav.File.Path}}" />
				{{range $,$nav:=$nav.Child}}
					<navPoint playOrder="{{$nav.Order}}" class="article" id="item-{{$nav.Order}}">
					<navLabel><text>{{index $nav.File.Prop "title"}}</text></navLabel>
					<content src="{{$nav.File.Path}}" />
					{{if (index $nav.File.Prop "description") }}<mbp:meta name="description">{{index $nav.File.Prop "description"}}</mbp:meta>{{end}}
					{{if (index $nav.File.Prop "author") }}<mbp:meta name="author">{{index $nav.File.Prop "author"}}</mbp:meta>{{end}}
					</navPoint>
				{{end}}
			</navPoint>
			{{end}}
		</navPoint>
	{{end}}
	</navMap>
</ncx>`)

	opfTmpl, _ = template.New("opf").Funcs(defaultFuncMap).Parse(`<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="BookId">
<metadata>
	<dc-metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
		<dc:title>{{index .Metadata "title"}}</dc:title>
		<dc:language>{{index .Metadata "language"}}</dc:language>
		<dc:Identifier id="BookId">urn:uuid:{{.UUID}}</dc:Identifier>
		<dc:creator>{{index .Metadata "creator"}}</dc:creator>
		<dc:publisher>{{index .Metadata "publisher"}}</dc:publisher>				
		<dc:date>{{fdate (index .Metadata "date") "2006-01-02"}}</dc:date>
		<dc:description>{{index .Metadata "description"}}</dc:description>	
	</dc-metadata>
	<x-metadata>
		<output content-type="application/x-mobipocket-subscription-magazine" encoding="utf-8"/>		
	</x-metadata>
</metadata>
<manifest>
{{range $,$f:=.Manifest}}
<item href="{{$f.Path}}" media-type="{{$f.MediaType}}" id="{{$f.Id}}"/>
{{end}}
</manifest>
<spine toc="{{.Ncx.Id}}">
{{range $,$f:=.Spine}}
<itemref idref="{{$f.Id}}"/>
{{end}}
</spine>
<guide>
<reference href="{{.Toc.Path}}" type="toc" title="{{index .Toc.Prop "title"}}" />
</guide>
</package>
`)
)
