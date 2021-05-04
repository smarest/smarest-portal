package persistence

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/util"
)

// URLRepository to get user info from cookie
type URLRepository struct {
	LoginURL string
}

// NewLoginService create LoginService
func NewURLRepository(loginURL string) *URLRepository {
	return &URLRepository{loginURL}
}

func (s *URLRepository) GetLoginUrl(c *gin.Context) string {
	return fmt.Sprintf("%s?frm=http://%s%s", s.LoginURL, c.Request.Host, c.Request.URL)
}

func (s *URLRepository) GetPageRedirectURL(from string, to string) string {
	return fmt.Sprintf("/portal/%s?%s=%s", to, util.PARAM_FROM_URL, from)
}

func (s *URLRepository) GetErrorURL() string {
	return "/portal/error"
}
