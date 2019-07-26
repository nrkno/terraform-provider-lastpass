package lastpass

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// Fetch record from upstream to update local record
func (c *Client) Read(id string) (Secret, error) {
	var s Secret
	err := c.login()
	if err != nil {
		return s, err
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
			s.ID = "0"
			return s, err
		}
		var err = errors.New(errbuf.String())
		return s, err
	}
	var secrets []Secret
	err = json.Unmarshal(outbuf.Bytes(), &secrets)
	if err != nil {
		return s, err
	}
	if secrets[0].URL == "http://" {
		secrets[0].URL = ""
	}
	secrets[0].Note = secrets[0].Note + "\n" // lastpass trims new line, provokes constant changes.
	return secrets[0], nil
}
