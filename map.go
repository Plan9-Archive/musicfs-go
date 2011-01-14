package main

import "fmt"

type Map interface{}
type MapSM map[string][]*AudioFile
type MapI  map[uint64]*AudioFile

func addToMapSM(key string, af *AudioFile, mr Map) {
	m := mr.(MapSM)
	old, present := m[key]
	if !present {
		m[key] = []*AudioFile{af}
	} else {
		m[key] = append(old, af)
	}
}
func addToMapI(key uint64, af *AudioFile, mr Map) {
	m := mr.(MapI)
	m[key] = af
}

func mapKeysS(m Map) []string {
	switch m:=m.(type) {
	case MapI:
		res := make([]string, len(m))
		i := 0
		for k, _ := range m {
			res[i] = fmt.Sprint(k)
			i++
		}
		return res
	case MapSM:
		res := make([]string, len(m))
		i := 0
		for k, _ := range m {
			res[i] = k
			i++
		}
		return res
	}
	return []string{}
}

