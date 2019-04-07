build_all: deps
	GOPATH=$(CURDIR)/.go CGO_LDFLAGS_ALLOW=".*" go build -o cas_server

test: deps
	GOPATH=$(CURDIR)/.go CGO_LDFLAGS_ALLOW=".*" go test

deps: cas_server.go utils.go login.go payment.go tc_wrap.h tc_wrap.c
	rm -Rf $(CURDIR)/.go/src/github.com/calblueprint/fsf-backend/go
	mkdir -p $(CURDIR)/.go/src/github.com/calblueprint/fsf-backend/go
	ln -s $(CURDIR) $(CURDIR)/.go/src/github.com/calblueprint/fsf-backend/go
	GOPATH=$(CURDIR)/.go go get -d

clean:
	rm cas_server
	rm -Rf $(CURDIR)/.go/src/github.com/calblueprint/fsf-backend/go
