package application

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

//https://github.com/gin-gonic/gin/issues/339#issuecomment-111694462

type PortalAPIService struct {
	Bean *Bean
}

func NewPortalAPIService(bean *Bean) *PortalAPIService {
	return &PortalAPIService{bean}
}

func (s *PortalAPIService) GetOrdersByAreaID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "areaId invalid."))
		return
	}

	var orders, err = s.Bean.APIRepository.GetRestaurantOrdersByAreaID(cookieCheckResult.Restaurant.ID, areaId)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "order not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetCategories(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	var orders, err = s.Bean.APIRepository.GetCategoriesByGroupID(cookieCheckResult.Restaurant.RestaurantGroupID)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "order not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetCommentsByProductID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	productID, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "categoryID invalid."))
		return
	}

	var results, err = s.Bean.APIRepository.GetRestaurantCommentsByProductID(cookieCheckResult.Restaurant.ID, productID)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "comments not found."))
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *PortalAPIService) GetProducts(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	categoryID, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "categoryID invalid."))
		return
	}

	var products, err = s.Bean.APIRepository.GetProductsByRestaurantIDAndCategoryID(cookieCheckResult.Restaurant.ID, categoryID)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "products not found."))
		return
	}
	c.JSON(http.StatusOK, products)
}

func (s *PortalAPIService) GetTablesByAreaID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "areaId invalid."))
		return
	}

	var orders, err = s.Bean.APIRepository.GetRestaurantTablesByAreaID(cookieCheckResult.Restaurant.ID, areaId)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "table not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) PutOrders(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsRedirect() {
		c.JSON(http.StatusUnauthorized, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	waiterID := cookieCheckResult.User.UserName

	orderRequest := entity.OrderRequest{}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, exception.GetError(exception.CodeValueInvalid))
		return
	}
	orderRequest.WaiterID = &waiterID
	orderRequest.RestaurantID = &cookieCheckResult.Restaurant.ID

	var orders, err = s.Bean.APIRepository.PutRestaurantOrders(cookieCheckResult.Restaurant.ID, orderRequest)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "table not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}
