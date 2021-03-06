package conn

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"goto/pkg/util"
)

var (
  Handler       util.ServerHandler = util.ServerHandler{Name: "connection", Middleware: Middleware}
  connectionKey *util.ContextKey   = &util.ContextKey{"connection"}
)

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {
  return context.WithValue(ctx, connectionKey, c)
}

func GetConn(r *http.Request) net.Conn {
  return r.Context().Value(connectionKey).(net.Conn)
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //log.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
    conn := GetConn(r)
    localAddr := conn.LocalAddr().String()
    remoteAddr := conn.RemoteAddr().String()
    util.AddLogMessage(fmt.Sprintf("LocalAddr: %s, RemoteAddr: %s", localAddr, remoteAddr), r)
    next.ServeHTTP(w, r)
  })
}
