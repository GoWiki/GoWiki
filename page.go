package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

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
				User    UserInfoType
			}{
				template.HTML(html),
				vars["page"],
				vars["page"],
				w.UserInfo(req),
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
			page.History.LoadUsers(tx)
			data := struct {
				Name   string
				Slug   string
				Events []*Event
				User   UserInfoType
			}{
				vars["page"],
				vars["page"],
				page.History.Events,
				w.UserInfo(req),
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
				User    UserInfoType
			}{
				template.HTML(html),
				vars["page"],
				vars["page"],
				w.UserInfo(req),
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
		s := w.store.Get(req)
		page, _ := GetPage(tx, vars["page"])
		data := req.FormValue("data")
		data = strings.Replace(data, "\r\n", "\n", -1)

		key, _ := SaveData(tx, []byte(data))
		if page != nil {
			page.History.AddEvent(page.Current)
		} else {
			page = &Page{}
		}
		page.Current = Event{DateTime: time.Now(), DataID: key, IP: req.RemoteAddr, AuthorID: s.User.ID}
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
				User    UserInfoType
			}{
				string(pagedata),
				vars["page"],
				vars["page"],
				w.UserInfo(req),
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
				User    UserInfoType
			}{
				"",
				vars["page"],
				vars["page"],
				w.UserInfo(req),
			}
			if err := w.tpl.ExecuteTemplate(rw, "edit.tpl", data); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
}
