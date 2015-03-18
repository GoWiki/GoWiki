// GoWiki project main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/andyleap/cajun"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

type Wiki struct {
	DB     *bolt.DB
	tpl    *template.Template
	router *mux.Router
	render *cajun.Cajun
	policy *bluemonday.Policy
}

func (w *Wiki) WikiLink(href string, text string) (link string) {
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, href)
		if page == nil {
			link = "<a href=\"" + href + "\" class=\"empty-link\">" + text + "</a>"
		} else {
			link = "<a href=\"" + href + "\" >" + text + "</a>"
		}
		return nil
	})
	return
}

func init() {

}

func main() {
	wiki := New()
	fmt.Println("it's working?")
	http.ListenAndServe(":3000", wiki.router)
}

func New() *Wiki {
	wiki := &Wiki{}
	db, err := bolt.Open("gowiki.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		SetupBuckets(tx)
		return nil
	})
	wiki.DB = db

	if err != nil {
		fmt.Println(err)
		return nil
	}

	tpl := template.Must(template.New("default").Funcs(template.FuncMap{
		"PageNav":    wiki.PageNav,
		"Route":      wiki.Route,
		"GetContent": wiki.GetContent,
	}).ParseGlob("templates/default/*"))

	wiki.tpl = tpl

	wiki.render = cajun.New()
	wiki.render.WikiLink = wiki
	wiki.policy = bluemonday.UGCPolicy()
	wiki.policy.AllowAttrs("class").OnElements("a")

	wiki.router = mux.NewRouter()
	wiki.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/default"))))
	wiki.router.HandleFunc("/{page}/edit", wiki.EditHandler).Methods("GET").Name("Edit")
	wiki.router.HandleFunc("/{page}", wiki.PageHandler).Methods("GET").Name("Read")
	wiki.router.HandleFunc("/{page}", wiki.UpdateHandler).Methods("POST").Name("Update")

	return wiki
}

type PageNavData struct {
	Read    string
	Edit    string
	Section string
}

func UrlToPath(url *url.URL, err error) string {
	if err != nil {
		panic(err)
	}
	return url.Path
}

func (w *Wiki) PageNav(Slug string, Section string) PageNavData {
	return PageNavData{
		Read:    UrlToPath(w.router.Get("Read").URLPath("page", Slug)),
		Edit:    UrlToPath(w.router.Get("Edit").URLPath("page", Slug)),
		Section: Section,
	}
}

func (w *Wiki) GetContent(Slug string) (Content template.HTML) {
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, Slug)
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe, _ := w.render.Transform(string(pagedata))
			fmt.Println(unsafe)
			html := w.policy.Sanitize(unsafe)
			Content = template.HTML(html)
		}
		return nil
	})
	return
}

func (w *Wiki) Route(Slug string, Route string) string {
	return UrlToPath(w.router.Get(Route).URLPath("page", Slug))
}

func (w *Wiki) PageHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe, _ := w.render.Transform(string(pagedata))
			html := w.policy.Sanitize(unsafe)
			rw.Header().Set("Content-Type", "text/html")

			data := struct {
				Content template.HTML
				Name    string
				Slug    string
			}{
				template.HTML(html),
				vars["page"],
				vars["page"],
			}

			if err := w.tpl.ExecuteTemplate(rw, "view.tpl", data); err != nil {
				fmt.Println(err)
			}
		} else {
			w.EditHandler(rw, req)
		}
		return nil
	})

}

func (w *Wiki) UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.Update(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		key, _ := SaveData(tx, []byte(req.FormValue("data")))
		if page != nil {
			//page.History.Events = append(page.History.Events, page.Current)
		} else {
			page = &Page{}
		}
		page.Current = Event{DataID: key, IP: req.RemoteAddr}
		page.Save(tx, vars["page"])
		return nil
	})
	w.PageHandler(rw, req)
}

func (w *Wiki) EditHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			pagedata := page.Current.GetData(tx)
			rw.Header().Add("Content-Type", "text/html")
			data := struct {
				Content string
				Name    string
				Slug    string
			}{
				string(pagedata),
				vars["page"],
				vars["page"],
			}

			if err := w.tpl.ExecuteTemplate(rw, "edit.tpl", data); err != nil {
				fmt.Println(err)
			}
		} else {
			rw.Header().Add("Content-Type", "text/html")
			data := struct {
				Content string
				Name    string
				Slug    string
			}{
				"",
				vars["page"],
				vars["page"],
			}
			if err := w.tpl.ExecuteTemplate(rw, "edit.tpl", data); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
}
