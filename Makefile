GO = go
GO_LDFLAGS = -s -w

all: 	clean phome

phome:
	$(GO) build $(GOFLAGS) -ldflags "$(GO_LDFLAGS)"

clean:
	rm -f phome

.PHONY: all clean