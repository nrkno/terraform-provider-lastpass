package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ansd/lastpass-go"
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

// Client is our Lastpass (lastpass-go) wrapper client.
type Client struct {
	Client    *lastpass.Client
	Accounts  []*lastpass.Account
	Username  string
	Password  string
	Trust     bool
	TwoFA     bool
	OnetimePW string
	ConfigDIR string
	BaseURL   string
}

func (c *Client) Login() error {
	var client *lastpass.Client
	// authenticate with LastPass servers
	basedir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	fullpath := filepath.Join(basedir, c.ConfigDIR)
	if !c.TwoFA {
		client, err = lastpass.NewClient(
			context.Background(),
			c.Username,
			c.Password,
			lastpass.WithBaseURL(c.BaseURL),
			lastpass.WithConfigDir(fullpath))
	}
	if c.Trust {
		client, err = lastpass.NewClient(
			context.Background(),
			c.Username,
			c.Password,
			lastpass.WithOneTimePassword(c.OnetimePW),
			lastpass.WithBaseURL(c.BaseURL),
			lastpass.WithTrust(),
			lastpass.WithConfigDir(fullpath))
	} else {
		client, err = lastpass.NewClient(
			context.Background(),
			c.Username,
			c.Password,
			lastpass.WithOneTimePassword(c.OnetimePW),
			lastpass.WithBaseURL(c.BaseURL),
			lastpass.WithConfigDir(fullpath))
	}
	if err != nil {
		return err
	}
	c.Client = client
	return nil
}

func (c *Client) Sync() error {
	// read all Accounts()
	accounts, err := c.Client.Accounts(context.Background())
	if err != nil {
		return err
	}
	c.Accounts = accounts
	return nil
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
			key := s.Note[ss[2]:ss[3]]
			value := s.Note[ss[4]:ss[5]]
			if key == "Notes" && strings.Contains(value, "\n") {
				notes[key] = strings.TrimSuffix(value, "\n")
			} else {
				notes[key] = value
			}
		}
	}
	s.CustomFields = notes
}
