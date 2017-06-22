package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Profile struct {
	Host      string   `json:"host"`
	Port      string   `json:"port"`
	Directory string   `json:"directory"`
	Command   []string `json:"command"`
}

func ValidateProfile(prof Profile) error {
	if prof.Host == "" {
		return errors.New("You must specify a host")
	}

	matched, err := regexp.MatchString(`:\d+`, prof.Port)
	if !matched || err != nil {
		return fmt.Errorf("Port %q does not match %q", prof.Port, ":[0-9]+")
	}

	if len(prof.Command) < 1 {
		return errors.New("You must specify a command")
	}
	return nil
}

func Load(fileName string) (Profile, error) {
	var prof Profile
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return prof, err
	}

	if err := json.Unmarshal(text, &prof); err != nil {
		return prof, err
	}

	return prof, nil
}
