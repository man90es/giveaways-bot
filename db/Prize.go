package db

type Prize struct {
	ID       string
	Name     string
	WinnerID string
}

func (prize *Prize) Create() {
	fsc, ctx := getClient()

	fsc.NewRequest().CreateEntities(ctx, prize)()
}

func (prize *Prize) AssignWinner(winnerID string) {
	fsc, ctx := getClient()

	fsc.NewRequest().GetEntities(ctx, prize)()
	prize.WinnerID = winnerID
	fsc.NewRequest().UpdateEntities(ctx, prize)()
}

func GetAvailablePrizes() []Prize {
	fsc, ctx := getClient()

	query := fsc.Client.Collection("Prize").Where("winnerid", "==", "")

	result := make([]Prize, 0)
	_ = fsc.NewRequest().QueryEntities(ctx, query, &result)()

	return result
}

func GetPastWinnerIDs() []string {
	fsc, ctx := getClient()

	prizesQuery := fsc.Client.Collection("Prize").Where("winnerid", "!=", "")
	wonPrizes := make([]Prize, 0)
	_ = fsc.NewRequest().QueryEntities(ctx, prizesQuery, &wonPrizes)()

	pastWinnerIDs := []string{}
	for _, wonPrize := range wonPrizes {
		pastWinnerIDs = append(pastWinnerIDs, wonPrize.WinnerID)
	}

	return pastWinnerIDs
}
