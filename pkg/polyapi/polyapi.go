package polyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Polyapi struct {
	BaseDataURL  string
	BaseGammaURL string
	HTTPClient   *http.Client

	// Data API limiters
	dataGeneralLimiter         *rate.Limiter // 1000 / 10s = 100 rps
	dataClosedPositionsLimiter *rate.Limiter // 150 / 10s = 15 rps

	// Gamma API limiters
	gammaGeneralLimiter *rate.Limiter // 4000 / 10s = 400 rps
	gammaMarketsLimiter *rate.Limiter // 300 / 10s = 30 rps
}

func NewPolyapi() *Polyapi {
	return &Polyapi{
		BaseDataURL:  "https://data-api.polymarket.com",
		BaseGammaURL: "https://gamma-api.polymarket.com",
		HTTPClient:   http.DefaultClient,

		// Data API
		dataGeneralLimiter:         rate.NewLimiter(rate.Every(time.Second/100), 100),
		dataClosedPositionsLimiter: rate.NewLimiter(rate.Every(time.Second/15), 15),

		// Gamma API
		gammaGeneralLimiter: rate.NewLimiter(rate.Every(time.Second/400), 400),
		gammaMarketsLimiter: rate.NewLimiter(rate.Every(time.Second/30), 30),
	}
}

// --------------------
// Internal helpers
// --------------------

func (c *Polyapi) doRequest(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {

	limiters := c.limitersFor(req)

	for _, l := range limiters {
		if err := l.Wait(ctx); err != nil {
			return nil, err
		}
	}

	return c.HTTPClient.Do(req)
}

func (c *Polyapi) limitersFor(req *http.Request) []*rate.Limiter {
	host := req.URL.Host
	path := req.URL.Path

	if strings.Contains(host, "data-api") {
		if strings.HasPrefix(path, "/closed-positions") {
			return []*rate.Limiter{
				c.dataGeneralLimiter,
				c.dataClosedPositionsLimiter,
			}
		}
		return []*rate.Limiter{
			c.dataGeneralLimiter,
		}
	}

	if strings.Contains(host, "gamma-api") {
		if strings.HasPrefix(path, "/markets") {
			return []*rate.Limiter{
				c.gammaGeneralLimiter,
				c.gammaMarketsLimiter,
			}
		}
		return []*rate.Limiter{
			c.gammaGeneralLimiter,
		}
	}

	return nil
}

// handleError tries to parse API error response, otherwise returns generic error with status code
func (c *Polyapi) handleError(resp *http.Response) error {
	var apiErr PolyError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error() != "" {
		return &apiErr
	}

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// --------------------
// API methods
// --------------------

// https://docs.polymarket.com/api-reference/core/get-top-holders-for-markets?playground=open
func (c *Polyapi) GetTopHolders(ctx context.Context, market string) (*[]TokenHolderGroup, error) {
	u, err := url.Parse(c.BaseDataURL + "/v1/holders")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("limit", "20")
	q.Set("minBalance", "1")
	q.Set("market", market)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp)
	}

	var tokenHolders []TokenHolderGroup
	if err := json.NewDecoder(resp.Body).Decode(&tokenHolders); err != nil {
		return nil, err
	}

	return &tokenHolders, nil
}

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

	resp, err := c.doRequest(ctx, req)
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
func (c *Polyapi) GetClosedPositions(ctx context.Context, user string, limit, offset int) ([]ClosedPosition, error) {
	u, err := url.Parse(c.BaseDataURL + "/v1/closed-positions")
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

	resp, err := c.doRequest(ctx, req)
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

// Auto-pagination
func (c *Polyapi) GetAllClosedPositions(ctx context.Context, user string) ([]ClosedPosition, error) {
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

		allPositions = append(allPositions, positions...)

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
