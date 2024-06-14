package v2

import (
	"context"
	"github.com/redis/go-redis/v9"
	. "rentServer/core/controller"
	"rentServer/initilization/db"
	"rentServer/pkg/oauth"
	"time"
)

type SignController interface {
	WxappSign(c *Context)
}

type signController struct {
	cacheClient *redis.Client
}

func (s signController) WxappSign(c *Context) {
	//TODO 前端会传过来一个 js_code
	jsCode := c.Query("code")
	resp, err := oauth.Code2Session(jsCode)
	if err != nil {
		//print(err)
		c.JSONE(500, err.Error(), err)
		return
	}

	//session.SessionKey
	//session.OpenId 可以对违规行为封号了
	//session.UnionId
	go func() {
		// 允许 1 小时内的访问
		s.cacheClient.Set(context.Background(), resp.SessionKey, resp.OpenId, 1*time.Hour)
	}()

	c.JSONOK(resp)
}

func NewSignController() SignController {
	return &signController{
		cacheClient: db.GetRedis(),
	}
}
