package domain

import (
	"errors"
	"time"
)

type Session struct {
	ID           string        `bson:"id"`
	Price        float64       `bson:"price"`
	Notes        string        `bson:"notes"`
	Date         time.Time     `bson:"date"`
	Duration     time.Duration `bson:"duration"`
	Patient      *Patient      `bson:"patient"`
	Professional *Professional `bson:"professional"`
	Plan         string        `bson:"plan"`
}

func NewSession(
	id string,
	price float64,
	notes string,
	date time.Time,
	duration time.Duration,
	patient *Patient,
	plan string,
	professional *Professional,
) (*Session, error) {
	session := &Session{
		ID:           id,
		Price:        price,
		Notes:        notes,
		Date:         date,
		Duration:     duration,
		Patient:      patient,
		Professional: professional,
		Plan:         plan,
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
	if s.Professional == nil {
		return errors.New("invalid professional")
	}
	return nil
}
