package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ansd/lastpass-go"
)

// Secret describes a Lastpass object.
type Secret struct {
	ID              string
	Name            string
	Username        string
	Password        string
	URL             string
	Share           string
	Group           string
	Notes           string
	LastModifiedGmt string
	LastTouch       string
	CustomFields    map[string]string
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
	if strings.HasPrefix(s.Notes, "NoteType:") {
		// fix notes so the regexp works
		s.Notes = "\n" + s.Notes

		// change '\n<words>:' to something more precise regexp can parse
		tokenizer := regexp.MustCompile(`\n([[:alnum:]][ -_[:alnum:]]+:)`)
		s.Notes = tokenizer.ReplaceAllString(s.Notes, "\a$1\a")

		// break up notes using '\n<word>:<multi-line-string without control character bell>'
		// - which implies that custom-fields values cannot include the bell character
		// - allows for an inexpensive parser using regexp
		splitter := regexp.MustCompile(`\a([ -_[:alnum:]]+):\a([^\a]*)`)
		splitted := splitter.FindAllStringSubmatchIndex(s.Notes, -1)
		fmt.Println(splitted)
		for _, ss := range splitted {
			fmt.Println("*>> ", string(s.Notes[ss[0]:ss[1]]))
			fmt.Println("[0] ", string(s.Notes[ss[2]:ss[3]]))
			fmt.Println("[1] ", string(s.Notes[ss[4]:ss[5]]))
		}

		for _, ss := range splitted {
			key := s.Notes[ss[2]:ss[3]]
			value := s.Notes[ss[4]:ss[5]]
			if key == "Notes" && strings.Contains(value, "\n") {
				notes[key] = strings.TrimSuffix(value, "\n")
			} else {
				notes[key] = value
			}
		}
	}
	s.CustomFields = notes
}

func epochToTime(s string) (string, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", err
	}
	return time.Unix(sec, 0).String(), nil
}
