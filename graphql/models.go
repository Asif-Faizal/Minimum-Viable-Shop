package graphql

type Account struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	UserType string  `json:"user_type"`
	Email    string  `json:"email"`
	Orders   []Order `json:"orders"`
}
