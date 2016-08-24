package main

import (
  "regexp"
  // "strings"

  "github.com/nlopes/slack"
)

type bot struct {
  account   *slack.UserDetails
  magicWord *regexp.Regexp
}

func (b *bot) msgInvolved(api *slack.Client, ev *slack.MessageEvent) bool {
  // // Check if the channel is DM type.
  // if strings.Index(ev.Channel, "D") >= 0 {
  //   return true
  // }

  // Check if magic word is hit.
  if b.magicWord.FindStringIndex(ev.Text) != nil {
    return true
  }

  return false
}