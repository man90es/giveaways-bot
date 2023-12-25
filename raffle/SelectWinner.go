package raffle

import (
	"fmt"
	"log"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/man90es/giveaways-bot/db"
	"github.com/man90es/giveaways-bot/utils"
)

func SelectWinner(session *discordgo.Session, dcEvent *discordgo.GuildScheduledEvent) {
	members, err := session.GuildMembers(db.Config[db.ConfigKeyGiveawayGuildlID], "", 1e3)
	if err != nil {
		log.Println("Error occured while trying to retrieve server members: ", err.Error())
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

	prize, err := db.GetPrizeSelectedForEvent(dcEvent.ID)
	if err != nil {
		log.Println("Error occured while trying to retrieve a selected prize: ", err.Error())
		return
	}

	err = prize.AssignWinner(winner.User.ID)
	if err != nil {
		log.Println("Error occured while trying to assign a winner to the prize: ", err.Error())
		return
	}

	announcementLn1 := fmt.Sprintf("%v winner is %v, congrats!", dcEvent.Name, winner.Mention())
	announcementLn2 := fmt.Sprintf("Please DM <@!%v> to claim your %v", db.Config[db.ConfigKeyGiveawayOrganiserID], prize.Name)
	session.ChannelMessageSend(db.Config[db.ConfigKeyGiveawayChannelID], announcementLn1+"\n"+announcementLn2)
}
