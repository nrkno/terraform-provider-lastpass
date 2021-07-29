package api

import (
	"context"
)

// Delete secret in upstream db
func (c *Client) Delete(id string) error {
	err := c.Client.Delete(context.Background(), id)
	if err != nil {
		return err
	}
	err = c.Sync()
	if err != nil {
		return err
	}
	return nil
}
