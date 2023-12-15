package db

import "github.com/bwmarrin/discordgo"

type Participant struct {
	ID       string
	Username string
}

func (participant *Participant) Upsert() (err error) {
	fsc, ctx := getClient()

	err = fsc.NewRequest().CreateEntities(ctx, participant)()

	return
}

func NewParticipantFromDiscordUser(dcUser *discordgo.User) (participant *Participant) {
	participant = &Participant{
		ID:       dcUser.ID,
		Username: dcUser.Username,
	}

	return
}
