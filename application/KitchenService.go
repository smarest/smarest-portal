package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type KitchenService struct {
	Bean *Bean
}

func NewKitchenService(bean *Bean) *KitchenService {
	return &KitchenService{bean}
}

func (s *KitchenService) Get(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, PageKitchen)
	if cookieCheckResult.IsRedirect() {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}

	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.IsKitchen = true
	resource.PageTitle = "Kitchen"

	c.HTML(http.StatusOK, PageKitchen, gin.H{
		"resource": resource,
	})
}
