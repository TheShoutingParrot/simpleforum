package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func newThread(c *gin.Context) {
	var n ForumThread

	if c.BindJSON(&n) != nil {
		c.String(http.StatusBadRequest, "Failed to bind json")
		return
	}

	previousID++
	threads[previousID] = ForumThread{
		ID:		previousID,
		OriginalPoster:	n.OriginalPoster,
		Title:		n.Title,
		Content:	n.Content,
		Date:		time.Now(),
		Votes:		0,
	}

	if err := saveThread(threads[previousID]); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "couldn't save")
		return
	}

	c.IndentedJSON(http.StatusCreated, threads[previousID])
}

func readThreads(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, threads)
}

func readThread(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "id not int!")
		return
	}

	c.IndentedJSON(http.StatusOK, threads[id])
}

func reply(c *gin.Context) {
	var reply struct {
		Poster		int64	`json: "poster"`
		Content		string	`json: "content"`
	}

	thread, err := strconv.ParseInt(c.Param("thread"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "id not int!")
		return
	}

	if err := c.BindJSON(&reply); err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "Failed to bind json")
		return
	}

	// Check that the thread does exist
	if _, ok := threads[thread]; !ok {
		c.String(http.StatusNotFound, "Thread doesn't exist (yet)")
		return
	}

	replies[thread] = ThreadReplies{
		ThreadID:	thread,
		Poster:		append(replies[thread].Poster, reply.Poster),
		Content:	append(replies[thread].Content, reply.Content),
		Votes:		append(replies[thread].Votes, 0),
	}

	if err := saveReply(replies[thread]); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "couldn't save")
		return
	}

	c.IndentedJSON(http.StatusCreated, replies[thread])
}

func readReplies(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("thread"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "id not int!")
		return
	}

	r, ok := replies[id]
	if !ok {
		c.String(http.StatusNotFound, "Doesn't exist!")
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

// TODO: Save users that have voted (so no one can vote twice) 
func voteThread(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("thread"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "id not int!")
		return
	}

	_, ok := threads[id]
	if !ok {
		c.String(http.StatusNotFound, "Doesn't exist!")
		return
	}

	var val int

	val = 0

	if c.Param("vote") == "down" {
		val = -1;
	} else if c.Param("vote") == "up" {
		val = 1;
	} else {
		c.String(http.StatusBadRequest, "unrecognized vote")
		return
	}

	threads[id] = ForumThread{
		ID:		threads[id].ID,
		OriginalPoster:	threads[id].OriginalPoster,
		Title:		threads[id].Title,
		Content:	threads[id].Content,
		Votes:		threads[id].Votes + val,
		Date:		threads[id].Date,
	}

	c.String(http.StatusAccepted, "voted!")
}

func startServer() {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	router.Use(cors.New(corsConfig))

	// This request can be used to test if it's working
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/api/new", newThread)
	router.GET("/api/read/:id", readThread)
	router.GET("/api/read", readThreads)
	//router.POST("/api/vote/:thread/:vote", voteThread)

	router.POST("/api/reply/:thread", reply)
	router.GET("/api/replies/:thread", readReplies)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
