// GoWiki project main.go
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/andyleap/cajun"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/gowiki/greentuesday"
	"github.com/justinas/alice"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

type Wiki struct {
	DB      *bolt.DB
	tpl     *template.Template
	router  *mux.Router
	render  *cajun.Cajun
	policy  *bluemonday.Policy
	gpolicy *greentuesday.Policy
	store   *MemoryStore
	config  *Config
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

func main() {
	wiki := New()
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

		wiki.config = GetConfig(tx)

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
		"EncodeID":   wiki.EncodeID,
	}).ParseGlob("templates/default/*"))

	wiki.tpl = tpl

	wiki.render = cajun.New()
	wiki.render.WikiLink = wiki
	wiki.policy = bluemonday.UGCPolicy()
	wiki.policy.AllowAttrs("class").Matching(regexp.MustCompile("empty-link")).OnElements("a")
	wiki.policy.RequireNoFollowOnLinks(false)

	wiki.gpolicy = &greentuesday.Policy{}
	wiki.gpolicy.Add = append(wiki.gpolicy.Add, greentuesday.AttrEle{Tag: "table", Attribute: html.Attribute{Key: "class", Val: "table"}})

	wiki.store = newMemoryStore()

	mainChain := alice.New(wiki.store.ContextClear)
	authChain := mainChain.Append(wiki.CheckAuth(AuthMember))
	adminChain := mainChain.Append(wiki.CheckAuth(AuthAdmin))

	wiki.router = mux.NewRouter()
	wiki.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/default"))))
	wiki.router.Handle("/", http.RedirectHandler("/Home", http.StatusMovedPermanently))
	wiki.router.Handle("/Setup", mainChain.ThenFunc(wiki.SetupFormHandler)).Methods("GET").Name("SetupForm")
	wiki.router.Handle("/Setup", mainChain.ThenFunc(wiki.SetupHandler)).Methods("POST").Name("Setup")
	wiki.router.Handle("/Login", mainChain.ThenFunc(wiki.LoginFormHandler)).Methods("GET").Name("LoginForm")
	wiki.router.Handle("/Login", mainChain.ThenFunc(wiki.LoginHandler)).Methods("POST").Name("Login")

	wiki.router.Handle("/favicon.ico", http.NotFoundHandler())

	wiki.router.Handle("/Admin", adminChain.ThenFunc(wiki.LoginHandler)).Methods("POST").Name("Admin")

	wiki.router.Handle("/{page:[^/]*}/edit", authChain.ThenFunc(wiki.EditHandler)).Methods("GET").Name("Edit")
	wiki.router.Handle("/{page:[^/]*}/history", mainChain.ThenFunc(wiki.HistoryHandler)).Methods("GET").Name("History")
	wiki.router.Handle("/{page:[^/]*}/version/{ver}", mainChain.ThenFunc(wiki.PageVersionHandler)).Methods("GET").Name("PageVersion")
	wiki.router.Handle("/{page:[^/]*}", mainChain.ThenFunc(wiki.PageHandler)).Methods("GET").Name("Read")
	wiki.router.Handle("/{page:[^/]*}", authChain.ThenFunc(wiki.UpdateHandler)).Methods("POST").Name("Update")

	return wiki
}

type PageNavData struct {
	Read    string
	Edit    string
	History string
	Section string
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

func (w *Wiki) PageNav(Slug string, Section string) PageNavData {
	return PageNavData{
		Read:    UrlToPath(w.router.Get("Read").URLPath("page", Slug)),
		Edit:    UrlToPath(w.router.Get("Edit").URLPath("page", Slug)),
		History: UrlToPath(w.router.Get("History").URLPath("page", Slug)),
		Section: Section,
	}
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

func (w *Wiki) PageHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			pagedata := page.Current.GetData(tx)
			unsafe, _ := w.render.Transform(string(pagedata))
			html := w.gpolicy.Massage(w.policy.Sanitize(unsafe))
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
			http.Redirect(rw, req, UrlToPath(w.router.Get("Edit").URLPath("page", vars["page"])), http.StatusTemporaryRedirect)
		}
		return nil
	})

}

func (w *Wiki) HistoryHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			rw.Header().Set("Content-Type", "text/html")

			data := struct {
				Name   string
				Slug   string
				Events []*Event
			}{
				vars["page"],
				vars["page"],
				page.History.Events,
			}

			if err := w.tpl.ExecuteTemplate(rw, "history.tpl", data); err != nil {
				fmt.Println(err)
			}
		} else {
			http.Redirect(rw, req, UrlToPath(w.router.Get("Edit").URLPath("page", vars["page"])), http.StatusTemporaryRedirect)
		}
		return nil
	})

}

func (w *Wiki) PageVersionHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.View(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		if page != nil {
			var pagedata []byte
			verid, _ := base64.URLEncoding.DecodeString(vars["ver"])
			for _, v := range page.History.Events {
				if bytes.Equal(v.DataID, verid) {
					pagedata = v.GetData(tx)
				}
			}
			if pagedata == nil {
				return nil
			}

			unsafe, _ := w.render.Transform(string(pagedata))
			html := w.gpolicy.Massage(w.policy.Sanitize(unsafe))
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
			http.Redirect(rw, req, UrlToPath(w.router.Get("Edit").URLPath("page", vars["page"])), http.StatusTemporaryRedirect)
		}
		return nil
	})

}

func (w *Wiki) UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.DB.Update(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		data := req.FormValue("data")
		data = strings.Replace(data, "\r\n", "\n", -1)

		key, _ := SaveData(tx, []byte(data))
		if page != nil {
			page.History.AddEvent(page.Current)
		} else {
			page = &Page{}
		}
		page.Current = Event{DateTime: time.Now(), DataID: key, IP: req.RemoteAddr}
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

func (w *Wiki) LoginFormHandler(rw http.ResponseWriter, req *http.Request) {
	if err := w.tpl.ExecuteTemplate(rw, "login.tpl", nil); err != nil {
		fmt.Println(err)
	}
}

func (w *Wiki) SetupFormHandler(rw http.ResponseWriter, req *http.Request) {
	if w.config.InitDone {
		http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusMovedPermanently)
		return
	}
	if err := w.tpl.ExecuteTemplate(rw, "setup.tpl", nil); err != nil {
		fmt.Println(err)
	}
}

func (w *Wiki) LoginHandler(rw http.ResponseWriter, req *http.Request) {
	w.DB.View(func(tx *bolt.Tx) error {
		u := GetUser(tx, req.FormValue("username"))
		if u != nil && u.CheckPassword(req.FormValue("password")) {
			s := w.store.Get(req)
			s.User = u
			w.store.Save(req, rw, s)
			http.Redirect(rw, req, s.PostLoginRedirect, http.StatusFound)
		} else {
			if err := w.tpl.ExecuteTemplate(rw, "login.tpl", nil); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
}

func (w *Wiki) SetupHandler(rw http.ResponseWriter, req *http.Request) {
	if w.config.InitDone {
		http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusMovedPermanently)
		return
	}
	w.DB.Update(func(tx *bolt.Tx) error {
		w.config.InitDone = true
		w.config.Save(tx)
		u := &User{}
		u.Name = req.FormValue("username")
		u.SetPassword(req.FormValue("password"))
		u.GiveAuth(AuthMember).GiveAuth(AuthModerator).GiveAuth(AuthAdmin)
		u.Save(tx)
		s := w.store.Get(req)
		s.User = u
		w.store.Save(req, rw, s)
		http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
		return nil
	})
}

func (w *Wiki) CheckAuth(auth Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			session := w.store.Get(req)
			if session.User != nil && session.User.HasAuth(auth) {
				next.ServeHTTP(rw, req)
			} else {
				session.PostLoginRedirect = req.URL.Path
				w.store.Save(req, rw, session)
				http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusTemporaryRedirect)
			}
		})
	}
}
