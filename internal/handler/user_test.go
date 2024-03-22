package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/service/mocks"
)

func TestUserRegister(t *testing.T) {
	t.Run("valid user age to register", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer([]byte(`{"username":"test", "email":"test@test.com", "password":"testtt", "dob":"2020-10-04"}`)))

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(rec)
		g.Request = req

		serviceMock := mocks.NewUserService(t)
		serviceMock.
			On("CheckIsAValidAge", "2020-10-04").
			Return(false, nil)

		userHandler := userHandlerImpl{svc: serviceMock}
		userHandler.UserRegister(g)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("error not all required data is provided", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer([]byte(`{"username":"", "email":"test@test.com", "password":"testtt", "dob":"2020-10-04"}`)))

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(rec)
		g.Request = req

		userHandler := userHandlerImpl{}
		userHandler.UserRegister(g)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("register with not valid date format", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer([]byte(`{"username":"test", "email":"test@test.com", "password":"test", "dob":"2010-10-044"}`)))

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(rec)
		g.Request = req

		serviceMock := mocks.NewUserService(t)
		serviceMock.
			On("CheckIsAValidAge", "2010-10-044").
			Return(false, errors.New(""))

		userHandler := userHandlerImpl{svc: serviceMock}
		userHandler.UserRegister(g)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("successfully register", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer([]byte(`{"username":"test", "email":"test@test.com", "password":"test", "dob":"2010-10-04"}`)))

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(rec)
		g.Request = req

		userRegData := model.UserSignUp{Username: "test", Email: "test@test.com", Password: "test", DOB: "2010-10-04"}

		serviceMock := mocks.NewUserService(t)
		serviceMock.
			On("CheckIsAValidAge", userRegData.DOB).
			Return(true, nil)

		serviceMock.
			On("UserRegister", g, userRegData).
			Return(&model.UserView{}, nil)

		userHandler := userHandlerImpl{svc: serviceMock}
		userHandler.UserRegister(g)

		assert.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	})
}
