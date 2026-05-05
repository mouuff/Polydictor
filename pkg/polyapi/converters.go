package polyapi

import (
	"encoding/json"
	"fmt"
)

func (m *Market) UnmarshalJSON(data []byte) error {
	type Alias Market // prevent recursion

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse Outcomes
	if m.OutcomesRaw != "" {
		if err := json.Unmarshal([]byte(m.OutcomesRaw), &m.Outcomes); err != nil {
			return fmt.Errorf("failed to parse outcomes: %w", err)
		}
	}

	// Parse ClobTokenIds
	if m.ClobTokenIdsRaw != "" {
		if err := json.Unmarshal([]byte(m.ClobTokenIdsRaw), &m.ClobTokenIds); err != nil {
			return fmt.Errorf("failed to parse clobTokenIds: %w", err)
		}
	}

	return nil
}
