package v2

import (
	"net/http"
	. "rentServer/core/controller"
	service2 "rentServer/core/service"
	"rentServer/pkg/service"
	"strconv"
)

type ClerkController interface {
	List(ctx *Context)
}

type clerkController struct {
	ClerkService service.ClerkService
	rentService  service.ReservationService
}

func NewClerkController() ClerkController {
	return &clerkController{
		ClerkService: service.NewClerkService(),
		rentService:  service.NewrentService(),
	}
}

func (_this clerkController) List(c *Context) {

	var param service2.BaseSearchParam

	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	filter := map[string]any{}
	list, err := _this.ClerkService.FilterList(&param, filter)
	if err != nil {
		return
	}

	timestamp := c.Query("timestamp")

	parseTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return
	}

	timeList := _this.rentService.GetServiceTime(parseTimestamp)
	for _, clerkModel := range list {
		clerkModel.ServiceTime = timeList
	}

	c.JSONOK(list)
}
