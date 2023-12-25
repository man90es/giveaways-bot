# Giveaways bot

A little Discord bot for making giveaways of random things to random server members. I quickly threw it together before Xmas to make a Steam keys giveaway for my friends, so it's rather hacky, but shouldn't have any breaking bugs. I may improve the code and add more features if I do any other giveaways in the future.

## Principle of work
The bot can currently work only with one server. It is triggered by an event in that server. When the event starts, the bot will randomly select and announce the prize from a list you provide. When the event ends, the bot will select a random member with a certain role and announce the winner.

## Detailed instructions
1. Create a Discord application and give it `server members` intent access
1. Add it to your Discord server with a `bot` scope and `Send Messages` permission
1. Create a Firebase project with a Firestore database and download credentials with access to it in .json format
1. In the database, create a collection named `ConfigItem`. Each item should have a strictly defined `id` and a single property called value `value`. All of them are mandatory and should be filled like this:
	- `discordapitoken`: Token for your Discord bot application created in step 1
	- `giveawaychannelid`: ID of a channel in which the bot will make announcements
	- `giveawayguildid`: ID of a guild in which the bot will conduct a giveaway
	- `giveawayorganiserid`: ID of a person that will be mentioned to winners like "Please DM @organiser to claim your prize"
	- `giveawayroleid`: ID of a role that all participants have. Server members without this role can't win
1. In the same database, create a collection named `Prize`. You can add as many prizes as you want here. IDs don't matter here but each item should have 3 fields:
	- eventid: empty string
	- name: non-empty string, put the name of the prize here
	- winnerid: empty string
1. Clone the repo and compile it by running `go build .`
1. Move the resulting binary file to your VPS
1. Place Firebase credentials .json file generated in step 3 alongside the binary and change its name to `firebase-credentials.json`
1. Run the binary and leave it like that
1. Create one or more Discord events in your server
1. Go back to the database, it should have an `Event` collection automatically created. Find the event(s) you created and switch theier `triggersraffle` property from false to true

## Contributing
All pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
