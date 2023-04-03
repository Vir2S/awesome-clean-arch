package profile

type Profile struct {
	ID        int    `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	School    string `json:"school"`
}
