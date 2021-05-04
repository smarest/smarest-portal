package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CashierService struct {
	Bean *Bean
}

func NewCashierService(bean *Bean) *CashierService {
	return &CashierService{bean}
}

func (s *CashierService) Get(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}

	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.IsCashier = true
	resource.PageTitle = "Cashier"

	areas, err := s.Bean.APIRepository.GetAreasByRestaurantID(cookieCheckResult.Restaurant.ID)
	if err != nil {
		s.Bean.ErrorService.HandlerError(c, err)
		return
	}

	resource.Areas = areas
	if len(areas) > 0 {
		areaID := int64(areas[0].(map[string]interface{})["id"].(float64))
		resource.Orders, err = s.Bean.APIRepository.GetRestaurantOrdersByAreaID(cookieCheckResult.Restaurant.ID, areaID)
		if err != nil {
			s.Bean.ErrorService.HandlerError(c, err)
			return
		}
		resource.AreaID = fmt.Sprint(areaID)
	}

	c.HTML(http.StatusOK, PageCashier, gin.H{
		"resource": resource,
	})
}
