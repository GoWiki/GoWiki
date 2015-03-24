package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/boltdb/bolt"
)

func (w *Wiki) SetupFormHandler(rw http.ResponseWriter, req *http.Request) {
	if w.config.InitDone {
		http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusFound)
		return
	}
	form := NewForm(w.tpl)
	form.NewString("Username", "Username", "Username")
	form.NewPassword("Password", "Password", "Password")
	form.NewButtons().AddButton("Finish Setup", "", "primary")

	data := struct {
		Name     string
		FormName string
		Form     template.HTML
	}{
		"Initial Setup",
		"Initial Setup",
		form.Render(nil),
	}

	if err := w.tpl.ExecuteTemplate(rw, "form.tpl", data); err != nil {
		fmt.Println(err)
	}
}

func (w *Wiki) SetupHandler(rw http.ResponseWriter, req *http.Request) {
	if w.config.InitDone {
		http.Redirect(rw, req, UrlToPath(w.router.Get("LoginForm").URLPath()), http.StatusFound)
		return
	}
	w.DB.Update(func(tx *bolt.Tx) error {
		w.config.InitDone = true
		w.config.Save(tx)
		u := &User{}
		form := NewForm(w.tpl)
		form.NewString("Username", "Username", "Username")
		form.NewPassword("Password", "Password", "Password")
		form.NewButtons().AddButton("Finish Setup", "", "primary")

		data := struct {
			Username string
			Password string
		}{}

		form.Parse(req.FormValue, &data)

		u.Name = data.Username
		u.SetPassword(data.Password)

		u.GiveAuth(AuthMember).GiveAuth(AuthModerator).GiveAuth(AuthAdmin)
		u.Save(tx)
		s := w.store.Get(req)
		s.User = u
		w.store.Save(req, rw, s)
		http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
		return nil
	})
}
