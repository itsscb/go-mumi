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
	Date    time.Time
	DateStr string
	Amount    int
}

type tplData struct {
	Data  []DB
	Date string
	Time  string
	Sum int
}

type DB struct {
	Datasets []ds
	Date        string
	Sum       int
}

var DataStruct = map[string]DB{}

var tpl *template.Template
var layout = "02.01.2006 15:04"
var laydt = "2006-01-02"

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
	http.ListenAndServe(":8080", nil)
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

	fm := req.FormValue("Amount")
	ft := req.FormValue("Date")
	fz := req.FormValue("Time")

	for _, fo := range DataStruct {
		for i, fi := range fo.Datasets {
			form := req.FormValue(fi.DateStr +"@" +strconv.Itoa(fi.Amount))
			ddate := fi.Date.Format(laydt) 
			if len(form) >= 1 {
				tmpslice := DataStruct[ddate]
				if len(tmpslice.Datasets) >= 2 {
					tmpslice.Datasets = append(tmpslice.Datasets[:i], tmpslice.Datasets[i+1:]...)
				} else if len(tmpslice.Datasets) == i {
					tmpslice.Datasets = tmpslice.Datasets[:i]
				} else {
					tmpslice.Datasets = tmpslice.Datasets[i+1:]
				}
				if len(tmpslice.Datasets) != 0 {
					tmpslice.Sum -= fi.Amount
					DataStruct[ddate] = tmpslice
				} else {
					delete(DataStruct, ddate)
				}
				file, _ := json.MarshalIndent(DataStruct, "", " ")
				_ = ioutil.WriteFile("data.json", file, 0644)
				break
			}
		}
	}

	if len(fm) >= 1 {
		m, _ = strconv.Atoi(fm)
		t = string2time((ft + "T" + fz), "2006-01-02T15:04")
		if dt, ok := DataStruct[t.Format(laydt)]; ok {
			dt.Datasets = append(dt.Datasets, ds{
				Date:    t,
				DateStr: t.Format(layout),
				Amount:    m,
			})
			dt.Sum = dt.Sum + m
			DataStruct[t.Format(laydt)] = dt
		} else {
			DataStruct[t.Format(laydt)] = DB{
				Datasets: []ds{{
					Date:    t,
					DateStr: t.Format(layout),
					Amount:    m,
				},
				},
				Date:  t.Format(laydt),
				Sum: m,
			}
		}
		file, _ := json.MarshalIndent(DataStruct, "", " ")
		_ = ioutil.WriteFile("data.json", file, 0644)
	} 

	t = time.Now()
		fz = t.Format("15:04")
		m = 0
	

	for _, d := range DataStruct[time.Now().Format(laydt)].Datasets {
		sum += d.Amount
	}
	var TplStruct []DB
	TplStruct = append(TplStruct, DataStruct[time.Now().Format(laydt)])
	td = tplData{
		Data:  TplStruct,
		Date: t.Format(laydt),
		Time:  fz,
		Sum: sum,
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
	for _, fo := range DataStruct {
		for i, fi := range fo.Datasets {
			form := req.FormValue(fi.DateStr +"@" +strconv.Itoa(fi.Amount))
			ddate := fi.Date.Format(laydt) 
			if len(form) >= 1 {
				tmpslice := DataStruct[ddate]
				if len(tmpslice.Datasets) >= 2 {
					tmpslice.Datasets = append(tmpslice.Datasets[:i], tmpslice.Datasets[i+1:]...)
				} else if len(tmpslice.Datasets) == i {
					tmpslice.Datasets = tmpslice.Datasets[:i]
				} else {
					tmpslice.Datasets = tmpslice.Datasets[i+1:]
				}
				if len(tmpslice.Datasets) != 0 {
					tmpslice.Sum -= fi.Amount
					DataStruct[ddate] = tmpslice
				} else {
					delete(DataStruct, ddate)
				}
				file, _ := json.MarshalIndent(DataStruct, "", " ")
				_ = ioutil.WriteFile("data.json", file, 0644)
				break
			}
		}
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
