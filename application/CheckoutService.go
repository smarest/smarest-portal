package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/infrastructure/persistence"
)

type CheckoutService struct {
	*application.LoginService
	apiRepository       persistence.APIRepository
	PageResourceFactory resource.PageResourceFactory
}

func NewCheckoutService(loginService *application.LoginService, APIRepository persistence.APIRepository, pageResourceFactory resource.PageResourceFactory) *CheckoutService {
	return &CheckoutService{loginService, APIRepository, pageResourceFactory}
}

func (s *CheckoutService) Get(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	resource := s.PageResourceFactory.CreateResource()
	resource.IsCashier = true
	resource.PageTitle = "Checkout"

	orderNumberIDStr := c.Query("orderNumberID")
	if orderNumberIDStr == "" {
		resource.ErrorMessage = "orderNumberID required."
		c.HTML(http.StatusOK, "checkout", gin.H{
			"resource": resource,
		})
		return
	}

	orderNumberIDInt, paramErr := strconv.ParseInt(orderNumberIDStr, 0, 64)
	if paramErr != nil {
		resource.ErrorMessage = "orderNumberID invalid."
		c.HTML(http.StatusOK, "checkout", gin.H{
			"resource": resource,
		})
		return
	}
	orders, err := s.apiRepository.GetOrdersByOrderNumberID(orderNumberIDInt)
	if err != nil {
		resource.ErrorMessage = err.ErrorMessage
		c.HTML(http.StatusOK, "checkout", gin.H{
			"resource": resource,
		})
		return
	}
	resource.Orders = orders
	c.HTML(http.StatusOK, "checkout", gin.H{
		"resource": resource,
	})
	return
}
