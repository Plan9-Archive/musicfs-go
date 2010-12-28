package main

import (
	"bitbucket.org/taruti/taglib.go"
	"go9p.googlecode.com/hg/p/srv"

	)

type AudioFile struct {
	taglib.Tags
	srv.File
	FileName, Suffix string
}


type TitleR []*AudioFile
func (r TitleR)LessThan(other interface{}) bool {
	return r[0].Title < other.(TitleR)[0].Title
}

type Query func(Map)
type Index chan Query


func indexLoop(ch chan Query, m Map) {
	for {
		f :=<- ch
		if nil==f { return }
		f(m)
	}
}

func spawnIndexLoop(tree Map) Index {
	ch := make(chan Query, 1)
	go indexLoop(ch,tree)
	return Index(ch)
}

func buildAudioIndexes() {
	if indexes == nil { indexes = &map[string]Index{} }

//	ititle := spawnIndexLoop(createIndexAF(w.ais, func(a *AudioFile)interface{}{return a.Title}))
//	_ = ititle
	iartist := spawnIndexLoop(Map{})
	(*indexes)["Artist"] = iartist
}

func audioAddToIndexes(af *AudioFile) {
	(*indexes)["Artist"] <- func(m Map) {
		if af.Artist == "" { af.Artist = "@" }
		addToMap(af.Artist, af, m) 
	}
}

var indexes *map[string]Index
