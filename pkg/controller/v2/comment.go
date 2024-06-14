package v2

import (
	"encoding/json"
	"net/http"
	. "rentServer/core/controller"
	core "rentServer/core/service"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/album"
	"rentServer/pkg/model/comment"
	"rentServer/pkg/service"
	"strconv"
)

type CommentController interface {
	Post(c *Context)
	List(ctx *Context)
}

type commentController struct {
	service.CommentService
	service.AlbumService
	service.UserService
}

func NewCommentController() CommentController {
	return &commentController{
		CommentService: service.NewCommentService(),
		AlbumService:   service.NewAlbumService(),
		UserService:    service.NewUserService(),
	}
}

func (_this commentController) Post(c *Context) {
	var commentModel comment.Comment

	err := c.ShouldBindJSON(&commentModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	commentModel.UserID = strconv.FormatUint(c.GetUint64(Uid), 10)

	picStr, err := json.Marshal(commentModel.P)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	commentModel.PicList = string(picStr)
	err = _this.CommentService.Save(&commentModel, 0)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 存入店铺相册
	for _, url := range commentModel.P {
		albumModel := album.Album{
			StoreID: commentModel.StoreID,
			From:    album.Customer,
			Url:     url,
		}
		err = _this.AlbumService.Save(&albumModel, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK(commentModel)
}

func (_this commentController) List(c *Context) {

	var param core.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	filter := map[string]any{
		"store_id": c.Query("store_id"),
	}
	clerkID := c.Query("clerk_id")
	if len(clerkID) != 0 {
		filter["clerk_id"] = clerkID
	}
	productID := c.Query("product_id")
	if len(productID) != 0 {
		filter["product_id"] = productID
	}
	list, err := _this.CommentService.FilterList(&param, filter)
	if err != nil {
		return
	}

	for _, commentItem := range list {
		err = json.Unmarshal([]byte(commentItem.PicList), &commentItem.P)
		if err != nil {
			logger.Info(err.Error())
			//c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = _this.UserService.FindByID(&commentItem.User, commentItem.UserID)
		if err != nil {
			logger.Info(err.Error())
			return
		}

	}
	c.JSONOK(list)
}
