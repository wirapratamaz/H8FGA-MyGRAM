package controller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/wirapratamaz/H8FGA-MyGRAM/app/auth"
	"github.com/wirapratamaz/H8FGA-MyGRAM/app/response"
	"github.com/wirapratamaz/H8FGA-MyGRAM/models"
	"github.com/wirapratamaz/H8FGA-MyGRAM/repository"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object to create"
// @Success 201 {object} UserCreateResponse
// @Router /users [post]
func (controller *UserController) CreateUser(ctx *gin.Context) {
	user := models.User{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	_, err = govalidator.ValidateStruct(&user)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&user).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusCreated, repository.UserCreateResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	})
}

// @Summary User login
// @Description Authenticate the user and generate a token
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object to authenticate"
// @Success 200 {object} UserLoginResponse
// @Router /users/login [post]
func (controller *UserController) UserLogin(ctx *gin.Context) {
	user := models.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	password := user.Password
	err = controller.db.Debug().Where("email = ?", user.Email).Take(&user).Error

	if err != nil {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "username / password is not match",
		})
		return
	}

	comparePass := auth.ComparePassword(user.Password, password)
	if !comparePass {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "username / password is not match",
		})
		return
	}
	token := auth.GenerateToken(user.Id, user.Email)
	ctx.JSON(http.StatusOK, repository.UserLoginResponse{
		Token: token,
	})
}

// UpdateUser updates an existing user's information
// @Summary Update an existing user
// @Description Update an existing user's information
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <access_token>"
// @Param user body UserUpdateRequest true "User object that needs to be updated"
// @Success 200 {object} UserUpdateResponse
// @Router /users/{id} [put]
func (controller *UserController) UpdateUser(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	userReq := repository.UserUpdateRequest{}
	user := models.User{}

	err := ctx.ShouldBindJSON(&userReq)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	updatedUser := models.User{
		Email:    userReq.Email,
		Username: userReq.Username,
	}

	_, err = govalidator.ValidateStruct(&userReq)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.First(&user, userId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "User data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	err = controller.db.Model(&user).Updates(updatedUser).Error
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, repository.UserUpdateResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		UpdatedAt: user.UpdatedAt,
	})
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Deletes the user account associated with the provided authentication token
// @Tags users
// @Produce json
// @Success 200 {object} gin.H
// @Router /users [delete]
func (controller *UserController) DeleteUser(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var user models.User

	err := controller.db.First(&user, userId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "User not found")
			return
		}
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Delete(&user).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
