#!/usr/bin/bash
go build main.go && zip dist/slack-ack-command.zip main && aws lambda update-function-code --function-name slack-ack-command --zip-file fileb://dist/slack-ack-command.zip && rm main
