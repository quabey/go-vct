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
			id integer not null primary key,
    	    match_id integer not null constraint messages_pkunique,
    		message_id        integer default 0,
    		announcement_sent integer default 0,
    		starting_sent     integer default 0,
    		result_sent       integer default 0,
    		timestamp         integer default 0
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
	fmt.Println("Getting sent messages...")
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

		err = rows.Scan(&message.Id, &message.MatchId, &message.MessageId, &message.AnnouncementSent, &message.StartingSent, &message.ResultSent, &message.Timestamp)
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

func AddSentMessage(matchId int, messageId int, timestamp int) error {
	db, err := sql.Open("sqlite3", common.DbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO messages (match_id, message_id, announcement_sent, starting_sent, result_sent, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
		matchId, messageId, true, false, false, timestamp)
	if err != nil {
		log.Fatal(err)
		return err
	}

	common.Messages, _ = GetSentMessages()
	return nil
}

func UpdateSentMessage(matchId string, field string) error {
	fmt.Println("Updating sent message")
	db, err := sql.Open("sqlite3", common.DbPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	if field == "starting_sent" {
		_, err = db.Exec("INSERT INTO messages (match_id, starting_sent) VALUES(?, true) ON CONFLICT(match_id) DO UPDATE SET starting_sent = true", matchId)
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else if field == "result_sent" {
		_, err = db.Exec("INSERT INTO messages (match_id, result_sent) VALUES(?, true) ON CONFLICT(match_id) DO UPDATE SET result_sent = true", matchId)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	common.Messages, _ = GetSentMessages()
	fmt.Println("=======================================")
	return nil
}
