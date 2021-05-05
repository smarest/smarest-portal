package application

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-common/util"
)

//https://github.com/gin-gonic/gin/issues/339#issuecomment-111694462

type PortalAPIService struct {
	Bean *Bean
}

func NewPortalAPIService(bean *Bean) *PortalAPIService {
	return &PortalAPIService{bean}
}

func (s *PortalAPIService) GetOrdersByAreaID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		s.HandlerError(c, exception.CreateErrorWithRootCause(exception.CodeValueInvalid, "areaId invalid.", paramErr))
		return
	}

	var orders, err = s.Bean.APIRepository.GetRestaurantOrdersByAreaID(cookieCheckResult.Restaurant.ID, areaId)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetOrders(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	orderBy := c.DefaultQuery("orderBy", util.ORDER_SORT_BY_PRODUCT)
	groupBy := c.Query("groupBy")
	var orders, err = s.Bean.APIRepository.GetRestaurantOrders(cookieCheckResult.Restaurant.ID, orderBy, groupBy)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetCategories(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	var orders, err = s.Bean.APIRepository.GetCategoriesByGroupID(cookieCheckResult.Restaurant.RestaurantGroupID)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetCommentsByProductID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	productID, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		s.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, "productID invalid."))
		return
	}

	var results, err = s.Bean.APIRepository.GetRestaurantCommentsByProductID(cookieCheckResult.Restaurant.ID, productID)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *PortalAPIService) GetProducts(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	categoryID, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		s.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, "categoryID invalid."))
		return
	}

	var products, err = s.Bean.APIRepository.GetProductsByRestaurantIDAndCategoryID(cookieCheckResult.Restaurant.ID, categoryID)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

func (s *PortalAPIService) GetTablesByAreaID(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		s.HandlerError(c, exception.CreateError(exception.CodeValueInvalid, "areaId invalid."))
		return
	}

	var orders, err = s.Bean.APIRepository.GetRestaurantTablesByAreaID(cookieCheckResult.Restaurant.ID, areaId)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) PutOrders(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c, "")
	if cookieCheckResult.IsRedirect() {
		s.HandlerError(c, exception.GetError(exception.CodeSignatureInvalid))
		return
	}

	waiterID := cookieCheckResult.User.UserName

	orderRequest := entity.OrderRequest{}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		s.HandlerError(c, exception.CreateErrorWithRootCause(exception.CodeValueInvalid, "orderRequest invalid", err))
		return
	}
	orderRequest.WaiterID = &waiterID
	orderRequest.RestaurantID = &cookieCheckResult.Restaurant.ID

	var orders, err = s.Bean.APIRepository.PutRestaurantOrders(cookieCheckResult.Restaurant.ID, orderRequest)
	if err != nil {
		s.HandlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) HandlerError(c *gin.Context, err *exception.Error) {
	log.Printf("Message=[%s], error=[%s]", err.ErrorMessage, err.RootCause)
	switch err.ErrorCode {
	case exception.CodeNotFound:
		c.JSON(http.StatusNotFound, exception.GetError(err.ErrorCode))
	case exception.CodeSignatureInvalid:
		c.JSON(http.StatusUnauthorized, exception.GetError(err.ErrorCode))
	case exception.CodeValueInvalid:
		c.JSON(http.StatusBadRequest, exception.GetError(err.ErrorCode))
	default:
		c.JSON(http.StatusInternalServerError, exception.GetError(exception.CodeSystemError))
	}
}
