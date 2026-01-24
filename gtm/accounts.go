package gtm

import (
	"context"

	tagmanager "google.golang.org/api/tagmanager/v2"
)

// Account is a simplified representation of a GTM account.
type Account struct {
	AccountID string `json:"accountId"`
	Name      string `json:"name"`
	Path      string `json:"path"`
}

// ListAccounts returns all GTM accounts accessible to the authenticated user.
// If account restriction is enabled (via GTM_ACCOUNT_ID), only the restricted account is returned.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	resp, err := c.Service.Accounts.List().Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	accounts := toAccounts(resp.Account)

	// Filter to restricted account if configured
	if c.RestrictedAccountID != "" {
		filtered := make([]Account, 0, 1)
		for _, a := range accounts {
			if a.AccountID == c.RestrictedAccountID {
				filtered = append(filtered, a)
				break
			}
		}
		return filtered, nil
	}

	return accounts, nil
}

func toAccounts(accounts []*tagmanager.Account) []Account {
	result := make([]Account, 0, len(accounts))
	for _, a := range accounts {
		result = append(result, Account{
			AccountID: a.AccountId,
			Name:      a.Name,
			Path:      a.Path,
		})
	}
	return result
}
