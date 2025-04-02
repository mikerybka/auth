package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mikerybka/twilio"
	"github.com/mikerybka/util"
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

func (s *Server) SendLoginCode(w http.ResponseWriter, r *http.Request) {
	// Decode input
	req := &SendLoginCodeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Phone = util.NormalizePhoneNumber(req.Phone)

	// Check the DB for the existing phone number
	phone, err := s.DB.Phone(req.Phone)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If not, create a user for that phone number.
			user := &User{
				ID: newUserID(),
			}
			err := s.DB.SaveUser(user)
			if err != nil {
				panic(err)
			}

			// And assign the new user to that phone.
			phone = &Phone{
				Number:     req.Phone,
				UserIDs:    []string{user.ID},
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

	// Create a new login code in the system
	loginCode := newLoginCode()
	phone.LoginCodes[loginCode] = true
	err = s.DB.SavePhone(phone)
	if err != nil {
		panic(err)
	}

	// Send login code to phone
	msg := fmt.Sprintf("Your login code is %s", loginCode)
	err = s.TwilioClient.SendSMS(phone.Number, msg)
	if err != nil {
		panic(err)
	}

	// Respond with a nil error.
	res := SendLoginCodeResponse{}
	json.NewEncoder(w).Encode(res)
}

type LoginRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	Token   string   `json:"token"`
	UserIDs []string `json:"userIDs"`
	Error   error    `json:"error"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Phone = util.NormalizePhoneNumber(req.Phone)

	// Read
	phone, err := s.DB.Phone(req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check validity of the login code.
	if !phone.LoginCodes[req.Code] {
		http.Error(w, "wrong code", http.StatusBadRequest)
		return
	}

	// Remove the login code from the DB since it's now been used.
	phone.LoginCodes[req.Code] = false
	err = s.DB.SavePhone(phone)
	if err != nil {
		panic(err)
	}

	// Create a new session in the DB.
	session := &Session{
		Token: newSessionToken(),
		Phone: phone.Number,
	}
	err = s.DB.SaveSession(session)
	if err != nil {
		panic(err)
	}

	// Respond with Session data.
	res := LoginResponse{
		Token:   session.Token,
		UserIDs: phone.UserIDs,
	}
	json.NewEncoder(w).Encode(res)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-login-code", s.SendLoginCode)
	mux.HandleFunc("/login", s.Login)
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
