include $(GOROOT)/src/Make.inc

TARG=musicfs
GOFILES=musicfs.go index.go map.go http.go

include $(GOROOT)/src/Make.cmd
