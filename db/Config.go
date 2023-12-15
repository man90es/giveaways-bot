package db

type configKey string

const (
	ConfigKeyDiscordAPIToken   configKey = "discordapitoken"
	ConfigKeyGiveawayChannelID configKey = "giveawaychannelid"
	ConfigKeyGiveawayGuildlID  configKey = "giveawayguildid"
	ConfigKeyGiveawayRoleID    configKey = "giveawayroleid"
)

type ConfigItem struct {
	ID    string
	Value string
}

var Config = map[configKey]string{}
var configKeys = []configKey{ConfigKeyDiscordAPIToken, ConfigKeyGiveawayChannelID, ConfigKeyGiveawayGuildlID, ConfigKeyGiveawayRoleID}

func LoadConfig() (err error) {
	fsc, ctx := getClient()

	for _, key := range configKeys {
		item := ConfigItem{ID: string(key)}
		_, err = fsc.NewRequest().GetEntities(ctx, &item)()
		if err != nil {
			return err
		}
		Config[key] = item.Value
	}

	return nil
}
