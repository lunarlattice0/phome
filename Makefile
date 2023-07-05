GO = go
GO_LDFLAGS = -s -w

all: 	clean phome

phome:
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)"

run:	all
	./phome

clean:
	rm -f phome

.PHONY: all clean
