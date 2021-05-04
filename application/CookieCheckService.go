package application

import (
	"log"

	"github.com/gin-gonic/gin"
	common_entity "github.com/smarest/smarest-common/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-portal/domain/entity"
)

type CookieCheckService struct {
	Bean *Bean
}

type CheckResult struct {
	IsRestaurantRedirect bool
	IsLoginRedirect      bool
	RedirectURL          string
	Restaurant           *entity.Restaurant
	User                 *common_entity.User
}

func (result *CheckResult) IsRedirect() bool {
	return result.IsLoginRedirect || result.IsRestaurantRedirect
}

func NewCookieCheckService(bean *Bean) *CookieCheckService {
	return &CookieCheckService{bean}
}

func (s *CookieCheckService) Check(c *gin.Context) *CheckResult {
	var err *exception.Error
	/*	_, err := s.CheckCookie(c)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, s.GetLoginUrl(c))
			return
		}
	*/
	/* */
	result := &CheckResult{IsLoginRedirect: false, IsRestaurantRedirect: false}
	result.User = &common_entity.User{
		UserName:          "test12345",
		RestaurantGroupID: 1,
		Name:              "test name",
		Role:              "A",
	}
	result.Restaurant, err = s.CheckRestaurantCookie(c)
	if err != nil {
		result.IsRestaurantRedirect = true
		result.RedirectURL = s.Bean.URLRepository.GetPageRedirectURL(PageCashier, PageRestaurant)
		return result
	}
	return result
}

func (s *CookieCheckService) CheckRestaurantCookie(c *gin.Context) (*entity.Restaurant, *exception.Error) {
	cookie, err := c.Cookie(s.Bean.COOKIE_TOKEN_RESTAURANT)
	if err != nil {
		log.Print(err)
		return nil, exception.CreateError(exception.CodeNotFound, "Cookie not found.")
	}
	claims, cErr := s.Bean.SecretService.CheckRestaurantToken(cookie)
	if cErr != nil {
		log.Print(cErr.ErrorMessage)
		return nil, exception.CreateError(exception.CodeUnknown, "client has error")
	}

	return claims.Restaurant, nil
}
