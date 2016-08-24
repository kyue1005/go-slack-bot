package main

import (
  "encoding/csv"
  "fmt"
  "net/http"
  "strings"
)

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
  sym = strings.ToUpper(sym)
  url := fmt.Sprintf("http://finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op", sym)
  resp, err := http.Get(url)
  if err != nil {
    return fmt.Sprintf("error: %v", err)
  }
  rows, err := csv.NewReader(resp.Body).ReadAll()
  if err != nil {
    return fmt.Sprintf("error: %v", err)
  }
  if len(rows) >= 1 && len(rows[0]) == 5 {
    return fmt.Sprintf("%s (<https://hk.finance.yahoo.com/q?s=%s&ql=1|%s>) is trading at $%s", rows[0][0], rows[0][1], rows[0][1], rows[0][2])
  }
  return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}