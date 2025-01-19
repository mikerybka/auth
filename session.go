package auth

type Session struct {
	Token string `json:"token"`
	Phone string `json:"phone"`
}
