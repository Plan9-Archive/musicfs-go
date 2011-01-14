package main

import (
	"fmt"
	"os"
	"syscall"
	)

type PL struct {
	Pid int
}

var plchan = openchanpl()

func openchanpl() chan func(*PL) {
	ch := make(chan func(*PL),4)
	pl := new(PL)
	go func() {
		for {
			(<-ch)(pl)
		}
	}()
	return ch
}


func kill(pl *PL) {
	if pl.Pid > 0 {
		syscall.Kill(pl.Pid, syscall.SIGTERM)
		pl.Pid = -1
	}
}

func Stop() {
	plchan <- kill
}

const MusicPlayer = "/usr/bin/mpg123"
const Amixer = "/usr/bin/amixer"

func PlayId(id int) string {
	reply := make(chan *AudioFile)
	(*indexes)["Id"] <- func(mr Map) {
		m := mr.(MapI)
		reply <- m[uint64(id)]
	}
	af :=<- reply
	if af!=nil {
		PlayFile(af.FileName)
	}
	return ""
}

func PlayFile(file string) {
	plchan <- func(pl *PL) {
		kill(pl)
		var e os.Error
		pl.Pid,e = os.ForkExec(MusicPlayer, []string{MusicPlayer, file}, nil, "", nil)
		if e!=nil { panic("ForkExec: "+e.String()) }
	}
}

func VolumeModPercent(percent int) {
	var minus string
	if percent < 0 {
		minus = "-"
		percent = 0-percent
	}
	os.ForkExec(Amixer, []string{Amixer,"set","Master",fmt.Sprintf("%d%%%s",percent,minus)}, nil, "", nil)
}
