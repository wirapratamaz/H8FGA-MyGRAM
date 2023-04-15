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

type CommentController struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{
		db: db,
	}
}

func (controller *CommentController) CreateComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentRequest := repository.CommentRequest{}

	err := ctx.ShouldBindJSON(&commentRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	comment := models.Comment{
		Message: commentRequest.Message,
		PhotoId: commentRequest.PhotoId,
		UserId:  uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&comment)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&comment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusCreated, repository.CommentCreateResponse{
		Id:        comment.Id,
		Message:   comment.Message,
		PhotoId:   comment.PhotoId,
		UserId:    comment.UserId,
		CreatedAt: comment.CreatedAt,
	})
}

func (controller *CommentController) FindAllComment(ctx *gin.Context) {
	var comments []models.Comment
	err := controller.db.Find(&comments).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	commentList := make([]repository.CommentCreateResponse, 0)
	for _, comment := range comments {
		commentList = append(commentList, repository.CommentCreateResponse{
			Id:        comment.Id,
			Message:   comment.Message,
			PhotoId:   comment.PhotoId,
			UserId:    comment.UserId,
			CreatedAt: comment.CreatedAt,
		})
	}

	response.WriteJsonResponse(ctx, http.StatusOK, commentList)
}

func (controller *CommentController) UpdateComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var comment models.Comment

	err := controller.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to update this comment",
		})
		return
	}

	commentRequest := repository.CommentRequest{}

	err = ctx.ShouldBindJSON(&commentRequest)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	comment.Message = commentRequest.Message

	_, err = govalidator.ValidateStruct(&comment)
	if err != nil {
		response.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Save(&comment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, err.Error())
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	response.WriteJsonResponse(ctx, http.StatusOK, repository.CommentCreateResponse{
		Id:        comment.Id,
		Message:   comment.Message,
		PhotoId:   comment.PhotoId,
		UserId:    comment.UserId,
		CreatedAt: comment.CreatedAt,
	})
}

func (controller *CommentController) DeleteComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var comment models.Comment

	err := controller.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			response.NotFoundResponse(ctx, "data not found")
			return
		}
		response.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		response.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this comment",
		})
		return
	}

	err = controller.db.Delete(&comment).Error
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
		"message": "Your comment has been successfully deleted",
	})
}
