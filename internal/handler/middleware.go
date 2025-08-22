package handler

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// logRequest 中间件记录请求基础信息、耗时等
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		clientIP := getClientIP(r)
		log.Printf("[REQ] %s %s ip=%s ua=%q took=%s",
			r.Method, r.URL.Path, clientIP, r.UserAgent(), time.Since(start))
	})
}

// getClientIP 尝试从 X-Forwarded-For, X-Real-IP 获取真实客户端 IP
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xr := r.Header.Get("X-Real-Ip"); xr != "" {
		return strings.TrimSpace(xr)
	}
	// strip port if present
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
