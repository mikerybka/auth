package auth

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	DB *DB
}

type SendLoginCodeRequest struct {
	Phone string `json:"phone"`
}

type SendLoginCodeResponse struct {
	UserIDs []string `json:"user_ids"`
	Error   error    `json:"error"`
}

func (s *Server) SendLoginCode(req *SendLoginCodeRequest) SendLoginCodeResponse {
	return SendLoginCodeResponse{}
}

type LoginRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error error  `json:"error"`
}

func (s *Server) Login(req *LoginRequest) LoginResponse {
	return LoginResponse{}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/send-login-code", func(w http.ResponseWriter, r *http.Request) {
		req := &SendLoginCodeRequest{}
		json.NewDecoder(r.Body).Decode(req)
		res := s.SendLoginCode(req)
		json.NewEncoder(w).Encode(res)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		req := &LoginRequest{}
		json.NewDecoder(r.Body).Decode(req)
		res := s.Login(req)
		json.NewEncoder(w).Encode(res)
	})

	mux.ServeHTTP(w, r)
}
