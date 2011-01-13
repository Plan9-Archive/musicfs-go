package main

import (
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
	}
}

func Pause() {
	plchan <- kill
}

const MusicPlayer = "/usr/bin/mpg123"

func Play(file string) {
	plchan <- func(pl *PL) {
		kill(pl)
		var e os.Error
		pl.Pid,e = os.ForkExec(MusicPlayer, []string{MusicPlayer, file}, nil, "", nil)
		if e!=nil { panic("ForkExec: "+e.String()) }
	}
}
