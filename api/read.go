package api

import (
	"errors"
)

// Fetch secrets from upstream
func (c *Client) Read(id string) (*Secret, error) {
	for _, account := range c.Accounts {
		if account.ID == id {
			secret := Secret{
				ID:       account.ID,
				Name:     account.Name,
				Username: account.Username,
				Password: account.Password,
				URL:      account.URL,
				Group:    account.Group,
				Note:     account.Notes,
			}
			secret.genCustomFields()
			return &secret, nil
		}
	}
	return nil, errors.New("Could not find specified secret")
}
