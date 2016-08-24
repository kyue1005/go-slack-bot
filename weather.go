package main

import (
  "encoding/json"
  "encoding/xml"
  "errors"
  "fmt"
  // "io/ioutil"
  "net/http"
  "os"
  // "strings"
  "time"
)

// openweathermap settings
const owmap_city_id = "1819729"
const owmap_city_appid = "c64b48956d2fb7650051b70880018ec4"

// Warining from HKO
type Channel struct {
  Title         string `xml:"title"`
  Link          string `xml:"link"`
  Description   string `xml:"description"`
  Items []WarningItem `xml:"item"`
}

type WarningItem struct {
  Title         string `xml:"title"`
  Link          string `xml:"link"`
  Description   string `xml:"description"`
}

// Weather from openweather map
type WeatherItem struct {
  Description string
  Icon string
}

type CurrentWeather struct {
  Name string
  Weather []WeatherItem `json:"weather"`
  Main Main `json:"main"`
  Cod int `json:"cod"`
}

type Main struct {
  Temp float64
  Humidity float64
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`
}

// Forecast from openweather map
type Forecast struct {
  Items []ForecastItem `json:"list"`
  Cod string `json:"cod"`
}

type ForecastItem struct {
  Dt int `json:"dt"`
  Temp ForecastTemp `json:"temp"`
  Weather []WeatherItem `json:"weather"`
}

type ForecastTemp struct {

  TempMin float64 `json:"min"`
  TempMax float64 `json:"max"`
}

// Final Report
type WeatherReport struct {
  Warnings []WarningItem
  CurrentWeather CurrentWeather
  Forecast []ForecastItem
}

func (w WarningItem) String() string {
  return "<" + w.Link + "|" + w.Title + ">"
}

func (f ForecastItem) String() string {
  tm := time.Unix(int64(f.Dt), 0)
  return fmt.Sprintf("*%d/%d:* %s ↑ %.1f°C ↓ %.1f°C %s\n", tm.Day(), tm.Month(), WeatherIcon(f.Weather[0].Icon), f.Temp.TempMax, f.Temp.TempMin, f.Weather[0].Description)
}

func WeatherIcon(iconCode string) string {
  switch iconCode {
  case "01d", "01n":
    return ":sunny:"
  case "02d", "02n":
    return ":mostly_sunny:"
  case "03d", "03n":
    return ":cloud:"
  case "04d", "04n":
    return ":cloud:"
  case "09d", "09n":
    return ":rain_cloud:"
  case "10d", "10n":
    return ":rain_cloud:"
  case "11d", "11n":
    return ":lightning:"
  case "13d", "13n":
    return ":snow_cloud:"
  case "50d", "50n":
    return ":fog:"
  }
  return ""
}

func getWeather(city string) (*WeatherReport, error) {
  var report *WeatherReport
  warningURL := "http://rss.weather.gov.hk/rss/WeatherWarningSummaryv2.xml"
  currentWeatherURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%s&units=metric&appid=%s", owmap_city_id, owmap_city_appid)
  forecastURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast/daily?id=%s&units=metric&cnt=9&appid=%s", owmap_city_id, owmap_city_appid)

  if city != "" {
    currentWeatherURL = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, owmap_city_appid)
    forecastURL = fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast/daily?q=%s&units=metric&cnt=9&appid=%s", city, owmap_city_appid)
  }

  // Get Current weather
  fmt.Println("Get current weather")
  fmt.Println("API Call: " + currentWeatherURL)
  resp, err := http.Get(currentWeatherURL)
  if err != nil {
    fmt.Fprintf(os.Stdout, "error: %v", err)
    return report, err
  }
  defer resp.Body.Close()

  jsonDecoder := json.NewDecoder(resp.Body)

  var curWeather CurrentWeather
  if err = jsonDecoder.Decode(&curWeather); err != nil {
    fmt.Fprintf(os.Stdout, "error: %v", err)
    return report, err
  }

  if (curWeather.Cod != 200) {
    return report, errors.New("Invalid Request")
  }

  // Get Forecast weather
  fmt.Println("Get Forecast weather")
  fmt.Println("API Call: " + forecastURL)
  resp, err = http.Get(forecastURL)
  if err != nil {
    fmt.Fprintf(os.Stdout, "error: %v", err)
    return report, err
  }
  defer resp.Body.Close()

  jsonDecoder = json.NewDecoder(resp.Body)

  var forecast Forecast
  if err = jsonDecoder.Decode(&forecast); err != nil {
    fmt.Fprintf(os.Stdout, "error: %v", err)
    return report, err
  }

  if (forecast.Cod != "200") {
    return report, errors.New("Invalid Request")
  }

  // Get warnings
  var warningRss struct {
    Channel Channel `xml:"channel"`
  }

  if (curWeather.Name == "Hong Kong") {
    fmt.Println("Get weather warnings")
    fmt.Println("API Call: " + warningURL)
    resp, err := http.Get(warningURL)
    if err != nil {
      fmt.Fprintf(os.Stdout, "error: %v", err)
      return report, err
    }
    defer resp.Body.Close()

    xmlDecoder := xml.NewDecoder(resp.Body)

    if err = xmlDecoder.Decode(&warningRss); err != nil {
      fmt.Fprintf(os.Stdout, "error: %v", err)
      return report, err
    }
  }

  report = &WeatherReport{
    Warnings: warningRss.Channel.Items,
    CurrentWeather: curWeather,
    Forecast: forecast.Items,
  }

  return report, nil
}
