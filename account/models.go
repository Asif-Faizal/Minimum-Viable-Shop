package account

type Account struct {
	ID       string `json:"id" validate:"required,uuidv4"`
	Name     string `json:"name" validate:"min=3,max=50"`
	Email    string `json:"email" validate:"required,email,normalizeemail"`
	Password string `json:"-" validate:"required,min=8,max=50"`
}
