include $(GOROOT)/src/Make.inc

TARG=github.com/fzzbt/neste
GOFILES=\
	manager.go\
	template.go\
	formatter.go\

include $(GOROOT)/src/Make.pkg
