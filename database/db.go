package database

import (
	"bey/go-vct/common"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath" // Correct import for filepath

	_ "github.com/mattn/go-sqlite3" // replace with the import path for your SQL driver
)

func InitDB(databasePath string) {
	// Check if the database file exists
	_, err := os.Stat(databasePath)
	if os.IsNotExist(err) {
		// The database file doesn't exist. Initialize it.
		log.Println("Initializing Database")
		// Create the directory if it doesn't exist
		dir := filepath.Dir(databasePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Printf("Database Does Not Exists. Creating DB Path: %s \n", databasePath)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		db, err := sql.Open("sqlite3", databasePath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		sqlStmt := `
		CREATE TABLE messages (
			id INTEGER NOT NULL PRIMARY KEY, 
			message_id TEXT, 
			match_id TEXT, 
			announcement_sent BOOLEAN, 
			starting_sent BOOLEAN, 
			result_sent BOOLEAN
		);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
		}
	} else if err != nil {
		// Some other error occurred when trying to check the file
		log.Fatal(err)
	}
}

func GetSentMessages() ([]common.Message, error) {
	db, err := sql.Open("sqlite3", common.DbPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM messages")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var messages []common.Message

	for rows.Next() {
		var message common.Message

		err = rows.Scan(&message.Id, &message.MessageId, &message.MatchId, &message.AnnouncementSent, &message.StartingSent, &message.ResultSent)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return messages, nil
}

func AddSentMessage(matchId int, messageId int) error {
	db, err := sql.Open("sqlite3", common.DbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO messages (match_id, message_id, announcement_sent, starting_sent, result_sent) VALUES (?, ?, ?, ?, ?)",
		matchId, messageId, true, false, false)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func UpdateSentMessage(messageId int, field string) error {
	db, err := sql.Open("sqlite3", common.DbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("UPDATE messages SET %s = ? WHERE message_id = ?", field), true, messageId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
