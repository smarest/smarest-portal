package application

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/infrastructure/persistence"
)

type OrderService struct {
	*application.LoginService
	apiRepository       persistence.APIRepository
	PageResourceFactory resource.PageResourceFactory
}

func NewOrderService(loginService *application.LoginService, APIRepository persistence.APIRepository, pageResourceFactory resource.PageResourceFactory) *OrderService {
	return &OrderService{loginService, APIRepository, pageResourceFactory}
}

func (s *OrderService) Get(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	resource := s.PageResourceFactory.CreateResource()
	resource.IsOrder = true
	resource.PageTitle = "Order"

	areas, err := s.apiRepository.GetAreasByRestaurantID(1)
	if err != nil {
		log.Print(err.ErrorMessage)
		resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
	} else {
		resource.Areas = areas
		if len(areas) > 0 {
			areaID := int64(areas[0].(map[string]interface{})["id"].(float64))
			resource.Tables, err = s.apiRepository.GetTablesByAreaID(areaID)
			if err != nil {
				log.Print(err.ErrorMessage)
				resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
			}
			resource.AreaID = fmt.Sprint(areaID)
		}
	}

	orderNumberIDStr := c.Query("orderNumberID")

	if orderNumberIDStr != "" {
		orderNumberIDInt, paramErr := strconv.ParseInt(orderNumberIDStr, 0, 64)
		if paramErr != nil {
			resource.ErrorMessage = "orderNumberID invalid."
		} else {
			orders, err := s.apiRepository.GetOrdersByOrderNumberID(orderNumberIDInt)
			if err != nil {
				log.Printf("GetOrdersByOrderNumberID: %s", err.ErrorMessage)
				resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
			}
			resource.Orders = orders
			resource.OrderNumberID = orderNumberIDStr
		}
	}

	resource.Categories, err = s.apiRepository.GetCategories()
	if err != nil {
		log.Printf("GetCategories: %s", err.ErrorMessage)
		resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
	} else {
		if len(resource.Categories) > 0 {
			categoryID := int64(resource.Categories[0].(map[string]interface{})["id"].(float64))
			resource.Products, err = s.apiRepository.GetProductsByRestaurantIDAndCategoryID(1, categoryID)
			if err != nil {
				log.Print(err.ErrorMessage)
				resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
			}
			resource.CategoryID = fmt.Sprint(categoryID)
		}
	}
	resource.Comments, err = s.apiRepository.GetComments()
	if err != nil {
		log.Printf("GetComments: %s", err.ErrorMessage)
		resource.ErrorMessage = "Co loi trong he thong, vui long lien he bo phan ky thuat"
	}

	c.HTML(http.StatusOK, "order", gin.H{
		"resource": resource,
	})
}
