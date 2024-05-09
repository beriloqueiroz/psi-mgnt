package domain

import (
	"testing"
	"time"
)

func TestValidSession(t *testing.T) {
	validPatient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	session := &Session{ID: "123", Price: 100, Notes: "Session notes", Date: time.Now(), PaymentDate: time.Now(), Duration: 30 * time.Minute, Patient: validPatient}
	err := session.IsValid()
	if err != nil {
		t.Errorf("Expected valid session, got error: %v", err)
	}
}

func TestInvalidSessionID(t *testing.T) {
	validPatient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	session := &Session{Price: 100, Notes: "Session notes", Date: time.Now(), PaymentDate: time.Now(), Duration: 30 * time.Minute, Patient: validPatient}
	err := session.IsValid()
	if err == nil {
		t.Errorf("Expected invalid session ID, got nil")
	}
}

func TestInvalidSessionPrice(t *testing.T) {
	validPatient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	session := &Session{ID: "123", Price: -100, Notes: "Session notes", Date: time.Now(), PaymentDate: time.Now(), Duration: 30 * time.Minute, Patient: validPatient}
	err := session.IsValid()
	if err == nil {
		t.Errorf("Expected invalid session price, got nil")
	}
}

func TestInvalidSessionDuration(t *testing.T) {
	validPatient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	session := &Session{ID: "123", Price: 100, Notes: "Session notes", Date: time.Now(), PaymentDate: time.Now(), Duration: 0, Patient: validPatient}
	err := session.IsValid()
	if err == nil {
		t.Errorf("Expected invalid session duration, got nil")
	}
}

func TestInvalidSessionNotes(t *testing.T) {
	validPatient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	session := &Session{ID: "123", Price: 100, Notes: "", Date: time.Now(), PaymentDate: time.Now(), Duration: 10 * time.Minute, Patient: validPatient}
	err := session.IsValid()
	if err == nil {
		t.Errorf("Expected invalid session notes, got nil")
	}
}
