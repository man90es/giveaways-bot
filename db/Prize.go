package db

type Prize struct {
	EventID  string
	ID       string
	Name     string
	WinnerID string
}

func (prize *Prize) Insert() {
	fsc, ctx := getClient()

	fsc.NewRequest().CreateEntities(ctx, prize)()
}

func (prize *Prize) AssignWinner(winnerID string) (err error) {
	fsc, ctx := getClient()

	_, err = fsc.NewRequest().GetEntities(ctx, prize)()
	if err != nil {
		return
	}

	prize.WinnerID = winnerID
	err = fsc.NewRequest().UpdateEntities(ctx, prize)()

	return
}

func (prize *Prize) AssignEvent(eventID string) (err error) {
	fsc, ctx := getClient()

	_, err = fsc.NewRequest().GetEntities(ctx, prize)()
	if err != nil {
		return
	}

	prize.EventID = eventID
	err = fsc.NewRequest().UpdateEntities(ctx, prize)()

	return
}

func NewPrize(name string) (prize *Prize) {
	prize = &Prize{Name: name}

	return
}

func GetAvailablePrizes() (prizes []Prize, err error) {
	fsc, ctx := getClient()

	query := fsc.Client.Collection("Prize").Where("winnerid", "==", "")
	err = fsc.NewRequest().QueryEntities(ctx, query, &prizes)()

	return
}

func GetPrizeSelectedForEvent(eventID string) (prize *Prize, err error) {
	fsc, ctx := getClient()

	result := make([]Prize, 0)
	query := fsc.Client.Collection("Prize").Where("eventid", "==", eventID)
	err = fsc.NewRequest().QueryEntities(ctx, query, &result)()

	if err != nil {
		return
	}

	foundPrize := result[0]
	return &foundPrize, nil
}

func GetPastWinnerIDs() (pastWinnerIDs []string, err error) {
	fsc, ctx := getClient()

	wonPrizes := make([]Prize, 0)
	prizesQuery := fsc.Client.Collection("Prize").Where("winnerid", "!=", "")
	err = fsc.NewRequest().QueryEntities(ctx, prizesQuery, &wonPrizes)()
	if err != nil {
		return
	}

	for _, wonPrize := range wonPrizes {
		pastWinnerIDs = append(pastWinnerIDs, wonPrize.WinnerID)
	}

	return
}
