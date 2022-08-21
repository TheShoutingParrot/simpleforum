package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func importThreads() (err error) {
	var data []ForumThread

	rows, err := db.Query("SELECT * FROM threads ORDER BY id")
	if err != nil {
		return err
	}

	fmt.Println(rows)

	defer rows.Close()
	for rows.Next() {
		var d ForumThread
		if err := rows.Scan(
				&d.ID, &d.OriginalPoster, &d.Title,
				&d.Content, &d.Date, &d.Votes);
			err != nil {
				return err
			}
		data = append(data, d)
	}

	fmt.Println("Imported data: ", data)

	err = db.QueryRow("select max(id) from threads").Scan(&previousID)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i++ {
		fmt.Println(data[i].ID)
		threads[data[i].ID] = ForumThread{
			ID:		data[i].ID,
			OriginalPoster:	data[i].OriginalPoster,
			Title:		data[i].Title,
			Content:	data[i].Content,
			Votes:		data[i].Votes,
			Date:		data[i].Date,
		}
	}

	return nil
}

func importReplies() (err error) {
	var data []ThreadReplies

	rows, err := db.Query("SELECT * FROM replies ORDER BY thread")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var d ThreadReplies
		if err := rows.Scan(
				&d.ThreadID, &d.Poster, &d.Content,
				&d.Votes);
			err != nil {
				return err
			}
		data = append(data, d)
	}

	fmt.Println("Imported data: ", data)

	for i := 0; i < len(data); i++ {
		fmt.Println(data[i].ThreadID)
		replies[data[i].ThreadID] = ThreadReplies{
			ThreadID:	data[i].ThreadID,
			Poster:		data[i].Poster,
			Content:	data[i].Content,
			Votes:		data[i].Votes,
		}
	}

	return nil
}

func saveThread(t ForumThread) (err error) {
	_, err = db.Exec("INSERT INTO threads (id, op, title, content, pubd, votes) values ($1, $2, $3, $4, $5, $6);", t.ID, t.OriginalPoster, t.Title, t.Content, t.Date, t.Votes)
	if err != nil {
		return err
	}

	return nil
}

func saveReply(t ThreadReplies) (err error) {
	if len(t.Content) <= 1 {
		_, err = db.Exec("INSERT INTO replies (thread, poster, content, votes) values ($1, $2, $3, $4);", t.ThreadID, t.Poster, t.Content, t.Votes)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = db.Exec("UPDATE replies set poster=$2, content=$3, votes=$4 where thread = $1", t.ThreadID, t.Poster, t.Content, t.Votes)
	if err != nil {
		return err
	}

	return nil
}

func initDb(username, password string) {
	var err error

	s := fmt.Sprintf("user=%v password=%v dbname=simpleforum sslmode=disable",
		username, password)

	db, err = sql.Open("postgres", s)
	if err != nil {
		log.Fatal("failed to access database ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database ", err)
	}

	if err = importThreads(); err != nil {
		log.Fatal("Failed to import threads from database ", err)
	}

	if err = importReplies(); err != nil {
		log.Fatal("Failed to import replies from database ", err)
	}
}
