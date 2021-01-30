# dogetracker

Simple discord bot that displays price of dogecoin and notifies in increments


## Configuration
Create file called config.json and fill this out (order matters)
```
{
    "Username": "",
    "DogeToken":"DTnt7VZqR5ofHhAxZuDy4m3PhSjKFXpw3e",
    "NomicsAPI":"",
    "DiscordWebhook":"",
    "Intervals": 120,
    "Thresholds": 0.01
}
```
Username: display name of bot
DogeToken: free api key to check amount of dogecoin in account  
NomicsAPI: API used to fetch the price of Dogecoin
DiscordWebhook: the integration from the server to allow bots sending messages
Intervals: The time in between calculations
Thresholds: TODO _ bolder notifications / only notify when the threshold change has passed (I.E., increase/ decrease by a whole cent / dollar, etc.)

## Stuck? heres a hand
https://nomics.com/

Get webhook: no need to worry about the github stuff
https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
