package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

type compressedWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (c *compressedWriter) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept-Encoding")

		switch {
		case strings.Contains(accept, "gzip"):
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			next.ServeHTTP(&compressedWriter{ResponseWriter: w, writer: gz}, r)

		case strings.Contains(accept, "deflate"):
			w.Header().Set("Content-Encoding", "deflate")
			w.Header().Del("Content-Length")
			zw := zlib.NewWriter(w)
			defer zw.Close()
			next.ServeHTTP(&compressedWriter{ResponseWriter: w, writer: zw}, r)

		default:
			next.ServeHTTP(w, r)
		}
	})
}
