package label

import (
	"fmt"
	"net/http"
	"os"

	"goto/pkg/http/server/listeners"
	"goto/pkg/util"

	"github.com/gorilla/mux"
)

var (
  Handler util.ServerHandler = util.ServerHandler{"delay", SetRoutes, Middleware}
)

func SetRoutes(r *mux.Router, parent *mux.Router, root *mux.Router) {
  labelRouter := r.PathPrefix("/label").Subrouter()
  util.AddRoute(labelRouter, "/set/{label}", setLabel, "PUT", "POST")
  util.AddRoute(labelRouter, "/clear", setLabel, "POST")
  util.AddRoute(labelRouter, "", getLabel)
}

func setLabel(w http.ResponseWriter, r *http.Request) {
  msg := ""
  if label := listeners.SetListenerLabel(r); label == "" {
    msg = "Label cleared"
  } else {
    msg = fmt.Sprintf("Will add header 'Server: %s' to all responses on port %s", label, util.GetListenerPort(r))
  }
  util.AddLogMessage(msg, r)
  w.WriteHeader(http.StatusAccepted)
  fmt.Fprintln(w, msg)
}

func getLabel(w http.ResponseWriter, r *http.Request) {
  label := listeners.GetListenerLabel(r)
  util.AddLogMessage("Server Label: "+label, r)
  w.WriteHeader(http.StatusOK)
  fmt.Fprintln(w, "Server Label: "+label)
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    label := listeners.GetListenerLabel(r)
    pod, present := os.LookupEnv("POD_NAME")
    if !present {
      pod, _ = os.Hostname()
    }
    ns, present := os.LookupEnv("NAMESPACE")
    if !present {
      ns = "local"
    }
    ip := util.GetHostIP()
    hostLabel := fmt.Sprintf("%s.%s@%s", pod, ns, ip)
    util.AddLogMessage(fmt.Sprintf("[%s] [%s]", hostLabel, label), r)
    w.Header().Add("Via-Goto", label)
    w.Header().Add("Goto-Host", hostLabel)
    next.ServeHTTP(w, r)
    util.AddLogMessage(util.GetResponseHeadersLog(w), r)
  })
}
