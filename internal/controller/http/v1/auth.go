package v1

import (
	"errors"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TODO: Add logger

func newAuthRoutes(handler *gin.RouterGroup) {
	usecase := usecase.NewUserUseCase()

	h := handler.Group("/auth")
	{
		h.POST("/register", func(ctx *gin.Context) {
			RegisterHandler(ctx, usecase)
		})
		h.POST("/login", func(ctx *gin.Context) {
			LoginHandler(ctx, usecase)
		})
	}
}

func RegisterHandler(ctx *gin.Context, usecase usecase.UserUseCase) {
	var body model.UserRegister

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body format"})
		return
	}

	user, err := usecase.CreateUser(body.FirstName, body.LastName, body.Username, body.Email, body.Password)

	if err == nil {
		if token, err := utils.GenerateJWT(user.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "User created", "token": token})
		}

		return
	}

	if validationErr, ok := err.(utils.ParamError); ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{validationErr.Param: validationErr.Message}})
		return
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "An user with that username or email already exists."})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
}

func LoginHandler(ctx *gin.Context, usecase usecase.UserUseCase) {
	var body model.UserLogin

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body format"})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(errs[0])
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	user, err := usecase.GetByUsernameOrEmail(body.Username, body.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user found with that username or email"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password mismatch"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if token, err := utils.GenerateJWT(user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "User created", "token": token})
	}
}