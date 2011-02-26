include $(GOROOT)/src/Make.inc

TARG=bitbucket.org/fzzbt/neste
GOFILES=\
	manager.go\
	template.go\
	formatter.go\

include $(GOROOT)/src/Make.pkg
