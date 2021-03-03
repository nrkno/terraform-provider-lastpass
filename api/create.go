package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"
)

// Create is used to create a new resource and generate ID.
func (c *Client) Create(s Secret) (Secret, error) {
    template := s.getTemplate()
    cmd := exec.Command("lpass", "add", s.Name, "--non-interactive", "--sync=now")
    return c.create(s.Name, template, cmd)
}

// Create a secret of type node-type
func (c *Client) CreateNodeType(name string, template string, nodetype string) (Secret, error) {
    cmd := exec.Command("lpass", "add", name, "--non-interactive", "--sync=now", "--note-type=" + nodetype)
    return c.create(name, template, cmd)
}

func (c *Client) create(name string, template string, cmd *exec.Cmd) (Secret, error) {
	var s Secret = Secret{ Name: name }
	err := c.login()
	if err != nil {
		return s, err
	}
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return s, err
	}
	var outbuf bytes.Buffer
	var secrets []Secret
	// because of the ridiculous way lpass sync works we will need to retry until we get our ID.
	// see open issue at https://github.com/lastpass/lastpass-cli/issues/450
	for i := 0; i < 10; i++ {
        // Still times out we creating several secrets at a time
		time.Sleep(time.Second * 4) // double the waiting time to 4s per iteration
		errbuf.Reset()
		outbuf.Reset()
		cmd = exec.Command("lpass", "sync")
		cmd.Stderr = &errbuf
		err = cmd.Run()
		if err != nil {
			var err = errors.New(errbuf.String())
			return s, err
		}
		cmd = exec.Command("lpass", "show", "--sync=now", s.Name, "--json", "-x")
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf
		err = cmd.Run()
		if err != nil {
			if !strings.Contains(errbuf.String(), "Could not find specified account") {
				var err = errors.New(errbuf.String())
				return s, err
			}
			continue
		}
		err = json.Unmarshal(outbuf.Bytes(), &secrets)
		if err != nil {
			return s, err
		}
		if len(secrets) > 1 {
			err := errors.New("more than one secret with same name, unable to determine ID")
			return s, err
		}
		if secrets[0].ID == "0" {
			// sync is still not done with upstream.
			continue
		}
		return secrets[0], nil
	}
	err = errors.New("timeout, unable to create new secret")
	return s, err
}
