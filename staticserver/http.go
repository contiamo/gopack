package staticserver

import "net/http"

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
	path := r.URL.Path
	if len(path) > 0 {
		if content, ok := s.content[path[1:]]; ok {
			w.Write(content)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
