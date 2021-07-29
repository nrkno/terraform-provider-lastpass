package api

import (
	"context"

	"github.com/ansd/lastpass-go"
)

// Create is used to create a new resource and generate ID.
func (c *Client) Create(s *Secret) error {
	a := &lastpass.Account{
		Name:     s.Name,
		Username: s.Username,
		Password: s.Password,
		URL:      s.URL,
		Group:    s.Group,
		Notes:    s.Note,
	}
	err := c.Client.Add(context.Background(), a)
	if err != nil {
		return err
	}
	s.ID = a.ID
	err = c.Sync()
	if err != nil {
		return err
	}
	return nil
}
