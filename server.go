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
	Error error `json:"error"`
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

	return SendLoginCodeResponse{}
}

type LoginRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	Token   string   `json:"token"`
	UserIDs []string `json:"user_ids"`
	Error   error    `json:"error"`
}

func (s *Server) Login(req *LoginRequest) LoginResponse {
	phone, err := s.DB.Phone(req.Phone)
	if err != nil {
		return LoginResponse{
			Error: err,
		}
	}

	if !phone.LoginCodes[req.Code] {
		return LoginResponse{
			Error: fmt.Errorf("bad code"),
		}
	}

	phone.LoginCodes[req.Code] = false
	err = s.DB.SavePhone(phone)
	if err != nil {
		panic(err)
	}

	session := &Session{
		Token: newSessionToken(),
		Phone: phone.Number,
	}
	err = s.DB.SaveSession(session)
	if err != nil {
		panic(err)
	}

	return LoginResponse{
		Token:   session.Token,
		UserIDs: phone.UserIDs,
	}
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

// GetUserID returns the ID of the requesting user
func (s *Server) GetUserID(r *http.Request) string {
	token := r.Header.Get("Token")
	session, err := s.DB.Session(token)
	if err != nil {
		return ""
	}
	phone, err := s.DB.Phone(session.Phone)
	if err != nil {
		return ""
	}
	userID := r.Header.Get("User")
	for _, id := range phone.UserIDs {
		if id == userID {
			return id
		}
	}
	return ""
}
