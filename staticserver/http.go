package staticserver

import (
	"mime"
	"net/http"
	"path/filepath"
)

// StaticContentServer is an http.Handler which serves static content
type StaticContentServer struct {
	content map[string][]byte
}

// NewStaticContentServer creates a new StaticContentServer instance
func New(content map[string][]byte) *StaticContentServer {
	return &StaticContentServer{content}
}

// ServeHTTP implements the http.Handler interface
func (s *StaticContentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Path
	if len(fileID) > 0 {
		fileID = fileID[1:] // strip leading slash
		if content, ok := s.content[fileID]; ok {
			mimeType := mime.TypeByExtension(filepath.Ext(fileID))
			if mimeType != "" {
				w.Header().Add(http.CanonicalHeaderKey("Content-Type"), mimeType)
			}
			w.Write(content)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
