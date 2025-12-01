package anconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("erreur ouverture config: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("erreur lecture JSON: %w", err)
	}

	return nil
}
