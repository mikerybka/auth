package auth

type Phone struct {
	Number  string   `json:"number"`
	UserIDs []string `json:"user_ids"`
}
