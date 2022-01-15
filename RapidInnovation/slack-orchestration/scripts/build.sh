#!/usr/bin/bash
go build *.go && zip dist/slack-orchestration.zip main && aws lambda update-function-code --function-name slack-orchestration --zip-file fileb://dist/slack-orchestration.zip && rm main
