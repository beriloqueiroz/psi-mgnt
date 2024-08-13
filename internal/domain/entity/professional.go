package domain

import "errors"

type Professional struct {
	ID       string
	Name     string `bson:"name"`
	Document string `bson:"document"`
	Email    string `bson:"email"`
}

func NewProfessional(id string, name string, document string, email string) (*Professional, error) {
	professional := &Professional{
		ID:       id,
		Name:     name,
		Document: document,
		Email:    email,
	}
	err := professional.IsValid()
	if err != nil {
		return nil, err
	}
	return professional, nil
}

func (p *Professional) IsValid() error {
	if p.ID == "" {
		return errors.New("invalid id")
	}
	if len(p.Name) <= 5 {
		return errors.New("invalid name")
	}
	return nil
}
