package application

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-common/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-portal/infrastructure/persistence"
)

//https://github.com/gin-gonic/gin/issues/339#issuecomment-111694462

type PortalAPIService struct {
	*application.LoginService
	apiRepository persistence.APIRepository
}

func NewPortalAPIService(loginService *application.LoginService, APIRepository persistence.APIRepository) *PortalAPIService {
	return &PortalAPIService{apiRepository: APIRepository}
}

func (s *PortalAPIService) GetOrdersByAreaID(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "areaId invalid."))
		return
	}

	var orders, err = s.apiRepository.GetOrdersByAreaID(areaId)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "order not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetCategories(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	var orders, err = s.apiRepository.GetCategories()
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "order not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) GetComments(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	var results, err = s.apiRepository.GetComments()
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "comments not found."))
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *PortalAPIService) GetProducts(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	categoryID, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "categoryID invalid."))
		return
	}

	var products, err = s.apiRepository.GetProductsByRestaurantIDAndCategoryID(1, categoryID)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "products not found."))
		return
	}
	c.JSON(http.StatusOK, products)
}

func (s *PortalAPIService) GetTablesByAreaID(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	areaId, paramErr := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if paramErr != nil {
		log.Print(paramErr)
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "areaId invalid."))
		return
	}

	var orders, err = s.apiRepository.GetTablesByAreaID(areaId)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "table not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) PutOrders(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	waiterID := "dienami"
	restaurantID := int64(1)

	orderRequest := entity.OrderRequest{}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, exception.GetError(exception.CodeValueInvalid))
		return
	}
	orderRequest.WaiterID = &waiterID
	orderRequest.RestaurantID = &restaurantID

	log.Printf("%+v", orderRequest)
	var orders, err = s.apiRepository.PutOrders(orderRequest)
	if err != nil {
		log.Print(err.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "table not found."))
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (s *PortalAPIService) DeleteOrders(c *gin.Context) {
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/

	c.JSON(http.StatusOK, nil)
}
