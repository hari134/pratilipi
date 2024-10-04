package serde

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Base64ToStruct decodes a base64 string and unmarshals it into the provided struct
func Base64ToStruct(base64Str string, v interface{}) error {
    // Decode the base64 string
    decodedData, err := base64.RawStdEncoding.DecodeString(base64Str)
    if err != nil {
        return fmt.Errorf("failed to decode base64 string: %w", err)
    }

    // Unmarshal the JSON into the provided struct
    err = json.Unmarshal(decodedData, v)
    if err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    return nil
}