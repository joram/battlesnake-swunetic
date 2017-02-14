package swu

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Copied from here: https://code.google.com/p/go/source/detail?r=5e03333d2dcf
// Remove this once it's merged into stdlib
func ParseBasicAuth(r *http.Request) (username, password string, ok bool) {

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}

	if !strings.HasPrefix(auth, "Basic ") {
		return
	}

	c, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}

	return cs[:s], cs[s+1:], true
}
