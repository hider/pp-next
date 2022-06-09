package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"okki.hu/garric/ppnext/config"
	"okki.hu/garric/ppnext/consts"
	"okki.hu/garric/ppnext/model"
)

func DisplayRoom(c *gin.Context) {
	user := c.MustGet("user")
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	h := gin.H{
		"room":    room,
		"user":    user,
		"options": model.VoteOptions,
		"support": consts.Support,
	}
	c.HTML(http.StatusOK, "room.html", h)
}

func UserList(c *gin.Context) {
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h := gin.H{
		"room":    room,
		"options": model.VoteOptions,
		"lookup":  model.VoteLookup,
	}
	c.HTML(http.StatusOK, "user-list.html", h)
}

func Results(c *gin.Context) {
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	h := gin.H{
		"room":    room,
		"options": model.VoteOptions,
		"lookup":  model.VoteLookup,
	}
	c.HTML(http.StatusOK, "results.html", h)
}

func AcceptVote(c *gin.Context) {
	var vote float64
	err := c.ShouldBind(&vote)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := c.MustGet("user").(string)
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	room.RegisterVote(model.NewVote(user, vote))
	err = config.Repository.Save(room)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func Reveal(c *gin.Context) {
	user := c.MustGet("user").(string)
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	room.Revealed = true
	room.RevealedBy = user
	err = config.Repository.Save(room)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func ResetRoom(c *gin.Context) {
	user := c.MustGet("user").(string)
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	room.Reset(user)
	err = config.Repository.Save(room)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func GetEvents(c *gin.Context) {
	name := c.Param("room")
	room, err := config.Repository.Load(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	e := &model.RoomEvent{
		Revealed: room.Revealed,
		ResetTs:  room.ResetTs.UnixMilli(),
	}
	c.JSON(http.StatusOK, e)
}
