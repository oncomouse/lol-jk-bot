package main
import (
  "database/sql"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "path"
  "path/filepath"
  "strings"

  "github.com/deckarep/golang-set"
  _ "github.com/mattn/go-sqlite3"
)

func main() {
  contents, err := ioutil.ReadFile("../acronyms.json")
  if err != nil {
    fmt.Println("IO Error")
    panic(err)
  }
  var acronyms []string
  if err := json.Unmarshal(contents, &acronyms); err != nil {
    fmt.Println("JSON Error")
    panic(err)
  }
  contents, err = ioutil.ReadFile("../stop-words.json")
  if err != nil {
    fmt.Println("IO Error")
    panic(err)
  }
  var stop_words_array []string
  if err := json.Unmarshal(contents, &stop_words_array); err != nil {
    fmt.Println("JSON Error")
    panic(err)
  }
  stop_words := mapset.NewSet()
  for _,sw := range stop_words_array {
    stop_words.Add(sw)
  }
  // Open DB connection:
  database,err := sql.Open("sqlite3", "../data.db")
  if err != nil {
    panic(err)
  }
  commands := make(map [string]*sql.Stmt)
  sets := make(map [string]mapset.Set)
  statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS acronyms (id INTEGER PRIMARY KEY, acronym TEXT)")
  statement.Exec()
  for _,acronym := range acronyms {
    // Create the table:
    statement, err := database.Prepare(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, phrase TEXT)", acronym))
    if err != nil {
      panic(err)
    }
    statement.Exec()
    // Add the acronym to the acronyms table:
    statement,err = database.Prepare("INSERT INTO acronyms (acronym) VALUES (?)")
    if err != nil {
      panic(err)
    }
    statement.Exec(acronym)
    // Cache a command to add to the new table:
    commands[acronym],err = database.Prepare(fmt.Sprintf("INSERT INTO %s (phrase) VALUES (?)", acronym))
    if err != nil {
      panic(err)
    }
    // Create a set for the acronym:
    sets[acronym] = mapset.NewSet()
  }
  // Process JSON files:
  json_files,_ := filepath.Glob("../data/*.json")
  for _,file := range json_files {
    fmt.Printf("Processing %s…", path.Base(file))
    var dat map[string]interface{}
    byt,_ := ioutil.ReadFile(file)
    if err := json.Unmarshal(byt, &dat); err != nil {
      panic(err)
    }
    for acronym, phrases := range dat {
      for _, phrase := range phrases.([]interface{}) {
        words := strings.Split(phrase.(string), " ")
        words_set := mapset.NewSet()
        for _,word := range words {
          words_set.Add(word)
        }
        if words_set.Union(stop_words).Cardinality() < words_set.Cardinality() - 1 {
          sets[acronym].Add(phrase)
        }
      }
    }
    fmt.Println("done")
  }
  for acronym, set := range sets {
    it := set.Iterator()
    for phrase := range it.C {
      fmt.Printf("Loading table %s with %s…", acronym, phrase)
      commands[acronym].Exec(phrase)
      fmt.Println("done")
    }
  }
  database.Close()
}
