package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUsecase,
	}
}

type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {

	var input RegisterUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := h.UserUseCase.Register(input.Username, input.Password, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, user)

}

// Handler para obtener un usuario por ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserUseCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
