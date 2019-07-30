package lastpass

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

// Create is used to create a new resource and generate ID.
func (c *Client) Create(s Secret) (Secret, error) {
	err := c.login()
	if err != nil {
		return s, err
	}
	template := s.getTemplate()
	cmd := exec.Command("lpass", "add", s.Name, "--non-interactive", "--sync=now")
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return s, err
	}
	time.Sleep(time.Second * 5) // Need to finish sync with upstream/lastpass before we get actual ID.
	cmd = exec.Command("lpass", "show", s.Name, "--json", "-x")
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return s, err
	}
	var secrets []Secret
	err = json.Unmarshal(outbuf.Bytes(), &secrets)
	if err != nil {
		return s, err
	}
	if len(secrets) > 1 {
		err := errors.New("more than one secret with same name, unable to determine ID")
		return s, err
	}
	if secrets[0].ID == "0" {
		err := errors.New("got invalid ID 0, problem with lastpass sync")
		return secrets[0], err
	}
	return secrets[0], nil
}
