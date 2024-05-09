package domain

import "testing"

func TestValidChatPhone(t *testing.T) {
	phone := Phone{Value: "123456789", IsChat: true}
	if !phone.IsChat {
		t.Error("Expected phone to be chat enabled, got false")
	}
}

func TestValidNonChatPhone(t *testing.T) {
	phone := Phone{Value: "987654321", IsChat: false}
	if phone.IsChat {
		t.Error("Expected phone to be non-chat enabled, got true")
	}
}
