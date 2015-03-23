package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

func (w *Wiki) SetupFormHandler(rw http.ResponseWriter, req *http.Request) {
	if w.config.InitDone {
		http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusMovedPermanently)
		return
	}
	if err := w.tpl.ExecuteTemplate(rw, "setup.tpl", nil); err != nil {
		fmt.Println(err)
	}
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
