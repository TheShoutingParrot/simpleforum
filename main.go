package main

import (
	"log"
	"time"
	"database/sql"
	"os"

	"github.com/lib/pq"
)

type User struct {
	ID			int64		`json:"id", db:"id"`
	Username		string		`json:"username", db:"username"`
	Password		string		`json:"password", db:"password"`
	Role			int		`json:"role", db:"role"`
}

type Session struct {
	SessionUser		User
	Expire			int
	ExpireTime		time.Time
}

type ForumThread struct {
	ID			int64		`json:"id", db:"id"`
	OriginalPoster		int64		`json:"op", db:"op"`
	Title			string		`json:"title", db:"title"`
	Content			string		`json:"content", db:"content"`
	Votes			int		`json:"votes", db:"votes"`
	Date			time.Time	`json:"pubd", db:"pubd"`
}

type ThreadReplies struct {
	ThreadID		int64		`json:"thread", db:"thread"`
	Poster			pq.Int64Array	`json:"user", json:"user"`
	Content			pq.StringArray	`json:"content", db:"content"`
	Votes			pq.Int64Array	`json:"votes", db:"votes"`
}

var threads map[int64]ForumThread
var replies map[int64]ThreadReplies
var sessions map[string]Session
var previousID int64

var db *sql.DB

func main() {
	log.Println("Starting the simpleforum server")

	threads = make(map[int64]ForumThread)
	replies = make(map[int64]ThreadReplies)
	sessions = make(map[string]Session)

	previousID = 0

	initDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"))

	startServer()
}
