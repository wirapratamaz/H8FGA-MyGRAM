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

type SocialController struct {
	db *gorm.DB
}

func NewSocialController(db *gorm.DB) *SocialController {
	return &SocialController{
		db: db,
	}
}

// CreateSocial godoc
// @Summary Create a social media account for a user
// @Description Create a new social media account for the authenticated user. Returns the created social media account.
// @Tags Social
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} repository.SocialCreateResponse
// @Router /social [post]
func (controller *SocialController) CreateSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialRequest := repository.SocialRequest{}

	err := ctx.ShouldBindJSON(&socialRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	social := models.Social{
		Name:           socialRequest.Name,
		SocialMediaUrl: socialRequest.SocialMediaUrl,
		UserId:         uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&social)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&social).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusCreated, repository.SocialCreateResponse{
		Id:             social.Id,
		Name:           social.Name,
		SocialMediaUrl: social.SocialMediaUrl,
		UserId:         social.UserId,
		CreatedAt:      social.CreatedAt,
	})
}

// FindAllSocial godoc
// @Summary Find all social media of the authenticated user
// @Description Get all social media data of the authenticated user
// @Tags Social
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} repository.SocialCreateResponse
// @Router /social [get]
func (controller *SocialController) FindAllSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var socials []models.Social

	err := controller.db.Where("user_id = ?", userId).Find(&socials).Error
	if err != nil {
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var socialResponses []repository.SocialCreateResponse
	for _, social := range socials {
		socialResponses = append(socialResponses, repository.SocialCreateResponse{
			Id:             social.Id,
			Name:           social.Name,
			SocialMediaUrl: social.SocialMediaUrl,
			UserId:         social.UserId,
			CreatedAt:      social.CreatedAt,
		})
	}

	response.WriteJsonResponse(ctx, http.StatusOK, socialResponses)
}

// UpdateSocial godoc
// @Summary Update social media data of the authenticated user
// @Description Update social media data of the authenticated user
// @Tags Social
// @Accept json
// @Produce json
// @Param socialMediaId path string true "Social Media ID"
// @Param socialMedia body repository.SocialRequest true "Social Media Data"
// @Security ApiKeyAuth
// @Success 200 {object} repository.SocialCreateResponse
// @Router /social/{socialMediaId} [put]
func (controller *SocialController) UpdateSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialMediaId := ctx.Param("socialMediaId")
	var social models.Social

	err := controller.db.First(&social, socialMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to update this social media",
		})
		return
	}

	socialRequest := repository.SocialRequest{}
	err = ctx.ShouldBindJSON(&socialRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	social.Name = socialRequest.Name
	social.SocialMediaUrl = socialRequest.SocialMediaUrl

	_, err = govalidator.ValidateStruct(&social)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Save(&social).Error
	if err != nil {
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, repository.SocialCreateResponse{
		Id:             social.Id,
		Name:           social.Name,
		SocialMediaUrl: social.SocialMediaUrl,
		UserId:         social.UserId,
		CreatedAt:      social.CreatedAt,
		UpdatedAt:      social.UpdatedAt,
	})
}

// DeleteSocial godoc
// @Summary Delete social media data of the authenticated user
// @Description Delete social media data of the authenticated user
// @Tags Social
// @Accept json
// @Produce json
// @Param socialMediaId path string true "Social Media ID"
// @Security ApiKeyAuth
// @Success 200 {object} gin.H
// @Router /social/{socialMediaId} [delete]
func (controller *SocialController) DeleteSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialMediaId := ctx.Param("socialMediaId")
	var social models.Social

	err := controller.db.First(&social, socialMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this social media",
		})
		return
	}

	err = controller.db.Delete(&social).Error
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
		"message": "Your social media has been successfully deleted",
	})
}
