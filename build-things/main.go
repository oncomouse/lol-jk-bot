package main

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "io/ioutil"

  _ "github.com/mattn/go-sqlite3"
)

func main() {
  contents, err := ioutil.ReadFile("../things.json")
  if err != nil {
    fmt.Println("IO Error")
    panic(err)
  }
  var things []string
  if err := json.Unmarshal(contents, &things); err != nil {
    fmt.Println("JSON Error")
    panic(err)
  }
  // Open DB connection:
  database,err := sql.Open("sqlite3", "../data.db")
  if err != nil {
    panic(err)
  }
  statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS things (id INTEGER PRIMARY KEY, thing TEXT)")
  statement.Exec()
  for _,thing := range things {
    // Add the thing to the things table:
    statement,err = database.Prepare("INSERT INTO things (thing) VALUES (?)")
    if err != nil {
      panic(err)
    }
    statement.Exec(thing)
  }
  database.Close()
}
