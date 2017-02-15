package swu

import (
	"fmt"
	"log"
	"os"
)

func GetEnvVariable(name string, required bool) string {
	variable, exists := os.LookupEnv(name)
	if required && !exists {
		panic(fmt.Sprintf("Unable to find required environment variable %v", name))
	} else if !required && !exists {
		log.Printf("Unable to find environment variable %v\n", name)
	}
	return variable
}
