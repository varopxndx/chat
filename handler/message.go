package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// limit of messages to be displayed
const maxLimit = 50

// Message handler
func (h Handler) Message(g *gin.Context) {
	limitString := g.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Default().Printf("error converting limit, wrong type %v", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if limit < 1 || limit > maxLimit {
		err := fmt.Sprintf("invalid limit: %d , out of range", limit)
		log.Default().Println(err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// get room from query param
	room, _ := g.GetQuery("room")

	msgs, err := h.usecase.GetMessages(g.Request.Context(), limit, room)
	if err != nil {
		log.Default().Printf("error: %v", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get messages"})
		return
	}
	g.JSON(http.StatusOK, msgs)
}
