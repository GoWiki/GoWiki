// GoWiki project main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var (
	DB     *bolt.DB
	tpl    *template.Template
	router *mux.Router
)

func init() {
	tpl = template.Must(template.New("default").Funcs(template.FuncMap{
		"PageNav":    PageNav,
		"Route":      Route,
		"GetContent": GetContent,
	}).ParseGlob("templates/default/*"))
}

func main() {
	db, err := bolt.Open("gowiki.db", 0600, nil)
	DB = db
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		SetupBuckets(tx)
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	router = mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/default"))))
	router.HandleFunc("/{page:.*}/edit", EditHandler).Methods("GET").Name("Edit")
	router.HandleFunc("/{page:.*}", PageHandler).Methods("GET").Name("Read")
	router.HandleFunc("/{page:.*}", UpdateHandler).Methods("POST").Name("Update")

	http.ListenAndServe(":3000", router)
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

func PageNav(Slug string, Section string) PageNavData {
	return PageNavData{
		Read:    UrlToPath(router.Get("Read").URLPath("page", Slug)),
		Edit:    UrlToPath(router.Get("Edit").URLPath("page", Slug)),
		Section: Section,
	}
}

func GetContent(Slug string) (Content template.HTML) {
	DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, Slug)
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe := blackfriday.MarkdownCommon(pagedata)
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			Content = template.HTML(html)
		}
		return nil
	})
	return
}

func Route(Slug string, Route string) string {
	return UrlToPath(router.Get(Route).URLPath("page", Slug))
}

func PageHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe := blackfriday.MarkdownCommon(pagedata)
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			rw.Header().Set("Content-Type", "text/html")

			data := struct {
				Content template.HTML
				Name    string
				Slug    string
			}{
				template.HTML(string(html)),
				vars["page"],
				vars["page"],
			}

			if err := tpl.ExecuteTemplate(rw, "view.tpl", data); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})

}

func UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	DB.Update(func(tx *bolt.Tx) error {
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
	PageHandler(rw, req)
}

func EditHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	DB.View(func(tx *bolt.Tx) error {
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

			if err := tpl.ExecuteTemplate(rw, "edit.tpl", data); err != nil {
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
			if err := tpl.ExecuteTemplate(rw, "edit.tpl", data); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
}
