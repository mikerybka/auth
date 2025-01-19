package auth

type Session struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
