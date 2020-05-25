package header

import (
	"net/http"

	"goto/pkg/http/server/request/header/tracking"
	"goto/pkg/util"

	"github.com/gorilla/mux"
)

var (
	Handler util.ServerHandler = util.ServerHandler{"header", SetRoutes, Middleware}
)

func SetRoutes(r *mux.Router, parent *mux.Router) {
	headersRouter := r.PathPrefix("/headers").Subrouter()
	tracking.SetRoutes(headersRouter, r)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.AddLogMessage(util.GetRequestHeadersLog(r), r)
		tracking.Middleware(next).ServeHTTP(w, r)
	})
}
