package api

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Secret describes a Lastpass object.
type Secret struct {
	Fullname        string            `json:"fullname"`
	Group           string            `json:"group"`
	ID              string            `json:"id"`
	LastModifiedGmt string            `json:"last_modified_gmt"`
	LastTouch       string            `json:"last_touch"`
	Name            string            `json:"name"`
	Note            string            `json:"note"`
	Password        string            `json:"password"`
	Share           string            `json:"share"`
	URL             string            `json:"url"`
	Username        string            `json:"username"`
	CustomFields    map[string]string `json:"custom_fields"`
}

// Client is our Lastpass (lpass) wrapper client.
type Client struct {
	Username string
	Password string
}

func (s *Secret) genCustomFields() {
	notes := make(map[string]string)
	if strings.HasPrefix(s.Note, "NoteType:") {
		splitted := strings.Split(s.Note, "\n")
		for _, split := range splitted {
			re := regexp.MustCompile(`:`)
			s := re.Split(split, 2)
			if s[0] == "Notes" {
				break
			}
			if len(s) == 2 {
				notes[s[0]] = s[1]
			}
		}
		// Fix for Notes with multiline. Always last in end of the string.
		n := strings.Split(s.Note, "\nNotes:")
		if len(n) == 2 {
			notes["Notes"] = n[1]
		}
	}
	s.CustomFields = notes
}

func (s *Secret) getTemplate() string {
	template := fmt.Sprintf(`Name: %s
URL: %s
Username: %s 
Password: %s
Notes:    # Add notes below this line.
%s
`, s.Name, s.URL, s.Username, s.Password, s.Note)
	return template
}

func (c *Client) login() error {
	cmd := exec.Command("lpass", "status", "-q")
	err := cmd.Run()
	if err != nil {
		if c.Username == "" {
			err := errors.New("Not logged in, please run 'lpass login' manually and try again")
			return err
		}
		cmd := exec.Command("lpass", "login", c.Username)
		var inbuf, errbuf bytes.Buffer
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "LPASS_DISABLE_PINENTRY=1")
		inbuf.Write([]byte(c.Password))
		cmd.Stdin = &inbuf
		cmd.Stderr = &errbuf
		err := cmd.Run()
		if err != nil {
			var err = errors.New(errbuf.String())
			return err
		}
	}
	return nil
}
