package utils

import (
	"fmt"
	"os"
)

func GetPort() string {
	return fmt.Sprintf(":%s", os.Getenv("PORT"))
}
