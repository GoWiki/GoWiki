// GoWiki project main.go
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

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
	theme   *Theme
	fb      *FormBuilder
}

var (
	GoLaunch = flag.Bool("GoLaunch", false, "Used for triggering GoLaunch functionality")
	Port     = flag.Int("port", 3000, "Port to serve the wiki on")
)

func main() {
	var socket net.Listener
	var db string

	server := &http.Server{}
	flag.Parse()

	db = os.Getenv("GOWIKIDB")
	if db == "" {
		db = "gowiki.db"
	}
	var err error
	if *GoLaunch {
		socket, _ = net.FileListener(os.NewFile(3, ""))
	} else {
		socket, err = net.Listen("tcp", fmt.Sprintf(":%d", *Port))
	}
	if err != nil {
		fmt.Println("Failed to open socket: ", err)
		os.Exit(1)
	}
	wiki := New(db)
	server.Handler = wiki.router
	server.Serve(socket)
}

func New(database string) *Wiki {
	wiki := &Wiki{}
	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		SetupBuckets(tx)

		wiki.config = GetConfig(tx)

		if !wiki.config.FilesLoaded {
			InitThemes(tx)
			wiki.config.FilesLoaded = true
			wiki.config.Save(tx)
		}

		return nil
	})
	wiki.DB = db

	if err != nil {
		fmt.Println(err)
		return nil
	}

	theme := &Theme{Name: "default"}
	wiki.theme = theme
	tpl := template.New("default").Funcs(template.FuncMap{
		"PageNav":    wiki.PageNav,
		"Route":      wiki.Route,
		"GetContent": wiki.GetContent,
		"EncodeID":   wiki.EncodeID,
	})
	db.View(func(tx *bolt.Tx) error {
		theme.ParseTemplates(tx, tpl)
		return nil
	})

	wiki.tpl = tpl

	wiki.fb = NewFormBuilder(wiki)

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
	wiki.router.PathPrefix("/static/").Handler(mainChain.ThenFunc(wiki.StaticHandler))
	wiki.router.Handle("/", http.RedirectHandler("/Home", http.StatusMovedPermanently))
	wiki.router.Handle("/Setup", mainChain.ThenFunc(wiki.SetupFormHandler)).Methods("GET").Name("SetupForm")
	wiki.router.Handle("/Setup", mainChain.ThenFunc(wiki.SetupHandler)).Methods("POST").Name("Setup")
	wiki.router.Handle("/Login", mainChain.ThenFunc(wiki.LoginFormHandler)).Methods("GET").Name("LoginForm")
	wiki.router.Handle("/Login", mainChain.ThenFunc(wiki.LoginHandler)).Methods("POST").Name("Login")
	wiki.router.Handle("/Create", mainChain.ThenFunc(wiki.UserCreateFormHandler)).Methods("GET").Name("UserCreateForm")
	wiki.router.Handle("/Create", mainChain.ThenFunc(wiki.UserCreateHandler)).Methods("POST").Name("UserCreate")
	wiki.router.Handle("/Logout", authChain.ThenFunc(wiki.LogoutHandler)).Methods("GET").Name("Logout")

	wiki.router.Handle("/favicon.ico", http.NotFoundHandler())

	wiki.router.Handle("/Admin", adminChain.ThenFunc(wiki.LoginHandler)).Methods("POST").Name("Admin")

	wiki.router.Handle("/{page:[^/]*}/edit", authChain.ThenFunc(wiki.EditHandler)).Methods("GET").Name("Edit")
	wiki.router.Handle("/{page:[^/]*}/history", mainChain.ThenFunc(wiki.HistoryHandler)).Methods("GET").Name("History")
	wiki.router.Handle("/{page:[^/]*}/version/{ver}", mainChain.ThenFunc(wiki.PageVersionHandler)).Methods("GET").Name("PageVersion")
	wiki.router.Handle("/{page:[^/]*}", mainChain.ThenFunc(wiki.PageHandler)).Methods("GET").Name("Read")
	wiki.router.Handle("/{page:[^/]*}", authChain.ThenFunc(wiki.UpdateHandler)).Methods("POST").Name("Update")

	setupform := wiki.fb.NewForm("SetupForm")
	setupform.NewString("Username", "Username", "Username", "text")
	setupform.NewPassword("Password", "Password", "Password")
	setupform.NewButtons().AddButton("Finish Setup", "", "primary")

	loginform := wiki.fb.NewForm("LoginForm")
	loginform.NewString("Username", "Username", "Username", "text")
	loginform.NewPassword("Password", "Password", "Password")
	buttons := loginform.NewButtons()
	buttons.AddButton("Login", "Login", "primary")
	buttons.AddButton("Create Account", "Create", "default").Link(wiki.Route("UserCreate"))

	usercreateform := wiki.fb.NewForm("UserCreateForm")
	usercreateform.NewString("Username", "Username", "Username", "text")
	usercreateform.NewString("Email", "Email", "Email", "email")
	usercreateform.NewPassword("Password", "Password", "Password")
	usercreatebuttons := usercreateform.NewButtons()
	usercreatebuttons.AddButton("Create Account", "Create", "default")

	return wiki
}

func (w *Wiki) StaticHandler(rw http.ResponseWriter, req *http.Request) {
	w.DB.View(func(tx *bolt.Tx) error {
		name := strings.TrimPrefix(req.RequestURI, "/static/")
		data := w.theme.GetFile(tx, name)
		if strings.HasSuffix(name, ".css") {
			rw.Header().Add("content-type", "text/css")
		}
		rw.Write(data)
		return nil
	})

}

func (w *Wiki) LoginFormHandler(rw http.ResponseWriter, req *http.Request) {
	form := w.fb.GetForm("LoginForm")

	data := struct {
		Name     string
		FormName string
		Form     template.HTML
	}{
		"Login",
		"Login",
		form.Render(nil, w.Route("Login"), "POST"),
	}

	if err := w.tpl.ExecuteTemplate(rw, "form.tpl", data); err != nil {
		fmt.Println(err)
	}
}

func (w *Wiki) LoginHandler(rw http.ResponseWriter, req *http.Request) {
	form := w.fb.GetForm("LoginForm")
	data := struct {
		Username string
		Password string
	}{}
	form.Parse(req.FormValue, &data)
	w.DB.View(func(tx *bolt.Tx) error {
		u := GetUser(tx, data.Username)

		if u != nil && u.CheckPassword(data.Password) {
			s := w.store.Get(req)
			s.User = u
			w.store.Save(req, rw, s)
			http.Redirect(rw, req, s.PostLoginRedirect, http.StatusFound)
		} else {
			w.LoginFormHandler(rw, req)
		}
		return nil
	})
}

func (w *Wiki) UserCreateFormHandler(rw http.ResponseWriter, req *http.Request) {
	form := w.fb.GetForm("UserCreateForm")

	data := struct {
		Name     string
		FormName string
		Form     template.HTML
	}{
		"Create User",
		"Create User",
		form.Render(nil, w.Route("UserCreate"), "POST"),
	}

	if err := w.tpl.ExecuteTemplate(rw, "form.tpl", data); err != nil {
		fmt.Println(err)
	}
}

func (w *Wiki) UserCreateHandler(rw http.ResponseWriter, req *http.Request) {
	form := w.fb.GetForm("UserCreateForm")
	data := struct {
		Username string
		Email    string
		Password string
	}{}
	form.Parse(req.FormValue, &data)

	w.DB.Update(func(tx *bolt.Tx) error {
		u := GetUser(tx, data.Username)
		if u == nil {
			u := &User{Name: data.Username, Email: data.Email}
			u.SetPassword(data.Password)
			u.GiveAuth(AuthMember)
			u.Save(tx)
			s := w.store.Get(req)
			s.User = u
			w.store.Save(req, rw, s)
			http.Redirect(rw, req, "/", http.StatusFound)
		} else {
			w.LoginFormHandler(rw, req)
		}
		return nil
	})

}

func (w *Wiki) LogoutHandler(rw http.ResponseWriter, req *http.Request) {
	s := w.store.Get(req)
	w.store.Destroy(req, rw, s)
	http.Redirect(rw, req, "/", http.StatusFound)
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
