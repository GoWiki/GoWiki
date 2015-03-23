package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"html/template"
	"io"
	"strings"

	"github.com/boltdb/bolt"
)

type Theme struct {
	Name string
}

func (theme *Theme) GetFile(t *bolt.Tx, file string) []byte {
	tx := &WikiTx{t}
	b_themes := tx.Themes()
	b_theme := b_themes.Bucket([]byte(theme.Name))
	b_static := b_theme.Bucket([]byte("static"))
	return b_static.Get([]byte(file))
}

func (theme *Theme) ParseTemplates(t *bolt.Tx, templates *template.Template) {
	tx := &WikiTx{t}
	var b_themes *bolt.Bucket = tx.Themes()
	b_theme := b_themes.Bucket([]byte(theme.Name))
	b_templates := b_theme.Bucket([]byte("templates"))
	c_templates := b_templates.Cursor()
	for name, temp := c_templates.First(); name != nil; name, temp = c_templates.Next() {
		newtemp := templates.New(string(name))
		newtemp.Parse(string(temp))
	}
}

func InsertTheme(t *bolt.Tx, name string, tarfile []byte) {
	tx := &WikiTx{t}
	gzipreader, _ := gzip.NewReader(bytes.NewReader(tarfile))
	tarreader := tar.NewReader(gzipreader)
	b_themes := tx.Themes()
	b_themes.DeleteBucket([]byte(name))
	b_theme, _ := b_themes.CreateBucket([]byte(name))
	b_static, _ := b_theme.CreateBucket([]byte("static"))
	b_templates, _ := b_theme.CreateBucket([]byte("templates"))
	for {
		header, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		buf := &bytes.Buffer{}
		io.Copy(buf, tarreader)
		name := header.Name
		if strings.HasPrefix(name, "static/") {
			name = strings.TrimPrefix(name, "static/")
			b_static.Put([]byte(name), buf.Bytes())
		} else if strings.HasPrefix(name, "templates/") {
			name = strings.TrimPrefix(name, "templates/")
			b_templates.Put([]byte(name), buf.Bytes())
		}
	}
}
