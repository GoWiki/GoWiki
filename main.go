// GoWiki project main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var (
	DB *bolt.DB
)

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
			rw.Write(html)
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
	rw.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(rw, "<form method=\"POST\" action=\"/%s\"><textarea name=\"data\"></textarea><input type=\"submit\"></form>", vars["page"])
}
