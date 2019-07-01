package main

import(
    "encoding/json"
    "fmt"
    "io/ioutil"
    "path"
    "path/filepath"
    "regexp"
    "strings"
)

func Map(vs []string, f func(string) string) []string {
  vsm := make([]string, len(vs))
  for i, v := range vs {
    vsm[i] = f(v)
  }
  return vsm
}

func main() {
  contents, err := ioutil.ReadFile("../acronyms.json")
  if err != nil {
    panic(err)
  }
  var acronyms []string
  if err := json.Unmarshal(contents, &acronyms); err != nil {
    panic(err)
  }
  matchers := make(map[string]*regexp.Regexp)
  for _, acronym := range acronyms {
    pieces := strings.Split(acronym, "")
    matchers[acronym] = regexp.MustCompile("\\s" + strings.Join(
      Map(pieces, func(v string) string {
        return "[" + v + strings.ToLower(v) + "]\\w+"
      }), "\\W+"))
  }
  dirname := "../gutenberg/data/text"
  files,_ := filepath.Glob(dirname + "/*.txt")
  cleaner, _ := regexp.Compile("\\W+")
  for _, file := range files {
  // if true {
    //file := "PG9_text.txt"
    fmt.Printf("Processing %s…", path.Base(file))
    contents, _ := ioutil.ReadFile(file)
    output := make(map[string][]string)
    for acronym, re := range matchers {
      matches := re.FindAll(contents, -1)
      str_matches := []string{}
      for i := uint32(0); i < uint32(len(matches)); i++ {
        str_matches = append(str_matches, cleaner.ReplaceAllString(strings.ToLower(string(matches[i])), " "))
      }
      output[acronym] = str_matches
    }
    jsonOutput,_ := json.Marshal(output)
    ioutil.WriteFile("../data/" + strings.Replace(path.Base(file), ".txt", ".json", -1), []byte(jsonOutput), 0644)
    fmt.Println("done")
  }
}
