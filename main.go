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
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.HTMLRender = bean.LoadTemplates()
	portal := router.Group("portal")
	{
		portal.GET("/cashier", bean.CashierService.Get)
		portal.GET("/order", bean.OrderService.Get)
		portal.GET("/checkout", bean.CheckoutService.Get)
		portal.GET("/restaurant", bean.RestaurantService.Get)
		portal.POST("/restaurant", bean.RestaurantService.Post)
		portal.GET("/error", bean.ErrorService.Get)
	}

	v1 := router.Group("v1")
	{
		v1.GET("/areas/:id/orders", bean.PortalAPIService.GetOrdersByAreaID)
		v1.GET("/areas/:id/tables", bean.PortalAPIService.GetTablesByAreaID)
		v1.GET("/categories", bean.PortalAPIService.GetCategories)
		v1.GET("/categories/:id/products", bean.PortalAPIService.GetProducts)
		v1.GET("/products/:id/comments", bean.PortalAPIService.GetCommentsByProductID)
		v1.PUT("/orders", bean.PortalAPIService.PutOrders)
	}
	router.Run(":9090")
}
