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
        // fix notes so the regexp works
        s.Note = "\n" + s.Note

        // change '\n<words>:' to something more precise regexp can parse
        tokenizer := regexp.MustCompile(`\n([[:alnum:]][ [:alnum:]]+:)`)
        s.Note = tokenizer.ReplaceAllString(s.Note, "\a$1\a")

        // break up notes using '\n<word>:<multi-line-string without control character bell>'
        // - which implies that custom-fields values cannot include the bell character
        // - allows for an inexpensive parser using regexp
        splitter := regexp.MustCompile(`\a([ [:alnum:]]+):\a([^\a]*)`)
        splitted := splitter.FindAllStringSubmatchIndex(s.Note, -1)
        fmt.Println(splitted)
        for _, ss := range splitted {
          fmt.Println("*>> ", string(s.Note[ss[0]:ss[1]]))
          fmt.Println("[0] ", string(s.Note[ss[2]:ss[3]]))
          fmt.Println("[1] ", string(s.Note[ss[4]:ss[5]]))
        }

		for _, ss := range splitted {
            key := s.Note[ ss[2]:ss[3] ]
            value := s.Note[ ss[4]:ss[5] ]
            if key == "Notes" && strings.Contains(value, "\n") {
              notes[ key ] = strings.TrimSuffix(value, "\n")
            } else {
              notes[ key ] = value
            }
		}
	} else {
		// Fix for Notes with multiline. Always last in end of the string.
		if strings.Contains(s.Note, "\n") {
			s.Note = s.Note + "\n" // lastpass trims new line, add back to multiline notes.
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
%s`, s.Name, s.URL, s.Username, s.Password, s.Note)
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

// perhaps we should call Login only once providerConfigure()
func (c *Client) Login() error {
	return c.login()
}
