package main

import (
  "fmt"
  "regexp"
  "strings"
  "time"

  "github.com/nlopes/slack"
)

func messageHandler(api *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent) bool {
  parts := strings.Fields(ev.Text)
  attachment := slack.Attachment{}

  params := slack.PostMessageParameters{}
  switch parts[1] {
  case "stock":
      quoteID := parts[2]
      r := regexp.MustCompile(`<http://.+\|(.+)>`)
      match := r.FindStringSubmatch(quoteID)
      if match != nil {
        quoteID = match[1]
      }
      attachment = slack.Attachment{
        Pretext: getQuote(quoteID),
        Text: "",
        ImageURL: fmt.Sprintf("https://chart.finance.yahoo.com/t?s=%s&lang=zh-Hant-HK&region=HK&width=300&height=180", quoteID),
      }
  case "weather":
      var city string
      if len(parts) > 2 {
        city = ev.Text[strings.Index(ev.Text, parts[2]):len(ev.Text)]
      }

      if weather,err := getWeather(city); err == nil {
        // Waring preprocess
        warnings := ""
        for _, warning := range weather.Warnings {
          warnings += warning.String()
        }
        // Forecast preprocess
        forecasts := ""
        cnt := 0
        curTimestamp := int(time.Now().Unix())
        for _, forecast := range weather.Forecast {
          if cnt < 7 && forecast.Dt > curTimestamp {
            forecasts += forecast.String()
            cnt++
          }
        }

        attachment = slack.Attachment{
          Color:    "#40E5F7",
          Pretext: weather.CurrentWeather.Name + "'s Current Weather",
          Text: WeatherIcon(weather.CurrentWeather.Weather[0].Icon) + " " + weather.CurrentWeather.Weather[0].Description + "\n" + warnings,
          // Uncomment the following part to send a field too
          Fields: []slack.AttachmentField{
            slack.AttachmentField{
              Title: "Temperture",
              Value: fmt.Sprintf("%.1f°C", weather.CurrentWeather.Main.Temp),
              Short: true,
            },
            slack.AttachmentField{
              Title: "Humidity",
              Value: fmt.Sprintf("%.1f%%", weather.CurrentWeather.Main.Humidity),
              Short: true,
            },
            slack.AttachmentField{
              Title: "Highest Temp.",
              Value: fmt.Sprintf("%.1f°C", weather.CurrentWeather.Main.TempMax),
              Short: true,
            },
            slack.AttachmentField{
              Title: "Lowest Temp.",
              Value: fmt.Sprintf("%.1f°C", weather.CurrentWeather.Main.TempMin),
              Short: true,
            },
            slack.AttachmentField{
              Title: "7 Day Forecast",
              Value: forecasts,
            },
          },
        }
      } else {
        attachment = slack.Attachment{
          Pretext: fmt.Sprintf("error: %v", err),
        }
      }
  case "marksix":
      attachment = slack.Attachment{
        Pretext: "Searching for marksix...",
        Text: "",
        // Uncomment the following part to send a field too
        // Fields: []slack.AttachmentField{
        //   slack.AttachmentField{
        //     Title: "TITLE",
        //     Value: "no",
        //   },
        // },
      }
  default:
    return false
  }
  params.Attachments = []slack.Attachment{attachment}
  api.PostMessage(ev.Channel, "", params)
  return true
}