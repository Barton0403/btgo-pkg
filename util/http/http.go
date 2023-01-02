package http

import (
	"net/http"
	"os"
	"strings"
)

func Get(req *http.Request) map[string]string {
	var result = make(map[string]string)
	keys := req.URL.Query()
	for k, v := range keys {
		result[k] = v[0]
	}

	return result
}

func PostForm(req *http.Request) map[string]string {
	var result = make(map[string]string)
	req.ParseForm()
	for k, v := range req.PostForm {
		if len(v) < 1 {
			continue
		}

		result[k] = v[0]
	}

	return result
}

func BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}

	return token, token != ""
}

func OutputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
