package polyapi

type User struct {
	UserId          string
	ClosedPositions []ClosedPosition
}

func (u *User) GetPredictionRate() float64 {
	var totalPositions float64
	var profitablePositions float64

	for _, p := range u.ClosedPositions {
		if p.RealizedPnl > 0 {
			profitablePositions += 1
		}

		totalPositions += 1
	}

	if totalPositions == 0 {
		return 0
	}

	return profitablePositions / totalPositions
}
