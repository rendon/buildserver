package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rendon/buildserver/profile"
)

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	prof, err := profile.Load(os.Args[1])
	if err != nil {
		fatalf("Failed to load profile: %s", err)
	}

	url := fmt.Sprintf("http://%s%s/build", prof.Host, prof.Port)
	resp, err := http.Get(url)
	if err != nil {
		fatalf("Build request failed: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		fatalf("Build request failed with status %d", resp.StatusCode)
	}

	fmt.Printf("Build request successfully sent\n")
}
