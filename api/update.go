package api

import (
	"context"

	"github.com/ansd/lastpass-go"
)

// Update is called to update secret with upstream
func (c *Client) Update(s *Secret) error {
	account := &lastpass.Account{
		ID:       s.ID,
		Name:     s.Name,
		Username: s.Username,
		Password: s.Password,
		URL:      s.URL,
		Group:    s.Group,
		Share:    s.Share,
		Notes:    s.Notes,
	}
	err := c.Client.Update(context.Background(), account)
	if err != nil {
		return err
	}
	err = c.Sync()
	if err != nil {
		return err
	}
	return nil
}
