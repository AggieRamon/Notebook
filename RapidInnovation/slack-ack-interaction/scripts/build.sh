#!/usr/bin/bash
go build *.go && zip dist/slack-ack-interaction.zip main && aws lambda update-function-code --function-name slack-ack-interaction --zip-file fileb://dist/slack-ack-interaction.zip && rm main
