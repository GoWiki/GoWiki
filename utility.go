package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"

	"github.com/boltdb/bolt"
)

func GetRandomID() string {
	data := make([]byte, 32)
	rand.Read(data)
	return base64.StdEncoding.EncodeToString(data)
}

func (w *Wiki) WikiLink(href string, text string) (link string) {
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, href)
		if page == nil {
			link = "<a href=\"/" + href + "\" class=\"empty-link\">" + text + "</a>"
		} else {
			link = "<a href=\"/" + href + "\" >" + text + "</a>"
		}
		return nil
	})
	return
}

type PageNavData struct {
	Read    string
	Edit    string
	History string
	Section string
}

func (w *Wiki) PageNav(Slug string, Section string) PageNavData {
	return PageNavData{
		Read:    UrlToPath(w.router.Get("Read").URLPath("page", Slug)),
		Edit:    UrlToPath(w.router.Get("Edit").URLPath("page", Slug)),
		History: UrlToPath(w.router.Get("History").URLPath("page", Slug)),
		Section: Section,
	}
}

func UrlToPath(url *url.URL, err error) string {
	if err != nil {
		panic(err)
	}
	return url.Path
}

func (w *Wiki) EncodeID(id []byte) string {
	return base64.URLEncoding.EncodeToString(id)
}

func (w *Wiki) GetContent(Slug string) (Content template.HTML) {
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, Slug)
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe, _ := w.render.Transform(string(pagedata))
			html := w.policy.Sanitize(unsafe)
			Content = template.HTML(html)
		}
		return nil
	})
	return
}

func (w *Wiki) Route(Route string, Params ...string) string {
	return UrlToPath(w.router.Get(Route).URLPath(Params...))
}

type UserInfoType struct {
	LoggedIn bool
	Name     string
}

func (w *Wiki) UserInfo(req *http.Request) UserInfoType {
	s := w.store.Get(req)
	ui := UserInfoType{}

	if s.User != nil {
		ui.LoggedIn = true
		ui.Name = s.User.Name
	}
	return ui
}
