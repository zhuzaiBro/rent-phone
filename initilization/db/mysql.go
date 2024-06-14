package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"rentServer/pkg/config"
	"rentServer/pkg/model/address"
	"rentServer/pkg/model/admin"
	"rentServer/pkg/model/album"
	"rentServer/pkg/model/attr"
	"rentServer/pkg/model/attr_value"
	"rentServer/pkg/model/cart"
	"rentServer/pkg/model/category"
	"rentServer/pkg/model/clerk"
	"rentServer/pkg/model/comment"
	"rentServer/pkg/model/coupon"
	"rentServer/pkg/model/customized_project"
	"rentServer/pkg/model/good"
	"rentServer/pkg/model/order"
	"rentServer/pkg/model/resource"
	"rentServer/pkg/model/resource_group"
	"rentServer/pkg/model/store"
	"rentServer/pkg/model/user"
	"rentServer/pkg/model/user_coupon_rel"
)

// DB 全局mysql数据库对象
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	DB = ConnMysql()

	dbAutoMigrate()

	ConnectRedis()
}

// 自动迁移表结构
func dbAutoMigrate() {
	_ = DB.AutoMigrate(&good.Product{}, &admin.Admin{}, &resource_group.ResourceGroup{}, &resource.Resource{}, &clerk.Clerk{},
		&user.User{}, &attr_value.AttrValue{}, &attr.Attr{}, &store.Store{}, &cart.Cart{}, &album.Album{},
		&order.Order{}, &address.Address{}, &category.Category{}, &coupon.Coupon{}, &user_coupon_rel.UserCouponRel{}, &customized_project.CProject{},
		&comment.Comment{},
	)
}

func ConnMysql() *gorm.DB {
	conf := config.GetConfig()
	Db, err := gorm.Open(mysql.Open(conf.MysqlConf.Dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := Db.DB()

	sqlDB.SetMaxOpenConns(conf.MysqlConf.MaxOpenConns)
	sqlDB.SetMaxOpenConns(conf.MysqlConf.MaxOpenConns)

	return Db
}
