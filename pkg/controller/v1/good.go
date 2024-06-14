package v1

import (
	"encoding/json"
	"net/http"
	. "rentServer/core/controller"
	"rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/category"
	"rentServer/pkg/model/good"
	good_serv "rentServer/pkg/service/v1/good"
)

var _ GoodController = (*goodController)(nil)

type GoodController interface {
	List(c *Context)
	Save(c *Context)
	Delete(c *Context)
	Detail(c *Context)
	GetGoodComments(c *Context)
	DeleteCategory(c *Context)
	SaveCategory(c *Context)
}

type goodController struct {
	GoodService good_serv.GoodService
	CategoryDao dao.CategoryDao
}

func GetGoodController() GoodController {
	return &goodController{
		GoodService: good_serv.NewGoodService(),
		CategoryDao: dao.NewCategoryDao(),
	}
}

// List 列表搜索产品
func (_this goodController) List(c *Context) {

	var param service.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}
	//arr := []string{"good_id=1", "user_id=1"}
	//f := map[string]any{}
	//for _, s := range arr {
	//	str := strings.Split(s, "=")
	//	f[str[0]] = str[1]
	//}

	list, err := _this.GoodService.List(&param)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(list)
}

func (_this goodController) Save(c *Context) {

	var goodModel good.Product

	err := c.ShouldBindJSON(&goodModel)
	if err != nil {
		c.JSONOK(err.Error())
		return
	}

	banner, err := json.Marshal(&goodModel.B)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}
	goodModel.Banners = string(banner)

	err = _this.GoodService.SaveGood(&goodModel)
	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(&goodModel)
}

func (_this goodController) Delete(c *Context) {

	var ids []uint64
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}
	err = _this.GoodService.Delete(ids)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK()
}

// Detail 查询商品详情
func (_this goodController) Detail(c *Context) {
	goodId := c.Query("good_id")

	var goodModel good.Product
	err := _this.GoodService.Detail(&goodModel, goodId)
	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrInternalServer)
		return
	}
	err = json.Unmarshal([]byte(goodModel.Banners), &goodModel.B)
	if err != nil {
		logger.Error(err.Error())

		return
	}

	c.JSONOK(&goodModel)
}

// GetGoodComments 查询商品评论
func (_this goodController) GetGoodComments(c *Context) {
	var param service.BaseSearchParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}
	goodId := c.Query("good_id")

	comments, err := _this.GoodService.GetGoodComments(goodId, &param)
	if err != nil {
		c.JSONOK(ErrInternalServer)
		return
	}
	c.JSONOK(comments)
}

func (_this goodController) DeleteCategory(c *Context) {
	var ids []uint64

	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	for _, id := range ids {
		err = _this.CategoryDao.Delete(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK()
}

func (_this goodController) SaveCategory(c *Context) {
	var categoryModel category.Category

	err := c.ShouldBindJSON(&categoryModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if categoryModel.ID == 0 {
		err = _this.CategoryDao.Insert(&categoryModel)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		err = _this.CategoryDao.Update(&categoryModel, categoryModel.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSONOK(categoryModel)

}
