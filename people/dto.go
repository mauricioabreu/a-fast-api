package people

type PersonDTO struct {
	ID        string   `json:"id,omitempty"`
	Nickname  string   `json:"nickname" validate:"required,max=32"`
	Name      string   `json:"name" validate:"required,max=100"`
	Birthdate string   `json:"birthdate" validate:"required,datetime=2006-01-02"`
	Stack     []string `json:"stack" validate:"dive,max=32"`
}
