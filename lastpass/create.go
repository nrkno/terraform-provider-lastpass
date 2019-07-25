package lastpass

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

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
