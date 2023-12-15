package raffle

import (
	"log"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/man90es/giveaways-bot/db"
	"github.com/man90es/giveaways-bot/utils"
)

func RunRaffle(session *discordgo.Session, eventID string) {
	availablePrizes, err := db.GetAvailablePrizes()
	if err != nil {
		log.Println("Error occured while trying to retrieve prizes from the DB: ", err.Error())
		return
	}

	prize, err := utils.RandomChoice(availablePrizes)
	if err != nil {
		log.Println("Error occured while trying to select a prize: ", err.Error())
		return
	}

	(&db.Prize{ID: prize.ID}).AssignEvent(eventID)

	members, err := session.GuildMembers(db.Config[db.ConfigKeyGiveawayGuildlID], "", 1e3)
	if err != nil {
		log.Println(err.Error())
		return
	}

	pastWinnerIDs, err := db.GetPastWinnerIDs()
	if err != nil {
		log.Println("Error occured while trying to retrieve past winners' IDs: ", err.Error())
		return
	}

	// Select participants
	participants := []*discordgo.Member{}
	for _, member := range members {
		// User already won something
		if slices.Contains(pastWinnerIDs, member.User.ID) {
			continue
		}

		// User doesn't have a required role
		if !slices.Contains(member.Roles, db.Config[db.ConfigKeyGiveawayRoleID]) {
			continue
		}

		db.NewParticipantFromDiscordUser(member.User).Upsert()
		participants = append(participants, member)
	}

	winner, err := utils.RandomChoice(participants)
	if err != nil {
		log.Println("Error occured while trying to select a winner: ", err.Error())
		return
	}

	(&db.Prize{ID: prize.ID}).AssignWinner(winner.User.ID)
	session.ChannelMessageSend(db.Config[db.ConfigKeyGiveawayChannelID], winner.Mention()+" won "+prize.Name)
}
