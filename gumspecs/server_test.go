package gumspecs

import (
	"os"
	"testing"
)

func TestReadServer(t *testing.T) {
	os.Clearenv()

	os.Setenv("HTTP_HOST", "127.0.0.1")
	os.Setenv("HTTP_PORT", "9090")

	srv := ReadHTTPServer()

	exp := "127.0.0.1:9090"
	if exp != srv.String() {
		t.Fatalf("Expect %v was %v", exp, srv)
	}

}
