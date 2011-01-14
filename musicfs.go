package main

import (
	"bitbucket.org/taruti/taglib.go"
	"flag"
	"fmt"
	"go9p.googlecode.com/hg/p"
	"go9p.googlecode.com/hg/p/srv"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var addr = flag.String("addr", ":5640", "network address")
var root = flag.String("root", ".", "root of filesystem to look for files")
var debug = flag.Int("debug", 0, "debug level")

type walker struct {
	matching map[string]int
	ignored  map[string]int
	//	ais         []AudioFile
	nfiles int
}

func newwalker() *walker {
	w := walker{map[string]int{}, map[string]int{}, 0}
	w.matching["mp3"] = 0
	w.matching["m4a"] = 0
	w.matching["ogg"] = 0
	w.matching["flac"] = 0
	return &w
}

func (*walker) VisitDir(string, *os.FileInfo) bool { return true }
func (w *walker) VisitFile(path string, fi *os.FileInfo) {
	w.nfiles++
	if w.nfiles%1000 == 0 {
		log.Println("Processed", w.nfiles, "files")
	}

	var suffix string
	idx := strings.LastIndex(path, ".")
	if idx > 0 {
		suffix = strings.ToLower(path[idx+1:])
	}
	val, present := w.matching[suffix]
	if !present {
		nval := w.ignored[suffix]
		w.ignored[suffix] = nval + 1
		return
	}
	w.matching[suffix] = val + 1

	// Collect Metadata 
	ai := taglib.GetTags(path)
	af := new(AudioFile)
	af.Tags = *ai
	af.FileName = path
	af.Suffix = suffix
	af.Length = uint64(fi.Size)
	audioAddToIndexes(af)
	// FIXME
	//	w.ais = append(w.ais, AudioFile{})

}
func (w *walker) PrintStats() {
	fmt.Printf("Matches:\n")
	for k, v := range w.matching {
		fmt.Printf("\t%-5s%6d\n", k, v)
	}
	fmt.Printf("Ignored:\n")
	for k, v := range w.ignored {
		fmt.Printf("\t%-5s%6d\n", k, v)
	}
}

type user string

func (u user) Id() int           { return 0 }
func (u user) Name() string      { return string(u) }
func (u user) Groups() []p.Group { return []p.Group{u} }
func (u user) Members() []p.User { return []p.User{u} }

type MapSIStats struct {
	srv.File
	m map[string]int
}

func (w *MapSIStats) Read(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	var s string
	for k, v := range w.m {
		s += fmt.Sprintf("\t%-5s%6d\n", k, v)
	}

	b := []byte(s)
	w.Length = uint64(len(b))
	start := int(offset)
	end := len(buf)
	if end > len(b) {
		end = len(b)
	}

	copy(buf, b[start:end])
	return end - start, nil
}

type IndexDir struct {
	srv.File
	Index
}

func (id *IndexDir) Populate(key string, u user) {
	log.Println("Populate key:", key)
	rc := make(chan []string)
	id.Index <- func(m Map) { rc <- mapKeysS(m) }
	for _, aname := range <-rc {
		ir := new(IndexRes)
		ir.F = func(m Map) []*AudioFile {
			switch m := m.(type) {
			case MapSM:
				return m[aname]
			case MapI:
				d,_ := strconv.Atoi64(aname)
				return []*AudioFile{m[uint64(d)]}
			}
			return []*AudioFile{}
		}
		ir.Index = id.Index
		ir.Add(&id.File, aname, u, u, p.DMDIR|0555, ir)
		ir.Populate(u)
	}
}

type IndexRes struct {
	srv.File
	Index
	F func(Map) []*AudioFile
}

func (id *IndexRes) Populate(u user) {
	rc := make(chan []*AudioFile)
	id.Index <- func(m Map) { rc <- id.F(m) }
	for _, af := range <-rc {
		fname := fmt.Sprintf("%s_%02d_%s_%d.%s", af.Album, af.Track, af.Title, af.Year, af.Suffix)
		af.Add(&id.File, fname, u, u, 0555, af)
	}
}

func (af *AudioFile) Read(fid *srv.FFid, buf []byte, offset uint64) (int, *p.Error) {
	f, e := os.Open(af.FileName, os.O_RDONLY, 0)
	if e != nil {
		return 0, &p.Error{e.String(), -1}
	}
	defer f.Close()
	n, e := f.ReadAt(buf, int64(offset))
	if e != nil {
		return 0, &p.Error{e.String(), -1}
	}
	return n, nil
}


func main() {
	flag.Parse()

	log.Println("Building indexes")
	buildAudioIndexes()
	log.Println("Creating filesystem")
	var err *p.Error
	log.Println("Starting walker")
	w := newwalker()
	path.Walk(*root, w, nil)
	w.PrintStats()

	u := user("root")
	root := new(srv.File)
	err = root.Add(nil, "/", u, u, p.DMDIR|0555, nil)
	if err != nil { goto error }
	// Status directory
	stats := new(srv.File)
	stats.Add(root, "stats", u, u, p.DMDIR|0555, stats)
	ms := new(MapSIStats)
	ms.m = w.matching
	err = ms.Add(stats, "matched", u, u, 0444, ms)
	if err != nil {
		goto error
	}
	is := new(MapSIStats)
	is.m = w.ignored
	err = is.Add(stats, "ignored", u, u, 0444, is)
	if err != nil { goto error }
	// Index directories
	log.Println("Top-level indexes", len(*indexes))
	for k, v := range *indexes {
		if k == "Id" { continue }
		d := new(IndexDir)
		d.Index = v
		d.Add(root, k, u, u, p.DMDIR|0555, d)
		d.Populate(k, u)
	}
	go httpmain()
	s := srv.NewFileSrv(root)
	s.Dotu = true
	s.Debuglevel = *debug
	s.Start(s)
	log.Println("Starting listener")
	s.StartNetListener("tcp", *addr)
	return
error:
	log.Println(fmt.Sprintf("Error: %s %d", err.Error, err.Errornum))
}
