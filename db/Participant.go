package db

type Participant struct {
	ID       string
	Username string
}

func (participant *Participant) Create() {
	fsc, ctx := getClient()

	fsc.NewRequest().CreateEntities(ctx, participant)()
}
