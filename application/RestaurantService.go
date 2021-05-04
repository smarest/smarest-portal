package application

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/util"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/domain/entity"
)

type RestaurantService struct {
	Bean *Bean
}

func NewRestaurantService(bean *Bean) *RestaurantService {
	return &RestaurantService{bean}
}

func (s *RestaurantService) Get(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsLoginRedirect {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}

	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.FromURL = c.DefaultQuery(util.PARAM_FROM_URL, PageCashier)

	s.GetRestaurantsByGroupID(resource, cookieCheckResult.User.RestaurantGroupID)
	c.HTML(http.StatusOK, "restaurant", gin.H{
		"resource": resource,
	})
}

func (s *RestaurantService) Post(c *gin.Context) {
	cookieCheckResult := s.Bean.CookieCheckService.Check(c)
	if cookieCheckResult.IsLoginRedirect {
		c.Redirect(http.StatusMovedPermanently, cookieCheckResult.RedirectURL)
		return
	}
	restaurantCodeStr := c.PostForm("code")
	fromURL := c.DefaultPostForm(util.PARAM_FROM_URL, PageCashier)
	resource := s.Bean.PageResourceFactory.CreateResource()
	resource.FromURL = fromURL

	if restaurantCodeStr == "" {
		resource.AddErrorMessage("RestaurantCode is required")
		s.GetRestaurantsByGroupID(resource, cookieCheckResult.User.RestaurantGroupID)
		c.HTML(http.StatusOK, "restaurant", gin.H{
			"resource": resource,
		})
		return
	}

	restaurantRaw, apiErr := s.Bean.APIRepository.GetRestaurantByCode(restaurantCodeStr)
	if apiErr != nil {
		log.Println(apiErr)
		s.GetRestaurantsByGroupID(resource, cookieCheckResult.User.RestaurantGroupID)
		c.HTML(http.StatusOK, "restaurant", gin.H{
			"resource": resource,
		})
		return
	}

	restaurant := entity.CreateRestaurantFromSlice(restaurantRaw)

	tokenStr, tErr := s.Bean.SecretService.GenerateRestaurantToken(restaurant)
	if tErr != nil {
		log.Println(tErr)
		resource.AddErrorMessage("Phat sinh loi! vui long lien he bo phan ky thuat")

		c.HTML(http.StatusOK, "restaurant", gin.H{
			"resource": resource,
		})
		return
	}
	c.SetCookie(s.Bean.COOKIE_TOKEN_RESTAURANT, tokenStr, 60*60*24, "", "", http.SameSiteDefaultMode, false, false)
	c.Redirect(http.StatusMovedPermanently, "/portal/"+fromURL)
}

func (s *RestaurantService) GetRestaurantsByGroupID(resource *resource.PageResource, restaurantGroupID int64) {
	results, err := s.Bean.APIRepository.GetRestaurantsByGroupID(restaurantGroupID)

	if err != nil {
		log.Println(err)
		resource.AddErrorMessage("Phat sinh loi! vui long lien he bo phan ky thuat")
	} else {
		resource.Restaurants = results
	}
}
