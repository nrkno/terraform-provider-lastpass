package api

import (
	"errors"
)

// Fetch secrets from upstream
func (c *Client) Read(id string) (*Secret, error) {
	for _, account := range c.Accounts {
		if account.ID == id {
			modifiedGMT, err := epochToTime(account.LastModifiedGMT)
			if err != nil {
				return nil, err
			}
			lastTouch, err := epochToTime(account.LastTouch)
			if err != nil {
				return nil, err
			}
			secret := Secret{
				ID:              account.ID,
				Name:            account.Name,
				Username:        account.Username,
				Password:        account.Password,
				URL:             account.URL,
				Group:           account.Group,
				Share:           account.Share,
				Notes:           account.Notes,
				LastModifiedGmt: modifiedGMT,
				LastTouch:       lastTouch,
			}
			secret.genCustomFields()
			return &secret, nil
		}
	}
	return nil, errors.New("Could not find specified secret")
}
