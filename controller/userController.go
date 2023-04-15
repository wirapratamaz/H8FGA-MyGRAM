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
