package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase usecase.IUserUseCase
}

func NewUserHandler(userUsecase usecase.IUserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUsecase,
	}
}

type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
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

	c.JSON(http.StatusCreated, user)

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

// Handler para actualizar un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserUseCase.UpdateUser(uint(id), input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)

}

// Handler para eliminar un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.UserUseCase.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
}
