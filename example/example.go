package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/aws/aws-sdk-go/aws/session"
	console "github.com/ozaaar/aws-go-console"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create aws session: %s", err)
		os.Exit(1)
	}

	con := console.New(sess)
	token, err := con.SignInTokenWithArn("example", "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly")
	if err != nil {
		log.Fatalf("failed to get sign-in token: %s", err)
		os.Exit(1)
	}

	url, err := token.SignInURL("https://console.aws.amazon.com/ecr")
	if err != nil {
		log.Fatalf("failed to create sign url: %s", err)
		os.Exit(1)
	}

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url.String()).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url.String()).Start()
	case "darwin":
		err = exec.Command("open", url.String()).Start()
	default:
		err = errors.New("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
