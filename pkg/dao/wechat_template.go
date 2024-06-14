package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/wechat_template"
)

var _ WechatTemplateDao = (*wechatTemplateDao)(nil)

type WechatTemplateDao interface {
	Dao[wechat_template.WechatTemplate]
}

type wechatTemplateDao struct {
	DaoIMPL[wechat_template.WechatTemplate]
}

func NewWechatTemplateDao() WechatTemplateDao {
	return &wechatTemplateDao{
		DaoIMPL: DaoIMPL[wechat_template.WechatTemplate]{
			Db:    db.DB,
			Model: &wechat_template.WechatTemplate{},
		},
	}
}
