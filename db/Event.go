package db

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Event struct {
	CreatorID          string
	ID                 string
	Name               string
	ScheduledEndTime   time.Time
	ScheduledStartTime time.Time
	Status             discordgo.GuildScheduledEventStatus
	TriggersRaffle     bool
}

func (event *Event) Upsert() (err error) {
	fsc, ctx := getClient()

	err = fsc.NewRequest().UpdateEntities(ctx, event)()

	return
}

func NewEventFromDiscordEvent(dcEvent *discordgo.GuildScheduledEvent) (event *Event) {
	fsc, ctx := getClient()

	event = &Event{ID: dcEvent.ID}
	fsc.NewRequest().GetEntities(ctx, event)()

	event.CreatorID = dcEvent.CreatorID
	event.Name = dcEvent.Name
	event.ScheduledEndTime = *dcEvent.ScheduledEndTime
	event.ScheduledStartTime = dcEvent.ScheduledStartTime
	event.Status = dcEvent.Status

	return
}

func GetEventByID(ID string) (event *Event, err error) {
	fsc, ctx := getClient()

	event = &Event{ID: ID}
	_, err = fsc.NewRequest().GetEntities(ctx, event)()

	return
}
