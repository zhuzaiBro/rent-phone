package v1

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	. "rentServer/core/controller"
	"rentServer/core/model"
	"rentServer/core/service"
	cosClient "rentServer/initilization/resource"
	"rentServer/pkg/config"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/resource"
	"rentServer/pkg/model/resource_group"
	v1 "rentServer/pkg/service/v1"
	"strconv"
	"strings"
	"time"
)

var _ UploadController = (*uploadController)(nil)

type UploadController interface {
	Upload(c *Context)
	SaveGroup(c *Context)
	DeleteGroup(c *Context)
	GetGroup(c *Context)
	GetUploadResource(c *Context)
}

type uploadController struct {
	ResourceService      v1.ResourceService
	ResourceGroupService v1.ResourceGroupService
}

func NewUploadController() UploadController {
	return &uploadController{
		ResourceGroupService: v1.NewResourceGroupService(),
		ResourceService:      v1.NewResourceService(),
	}
}

func (_this uploadController) Upload(c *Context) {
	_file, err := c.FormFile(Upload)

	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrBadRequest)
		return
	}
	n := time.Now()

	ns := strconv.Itoa(int(n.UnixMicro()))

	r := rand.Int()
	rs := strconv.Itoa(r)

	fileList := strings.Split(_file.Filename, ".")

	key := "/upload/" + ns + rs[0:6] + "." + fileList[1]
	if err != nil {
		print(err.Error(), 123)
		return
	}

	file, err3 := _file.Open()
	if err3 != nil {
		print(err3.Error())
		return
	}
	logger.Info("开始上传图片资源")
	_, err2 := cosClient.CosClient.Object.Put(context.Background(), key, file, &cos.ObjectPutOptions{})

	if err2 != nil {
		logger.Error(err2.Error())
		c.JSONOK(ErrInternalServer)
		return
	}
	c.JSONOK(config.GetConfig().CosConfig.CosUrl + key)
	// 上传成功之后, 开一个子进程存储
	go func() {
		resourceModel := resource.Resource{
			BaseModel: model.BaseModel{},
			Name:      key,
			Url:       config.GetConfig().CosConfig.CosUrl + key,
			GroupID:   c.Query("group_id"),
			Type:      fileList[1],
		}
		err := _this.ResourceService.Save(&resourceModel, 0)
		if err != nil {
			return
		}
	}()

}

func (_this uploadController) SaveGroup(c *Context) {
	var groupModel resource_group.ResourceGroup

	err := c.ShouldBindJSON(&groupModel)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	err = _this.ResourceGroupService.Save(&groupModel, groupModel.ID)
	if err != nil {

		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(&groupModel)
}

func (_this uploadController) DeleteGroup(c *Context) {

	var ids []uint64

	err := _this.ResourceGroupService.Delete(ids)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	c.JSONOK()
}

func (_this uploadController) GetGroup(c *Context) {

	var param service.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrBadRequest)
		return
	}

	logger.Info(param, "传入的param")
	list, err := _this.ResourceGroupService.List(&param)
	if err != nil {

		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(list)
}

func (_this uploadController) GetUploadResource(c *Context) {

	var param service.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSONOK(ErrBadRequest)
		return
	}

	groupId := c.Query("group_id")
	filter := map[string]any{
		"group_id": groupId,
	}
	list, err := _this.ResourceService.FilterList(&param, filter)
	if err != nil {
		logger.Error(err.Error())
		c.JSONOK(ErrInternalServer)
		return
	}

	c.JSONOK(list)
}
