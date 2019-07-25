package lastpass

import (
	"bytes"
	"errors"
	"os/exec"
)

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
