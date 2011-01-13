include $(GOROOT)/src/Make.inc

TARG=musicfs
GOFILES=musicfs.go index.go map.go http.go playlocal.go

include $(GOROOT)/src/Make.cmd
