package lastpass

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Record describes a Lastpass json response
type Record struct {
	Fullname        string `json:"fullname"`
	Group           string `json:"group"`
	ID              string `json:"id"`
	LastModifiedGmt string `json:"last_modified_gmt"`
	LastTouch       string `json:"last_touch"`
	Name            string `json:"name"`
	Note            string `json:"note"`
	Password        string `json:"password"`
	Share           string `json:"share"`
	URL             string `json:"url"`
	Username        string `json:"username"`
}

// Client is our Lastpass (lpass) wrapper client.
type Client struct {
	Username string
	Password string
}

func (r *Record) getTemplate() string {
	template := fmt.Sprintf(`Name: %s
URL: %s
Username: %s 
Password: %s
Notes:    # Add notes below this line.
%s
`, r.Name, r.URL, r.Username, r.Password, r.Note)
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
