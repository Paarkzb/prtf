package handler

import (
	"net/http"
	"prtf"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) saveQuiz(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input prtf.SaveQuizInput
	if err := c.BindJSON(&input); err != nil {
		logrus.Println("HERE")
		newErrorResponse(c, http.StatusBadRequest, "binding quiz "+err.Error())
		return
	}

	quizId, err := h.service.Quiz.Save(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "saving quiz"+err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": quizId,
	})
}

type getAllQuizResponse struct {
	Data []prtf.QuizResponse `json:"data"`
}

func (h *Handler) getAllQuiz(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		return
	}

	quizes, err := h.service.Quiz.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllQuizResponse{
		Data: quizes,
	})
}

func (h *Handler) getQuizById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	quiz, err := h.service.Quiz.GetById(userId, uid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *Handler) updateQuiz(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	var input prtf.UpdateQuizInput
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Quiz.Update(userId, uid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteQuiz(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	err = h.service.Quiz.DeleteById(userId, uid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
