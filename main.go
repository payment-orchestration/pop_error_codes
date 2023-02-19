package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type StringList map[string]struct{}

func (sl StringList) Add(value string) {
	sl[value] = struct{}{}
}

func (sl StringList) Has(value string) bool {
	_, ok := sl[value]
	return ok
}

func main() {
	errCodes, err := validateCoreConfig("./core")
	if err != nil {
		panic(err)
	}

	err = validateConnectorConfig("./connector", errCodes)
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}

func validateCoreConfig(configDir string) (StringList, error) {
	format := struct {
		Errors []struct {
			Code        string
			Description string
		}
	}{}

	codes := StringList{}

	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("reading config directory error: %s", err)
	}

	for _, file := range files {
		filePath := filepath.Join(configDir, file.Name())
		f, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("file reading error: %s.\n%s", filePath, err)
		}

		err = yaml.Unmarshal(f, &format)
		if err != nil {
			return nil, fmt.Errorf("invalid error yaml format: %s.\n%s", filePath, err)
		}

		if len(format.Errors) == 0 {
			return nil, fmt.Errorf("`errors` key is not in core config: %s", filePath)
		}

		for _, e := range format.Errors {
			if e.Code == "" {
				return nil, fmt.Errorf("invalid error yaml format: value of `code` is empty")
			}

			codes.Add(e.Code)
		}
	}

	return codes, nil
}

func validateConnectorConfig(configDir string, errCodes StringList) error {
	format := map[string][]string{}
	files, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("reading config directory error: %s", err)
	}

	for _, file := range files {
		filePath := filepath.Join(configDir, file.Name())
		f, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("file reading error: %s.\n%s", filePath, err)
		}

		err = yaml.Unmarshal(f, &format)
		if err != nil {
			return fmt.Errorf("invalid error yaml format: %s.\n%s", filePath, err)
		}

		for k := range format {
			if !errCodes.Has(k) {
				return fmt.Errorf("unkown error code `%s` in %s", k, filePath)
			}
		}
	}

	return nil
}
