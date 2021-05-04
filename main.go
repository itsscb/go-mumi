package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

type ds struct {
	Datum    time.Time
	DatumStr string
	Menge    int
}

type tplData struct {
	Data  map[string]DB
	Datum string
	Zeit  string
	Summe int
}

type DB struct {
	Datensaetze []ds
	Date        string
	Summe       int
}

var DataStruct = map[string]DB{}

var tpl *template.Template
var layout = "02.01.2006 15:04"

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func sfbs(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/assets/bootstrap/css/bootstrap.min.css")
}

func sfbj(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/assets/bootstrap/js/bootstrap.min.js")
}

func sfnv(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/assets/css/Navigation-Clean.css")
}
func sfst(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/assets/css/styles.css")
}
func sfjq(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/assets/js/jquery.min.js")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/details", details)
	http.HandleFunc("/assets/bootstrap/css/bootstrap.min.css", sfbs)
	http.HandleFunc("/assets/bootstrap/js/bootstrap.min.js", sfbj)
	http.HandleFunc("/assets/css/Navigation-Clean.css", sfnv)
	http.HandleFunc("/assets/css/styles.css", sfst)
	http.HandleFunc("/assets/js/jquery.min.js", sfjq)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8081", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	var td tplData
	var m int
	var t time.Time
	var sum int = 0
	_, err := os.Stat("data.json")
	if !os.IsNotExist(err) {
		fdata, err := ioutil.ReadFile("data.json")
		checkErr(err)
		err = json.Unmarshal(fdata, &DataStruct)
		checkErr(err)
	}

	fm := req.FormValue("Menge")
	ft := req.FormValue("Date")
	fz := req.FormValue("Time")

	if len(fm) >= 1 {
		m, _ = strconv.Atoi(fm)
		t = string2time((ft + "T" + fz), "2006-01-02T15:04")
		if dt, ok := DataStruct[t.Format("2006-01-02")]; ok {
			dt.Datensaetze = append(dt.Datensaetze, ds{
				Datum:    t,
				DatumStr: t.Format(layout),
				Menge:    m,
			})
			dt.Summe = dt.Summe + m
			DataStruct[t.Format("2006-01-02")] = dt
		} else {
			DataStruct[t.Format("2006-01-02")] = DB{
				Datensaetze: []ds{{
					Datum:    t,
					DatumStr: t.Format(layout),
					Menge:    m,
				},
				},
				Date:  t.Format("2006-01-02"),
				Summe: m,
			}
		}
	} else {
		t = time.Now()
		fz = t.Format("15:04")
		m = 0
	}

	file, _ := json.MarshalIndent(DataStruct, "", " ")
	_ = ioutil.WriteFile("data.json", file, 0644)

	for _, d := range DataStruct[time.Now().Format("2006-01-02")].Datensaetze {
		sum += d.Menge
	}
	td = tplData{
		Data:  DataStruct,
		Datum: t.Format("2006-01-02"),
		Zeit:  fz,
		Summe: sum,
	}
	err = tpl.ExecuteTemplate(w, "index.gohtml", td)
	checkErr(err)
}

func details(w http.ResponseWriter, req *http.Request) {
	_, err := os.Stat("data.json")
	if !os.IsNotExist(err) {
		fdata, err := ioutil.ReadFile("data.json")
		checkErr(err)
		err = json.Unmarshal(fdata, &DataStruct)
		checkErr(err)
	}
	err = tpl.ExecuteTemplate(w, "details.gohtml", DataStruct)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func string2time(s string, layout string) time.Time {
	t, err := time.Parse(layout, s)
	checkErr(err)
	return t
}
