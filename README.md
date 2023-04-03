## Ridepanda Notifier

This is an automatic web scraper I used to check for bike availability on RidePanda's website - so I could get a bike that I wanted.

It uses Gotify (https://gotify.net/) to send notifications to my phone.

Requires a gotify server running on the same device as the notifier.

Also requires a `config.json` file in the following form:

```json5
{
  "ridePandaAPI": { // Browse the website and then inspect browser logs to get these
    "hubId": "",
    "rpFfId": "",
    "rpOrg": "",
    "userAgent": "",
    "serverURL": ""
  },
  "notifications": { // Take from your gotify instance
    "serverURL": "",
    "appToken": ""
  },
  "worker": {
    "checkFrequency": "6h", // how often to check for bikes
    "jitter": "30m", // max random jitter to apply after checkFrequency has elapsed
    "startHour": 9, // hour of time of day to start
    "startMinute": 0, // minute of time of day to start
    "notifyNothing": 28, // Every x attempts, send a notification even if no bikes match rather than being silent like normal to re-assure that the system is still working
  },
  "matcher": {
    "targets": [
      {
        "name": "Diamondback Union 1", // Arbitrary name shown in the notification
        "criteria": {
          "name": "Diamondback Union 1", // Bike name
          "size": "L", // Bike size
          "color": "*" // Bike color
        }
      }
    ]
  }
}
```
