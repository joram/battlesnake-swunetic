package swu

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	_routeLogging bool
	_awsSession   *session.Session = session.New(&aws.Config{Region: aws.String("us-east-1")})
)

func getEnv(name string, required bool) string {
	v := os.Getenv(name)
	if required && v == "" {
		panic(fmt.Sprintf("$%s not set", name))
	}

	return v
}

func initConfig() {
	_routeLogging = getEnv("ROUTE_LOGGING", false) != ""
}

func getAwsSession() *session.Session {
	return _awsSession
}
