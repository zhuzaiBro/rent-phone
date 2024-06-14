package v2

import (
	"net/http"
	. "rentServer/core/controller"
	"rentServer/pkg/model/order"
	"rentServer/pkg/model/reservation"
	"rentServer/pkg/service"
)

var _ ReservationController = (*reservationController)(nil)

type ReservationController interface {
	Handlerent(c *Context)
}

type reservationController struct {
	rentService service.ReservationService
}

func NewrentController() ReservationController {
	return &reservationController{
		rentService: service.NewrentService(),
	}
}

func (_this reservationController) Handlerent(c *Context) {

	var (
		rentModel  reservation.Reservation
		orderModel order.Order
	)
	err := c.ShouldBindJSON(&rentModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = _this.rentService.Handlerent(&rentModel, &orderModel, c.GetUint64(Uid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONOK(orderModel)
}
