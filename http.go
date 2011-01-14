package main

import (
	"bitbucket.org/taruti/http_jsonrpc.go"
	"bitbucket.org/taruti/taglib.go"
    "http"
    "os"
	"regexp"
)

type afhtml struct {
	taglib.Tags
	Id uint64
}

func Search(regex string) ([]afhtml,os.Error) {
	re,err := regexp.Compile(regex)
	if err!=nil {
		return nil,err
	}
	ch := make(chan []afhtml)
	(*indexes)["Artist"] <- func(mr Map) {
		res := []afhtml{}
		m := mr.(MapSM)
		for artist,afs := range m {
			if re.MatchString(artist) {
				for i:=0; i<len(afs); i++ {
					res = append(res, afhtml{afs[i].Tags, afs[i].Path})
				}
			}
		}
		ch <- res
	}
	return <-ch,nil
}

func httpmain() {
	s := http_jsonrpc.New()
	s.Register("search", Search)
	s.Register("play", PlayId)
	s.Register("stop", Stop)
	s.Register("volume", VolumeModPercent)
	http.Handle("/post", s)
	http.Handle("/", http.FileServer("static", ""))
	http.ListenAndServe(":8080", nil)
}
