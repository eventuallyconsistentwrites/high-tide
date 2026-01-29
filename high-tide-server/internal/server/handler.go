package server

import (
	"net/http"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/post"
)

type Server struct {
	Store post.Repository
}

func (s *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	// Logic to call s.Store.Get(...)
	w.Write([]byte("Post data here"))
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Logic to call s.Store.Create(...)
}

func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Logic to call s.Store.Create(...)
}
