package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CamiloLeonP/parking-radar/internal/test/parking/mockgen"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Configuración de prueba para el handler de usuario
func setupUserHandler(mockUseCase *mockgen.MockIUserUseCase) *gin.Engine {
	userHandler := NewUserHandler(mockUseCase)

	r := gin.Default()
	users := r.Group("/users")
	{
		users.POST("/register", userHandler.Register)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}

// test para el registro de usuario
func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUser := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	mockUseCase.EXPECT().Register("testuser", "password", "test@example.com").Return(mockUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(RegisterUserInput{Username: "testuser", Password: "password", Email: "test@example.com"})
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseUser domain.User
	json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.Equal(t, mockUser, &responseUser)
}

// test para obtener un usuario por ID
func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUser := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	mockUseCase.EXPECT().FindByID(uint(1)).Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseUser domain.User
	json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.Equal(t, mockUser, &responseUser)
}

// test para actualizar un usuario
func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUser := &domain.User{ID: 1, Username: "updateduser", Email: "updated@example.com"}
	mockUseCase.EXPECT().UpdateUser(uint(1), "updateduser", "updated@example.com", "newpassword").Return(mockUser, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(UpdateUserInput{Username: "updateduser", Email: "updated@example.com", Password: "newpassword"})
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseUser domain.User
	json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.Equal(t, mockUser, &responseUser)
}

// test para eliminar un usuario
func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUseCase.EXPECT().DeleteUser(uint(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User deleted", response["status"])
}

// test para error en el registro de usuario
func TestRegisterUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUseCase.EXPECT().Register("testuser", "password", "test@example.com").Return(nil, errors.New("registration failed"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(RegisterUserInput{Username: "testuser", Password: "password", Email: "test@example.com"})
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to register user", response["error"])
}

// test para error en obtener usuario por ID
func TestGetUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUseCase.EXPECT().FindByID(uint(1)).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User not found", response["error"])
}

// test para error en actualizar usuario
func TestUpdateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUseCase.EXPECT().UpdateUser(uint(1), "updateduser", "updated@example.com", "newpassword").Return(nil, errors.New("update failed"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(UpdateUserInput{Username: "updateduser", Email: "updated@example.com", Password: "newpassword"})
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to update user", response["error"])
}

// test para error en eliminar usuario
func TestDeleteUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIUserUseCase(ctrl)
	r := setupUserHandler(mockUseCase)

	mockUseCase.EXPECT().DeleteUser(uint(1)).Return(errors.New("delete failed"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to delete user", response["error"])
}
