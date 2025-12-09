package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
)

var errInvalidMapPointer = errors.New("map pointer cannot be nil")

// JsonUnmarshalToMapAndStruct unmarshal JSON data into a map and a struct.
func JsonUnmarshalToMapAndStruct(data []byte, dest any, m *map[string]any) error {
	if m == nil {
		return errInvalidMapPointer
	}

	err := json.Unmarshal(data, m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	err = decoder.Decode(*m)
	if err != nil {
		return fmt.Errorf("failed to decode data to struct: %w", err)
	}

	return nil
}
