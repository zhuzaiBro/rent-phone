package service

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/attr_value"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/good"
	"strconv"
)

var _ CartService = (*cartService)(nil)

type CartService interface {
	Service[cart.Cart]
	AddGoods2Cart(userID uint64, addItems cart.AddRequest) (list []*cart.Cart, err error)
}

type cartService struct {
	ServiceIMPL[cart.Cart]
	AttrValueDao dao.AttrValueDao
	GoodDao      dao.GoodDao
}

func NewCartService() CartService {
	cartDao := dao.NewCartDao()
	return &cartService{
		ServiceIMPL: ServiceIMPL[cart.Cart]{
			BaseDao: cartDao,
		},
		AttrValueDao: dao.NewAttrValueDao(),
		GoodDao:      dao.NewGoodDao(),
	}
}

// AddGoods2Cart Add 加入商品购物车
func (_this cartService) AddGoods2Cart(userID uint64, addItems cart.AddRequest) (list []*cart.Cart, err error) {
	for _, item := range addItems {
		filter := map[string]any{
			"unique": item.UniqueID,
		}
		var (
			goodModel      good.Product
			attrValueModel attr_value.AttrValue
		)

		err := _this.AttrValueDao.FilterFindOne(&attrValueModel, filter)
		if err != nil {
			return nil, err
		}

		err = _this.GoodDao.FindByID(&goodModel, attrValueModel.ProductID)
		if err != nil {
			return nil, err
		}

		price, err := decimal.NewFromString(attrValueModel.Price)
		if err != nil {
			return nil, err
		}
		product := cart.CartProduct{
			Cover: attrValueModel.Image,
			Title: goodModel.Title + "-" + attrValueModel.Sku,
			Price: price,
			Num:   item.Num,
			// todo 多商户预留字段
			Sku:         attrValueModel.Sku,
			MerID:       "",
			ProductType: GOOD,
			ProductID:   attrValueModel.ProductID,
		}
		productStr, err := json.Marshal(&product)
		cartModel := cart.Cart{
			Product:   string(productStr),
			Total:     price.Mul(decimal.New(int64(item.Num), 0)),
			P:         product,
			Type:      GOOD,
			ProductID: attrValueModel.ProductID,
			UserID:    strconv.FormatUint(userID, 10),
		}

		err = _this.Save(&cartModel, 0)
		if err != nil {
			return nil, err
		}
		list = append(list, &cartModel)
	}
	return

}
