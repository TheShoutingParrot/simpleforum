package main

import (
	"net/http"
	"log"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

func checkCookie(c *gin.Context) (user User, err error) {
	cookie, err := c.Cookie("session_token")

	if err != nil {
		return user, err
	}

	s, exists := sessions[cookie]
	if !exists {
		return user, sql.ErrNoRows
	}

	return s.SessionUser, nil
}

func signup(c *gin.Context) {
	var n User

	if err := c.BindJSON(&n); err != nil {
		c.String(http.StatusBadRequest, "failed to bind json")

		return
	}

	var temp int
	res := db.QueryRow("select role from users where username=$1", n.Username)

	err := res.Scan(&temp)
	if err != nil && err != sql.ErrNoRows {
		log.Println("error occured: ", err)
		c.String(http.StatusInternalServerError, "Failed!")
		return
	} else if err == nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "username in use")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), 12)
	if err != nil {
		log.Println("error occured: ", err)
		c.IndentedJSON(http.StatusInternalServerError, "failed to generate password")

		return
	}

	_, err = db.Exec("insert into users (username, password, role) values ($1, $2, $3)", n.Username, string(hash), 0)
	if err != nil {
		log.Println("error occured: ", err)
		c.String(http.StatusInternalServerError, "couldn't save")
		return
	}

	log.Println("New user created %v", n.Username)
	c.String(http.StatusCreated, "User created!")
}

func signin(c *gin.Context) {
	var user, stored User

	if err := c.BindJSON(&user); err != nil {
		log.Println("signin: ", err)
		c.String(http.StatusBadRequest, "bad request")
	}

	res := db.QueryRow("select id, password, role from users where username=$1", user.Username)

	err := res.Scan(&stored.ID, &stored.Password, &stored.Role)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(user.Password))
	if err != nil {
		// We actually know that the password is incorrect
		c.String(http.StatusUnauthorized, "Incorrect username or password")
		return
	}

	secs := 7200 // 2 hours

	token := uuid.NewString()
	expiresAt := time.Now().Add(time.Duration(secs) * time.Second)

	stored.Password = ""
	stored.Username = user.Username

	sessions[token] = Session {
		SessionUser:	stored,
		Expire:		secs,
		ExpireTime:	expiresAt,
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("session_token", token, sessions[token].Expire, "/", "", true, false)
	c.String(http.StatusAccepted, "logged in!")
	log.Println("added cookie ", token)
}
