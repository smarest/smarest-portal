package application

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type ErrorService struct {
	Bean *Bean
}

func NewErrorService(bean *Bean) *ErrorService {
	return &ErrorService{bean}
}

func (s *ErrorService) Get(c *gin.Context) {
	c.HTML(http.StatusOK, PageError, gin.H{
		"resource": s.Bean.PageResourceFactory.CreateResource(),
	})
}

func (s *ErrorService) HandlerError(c *gin.Context, err *exception.Error) {
	log.Print(err.ErrorMessage)
	c.Redirect(http.StatusMovedPermanently, s.Bean.URLRepository.GetErrorURL())
}
