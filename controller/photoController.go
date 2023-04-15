package controller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/wirapratamaz/H8FGA-MyGRAM/app/response"
	"github.com/wirapratamaz/H8FGA-MyGRAM/models"
	"github.com/wirapratamaz/H8FGA-MyGRAM/repository"
	"gorm.io/gorm"
)

type PhotoController struct {
	db *gorm.DB
}

func NewPhotoController(db *gorm.DB) *PhotoController {
	return &PhotoController{
		db: db,
	}
}

func (controller *PhotoController) CreatePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoRequest := repository.PhotoRequest{}

	err := ctx.ShouldBindJSON(&photoRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	photo := models.Photo{
		Title:    photoRequest.Title,
		Caption:  photoRequest.Caption,
		PhotoUrl: photoRequest.PhotoUrl,
		UserId:   uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&photo)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&photo).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusCreated, repository.PhotoCreateResponse{
		Id:        photo.Id,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserId:    photo.UserId,
		CreatedAt: photo.CreatedAt,
	})
}

func (controller *PhotoController) FindAllPhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var photos []models.Photo

	err := controller.db.Where("user_id = ?", userId).Find(&photos).Error
	if err != nil {
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, photos)
}

func (controller *PhotoController) UpdatePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoId := ctx.Param("photoId")

	var photo models.Photo
	err := controller.db.First(&photo, photoId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to update this photo",
		})
		return
	}

	photoRequest := repository.PhotoRequest{}
	err = ctx.ShouldBindJSON(&photoRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	photo.Title = photoRequest.Title
	photo.Caption = photoRequest.Caption
	photo.PhotoUrl = photoRequest.PhotoUrl

	_, err = govalidator.ValidateStruct(&photo)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Save(&photo).Error
	if err != nil {
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, repository.PhotoCreateResponse{
		Id:        photo.Id,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserId:    photo.UserId,
		CreatedAt: photo.CreatedAt,
	})
}

func (controller *PhotoController) DeletePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoId := ctx.Param("photoId")
	var photo models.Photo

	err := controller.db.First(&photo, photoId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this photo",
		})
		return
	}

	err = controller.db.Delete(&photo).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"error":   false,
		"message": "Your photo has been successfully deleted",
	})
}
