package http

import (
	"github.com/gin-gonic/gin"
	. "rentServer/core/controller"
	v1 "rentServer/pkg/controller/v1"
	"rentServer/pkg/controller/v2"
	"rentServer/pkg/middle_ware"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	//router.Use(middle_ware.IPMiddleware())
	router.Static("/static", "./static") // Gin 1.7.0之前版本使用这种方式
	sign_controller := v2.NewSignController()
	router.GET("/spa/sign", Handle(sign_controller.WxappSign))
	api := router.Group("/spa")
	api.Use(middle_ware.WxappSignMiddleware())
	{
		api.Any("/", func(c *gin.Context) {
			c.JSON(0, "你好")
		})
		wxapp_img_controller := v2.NewWxappImgController()

		api.Any("/wxapp_img", Handle(wxapp_img_controller.GetImages))

	}

	admin := api.Group("/v1")

	{
		adminController := v1.NewAdminController()
		admin.POST("/users/login", Handle(adminController.Login))
		admin.GET("/users/info", AdminAuth(adminController.GetUserInfo))

		userRoot := admin.Group("/user")
		{
			userController := v1.NewUserController()
			userRoot.GET("/", AdminAuth(userController.List))
			userRoot.POST("/", AdminAuth(userController.Update))
		}

		good := admin.Group("/good")
		{
			goodController := v1.GetGoodController()
			_goodController := v2.NewGoodController()
			good.GET("/list", AdminAuth(goodController.List))
			good.GET("/detail", AdminAuth(goodController.Detail))
			good.POST("/", AdminAuth(goodController.Save))
			good.DELETE("/", AdminAuth(goodController.Delete))

			// 商品类目
			good.GET("/category", Handle(_goodController.GetCategoryList))
			good.POST("/category", AdminAuth(goodController.SaveCategory))
			good.DELETE("/category", AdminAuth(goodController.DeleteCategory))
		}

		uploadRoot := admin.Group("/upload")

		{
			uploadController := v1.NewUploadController()

			// 获取资源
			uploadRoot.GET("/", AdminAuth(uploadController.GetUploadResource))
			uploadRoot.POST("/", AdminAuth(uploadController.Upload))
			// 获取资源分组目录
			uploadRoot.GET("/group", AdminAuth(uploadController.GetGroup))
			// 添加或者修改分组
			uploadRoot.POST("/group", AdminAuth(uploadController.SaveGroup))
			// 删除资源分组
			uploadRoot.DELETE("/group", AdminAuth(uploadController.DeleteGroup))
		}

		clerkRoot := admin.Group("/clerk")
		{
			clerkController := v1.NewClerkController()
			clerkRoot.GET("/", AdminAuth(clerkController.GetClerkList))
			clerkRoot.POST("/", AdminAuth(clerkController.Save))
			clerkRoot.DELETE("/", AdminAuth(clerkController.Delete))
		}

		storeRoot := admin.Group("/store")
		{
			storeController := v1.NewStoreController()
			storeRoot.GET("/", AdminAuth(storeController.StoreInfo))
			storeRoot.POST("/", AdminAuth(storeController.Save))
		}

		orderRoot := admin.Group("/order")

		{
			orderController := v1.NewOrderController()

			orderRoot.GET("/", AdminAuth(orderController.GetOrders))
			orderRoot.PUT("/", AdminAuth(orderController.Update))
			orderRoot.GET("/detail", AdminAuth(orderController.Detail))

		}

		customProjRoot := admin.Group("/custom")
		{
			customProjController := v1.NewCustomProjController()
			customProjRoot.GET("/", Handle(customProjController.Get))
			customProjRoot.POST("/", AdminAuth(customProjController.Save))
			customProjRoot.DELETE("/", AdminAuth(customProjController.Delete))
		}

		couponRoot := admin.Group("/coupon")
		{
			couponController := v1.NewCouponController()
			couponRoot.GET("/", AdminAuth(couponController.GetCouponRecord))
			couponRoot.POST("/", AdminAuth(couponController.SaveCoupon))
			couponRoot.DELETE("/", AdminAuth(couponController.DeleteCoupon))
		}

		albumRoot := admin.Group("/album")
		{
			albumController := v1.NewAlbumController()
			{
				albumRoot.GET("/", AdminAuth(albumController.List))
				albumRoot.POST("/", AdminAuth(albumController.Save))
				albumRoot.DELETE("/", AdminAuth(albumController.Delete))
			}
		}

	}

	user := api.Group("/v2")

	{
		userController := v2.NewUserController()
		user.POST("/login", Handle(userController.MpLogin))
		user.GET("/info", Auth(userController.GetUserInfo))
		user.POST("/info", Auth(userController.UpdateUserInfo))
		user.POST("/wx/phone", Auth(userController.BindUserWxPhone))
		user.POST("/upload", Auth(userController.UploadAvatar))

		order := user.Group("/order")
		{
			orderController := v2.NewOrderController()
			order.POST("/wxPay", Auth(orderController.WxPayOrder))
			order.POST("/", Auth(orderController.GenerateOrder))
			order.GET("/", Auth(orderController.GetMyOrder))
			order.POST("/compute", Auth(orderController.ComputeOrder))
			order.GET("/detail", Auth(orderController.Detail))
			order.GET("/count", Auth(orderController.GetOrderCount))
		}

		good := user.Group("/good")
		{
			goodController := v2.NewGoodController()
			good.GET("/", Handle(goodController.List))
			good.GET("/category", Handle(goodController.GetCategoryList))
			good.GET("/detail", Handle(goodController.Detail))
		}

		clerk := user.Group("/clerk")
		{
			clerkController := v2.NewClerkController()
			clerk.GET("/", Handle(clerkController.List))
		}

		rent := user.Group("/rent")
		{
			rentController := v2.NewrentController()
			rent.POST("/", Auth(rentController.Handlerent))
		}

		storeRoot := user.Group("/store")
		{
			storeController := v1.NewStoreController()
			storeRoot.GET("/", Handle(storeController.StoreInfo))
		}

		cartRoot := user.Group("/cart")
		{
			cartController := v2.NewCartController()
			cartRoot.POST("/", Auth(cartController.Add))
			cartRoot.GET("/", Auth(cartController.List))
			cartRoot.DELETE("/", Auth(cartController.Delete))
		}

		wxRoot := user.Group("/wx")
		{
			wxController := v2.NewPayController()
			{
				wxRoot.POST("/callback", Handle(wxController.WxPayCallback))
			}
		}

		addressRoot := user.Group("/address")
		{
			addressController := v2.NewAddressController()
			addressRoot.GET("/", Auth(addressController.GetMyAddress))
			addressRoot.POST("/", Auth(addressController.Save))
			addressRoot.DELETE("/", Auth(addressController.Delete))
		}

		couponRoot := user.Group("/coupon")
		{
			couponController := v2.NewCouponController()
			couponRoot.GET("/", Auth(couponController.UserGetCouponList))
			couponRoot.POST("/", Auth(couponController.FetchCoupon))
		}

		albumRoot := user.Group("/album")
		{
			albumController := v2.NewAlbumController()
			{
				albumRoot.GET("/", AdminAuth(albumController.List))
			}
		}

		commentRoot := user.Group("/comment")
		{
			commentController := v2.NewCommentController()
			commentRoot.POST("/", Auth(commentController.Post))
			commentRoot.GET("/", Handle(commentController.List))
		}

		uploadRoot := user.Group("/upload")

		{
			uploadController := v1.NewUploadController()
			uploadRoot.POST("/", Auth(uploadController.Upload))
		}
	}

	return router
}
