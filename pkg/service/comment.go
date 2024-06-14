package service

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/comment"
	"rentServer/pkg/model/order"
)

type CommentService interface {
	Service[comment.Comment]
}

type commentService struct {
	ServiceIMPL[comment.Comment]
	CommentDao   dao.CommentDao
	OrderService OrderService
}

func NewCommentService() CommentService {
	commentDao := dao.NewCommentDao()
	return &commentService{
		ServiceIMPL: ServiceIMPL[comment.Comment]{
			BaseDao: commentDao,
		},
		CommentDao:   commentDao,
		OrderService: NewOrderService(),
	}
}

// PostUserComment 用户发送评论
func (_this commentService) PostUserComment(commentModel *comment.Comment, ) error {

	var orderModel order.Order
	err := _this.OrderService.FilterFindOne(&orderModel, map[string]any{
		"id": commentModel.OrderID,
	})
	if err != nil {
		return err
	}

	//productStr, err := json.Marshal(orderModel.ProductList)
	//commentModel.Pro = string(productStr)

	err = _this.CommentDao.Insert(commentModel)
	if err != nil {
		return err
	}

	return nil
}
