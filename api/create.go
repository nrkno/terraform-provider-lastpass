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
	secrets, err := c.read(s.Name)
	if err == nil {
	  // We find a secret by the same name - don't create it again.
	  // Note: this deviates from default lastpass behavior which allows multiple secrets of the same name
	  var err = errors.New("Secret " + s.Name + " exists already - try import - aborting")
	  var secret = secrets[0]
	  return secret, err
	}
	cmd := exec.Command("lpass", "add", s.Name, "--non-interactive", "--sync=now")
	return c.create(s.Name, template, cmd)
}

// Create a secret of type node-type
func (c *Client) CreateNodeType(name string, template string, nodetype string) (Secret, error) {
	secrets, err := c.read(name)
	if err == nil {
	  // We find a secret by the same name - don't create it again.
	  // Note: this deviates from default lastpass behavior which allows multiple secrets of the same name
	  var err = errors.New("Secret " + name + " exists already - try import - aborting")
	  var secret = secrets[0]
	  return secret, err
	}
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
	// 7/22/21: Create continues to timeout during bulk operations, or when client machine is
	// busy with high CPU/network - for example, when creating 10 secrets, even with terraform
	// -parallelism=1, we always timeout. The secret is typically eventually created in lastpass
	// but not registered in the terraform state-file, which will give us an error on the next
	// try because when we look it up we get two secrets with the same name. I will increase
	// the timeout retry, and I will disable creation of secrets with the same name.
	const retry_count = 20 // a number less than infinite, but large enough to not timeout
	for i := 0; i < retry_count; i++ {
		// Still times out we creating several secrets at a time
		// increase wait time between calls from 4s to 30s per iteration to
		// avoid potential rate-limiting by lastpass server
		time.Sleep(time.Second * 30)
		errbuf.Reset()
		outbuf.Reset()
		cmd = exec.Command("lpass", "sync")
		cmd.Stderr = &errbuf
		err = cmd.Run()
		if err != nil {
			var err = errors.New(errbuf.String())
			return s, err
		}
		cmd = exec.Command("lpass", "show", "--sync=now", "-G", s.Name, "--json", "-x")
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
	err = errors.New("timeout, unable to create new secret " + s.Name)
	return s, err
}
