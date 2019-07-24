package lastpass

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

// Client is our Lastpass (lpass) wrapper client.
type Client struct {
	username string
	password string
}

// NewClient creates a new Lastpass client
func NewClient(u, p string) *Client {
	c := new(Client)
	c.username = u
	c.password = p
	return c
}

func (c *Client) login() error {
	cmd := exec.Command("lpass", "status", "-q")
	err := cmd.Run()
	if err != nil {
		if c.username == "" {
			err := errors.New("Not logged in, please run 'lpass login' manually and try again")
			return err
		}
		cmd := exec.Command("lpass", "login", c.username)
		var inbuf, errbuf bytes.Buffer
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "LPASS_DISABLE_PINENTRY=1")
		inbuf.Write([]byte(c.password))
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

/* Our client CRUD methods */

// Create is used to create a new resource and generate ID.
func (c *Client) Create(r Record) (Record, error) {
	err := c.login()
	if err != nil {
		return r, err
	}
	template := r.getTemplate()
	cmd := exec.Command("lpass", "add", r.Name, "--non-interactive", "--sync=now")
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return r, err
	}
	time.Sleep(time.Second * 5) // Need to finish sync with upstream/lastpass before we get actual ID.
	cmd = exec.Command("lpass", "show", r.Name, "--json", "-x")
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return r, err
	}
	var records []Record
	err = json.Unmarshal(outbuf.Bytes(), &records)
	if err != nil {
		return r, err
	}
	if len(records) > 1 {
		err := errors.New("more than one record with same name")
		return r, err
	}
	if records[0].ID == "0" {
		err := errors.New("got invalid ID 0, possible problem with lpass sync")
		return records[0], err
	}
	return records[0], nil
}

// Update is called to update record with upstream
func (c *Client) Update(r Record) error {
	err := c.login()
	if err != nil {
		return err
	}
	template := r.getTemplate()
	cmd := exec.Command("lpass", "edit", r.ID, "--non-interactive", "--sync=now")
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return err
	}
	return nil
}

// Fetch record from upstream to update local record
func (c *Client) Read(id string) (Record, error) {
	var r Record
	err := c.login()
	if err != nil {
		return r, err
	}
	cmd := exec.Command("lpass", "show", id, "--json", "-x")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		// Make sure the record is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			// If no record is found, set to 0 for deletion.
			r.ID = "0"
			return r, err
		}
		var err = errors.New(errbuf.String())
		return r, err
	}
	var records []Record
	err = json.Unmarshal(outbuf.Bytes(), &records)
	if err != nil {
		return r, err
	}
	if records[0].URL == "http://" {
		records[0].URL = ""
	}
	records[0].Note = records[0].Note + "\n" // lastpass trims new line, provokes constant changes.
	return records[0], nil
}

// Delete record in upstream db
func (c *Client) Delete(id string) error {
	err := c.login()
	if err != nil {
		return err
	}
	var errbuf bytes.Buffer
	cmd := exec.Command("lpass", "rm", id, "--sync=now")
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		// Make sure the record is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			return nil
		}
		var err = errors.New(errbuf.String())
		return err
	}
	return nil
}
