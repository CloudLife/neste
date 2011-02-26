A=6
FILES=http_dirnames
GOFILES=$(wildcard $(FILES).go)

all: $(GOFILES:.go=)

$(GOFILES:.go=): %: %.$(A)
	$(A)l -o $@ $<

$(GOFILES:.go=.$(A)): %.$(A): %.go
	$(A)g -o $@ $<
	
	
.PHONY: clean

clean:
	rm -f *.a *.$(A) $(FILES)
