package student

import "net/http"

func NewStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Student API"))
	}
}
