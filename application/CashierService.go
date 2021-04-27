package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/infrastructure/persistence"
)

type CashierService struct {
	*application.LoginService
	apiRepository       persistence.APIRepository
	PageResourceFactory resource.PageResourceFactory
}

func NewCashierService(loginService *application.LoginService, APIRepository persistence.APIRepository, pageResourceFactory resource.PageResourceFactory) *CashierService {
	return &CashierService{loginService, APIRepository, pageResourceFactory}
}

func (s *CashierService) Get(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	resource := s.PageResourceFactory.CreateResource()
	resource.IsCashier = true
	resource.PageTitle = "Cashier"

	areas, err := s.apiRepository.GetAreasByRestaurantID(1)
	if err != nil {
		log.Print(err.ErrorMessage)
		resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
	} else {
		resource.Areas = areas
		if len(areas) > 0 {
			areaID := int64(areas[0].(map[string]interface{})["id"].(float64))
			resource.Orders, err = s.apiRepository.GetOrdersByAreaID(areaID)
			if err != nil {
				log.Print(err.ErrorMessage)
				resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
			}
			resource.AreaID = fmt.Sprint(areaID)
		}
	}

	c.HTML(http.StatusOK, "cashier", gin.H{
		"resource": resource,
	})
}
