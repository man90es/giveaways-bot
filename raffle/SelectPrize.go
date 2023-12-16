package raffle

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/man90es/giveaways-bot/db"
	"github.com/man90es/giveaways-bot/utils"
)

func SelectPrize(session *discordgo.Session, dcEvent *discordgo.GuildScheduledEvent) {
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

	prize.AssignEvent(dcEvent.ID)

	delay := time.Until(*dcEvent.ScheduledEndTime)
	announcement := fmt.Sprintf("%v starts now, the prize is: %v\nThe winner will be announced in %v", dcEvent.Name, prize.Name, delay)
	session.ChannelMessageSend(db.Config[db.ConfigKeyGiveawayChannelID], announcement)
}
