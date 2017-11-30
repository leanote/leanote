include $(GOROOT)/src/Make.inc

TARG=cascadia

GOFILES= \
	parser.go \
	selector.go \
	
include $(GOROOT)/src/Make.pkg

format:
	gofmt -w ${GOFILES} *_test.go
