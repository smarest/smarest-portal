package application

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	loginDomainService "github.com/smarest/smarest-account/domain/service"
	loginResource "github.com/smarest/smarest-common/client/resource"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-common/util"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/infrastructure/persistence"
)

type RestaurantService struct {
	*loginDomainService.RestaurantCookieService
	apiRepository       persistence.APIRepository
	PageResourceFactory resource.PageResourceFactory
}

func NewRestaurantService(APIRepository persistence.APIRepository, pageResourceFactory resource.PageResourceFactory) *RestaurantService {
	return &RestaurantService{loginDomainService.NewRestaurantCookieService("b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABD1WGaxt2"),
		APIRepository, pageResourceFactory}
}

func (s *RestaurantService) Get(c *gin.Context) {
	/*if !s.IsLogin(c) {
		log.Print("host %q", c.Request.URL.Host)
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?origin=%s", s.LoginUrl, c.Request.URL))
		return
	}*/

	resource := s.PageResourceFactory.CreateResource()
	resource.FromURL = c.DefaultQuery(util.PARAM_FROM_URL, "cashier")

	results, err := s.apiRepository.GetRestaurants()
	if err != nil {
		resource.ErrorMessage = err.ErrorMessage
	} else {
		resource.Restaurants = results
	}

	c.HTML(http.StatusOK, "restaurant", gin.H{
		"resource": resource,
	})
}

func (s *RestaurantService) Post(c *gin.Context) {
	/*if !s.IsLogin(c) {
		c.JSON(http.StatusUnauthorized, exception.CreateError(exception.CodeSignatureInvalid, "Access denied."))
		return
	}*/
	restaurantIdString := c.PostForm("restaurant_id")
	restaurantAccessKey := c.PostForm("access_key")
	fromURL := c.DefaultPostForm(util.PARAM_FROM_URL, "cashier")
	_, err := strconv.ParseInt(restaurantIdString, 10, 64)
	if err != nil {
		log.Print(err)
		resource := s.PageResourceFactory.CreateResource()
		resource.ErrorMessage = "restaurantId not found."
		resource.FromURL = fromURL
		results, err := s.apiRepository.GetRestaurants()
		if err != nil {
			resource.ErrorMessage = resource.ErrorMessage + " AND " + err.ErrorMessage
		} else {
			resource.Restaurants = results
		}
		c.HTML(http.StatusOK, "restaurant", gin.H{
			"resource": resource,
		})
		return
	}

	restaurant, apiErr := s.apiRepository.PostRestaurant(&loginResource.RestaurantRequestResource{
		ID:        restaurantIdString,
		AccessKey: restaurantAccessKey,
	})

	if apiErr != nil {
		log.Print(apiErr.ErrorMessage)
		c.JSON(http.StatusNotFound, exception.CreateError(apiErr.ErrorCode, "Phat sinh loi."))
		return
	}

	cookieString := restaurant.(map[string]interface{})["cookie"].(string)
	if cookieString == "" {
		log.Print(apiErr.ErrorMessage)
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeSystemError, "cookie is null"))
		return
	}
	claims, cookieErr := s.CheckRestaurantToken(cookieString)
	if cookieErr != nil {
		log.Print(cookieErr.ErrorMessage)
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeSystemError, "cookie is null"))
		return
	}

	fmt.Printf("%+v", claims.Restaurant)
	//c.SetCookie("test", restaurant.(*Cookie), 0, "/", "", http.SameSiteNoneMode, false, true)
	//c.Redirect(http.StatusMovedPermanently, "/portal/"+fromURL)
	//	c.JSON(http.StatusOK, gin.H{"filePath": c.Request.URL.String() + "/view/" + directory + "/" + baseFileName, "message": fmt.Sprintf("'%s' uploaded!", file.Filename)})
}
