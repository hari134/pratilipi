package serde

import (
	"testing"

	"github.com/hari134/pratilipi/pkg/messaging"
)

func TestBase64ToStruct(t *testing.T) {
	// Base64-encoded JSON string
	base64Str := "eyJ1c2VyX2lkIjoiMTQiLCJlbWFpbCI6ImhhcmkxMi5qYW1lc0BleGFtcGxlLmNvbSIsInBob25lX25vIjoiMTIzNDU2Nzg5MSJ9"

	expected := messaging.UserRegistered{
		UserID:  "14",
		Email:   "hari12.james@example.com",
		PhoneNo: "1234567891",
	}

	var actual messaging.UserRegistered
	err := Base64ToStruct(base64Str, &actual)
	if err != nil {
		t.Fatalf("Base64ToStruct failed: %v", err)
	}

	// Compare the expected and actual struct values
	if actual != expected {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}