package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-portal/application"
)

func main() {
	bean, err := application.InitBean()
	if err != nil {
		log.Fatalln("can not create bean", err)
	}

	router := gin.Default()
	router.HTMLRender = bean.LoadTemplates()
	portal := router.Group("portal")
	{
		portal.GET("/cashier", bean.CashierService.Get)
		portal.GET("/order", bean.OrderService.Get)
		portal.GET("/checkout", bean.CheckoutService.Get)
		portal.GET("/restaurant", bean.RestaurantService.Get)
		portal.POST("/restaurant", bean.RestaurantService.Post)
	}

	v1 := router.Group("v1")
	{
		v1.GET("/areas/:id/orders", bean.PortalAPIService.GetOrdersByAreaID)
		v1.GET("/areas/:id/tables", bean.PortalAPIService.GetTablesByAreaID)
		v1.GET("/categories", bean.PortalAPIService.GetCategories)
		v1.GET("/categories/:id/products", bean.PortalAPIService.GetProducts)
		v1.GET("/comments", bean.PortalAPIService.GetComments)
		v1.PUT("/orders", bean.PortalAPIService.PutOrders)
		v1.DELETE("/orders", bean.PortalAPIService.DeleteOrders)
	}
	router.Run(":9090")
}
