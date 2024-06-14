package middle_ware

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rentServer/initilization/db"
)

func WxappSignMiddleware() gin.HandlerFunc {

	cacheClient := db.GetRedis()
	return func(c *gin.Context) {

		sign := c.GetHeader("w-sign")
		err := cacheClient.Get(context.Background(), sign).Err()

		log.Println(err)

		if err != nil {
			c.JSON(200, map[string]any{
				"code": http.StatusProxyAuthRequired,
				"msg":  "UnSigned!",
				//"err":  err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}

}
