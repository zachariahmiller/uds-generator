// handler.go

package api

import (
	"net/http"
)

// APIHandler handles API requests
func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the API!"))
}
