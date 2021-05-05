package application

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type OrderService struct {
	Bean *Bean
}

func NewOrderService(bean *Bean) *OrderService {
	return &OrderService{bean}
}

func (s *OrderService) Get(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, PageOrder)
	if cookieCheckResult.IsRedirect() {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}

	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.IsOrder = true
	resource.PageTitle = "Order"

	areas, err := s.Bean.APIRepository.GetAreasByRestaurantID(cookieCheckResult.Restaurant.ID)
	if err != nil {
		s.Bean.ErrorService.HandlerError(c, err)
		return
	} else {
		resource.Areas = areas
		if len(areas) > 0 {
			areaID := int64(areas[0].(map[string]interface{})["id"].(float64))
			resource.Tables, err = s.Bean.APIRepository.GetRestaurantTablesByAreaID(cookieCheckResult.Restaurant.ID, areaID)
			if err != nil {
				s.Bean.ErrorService.HandlerError(c, err)
				return
			}
			resource.AreaID = fmt.Sprint(areaID)
		}
	}

	orderNumberIDStr := c.Query("orderNumberID")

	if orderNumberIDStr != "" {
		orderNumberIDInt, paramErr := strconv.ParseInt(orderNumberIDStr, 0, 64)
		if paramErr != nil {
			s.Bean.ErrorService.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, paramErr.Error()))
			return
		} else {
			orders, err := s.Bean.APIRepository.GetRestaurantOrdersByOrderNumberID(cookieCheckResult.Restaurant.ID, orderNumberIDInt)
			if err != nil {
				s.Bean.ErrorService.HandlerError(c, err)
				return
			}
			resource.Orders = orders
			resource.OrderNumberID = orderNumberIDStr
		}
	}

	resource.Categories, err = s.Bean.APIRepository.GetCategoriesByGroupID(cookieCheckResult.Restaurant.RestaurantGroupID)
	if err != nil {
		s.Bean.ErrorService.HandlerError(c, err)
		return
	} else {
		if len(resource.Categories) > 0 {
			categoryID := int64(resource.Categories[0].(map[string]interface{})["id"].(float64))
			resource.Products, err = s.Bean.APIRepository.GetProductsByRestaurantIDAndCategoryID(cookieCheckResult.Restaurant.ID, categoryID)
			if err != nil {
				s.Bean.ErrorService.HandlerError(c, err)
				return
			}
			resource.CategoryID = fmt.Sprint(categoryID)
		}
	}

	c.HTML(http.StatusOK, "order", gin.H{
		"resource": resource,
	})
}
