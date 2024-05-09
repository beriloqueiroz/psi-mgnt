package domain

import (
	"testing"
)

func TestValidPatient(t *testing.T) {
	patient := &Patient{ID: "12345", Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	err := patient.IsValid()
	if err != nil {
		t.Errorf("Expected valid patient, got error: %v", err)
	}
}

func TestInvalidPatientID(t *testing.T) {
	patient := &Patient{Name: "John Doe", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	err := patient.IsValid()
	if err == nil {
		t.Errorf("Expected invalid patient ID, got nil")
	}
}

func TestInvalidPatientName(t *testing.T) {
	patient := &Patient{ID: "12345", Name: "John", Document: "1234567890", Email: "john@example.com", Phone: Phone{Value: "123456789", IsChat: true}}
	err := patient.IsValid()
	if err == nil {
		t.Errorf("Expected invalid patient name, got nil")
	}
}
