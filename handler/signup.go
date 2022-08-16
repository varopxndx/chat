package handler

import (
	"log"
	"net/http"

	"github.com/varopxndx/chat/model"

	"github.com/gin-gonic/gin"
)

// SignUp handler
func (h Handler) SignUp(g *gin.Context) {
	var data model.User
	if err := g.ShouldBind(&data); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.SignUp(g.Request.Context(), data); err != nil {
		log.Default().Printf("unable to signup: %v", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.HTML(http.StatusOK, "index.html", nil)
}
