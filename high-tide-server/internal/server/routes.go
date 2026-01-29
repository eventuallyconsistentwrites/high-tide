package server

import "net/http"

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts", s.GetPost)
	mux.HandleFunc("POST /posts", s.CreatePost)
	mux.HandleFunc("DELETE /posts", s.DeletePost)

	return mux
}
