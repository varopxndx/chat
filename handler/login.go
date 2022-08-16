package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/varopxndx/chat/model"

	"github.com/gin-gonic/gin"
)

// Login handler
func (h Handler) Login(g *gin.Context) {
	var data model.LoginData
	if err := g.ShouldBind(&data); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usecase.Login(g.Request.Context(), data)
	if err != nil {
		log.Default().Printf("unable to login: %s", err.Error())
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenizer.GenerateToken(*user)
	if err != nil {
		log.Default().Printf("error to generate token: %s", err.Error())
	}

	g.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	g.Header("room", data.Room)
	g.JSON(http.StatusOK, user)
}
