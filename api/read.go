package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// Fetch secrets from upstream
func (c *Client) Read(id string) ([]Secret, error) {
	var secrets []Secret
	err := c.login()
	if err != nil {
		return secrets, err
	}
	cmd := exec.Command("lpass", "show", "--sync=auto", "-G", id, "--json", "-x")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		// Make sure the secret is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			// return empty secret list
			return secrets, nil
		}
		var err = errors.New(errbuf.String())
		return secrets, err
	}
	err = json.Unmarshal(outbuf.Bytes(), &secrets)
	if err != nil {
		return secrets, err
	}
	for i := range secrets {
		secrets[i].genCustomFields()
		secrets[i].Name = secrets[i].Fullname // lastpass trims path from name, so we need to copy fullname
	}
	return secrets, nil
}
