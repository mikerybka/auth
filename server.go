package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mikerybka/twilio"
)

type Server struct {
	DB           *DB
	TwilioClient *twilio.Client
}

type SendLoginCodeRequest struct {
	Phone string `json:"phone"`
}

type SendLoginCodeResponse struct {
	UserIDs []string `json:"user_ids"`
	Error   error    `json:"error"`
}

func (s *Server) SendLoginCode(req *SendLoginCodeRequest) SendLoginCodeResponse {
	phone, err := s.DB.Phone(req.Phone)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			userID := newUserID()

			user := &User{
				ID: userID,
			}
			err := s.DB.SaveUser(user)
			if err != nil {
				panic(err)
			}

			phone = &Phone{
				Number:     req.Phone,
				UserIDs:    []string{userID},
				LoginCodes: map[string]bool{},
			}
			err = s.DB.SavePhone(phone)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	loginCode := newLoginCode()
	phone.LoginCodes[loginCode] = true
	err = s.DB.SavePhone(phone)
	if err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("Your login code is %s", loginCode)
	err = s.TwilioClient.SendSMS(phone.Number, msg)
	if err != nil {
		panic(err)
	}

	return SendLoginCodeResponse{
		UserIDs: phone.UserIDs,
	}
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

func (s *Server) GetUserID(r *http.Request) string
