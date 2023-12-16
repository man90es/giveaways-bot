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

	announcementLn1 := fmt.Sprintf("%v has started, the prize is: %v", dcEvent.Name, prize.Name)
	announcementLn2 := fmt.Sprintf("The winner will be selected in %v", time.Until(*dcEvent.ScheduledEndTime).Round(time.Minute))
	session.ChannelMessageSend(db.Config[db.ConfigKeyGiveawayChannelID], announcementLn1+"\n"+announcementLn2)
}
