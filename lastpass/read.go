package lastpass

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

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
