#!/bin/zsh

GOOS=linux GOARCH=arm go build -o bikep cmd/main.go

scp bikep pi:/home/pi/bike/
scp config.json pi:/home/pi/bike/
scp start.sh pi:/home/pi/bike
scp stop.sh pi:/home/pi/bike