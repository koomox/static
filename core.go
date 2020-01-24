package static

import (
	"bytes"
	"encoding/base64"
	"github.com/koomox/ext"
	"io/ioutil"
	"path"
	"strings"
)

type File struct {
	Name string
	ContentType string
	Content []byte
}

type Item struct {
	root string
	filter []string
	prefix string
}

type Encoding struct {
	buffer bytes.Buffer
	item []*Item
}

func NewEncoding() *Encoding {
	return &Encoding{}
}

func DeCompress(data []byte) (buf []byte, err error) {
	if data, err = base64.RawStdEncoding.DecodeString(string(data)); err != nil {
		return
	}

	return ext.NewEncoding().DeCompress(data)
}

func (this *Encoding)CompressAllFile() (buffer []byte, err error) {
	var (
		b bytes.Buffer
		buf []byte
	)

	b.Write([]byte("var files = []File{\n"))
	for _, item := range this.item {
		if buf, err = item.Compress(); err != nil {
			return
		}
		b.Write(buf)
	}
	b.Write([]byte("}"))
	buffer = b.Bytes()
	return
}

func (this *Encoding)New() (*Item) {
	item := &Item{}
	this.item = append(this.item, item)
	return item
}

func (this *Item)Compress() (buffer []byte, err error){
	var (
		fs []string
		b bytes.Buffer
		buf []byte
	)
	if this.root == "" {
		if this.root, err = ext.GetCurrentDirectory(); err != nil {
			return
		}
	}
	this.root = path.Join(this.root, "")
	if fs, err = ext.GetCustomDirectoryAllFile(this.root); err != nil {
		return
	}

	for _, f := range fs {
		filter := false
		for _, v := range this.filter {
			if strings.EqualFold(f, v) || strings.HasPrefix(f, v) {
				filter = true
				break
			}
		}
		if filter {
			continue
		}
		b.Write([]byte("\tFile{\n\t\tName: \""))
		b.WriteString(path.Join(this.prefix, strings.TrimPrefix(f, this.root)))
		b.Write([]byte("\",\n\t\tContentType: \""))
		switch path.Ext(f) {
		case ".js":
			b.WriteString("application/javascript; charset=utf-8")
		case ".css":
			b.WriteString("text/css; charset=utf-8")
		case ".html", ".htm", ".php":
			b.WriteString("text/html; charset=utf-8")
		case ".jpg", "jpeg":
			b.WriteString("image/jpeg")
		case ".gif":
			b.WriteString("image/gif")
		case ".png":
			b.WriteString("image/png")
		case ".svg":
			b.WriteString("image/svg+xml")
		case ".webp":
			b.WriteString("image/webp")
		case ".xml":
			b.WriteString("text/xml; charset=utf-8")
		case ".pdf":
			b.WriteString("application/pdf")
		case ".otf":
			b.WriteString("font/otf")
		case ".ttf":
			b.WriteString("font/ttf")
		case ".woff":
			b.WriteString("font/woff")
		case ".woff2":
			b.WriteString("font/woff2")
		default:
			b.WriteString("text/plain; charset=utf-8")
		}
		b.Write([]byte("\",\n\t\tContent: []byte(\""))
		if buf, err = ioutil.ReadFile(f);err != nil {
			return
		}
		if buf, err = ext.NewEncoding().Compress(buf); err != nil {
			return
		}
		b.WriteString(base64.RawStdEncoding.EncodeToString(buf))
		b.Write([]byte("\"),\n\t},\n"))
	}
	buffer = b.Bytes()
	return
}

func (this *Item)Filter(elem ...string) (*Item){
	for i, e := range elem {
		if e != "" {
			this.filter = append(this.filter, elem[i])
		}
	}

	return this
}

func (this *Item)Root(root string) (*Item){
	root = strings.Replace(root, "\\", "/", -1)
	this.root = root

	return this
}

func (this *Item)Prefix(prefix string) (*Item){
	prefix = strings.Replace(prefix, "\\", "/", -1)
	this.prefix = path.Join("/", prefix)
	return this
}
