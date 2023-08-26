package people

type PersonDTO struct {
	ID        string   `json:"id,omitempty"`
	Nickname  *string  `json:"apelido" validate:"required,max=32"`
	Name      *string  `json:"nome" validate:"required,max=100"`
	Birthdate string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack     []string `json:"stack,omitempty" validate:"omitempty,dive,gt=0,max=32"`
}
