package v2

import (
	"encoding/json"
	"net/http"
	. "rentServer/core/controller"
	"rentServer/core/service"
	good2 "rentServer/pkg/model/good"
	biz "rentServer/pkg/service"
	"rentServer/pkg/service/v1/good"
)

type GoodController interface {
	List(ctx *Context)
	Detail(*Context)
	GetCategoryList(c *Context)
}

type goodController struct {
	GoodService good.GoodService
	biz.CommentService
}

func NewGoodController() GoodController {
	return &goodController{
		GoodService:    good.NewGoodService(),
		CommentService: biz.NewCommentService(),
	}
}

func (_this goodController) List(c *Context) {

	var param service.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}
	categoryID := c.Query("category_id")

	filter := map[string]any{}

	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	list, err := _this.GoodService.FilterList(&param, filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	for i := 0; i < len(list); i++ {
		//goodItem := list[i]
		var banners []string
		err = json.Unmarshal([]byte(list[i].Banners), &banners)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		list[i].B = banners
	}

	c.JSONOK(list)
}

func (_this goodController) Detail(c *Context) {

	var goodModel good2.Product

	err := _this.GoodService.Detail(&goodModel, c.Query("good_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal([]byte(goodModel.Banners), &goodModel.B)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSONOK(goodModel)
}

func (_this goodController) GetCategoryList(c *Context) {
	list, err := _this.GoodService.GetCategoryList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(list)
}
