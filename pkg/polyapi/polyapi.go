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
	BaseDataURL  string
	BaseGammaURL string
	HTTPClient   *http.Client
}

func NewPolyapi() *Polyapi {
	return &Polyapi{
		BaseDataURL:  "https://data-api.polymarket.com",
		BaseGammaURL: "https://gamma-api.polymarket.com",
		HTTPClient:   http.DefaultClient,
	}
}

// https://docs.polymarket.com/api-reference/core/get-top-holders-for-markets?playground=open
// https://data-api.polymarket.com/holders?limit=20&minBalance=1&market=0xfc613d67a16f3a9a10b63baa0f48cee855d49310b33643112e43f769d68b80a5
func (c *Polyapi) GetTopHolders(ctx context.Context, market string) (*TokenHolders, error) {
	u, err := url.Parse(c.BaseDataURL + "/holders")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("limit", "20") // 20 is the max allowed by the API
	q.Set("minBalance", "1")
	q.Set("market", market)
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
		return nil, c.handleError(resp)
	}

	var tokenHolders TokenHolders
	if err := json.NewDecoder(resp.Body).Decode(&tokenHolders); err != nil {
		return nil, err
	}

	return &tokenHolders, nil
}

// https://docs.polymarket.com/api-reference/core/get-top-holders-for-markets?playground=open
// https://gamma-api.polymarket.com/markets/slug/will-the-steam-machine-cost-700-or-more-at-release
func (c *Polyapi) GetMarketBySlug(ctx context.Context, slug string) (*Market, error) {
	u, err := url.Parse(c.BaseGammaURL + "/markets/slug/" + slug)
	if err != nil {
		return nil, err
	}

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
		return nil, c.handleError(resp)
	}

	var market Market
	if err := json.NewDecoder(resp.Body).Decode(&market); err != nil {
		return nil, err
	}

	return &market, nil
}

// https://docs.polymarket.com/api-reference/core/get-closed-positions-for-a-user
// https://data-api.polymarket.com/closed-positions?limit=10&sortBy=REALIZEDPNL&sortDirection=DESC&user=0x3a6EFc8104f17068a8B08360518B0618c4e53291
func (c *Polyapi) GetClosedPositions(ctx context.Context, user string, limit, offset int) ([]ClosedPosition, error) {
	u, err := url.Parse(c.BaseDataURL + "/closed-positions")
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
		return nil, c.handleError(resp)
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

func (c *Polyapi) handleError(resp *http.Response) error {
	var apiErr PolyError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error() != "" {
		return &apiErr
	}

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
