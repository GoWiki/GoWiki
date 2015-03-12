// GoWiki project main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var (
	DB *bolt.DB
)

type PageInfo struct {
	DataID []byte
}

func main() {
	db, err := bolt.Open("gowiki.db", 0600, nil)
	DB = db
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("pages"))
		pages := tx.Bucket([]byte("pages"))
		pages.CreateBucketIfNotExists([]byte("data"))
		pages.CreateBucketIfNotExists([]byte("history"))
		pages.CreateBucketIfNotExists([]byte("names"))
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
		pages := tx.Bucket([]byte("pages"))
		pageinfo := pages.Bucket([]byte("names")).Get([]byte(vars["page"]))
		if pageinfo != nil {
			var pi PageInfo
			json.Unmarshal(pageinfo, &pi)
			pagedata := pages.Bucket([]byte("data")).Get(pi.DataID)
			rw.Write(pagedata)
		}
		return nil
	})

}

func UpdateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	DB.Update(func(tx *bolt.Tx) error {
		pages := tx.Bucket([]byte("pages"))
		names := pages.Bucket([]byte("names"))
		pageinfo := names.Get([]byte(vars["page"]))
		var pi PageInfo
		if pageinfo != nil {
			json.Unmarshal(pageinfo, &pi)
		}

		data := pages.Bucket([]byte("data"))
		dataid := NextKey(data)

		pi.DataID = dataid

		data.Put(dataid, []byte(req.FormValue("data")))

		pageinfo, _ = json.Marshal(&pi)

		names.Put([]byte(vars["page"]), pageinfo)
		return nil
	})
	PageHandler(rw, req)
}

func NextKey(b *bolt.Bucket) []byte {
	i, _ := b.NextSequence()
	key, _ := json.Marshal(i)
	return key
}

func EditHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rw.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(rw, "<form method=\"POST\" action=\"/%s\"><textarea name=\"data\"></textarea><input type=\"submit\"></form>", vars["page"])
}
