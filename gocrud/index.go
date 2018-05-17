package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"./mycard"
	_ "github.com/go-sql-driver/mysql"
)

type ResultData struct {
	No    int
	Fname string
	Lname string
}
type Test struct {
	ID   string
	Name string
}
type Person struct {
	ID        string `json:"id"`
	THprefix  string `json:"thprefix"`
	THfname   string `json:"thfname"`
	THlname   string `json:"thlname"`
	ENprefix  string `json:"enprefix"`
	ENfname   string `json:"enfname"`
	ENlname   string `json:"enlname"`
	Addr      string `json:"addr"`
	Birthdate string `json:"birthdate"`
	Age       string `json:"age"`
	Sex       string `json:"sex"`
}

var templates = template.Must(template.ParseFiles("index.html"))

//var tmpl = template.Must(template.ParseGlob("tmpl/*"))
func readjson(res http.ResponseWriter, req *http.Request) {
	p := Test{"123456", "Songwut"}
	ijson, _ := json.Marshal(p)
	fmt.Fprintln(res, string(ijson))
}
func readcard(res http.ResponseWriter, req *http.Request) {
	//return mycard.ReadCard()
	p := mycard.ReadCard()
	//a := Person{"a", "b", "c", "d", "e", "g", "g", "h", "i", "10", "1"}
	pjson, _ := json.Marshal(p)
	//fmt.Fprintln(res, p)
	fmt.Fprintln(res, string(pjson))

	//fmt.Println(mycard.ReadCard())
	//templates.Execute(res, mycard.ReadCard())
}
func delete(res http.ResponseWriter, req *http.Request) {
	var db, err = sql.Open("mysql", "root:12345678@/mywind")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("DELETE FROM employees WHERE emp_no=?")
	stmt.Exec(req.URL.Query().Get("no"))
	if err != nil {
		panic(err)
	}
	http.Redirect(res, req, "/", 301)
}
func result(res http.ResponseWriter, req *http.Request) {
	var db, err = sql.Open("mysql", "root:12345678@/mywind")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT emp_no,first_name,last_name FROM employees")
	if err != nil {
		panic(err)
	}
	iData := ResultData{}
	var results []ResultData
	for rows.Next() {
		var no int
		var fname, lname string
		err = rows.Scan(&no, &fname, &lname)
		iData.No = no
		iData.Fname = fname
		iData.Lname = lname
		results = append(results, iData)
		if err != nil {
			panic(err)
		}
	}
	templates.Execute(res, results)
	//tmpl.ExecuteTemplate(res, "Index", results)
	fmt.Println(results)
}
func main() {
	http.HandleFunc("/", result)
	http.HandleFunc("/readcard", readcard)
	http.HandleFunc("/readjson", readjson)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":3000", nil)
}
