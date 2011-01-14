package main

import (
	"bitbucket.org/taruti/http_jsonrpc.go"
    "http"
    "io/ioutil"
    "os"
)

type page struct {
    title string
    body  []byte
}

func (p *page) save() os.Error {
    filename := p.title + ".txt"
    return ioutil.WriteFile(filename, p.body, 0600)
}

func loadPage(title string) (*page, os.Error) {
    filename := "test.html"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &page{title: title, body: body}, nil
}

const lenPath = len("/view/")

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    p, _ := loadPage(title)
    w.Write(p.body)
	switch r.Method {
	case "GET":
	case "POST":
	default:
	}
}

func httpmain() {
	s := http_jsonrpc.New()
	s.Register("play", PlayId)
	s.Register("stop", Stop)
	s.Register("volume", VolumeModPercent)
	http.Handle("/post", s)
	http.Handle("/", http.FileServer("static", ""))
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
