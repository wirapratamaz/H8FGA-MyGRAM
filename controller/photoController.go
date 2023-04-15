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

// CreatePhoto godoc
// @Summary Create a new photo data for authenticated user
// @Description Create a new photo data for authenticated user
// @Tags Photos
// @Accept json
// @Produce json
// @Param photo body repository.PhotoRequest true "Photo data"
// @Security ApiKeyAuth
// @Success 201 {object} repository.PhotoCreateResponse
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 500 {object} response.InternalServerError
// @Router /photos [post]

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

// FindAllPhoto godoc
// @Summary Find all photos of the authenticated user
// @Description Get all photos of the authenticated user
// @Tags Photo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Photo
// @Failure 401 {object} response.Unauthorized
// @Failure 500 {object} response.InternalServerError
// @Router /photo [get]
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

// UpdatePhoto godoc
// @Summary Update photo data of the authenticated user
// @Description Update photo data of the authenticated user
// @Tags Photo
// @Accept json
// @Produce json
// @Param photoId path string true "Photo ID"
// @Param photoRequest body repository.PhotoRequest true "Photo Data"
// @Security ApiKeyAuth
// @Success 200 {object} repository.PhotoCreateResponse
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.NotFound
// @Failure 500 {object} response.InternalServerError
// @Router /photo/{photoId} [put]
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

// DeletePhoto godoc
// @Summary Delete photo data of the authenticated user
// @Description Delete photo data of the authenticated user
// @Tags Photos
// @Accept json
// @Produce json
// @Param photoId path string true "Photo ID"
// @Security ApiKeyAuth
// @Success 200 {object} gin.H
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.NotFound
// @Failure 500 {object} response.InternalServerError
// @Router /photos/{photoId} [delete]
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
