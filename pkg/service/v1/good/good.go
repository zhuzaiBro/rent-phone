package good

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/category"
	"rentServer/pkg/model/good"
	"rentServer/pkg/model/order"
	"strconv"
)

var _ GoodService = (*goodServiceIMPL)(nil)

type GoodService interface {
	Service[good.Product]
	PostGoodComment(orderID string, commentDto *good.GoodCommentDto) error
	GetGoodComments(goodID string, param *BaseSearchParam) (list []*good.GoodCommentDto, err error)
	SaveGood(productModel *good.Product) error
	Detail(in *good.Product, goodID string) error
	GetCategoryList() (list []*category.Category, err error)
}

type goodServiceIMPL struct {
	ServiceIMPL[good.Product]
	GoodDao        dao.GoodDao
	GoodCommentDao dao.GoodCommentDao
	OrderDao       dao.OrderDao
	AttrValueDao   dao.AttrValueDao
	AttrDao        dao.AttrDao
	CategoryDao    dao.CategoryDao
}

func NewGoodService() GoodService {
	goodCommentDao := dao.NewGoodCommentDao()
	goodDao := dao.NewGoodDao()
	orderDao := dao.NewOrderDao()
	return &goodServiceIMPL{
		ServiceIMPL: ServiceIMPL[good.Product]{
			BaseDao: goodDao,
		},
		GoodDao:        goodDao,
		GoodCommentDao: goodCommentDao,
		OrderDao:       orderDao,
		AttrValueDao:   dao.NewAttrValueDao(),
		AttrDao:        dao.NewAttrDao(),
		CategoryDao:    dao.NewCategoryDao(),
	}
}

func (_this goodServiceIMPL) SaveGood(productModel *good.Product) error {

	err := _this.Save(productModel, productModel.ID)
	if err != nil {
		return err
	}
	isUpdate := productModel.ID > 0
	if isUpdate {
		// 如果有id那就是修改，删除之前的attr_value
		attrValueIDS, err := _this.AttrValueDao.FindAttrIDSByProductID(productModel.ID)
		if err != nil {
			return err
		}
		for _, attrValueID := range attrValueIDS {
			err = _this.AttrValueDao.Delete(attrValueID)
			if err != nil {
				return err
			}
		}

		attrIds, err := _this.AttrDao.FindAttrIDSByProductID(productModel.ID)
		if err != nil {
			return err
		}

		for _, attrId := range attrIds {
			err = _this.AttrDao.Delete(attrId)
			if err != nil {
				return err
			}
		}

	}

	// 雪花算法生成一个唯一的分布式ID
	var node *snowflake.Node
	node, err = snowflake.NewNode(1)
	if err != nil {
		return err
	}
	// 挨个存储 attr_value
	for _, attrValue := range productModel.Values {
		// 先存储Attr
		attrValue.ProductID = strconv.FormatUint(productModel.ID, 10)
		attrValue.Unique = node.Generate().String()
		err := _this.AttrValueDao.Insert(attrValue)
		if err != nil {
			return err
		}
	}

	for _, attr := range productModel.Attrs {
		attr.ProductID = strconv.FormatUint(productModel.ID, 10)
		attr.D, err = attr.Detail2D()
		err := _this.AttrDao.Insert(attr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (_this goodServiceIMPL) Detail(in *good.Product, goodID string) error {

	err := _this.GoodDao.FindByID(in, goodID)
	if err != nil {
		return err
	}
	filter := map[string]any{
		"product_id": goodID,
	}
	in.Values, err = _this.AttrValueDao.ToList(0, 100, filter, "", "")
	if err != nil {
		return err
	}
	in.Attrs, err = _this.AttrDao.ToList(0, 100, filter, "", "")
	if err != nil {
		return err
	}
	for _, attr := range in.Attrs {
		attr.Detail, err = attr.D2Detail()
		if err != nil {
			return err
		}
	}

	return nil
}

// 获取商品的评论
func (_this goodServiceIMPL) GetGoodComments(goodID string, param *BaseSearchParam) (list []*good.GoodCommentDto, err error) {
	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Size)
	if err != nil {
		return nil, err
	}

	start := (page - 1) * size
	filter := map[string]any{
		"good_id": goodID,
	}

	data, err := _this.GoodCommentDao.ToList(start, size, filter, "", "")
	if err != nil {
		return nil, err
	}

	for _, comment := range data {
		var dto good.GoodCommentDto
		comment.BuildDto(&dto)

		list = append(list, &dto)
	}

	return
}

// PostGoodComment 发表商品评论 默认只有购买了商品的用户才可以评论
func (_this goodServiceIMPL) PostGoodComment(orderID string, commentDto *good.GoodCommentDto) error {

	var orderModel order.Order

	err := _this.OrderDao.FindByID(&orderModel, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("没有相关的购买记录不能评论！")
		}
		return err
	}

	// 查看order中是否含有这件商品
	var (
		isOK  = false
		carts []*cart.CartProduct
	)
	err = json.Unmarshal([]byte(orderModel.Products), &carts)
	if err != nil {
		return err
	}

	for _, product := range carts {
		isOK = product.ProductID == strconv.FormatUint(commentDto.Good.ID, 10)
		if isOK {
			break
		}
	}

	if !isOK {
		// 用户没有购买
		return errors.New("用户没有购买改商品，不能评论！")
	}

	// 插入评论 评论里面包含了goodID
	var commentModel good.GoodComment
	commentModel = commentModel.BuildFromDto(commentDto)

	// 订单状态修改

	// 插入评论
	err = _this.GoodCommentDao.Insert(&commentModel)
	if err != nil {
		return err
	}

	// 商品修改？

	return nil
}

func (_this goodServiceIMPL) GetCategoryList() (list []*category.Category, err error) {
	return _this.getCategory(0, 1000, 0)

}

func (_this goodServiceIMPL) getCategory(start, size int, parentID uint64) (list []*category.Category, err error) {
	filter := map[string]any{}

	if parentID != 0 {
		filter["parent_id"] = parentID
	}

	list, err = _this.CategoryDao.ToList(start, size, filter, "", "")

	// 当没有子节点
	if len(list) == 0 {
		return
	}

	for _, _category := range list {
		// todo 默认先全部查出来
		_category.Children, err = _this.getCategory(0, 1000, _category.ID)
	}

	return
}
