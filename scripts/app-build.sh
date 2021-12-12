#!/bin/bash
REPO="github.com/duyquang6/git-watchdog"
NOW=$(date +'%Y-%m-%d_%T')

go build -ldflags "-X $REPO/internal/buildinfo.buildID=`git rev-parse --short HEAD` -X $REPO/internal/buildinfo.buildTime=$NOW" -o bin/app cmd/app/main.go
