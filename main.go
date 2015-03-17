// GoWiki project main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var (
	DB  *bolt.DB
	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/default/*"))
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

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/default"))))
	r.HandleFunc("/{page:.*}/edit", EditHandler).Methods("GET")
	r.HandleFunc("/{page:.*}", PageHandler).Methods("GET")
	r.HandleFunc("/{page:.*}", UpdateHandler).Methods("POST")

	http.ListenAndServe(":3000", r)
	fmt.Println("Hello World!")
}

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "Test")
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
			}{
				template.HTML(string(html)),
				vars["page"],
			}

			tpl.ExecuteTemplate(rw, "view.tpl", data)
		}
		return nil
	})

}

func UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Println(DB.Update(func(tx *bolt.Tx) error {
		page, _ := GetPage(tx, vars["page"])
		key, _ := SaveData(tx, []byte(req.FormValue("data")))
		if page != nil {
			//page.History.Events = append(page.History.Events, page.Current)
		} else {
			page = &Page{}
		}
		page.Current = Event{DataID: key, IP: req.RemoteAddr}
		fmt.Println(page.Save(tx, vars["page"]))
		return nil
	}))
	PageHandler(rw, req)
}

func EditHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rw.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(rw, "<form method=\"POST\" action=\"/%s\"><textarea name=\"data\"></textarea><input type=\"submit\"></form>", vars["page"])
}
