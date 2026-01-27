package rest

import pb "github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	UserType string `json:"user_type"`
	Email    string `json:"email"`
}

type AuthenticatedResponse struct {
	Account      *Account `json:"account"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

func toAccount(a *pb.Account) *Account {
	if a == nil {
		return nil
	}
	return &Account{
		ID:       a.Id,
		Name:     a.Name,
		UserType: a.Usertype,
		Email:    a.Email,
	}
}
