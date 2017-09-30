package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/rendon/buildserver/profile"
)

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func printOutput(rc io.ReadCloser, done chan bool) {
	defer func() {
		rc.Close()
		done <- true
	}()

	buf := make([]byte, 1024)
	for {
		n, err := rc.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Something went wrong while reading: %s", err)
			return
		}
		fmt.Printf("%s", buf[0:n])
		time.Sleep(20 * time.Millisecond)
	}
}

func build(prof profile.Profile) {
	if prof.Directory != "" {
		log.Printf("cd %s", prof.Directory)
		if err := os.Chdir(prof.Directory); err != nil {
			fatalf("Failed to change directory: %s", err)
		}
	}

	log.Printf("Executing command %v", prof.Command)
	cmd := exec.Command(prof.Command[0], prof.Command[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to create STDOUT pipe: %s", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Failed to create STDERR pipe: %s", err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start command: %s", err)
		return
	}

	done := make(chan bool)
	go printOutput(stderr, done)
	go printOutput(stdout, done)

	// Wait for print STDERR and print STDOUT to finish
	<-done
	<-done

	if err := cmd.Wait(); err != nil {
		log.Printf("%s", err)
		return
	}
}

func main() {
	if len(os.Args) != 2 {
		fatalf("Usage: %s <profile_file>\n", os.Args[0])
	}

	prof, err := profile.Load(os.Args[1])
	if err != nil {
		fatalf("Failed to load profile: %s\n", err)
	}

	if err := profile.ValidateProfile(prof); err != nil {
		fatalf("Invalid profile: %s", err)
	}

	http.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		go build(prof)
	})
	log.Printf("Listening at localhost%s...", prof.Port)
	log.Fatal(http.ListenAndServe(prof.Port, nil))
}
