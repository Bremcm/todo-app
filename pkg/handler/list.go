package handler

import (
	"net/http"

	"github.com/bremcm/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		newErrorResponce(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

}
func (h *Handler) getAllList(c *gin.Context) {

}
func (h *Handler) getListById(c *gin.Context) {

}
func (h *Handler) updateList(c *gin.Context) {

}
func (h *Handler) deleteList(c *gin.Context) {

}
