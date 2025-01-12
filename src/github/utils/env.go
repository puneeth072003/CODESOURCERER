package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func Loadenv(filePath string) (map[string]string, error) {
	envs := make(map[string]string)

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // skip lines that don't have one '=' sign
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envs[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envs, nil
}
