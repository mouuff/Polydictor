package polyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Polyapi struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewPolyapi() *Polyapi {
	return &Polyapi{
		BaseURL:    "https://data-api.polymarket.com",
		HTTPClient: http.DefaultClient,
	}
}

// TODO
// https://docs.polymarket.com/api-reference/core/get-top-holders-for-markets?playground=open

// https://docs.polymarket.com/api-reference/core/get-closed-positions-for-a-user
// https://data-api.polymarket.com/closed-positions?limit=10&sortBy=REALIZEDPNL&sortDirection=DESC&user=0x3a6EFc8104f17068a8B08360518B0618c4e53291
func (c *Polyapi) GetClosedPositions(ctx context.Context, user string, limit, offset int) ([]ClosedPosition, error) {
	u, err := url.Parse(c.BaseURL + "/closed-positions")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("user", user)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("sortBy", "TIMESTAMP")
	q.Set("sortDirection", "DESC")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr PolyError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error() != "" {
			return nil, &apiErr
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var positions []ClosedPosition
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

// Fetch ALL closed positions (auto-pagination)
func (c *Polyapi) GetAllClosedPositions(ctx context.Context, user string) ([]ClosedPosition, error) {
	// Required range: 0 <= x <= 50
	const limit = 50

	var (
		allPositions []ClosedPosition
		offset       = 0
	)

	for {
		positions, err := c.GetClosedPositions(ctx, user, limit, offset)
		if err != nil {
			return nil, err
		}

		// Append results
		allPositions = append(allPositions, positions...)

		// Stop condition: last page
		if len(positions) < limit {
			break
		}

		offset += limit
	}

	return allPositions, nil
}

func (c *Polyapi) GetUser(ctx context.Context, user string) (*User, error) {
	closedPositions, err := c.GetAllClosedPositions(ctx, user)
	if err != nil {
		return nil, err
	}

	return &User{
		UserId:          user,
		ClosedPositions: closedPositions,
	}, nil
}
