package main

import "fmt"

type Map map[interface{}][]*AudioFile

func addToMap(key interface{}, af *AudioFile, m Map) {
	old, present := m[key]
	if !present {
		m[key] = []*AudioFile{af}
	} else {
		m[key] = append(old, af)
	}
}

func mapKeysS(m Map) []string {
	res := make([]string, len(m))
	i := 0
	for k, _ := range m {
		res[i] = fmt.Sprint(k)
		i++
	}
	return res
}
