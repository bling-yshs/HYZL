package app_cron

import "github.com/robfig/cron/v3"

var AppCronInstance = cron.New()

var TaskIdEntryIdMap = make(map[int]cron.EntryID)
