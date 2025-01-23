package auth

type Team struct {
	ID      string   `json:"id"`
	Members []string `json:"members"`
}

func (t *Team) Owner() string {
	return t.Members[0]
}
