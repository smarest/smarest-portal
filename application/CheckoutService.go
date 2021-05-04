package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type CheckoutService struct {
	Bean *Bean
}

func NewCheckoutService(bean *Bean) *CheckoutService {
	return &CheckoutService{bean}
}

func (s *CheckoutService) Get(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}

	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.IsCashier = true
	resource.PageTitle = "Checkout"

	orderNumberIDStr := c.Query("orderNumberID")
	if orderNumberIDStr == "" {
		s.Bean.ErrorService.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, "orderNumberID is required"))
		return
	}

	orderNumberIDInt, paramErr := strconv.ParseInt(orderNumberIDStr, 0, 64)
	if paramErr != nil {
		s.Bean.ErrorService.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, "orderNumberID invalid."))
		return
	}
	orders, err := s.Bean.APIRepository.GetRestaurantOrdersByOrderNumberID(cookieCheckResult.Restaurant.ID, orderNumberIDInt)
	if err != nil {
		s.Bean.ErrorService.HandlerError(c, err)
		return
	}

	resource.Orders = orders
	c.HTML(http.StatusOK, PageCheckout, gin.H{
		"resource": resource,
	})
}
