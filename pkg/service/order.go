package service

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	. "rentServer/core/service"
	"rentServer/initilization/db"
	"rentServer/initilization/wxpay"
	"rentServer/pkg/config"
	"rentServer/pkg/dao"
	"rentServer/pkg/logger"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/clerk"
	"rentServer/pkg/model/order"
	"rentServer/pkg/model/reservation"
	"rentServer/pkg/model/user"
	"rentServer/pkg/model/wx_api"
	"strconv"
	"time"
)

/*** 通过这个字段查询用户名下的订单
*** @param cacheKey := "order:" + strconv.FormatUint(userID, 10)
***
*** 通过这个字段查询用户具体某一个订单
*** @param orderKey := in.TradeNo + strconv.FormatUint(userID, 10)
 */

type OrderService interface {
	Service[order.Order]
	GenerateGoodOrder(userID uint64, cartIds []uint64, in *order.Order) error
	//GenerateGoodOrder(cartIds []string, in *order.Order) error
	PayOrder(wxCallback *wx_api.WxPayCallback) error
	GenerateWxMpPayOrder(TradeNo string, wxOpenID string, userID uint64, couponRelID string) (res *jsapi.PrepayWithRequestPaymentResponse, err error)
	GenerateServiceOrder(in *order.Order, rent *reservation.Reservation) error
	GetMyOrder(param *BaseSearchParam, userID uint64) (list []*order.Order, err error)
	GetMyOrderCount(filter map[string]any) (int64, error)
	HandleDateOrder(TradeNo string, in *order.Order, userID uint64) error
	ComputeOrder(orderModel *order.Order, userID uint64) error
	GenerateDepositOrder(in *order.Order, depositNum string, userID uint64) (err error)
	OrderTask() error
}

type orderService struct {
	ServiceIMPL[order.Order]
	OrderDao      dao.OrderDao
	CartDao       dao.CartDao
	StoreDao      dao.StoreDao
	CouponService CouponService
	AddressDao    dao.AddressDao
	UserDao       dao.UserDao
	ClerkDao      dao.ClerkDao
}

func NewOrderService() OrderService {
	orderDao := dao.NewOrderDao()
	cartDao := dao.NewCartDao()
	return &orderService{
		ServiceIMPL: ServiceIMPL[order.Order]{
			BaseDao: orderDao,
		},
		OrderDao:      orderDao,
		CartDao:       cartDao,
		StoreDao:      dao.NewStoreDao(),
		CouponService: NewCouponService(),
		AddressDao:    dao.NewAddressDao(),
		UserDao:       dao.NewUserDao(),
		ClerkDao:      dao.NewClerkDao(),
	}
}

const (
	SERVICE = "service"
	GOOD    = "good"
	DEPOSIT = "deposit"
)

var Node *snowflake.Node

func init() {

	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}

	Node = node
}

// GenerateGoodOrder 生成产品订单 *需要生成地址
func (_this orderService) GenerateGoodOrder(userID uint64, cartIds []uint64, in *order.Order) error {

	// 处理 ProductData
	for _, cartId := range cartIds {
		var cartModel cart.Cart
		err := _this.CartDao.FindByID(&cartModel, cartId)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(cartModel.Product), &cartModel.P)
		if err != nil {
			return err
		}
		in.ProductList = append(in.ProductList, &cartModel.P)
		in.Total = in.Total.Add(cartModel.Total)

	}
	in.Status = order.PENDING
	in.Type = GOOD

	// 处理商品订单
	pb, err := json.Marshal(in.ProductList)

	if err != nil {
		return err
	}

	in.TradeNo = Node.Generate().String()
	in.UserID = strconv.FormatUint(userID, 10)
	in.Products = string(pb)

	err = _this.StoreDao.FindByID(&in.Store, 1)
	if err != nil {
		return err
	}
	in.StoreID = strconv.FormatUint(in.Store.ID, 10)

	// 提交 cartInfo 后生成订单，使用用户默认地址进行计算
	var userModel user.User
	err = _this.UserDao.FindByID(&userModel, userID)
	if err != nil {
		return err
	}
	// 此时地址字段已填充
	err = _this.getAddress(in, userModel.AddressID)
	if err != nil {
		return err
	}

	cacheClient := db.GetRedis()
	var orderTradeNos []string
	cacheKey := "order:" + strconv.FormatUint(userID, 10)
	// 先给这个用户存一个order的记录
	cacheStr, err := cacheClient.Get(context.Background(), cacheKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// 如果查询到了
	if err == nil && len(cacheStr) > 0 {
		err = json.Unmarshal([]byte(cacheStr), &orderTradeNos)
		if err != nil {
			return err
		}
	}

	orderTradeNos = append(orderTradeNos, in.TradeNo)
	cacheValue, err := json.Marshal(&orderTradeNos)

	err = cacheClient.Set(context.Background(), cacheKey, cacheValue, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	orderKey := in.TradeNo + strconv.FormatUint(userID, 10)
	return _this.OrderDao.CacheSave(in, orderKey, 15*time.Minute)
}

// GenerateServiceOrder 服务订单不需要生成地址
func (_this orderService) GenerateServiceOrder(in *order.Order, rent *reservation.Reservation) error {

	// 处理服务订单, 不需要查询用户之前的购物车，直接NEW 一个购物车
	var cartProduct cart.CartProduct

	cartProduct.ProductType = SERVICE
	cartProduct.Cover = rent.Clerk.Avatar
	var str = "需要淋浴"
	if !rent.NeedShower {
		str = "不" + str
	}
	cartProduct.Sku = str + "," + strconv.Itoa(rent.People) + "人"
	cartProduct.Num = 1
	cartProduct.Title = "线下服务，愈疗师：" + rent.Clerk.Name
	cartProduct.ArriveTime = rent.ArriveTime

	cartProduct.Price = decimal.NewFromFloat(rent.TotalTime).Mul(decimal.NewFromFloat(rent.Clerk.Commission)).Mul(decimal.NewFromInt(int64(rent.People)))

	in.Status = order.PENDING
	in.ProductList = append(in.ProductList, &cartProduct)
	in.Type = SERVICE
	in.ClerkID = strconv.FormatUint(rent.Clerk.ID, 10)
	pb, err := json.Marshal(&in.ProductList)
	if err != nil {
		return err
	}
	in.Products = string(pb)

	in.TradeNo = Node.Generate().String()
	in.Total = cartProduct.Price

	err = _this.StoreDao.FindByID(&in.Store, 1)
	if err != nil {
		return err
	}

	in.StoreID = strconv.FormatUint(in.Store.ID, 10)

	cacheClient := db.GetRedis()
	var orderTradeNos []string
	cacheKey := "order:" + in.UserID
	// 先给这个用户存一个order的记录
	cacheStr, err := cacheClient.Get(context.Background(), cacheKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// 如果查询到了
	if err == nil && len(cacheStr) > 0 {
		err = json.Unmarshal([]byte(cacheStr), &orderTradeNos)
		if err != nil {
			return err
		}
	}

	orderTradeNos = append(orderTradeNos, in.TradeNo)

	cacheValue, err := json.Marshal(&orderTradeNos)
	err = cacheClient.Set(context.Background(), cacheKey, cacheValue, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	orderKey := in.TradeNo + in.UserID
	return _this.OrderDao.CacheSave(in, orderKey, 15*time.Minute)
}

func (_this orderService) GenerateWxMpPayOrder(TradeNo string, wxOpenID string, userID uint64, couponRelID string) (res *jsapi.PrepayWithRequestPaymentResponse, err error) {

	// 查询缓存中的Order字段
	var (
		orderKey   = TradeNo + strconv.FormatUint(userID, 10)
		orderModel order.Order
	)

	err = _this.OrderDao.CacheGetValue(&orderModel, orderKey)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	productStr, err := json.Marshal(&orderModel.ProductList)

	orderModel.Products = string(productStr)
	addressStr, err := json.Marshal(&orderModel.Address)
	orderModel.Addr = string(addressStr)
	orderModel.PayType = order.WXPAY

	if len(couponRelID) > 0 {
		// 有使用优惠券
		err = _this.CouponService.UserUseCoupon(couponRelID)
		if err != nil {
			return nil, err
		}
	}

	count, err := _this.OrderDao.GetOrderCount(map[string]any{
		"trade_no": TradeNo,
	})
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		// 没有这条记录
		err = _this.Save(&orderModel, 0)
		if err != nil {
			return
		}
		// 开启一个定时任务， 15分钟后执行
		ticker := time.NewTicker(15 * time.Minute)

		// 在一个独立的 goroutine 中执行无效订单删除任务
		go func() {
			for {
				<-ticker.C // 等待定时器触发
				var _orderModel order.Order
				filter := map[string]any{
					"trade_no": orderModel.TradeNo,
				}
				err := _this.OrderDao.FilterFindOne(&_orderModel, filter)
				if err != nil {
					logger.Error(err.Error())
					return
				}

				if !_orderModel.IsPayed {
					// 没有支付，无效订单
					err = _this.OrderDao.Delete(_orderModel.ID)
					if err != nil {
						return
					}
				}
			}
		}()

	}

	svc := jsapi.JsapiApiService{
		Client: wxpay.WxPayClient,
	}
	conf := config.GetConfig()

	res, _, err = svc.PrepayWithRequestPayment(wxpay.WxPayCtx, jsapi.PrepayRequest{
		Appid:       core.String(conf.WechatAppConfiguration.AppId),
		Mchid:       core.String(conf.WxPayConfig.MchId),
		Description: core.String(conf.WxPayConfig.Description),
		OutTradeNo:  core.String(TradeNo),
		Attach:      core.String("自定义数据说明"),
		NotifyUrl:   core.String(conf.WxPayConfig.NotifyUrl),
		Amount: &jsapi.Amount{
			Total:    core.Int64(int64(orderModel.Total.InexactFloat64() * 100)),
			Currency: core.String(conf.WxPayConfig.Currency),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(wxOpenID), //
		},
	})

	return
}

// PayOrder 来自微信支付的回调
func (_this orderService) PayOrder(wxCallback *wx_api.WxPayCallback) error {

	filter := map[string]any{
		"trade_no": wxCallback.TransIdo,
	}

	// start 修改订单状态
	var orderModel order.Order
	err := _this.OrderDao.FilterFindOne(&orderModel, filter)
	if err != nil {
		return err
	}
	orderModel.IsPayed = true
	orderModel.Status = order.ING
	orderModel.PayType = order.WXPAY
	err = _this.OrderDao.Update(&orderModel, orderModel.ID)
	if err != nil {
		return err
	}
	// end 修改订单状态

	switch orderModel.Type {
	case SERVICE:
		{
			return _this.handleService(wxCallback.PayUser.Openid, &orderModel)
		}
	case GOOD:
		{

			return _this.handleGood(wxCallback.PayUser.Openid, &orderModel)
		}
	case DEPOSIT:
		{
			return _this.handleDeposit(wxCallback.PayUser.Openid, &orderModel)
		}
	}

	return nil
}

// 处理充值回调
func (_this orderService) handleDeposit(wxOpenID string, orderModel *order.Order) error {
	// start 这里需要一个方法来检查会员充值
	// 通过 openID 找到用户
	var userModel user.User
	userFilter := map[string]any{
		"wx_open_id": wxOpenID,
	}
	err := _this.UserDao.FilterFindOne(&userModel, userFilter)
	if err != nil {
		return err
	}
	// end 这里需要一个方法来检查会员充值

	// start 修改会员余额
	userModel.Account = userModel.Account.Add(orderModel.ProductList[0].Price)
	return _this.UserDao.Update(&userModel, userModel.ID)

	// end 修改会员余额
}

func (_this orderService) handleService(wxOpenID string, orderModel *order.Order) error {
	// todo 发送提醒消息给用户
	// 店员状态锁定
	var clerkModel clerk.Clerk
	err := _this.ClerkDao.FindByID(&clerkModel, orderModel.ClerkID)
	if err != nil {
		return err
	}
	return nil
}

func (_this orderService) handleGood(wxOpenID string, orderModel *order.Order) error {
	// todo

	return nil
}

func (_this orderService) SetTimeTask(d time.Duration, orderID string, orderStatus order.OrderStatus) {
	// 创建一个每隔 1 小时触发一次的定时器

	ticker := time.NewTicker(d)

	// 在一个独立的 goroutine 中执行订单状态更新任务
	go func() {
		for {
			<-ticker.C // 等待定时器触发
			// todo 操作订单状态
			fields := map[string]any{
				"status": orderStatus,
			}

			filter := map[string]any{
				"id": orderID,
			}
			err := _this.OrderDao.UpdateFromMap(fields, filter)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			//orderStatus = "已发货"
			//fmt.Println("订单状态已更新为：", orderStatus)
		}
	}()
}

func (_this orderService) GetMyOrder(param *BaseSearchParam, userID uint64) (list []*order.Order, err error) {
	filter := map[string]any{
		"user_id": userID,
	}
	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Size)
	if err != nil {
		return
	}
	start := (page - 1) * size
	list, err = _this.OrderDao.ToList(start, size, filter, "", "")

	return
}

func (_this orderService) GetMyOrderCount(filter map[string]any) (int64, error) {
	return _this.OrderDao.GetOrderCount(filter)
}

func (_this orderService) HandleDateOrder(TradeNo string, in *order.Order, userID uint64) error {
	var (
		orderKey = TradeNo + strconv.FormatUint(userID, 10)
	)
	err := _this.OrderDao.CacheGetValue(in, orderKey)
	if err != nil && err != redis.Nil {
		return err
	}

	productStr, err := json.Marshal(in.ProductList)

	in.Products = string(productStr)
	in.IsPayed = true
	in.PayType = order.DATE
	in.Total = decimal.NewFromFloat(0.00)
	in.Status = order.ING
	return _this.Save(in, 0)

}

func (_this orderService) ComputeOrder(orderModel *order.Order, userID uint64) error {
	orderKey := orderModel.TradeNo + strconv.FormatUint(userID, 10)
	return _this.OrderDao.CacheChangeValue(orderKey, orderModel)
}

// GenerateDepositOrder 处理充值
func (_this orderService) GenerateDepositOrder(in *order.Order, depositNum string, userID uint64) (err error) {

	total, err := decimal.NewFromString(depositNum)

	if err != nil {
		return
	}

	var cartItem = cart.CartProduct{
		Cover:       "https://uls.yishares.com/web/statics/images/vip-card/bg-3.png",
		Title:       "会员卡充值",
		Price:       total,
		Num:         1,
		MerID:       "",
		ProductType: "deposit",
		ProductID:   "",
		Sku:         total.String(),
		ArriveTime:  "",
	}
	in.ProductList = append(in.ProductList, &cartItem)
	in.TradeNo = Node.Generate().String()
	in.Status = order.PENDING
	in.PayType = order.WXPAY
	products, err := json.Marshal(&cartItem)
	in.Products = string(products)
	in.IsPayed = false
	in.UserID = strconv.FormatUint(userID, 10)
	in.Total = total
	in.Type = "deposit"

	orderKey := in.TradeNo + strconv.FormatUint(userID, 10)
	err = _this.OrderDao.CacheSave(in, orderKey, 15*time.Minute)
	return nil
}

func (_this orderService) getAddress(in *order.Order, addressID any) error {
	err := _this.AddressDao.FilterFindOne(&in.Address, map[string]any{
		"id": addressID,
	})
	addr, err := json.Marshal(in.Address)
	if err != nil {
		return err
	}
	in.Addr = string(addr)
	return err
}

func (_this orderService) BalancePay(in *order.Order) error {

	return nil
}

func (_this orderService) OrderTask() error {

	// 首先查询数据库创建时间超过15分钟，并且未支付的订单

	logger.Info("list", "=======================")
	now := time.Now()
	before15min := now.Add((-15) * time.Minute)

	var list []*order.Order

	err := db.DB.Model(&order.Order{}).Where("create_at <= ?", before15min).Where("is_payed = 0").Find(&list).Error
	if err != nil {
		return err
	}

	for _, item := range list {
		err = _this.OrderDao.Delete(item.ID)
		if err != nil {
			return err
		}
	}
	//_this.OrderDao.

	return nil
}
