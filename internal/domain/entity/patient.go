package domain

import "errors"

type Patient struct {
	ID       string
	Name     string
	Document string
	Email    string
	Phone    Phone
}

func NewPatient(id string, name string, document string, email string, phone Phone) (*Patient, error) {
	patient := &Patient{
		ID:       id,
		Name:     name,
		Document: document,
		Email:    email,
		Phone:    phone,
	}
	err := patient.IsValid()
	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (p *Patient) IsValid() error {
	if p.ID == "" {
		return errors.New("invalid id")
	}
	if len(p.Name) <= 5 {
		return errors.New("invalid name")
	}
	return nil
}
