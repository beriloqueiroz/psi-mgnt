package domain

import (
	"errors"
	"time"
)

type Session struct {
	ID          string
	Price       float64
	Notes       string
	Date        time.Time
	PaymentDate time.Time
	Duration    time.Duration
	Patient     *Patient
}

func NewSession(
	id string,
	price float64,
	notes string,
	date time.Time,
	paymentDate time.Time,
	duration time.Duration,
	patient *Patient,
) (*Session, error) {
	session := &Session{
		ID:          id,
		Price:       price,
		Notes:       notes,
		Date:        date,
		PaymentDate: paymentDate,
		Duration:    duration,
		Patient:     patient,
	}
	err := session.IsValid()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Session) IsValid() error {
	if s.ID == "" {
		return errors.New("invalid id")
	}
	if s.Price <= 0 {
		return errors.New("invalid price")
	}
	if s.Duration <= 0 {
		return errors.New("invalid duration")
	}
	if len(s.Notes) <= 4 {
		return errors.New("invalid duration")
	}
	return nil
}
