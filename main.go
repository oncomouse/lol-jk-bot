package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func randomRow(db *sql.DB, table string, keyOtional ...string) *sql.Rows {
	key := "phrase"
	if len(keyOtional) > 0 {
		key = keyOtional[0]
	}
	results, _ := db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE ROWID = (ABS(RANDOM()) %% (SELECT (SELECT MAX(ROWID) FROM %s)+1));", key, table, table))
	return results
}

func makeTweet(subject string, dictionary map[string]string) string {
	// Build the tweet from subject and dictionary:
	tweet := fmt.Sprintf("Is your child texting about %s? Know the signs:\n\n", subject)
	var definitions []string
	for acronym, definition := range dictionary {
		definitions = append(definitions, fmt.Sprintf("%s:%s", acronym, definition))
	}
	tweet = tweet + strings.Join(definitions, "\n")
	return tweet
}

func main() {
	database, _ := sql.Open("sqlite3", "./data.db")

	// Read the acronyms from the database:
	rows, _ := database.Query("SELECT acronym FROM acronyms")
	var acronym string
	acronyms := []string{}
	for rows.Next() {
		rows.Scan(&acronym)
		acronyms = append(acronyms, acronym)
	}

	// Make random choices:
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(acronyms), func(i, j int) { acronyms[i], acronyms[j] = acronyms[j], acronyms[i] })
	numDefinitions := rand.Intn(5) + 5

	// Make sure we always get a tweet less than 280
	var tweet string
	for {
		// Build the dictionary components:
		dictionary := make(map[string]string)
		for _, acronym := range acronyms[:numDefinitions] {
			row := randomRow(database, acronym)
			var definition string
			for row.Next() {
				row.Scan(&definition)
				dictionary[acronym] = definition
			}
		}
		// Choose the subject here:
		var subject string
		row := randomRow(database, "things", "thing")
		for row.Next() {
			row.Scan(&subject)
		}

		tweet = makeTweet(subject, dictionary)
		if len(tweet) < 280 {
			break
		}
	}
	database.Close()

	// We'd wire in Twitter here:
	fmt.Println(tweet)
}
