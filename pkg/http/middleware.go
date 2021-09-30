package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HTTPLogger(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := SmartProxyResponseWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)
		sw.Logger.Info("HTTPLogger", zap.String("Host", r.Host), zap.String("RemoteAddr", r.RemoteAddr), zap.String("Method", r.Method), zap.String("RequestURI", r.RequestURI), zap.String("Proto", r.Proto), zap.Any("Status", sw.HttpRetCode), zap.Any("ContentLen", sw.Bytes), zap.String("UserAgent", r.Header.Get("User-Agent")), zap.Any("Duration", duration))
	}
}
