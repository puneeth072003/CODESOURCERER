package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GenerateTestFiles(tests []struct {
	TestName    string
	TestPath    string
	ParentPath  string
	CodeContent string
}) error {
	for _, test := range tests {
		dir := filepath.Dir(test.TestPath)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("unable to create directory %s: %w", dir, err)
		}
		file, err := os.Create(test.TestPath)
		if err != nil {
			return fmt.Errorf("unable to create test file %s: %w", test.TestPath, err)
		}
		defer file.Close()
		_, err = file.WriteString(test.CodeContent)
		if err != nil {
			return fmt.Errorf("unable to write to test file %s: %w", test.TestPath, err)
		}
		log.Printf("Generated test file: %s", test.TestPath)
	}
	return nil
}
