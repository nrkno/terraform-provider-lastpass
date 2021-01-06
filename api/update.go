package api

import (
	"bytes"
	"errors"
	"os/exec"
)

// Update is called to update secret with upstream
func (c *Client) Update(s Secret) error {
    template := s.getTemplate()
    return c.update(s.ID, template)
}

// Update is called to update secret with upstream
func (c *Client) UpdateNodeType(id string, template string) error {
    return c.update(id, template)
}

func (c *Client) update(id string, template string) error {
	err := c.login()
	if err != nil {
		return err
	}
	cmd := exec.Command("lpass", "edit", id, "--non-interactive", "--sync=now")
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
