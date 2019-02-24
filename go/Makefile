cas_server: cas_server.go utils.go login.go payment.go tc_wrap.h tc_wrap.c
	CGO_LDFLAGS_ALLOW=".*" go build -o cas_server

test: cas_server
	CGO_LDFLAGS_ALLOW=".*" go test

clean:
	rm cas_server
