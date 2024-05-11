package domain

import (
	"errors"
	"time"
)

type Session struct {
	ID          string        `bson:"id"`
	Price       float64       `bson:"price"`
	Notes       string        `bson:"notes"`
	Date        time.Time     `bson:"date"`
	PaymentDate time.Time     `bson:"payment_date"`
	Duration    time.Duration `bson:"duration"`
	Patient     *Patient      `bson:"patient"`
	OwnerId     string        `bson:"owner_id"`
}

func NewSession(
	id string,
	price float64,
	notes string,
	date time.Time,
	paymentDate time.Time,
	duration time.Duration,
	patient *Patient,
	ownerId string,
) (*Session, error) {
	session := &Session{
		ID:          id,
		Price:       price,
		Notes:       notes,
		Date:        date,
		PaymentDate: paymentDate,
		Duration:    duration,
		Patient:     patient,
		OwnerId:     ownerId,
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
		return errors.New("invalid notes")
	}
	if s.OwnerId == "" {
		return errors.New("invalid owner")
	}
	return nil
}
