package service

import (
	"errors"
	"gorm.io/gorm"
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/coupon"
	"rentServer/pkg/model/user_coupon_rel"
	"strconv"
)

type CouponService interface {
	Service[coupon.Coupon]
	Dispatch(userID uint64, couponID string) error
	UserUseCoupon(userCouponRelID string) error
	UserList(userID any, status coupon.CouponStatus, param *BaseSearchParam) (list []*user_coupon_rel.UserCouponRel, err error)
}

type couponServiceIMPL struct {
	ServiceIMPL[coupon.Coupon]
	CouponDao        dao.CouponDao
	UserCouponRelDao dao.UserCouponRelDao
	UserDao          dao.UserDao
}

func NewCouponService() CouponService {
	couponDao := dao.NewCouponDao()
	return &couponServiceIMPL{
		ServiceIMPL: ServiceIMPL[coupon.Coupon]{
			BaseDao: couponDao,
		},
		CouponDao:        couponDao,
		UserCouponRelDao: dao.NewUserCouponRelDao(),
		UserDao:          dao.NewUserDao(),
	}
}

// Dispatch 发放优惠券到用户
func (_this couponServiceIMPL) Dispatch(userID uint64, couponID string) error {

	var (
		//userCouponRelModel  user_coupon_rel.UserCouponRel
		couponModel coupon.Coupon
	)
	err := _this.CouponDao.FindByID(&couponModel, couponID)
	if err != nil {
		return err
	}

	if couponModel.ID <= 0 {
		return errors.New("优惠券不存在！")
	}

	// 判断用户是否可以领取
	err = _this.check(userID, couponID)
	if err != nil {
		return err
	}

	var userCouponRelModel = user_coupon_rel.UserCouponRel{
		Category:       couponModel.ProductCategory,
		CouponID:       strconv.FormatUint(couponModel.ID, 10),
		Discount:       couponModel.DiscountValue,
		ExpirationDate: couponModel.EndDate,
		IsActive:       true,
		Product:        couponModel.ProductID,
		UserID:         strconv.FormatUint(userID, 10),
		CouponCode:     Node.Generate().String(),
	}

	err = _this.UserCouponRelDao.Insert(&userCouponRelModel)
	if err != nil {
		return err
	}

	// todo 微信通知消息

	return nil

}

// UserUseCoupon  用户消费优惠券
func (_this couponServiceIMPL) UserUseCoupon(userCouponRelID string) error {

	var userCouponRel user_coupon_rel.UserCouponRel
	filter := map[string]any{
		"id":        userCouponRelID,
		"is_active": true,
	}
	err := _this.UserCouponRelDao.FilterFindOne(&userCouponRel, filter)
	if err != nil {
		return errors.New("没有这张优惠券")
	}

	userCouponRel.IsActive = false
	err = _this.UserCouponRelDao.Update(&userCouponRel, userCouponRelID)
	if err != nil {
		return err
	}

	return nil
}

// UserList 用户查询自己的优惠券记录
func (_this couponServiceIMPL) UserList(userID any, status coupon.CouponStatus, param *BaseSearchParam) (list []*user_coupon_rel.UserCouponRel, err error) {

	page, err := strconv.Atoi(param.Page)
	size, err := strconv.Atoi(param.Size)
	if err != nil {
		return
	}

	start := (page - 1) * size

	list, err = _this.UserCouponRelDao.ToList(start, size, map[string]any{
		"user_id": userID,
	}, "", "")
	if err != nil {
		return nil, err
	}

	for _, rel := range list {

		err = _this.CouponDao.FindByID(&rel.Coupon, rel.CouponID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	}

	return
}

// check
func (_this couponServiceIMPL) check(userID uint64, couponID string) error {

	var couponRel user_coupon_rel.UserCouponRel
	// 理论上同样的优惠券，用户领取过一次就不能在领取了
	err := _this.UserCouponRelDao.FilterFindOne(&couponRel, map[string]any{
		"coupon_id": couponID,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if couponRel.ID > 0 {
		return errors.New("已经领取过了")
	}

	return nil
}
