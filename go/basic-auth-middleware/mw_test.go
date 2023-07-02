package basicauthmiddleware_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"

	basicauthmiddleware "github.com/michurin/warehouse/go/basic-auth-middleware"
)

func ExampleAuthBasic() {
	nakedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK: user=" + basicauthmiddleware.UserName(r.Context())))
	})

	hmacKey := []byte{1, 2, 3, 4}
	user := "one"
	password := []byte("secret")
	passwd := map[string][]byte{
		user: hmac.New(sha256.New, hmacKey).Sum(password),
	}

	checker := basicauthmiddleware.StaticAuth(passwd, hmacKey)

	wrappedHandler := basicauthmiddleware.AuthBasic(nakedHandler, "test", checker)

	server := httptest.NewServer(wrappedHandler)
	defer server.Close()

	url := server.URL
	curl := func(args ...string) {
		fmt.Println("curl " + strings.Join(append(args, "http://testserver/"), " "))
		cmd := exec.Command("curl", append([]string{"-qs"}, append(args, url)...)...) //nolint:gosec // -q must be first arg
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	curl()
	curl("-u", "xxx:secret")
	curl("-u", "one:xxxxxx")
	curl("-u", "one:secret")

	// output:
	// curl http://testserver/
	// Unauthorized
	// curl -u xxx:secret http://testserver/
	// Unauthorized
	// curl -u one:xxxxxx http://testserver/
	// Unauthorized
	// curl -u one:secret http://testserver/
	// OK: user=one
}
