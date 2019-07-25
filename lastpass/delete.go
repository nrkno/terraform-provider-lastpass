package lastpass

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

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
