package application

import (
	"strconv"

	"github.com/gin-contrib/multitemplate"
	common_application "github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-common/util"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/domain/service"
	"github.com/smarest/smarest-portal/infrastructure/client"
	"github.com/smarest/smarest-portal/infrastructure/persistence"

	common_client "github.com/smarest/smarest-common/client"
	common_repo "github.com/smarest/smarest-common/infrastructure/persistence"
)

const (
	PageCashier    = "cashier"
	PageKitchen    = "kitchen"
	PageOrder      = "order"
	PageCheckout   = "checkout"
	PageRestaurant = "restaurant"
	PageError      = "error"
	CommonTitle    = "title"
	CommonHeader   = "header"
	CommonFooter   = "footer"
	TemplatePath   = "./templates"
)

type Bean struct {
	COOKIE_TOKEN_RESTAURANT string
	PageLayouts             map[string]string
	CommonLayouts           map[string]string
	CookieCheckService      *CookieCheckService
	CashierService          *CashierService
	OrderService            *OrderService
	KitchenService          *KitchenService
	CheckoutService         *CheckoutService
	RestaurantService       *RestaurantService
	PortalAPIService        *PortalAPIService
	ErrorService            *ErrorService
	LoginService            *common_application.LoginService
	SecretService           *service.SecretService
	URLRepository           *persistence.URLRepository
	APIRepository           *persistence.APIRepository
	LoginRepository         *common_repo.LoginRepository
	APIClient               *client.APIClient
	LoginClient             *common_client.LoginClient
	PageResourceFactory     *resource.PageResourceFactory
}

func InitBean() (*Bean, error) {
	userTimeout, err := strconv.Atoi(util.GetEnvDefault("SMAREST_ACCOUNT_TIMEOUT", "5000"))
	if err != nil {
		return nil, err
	}

	apiTimeout, err := strconv.Atoi(util.GetEnvDefault("SMAREST_API_TIMEOUT", "5000"))
	if err != nil {
		return nil, err
	}

	bean := &Bean{}
	bean.COOKIE_TOKEN_RESTAURANT = util.GetEnvDefault("COOKIE_TOKEN_RESTAURANT", "res_token")
	bean.APIClient = client.NewAPIClient(
		util.GetEnvDefault("SMAREST_API_HOST", "http://localhost:8081"),
		apiTimeout,
	)

	bean.LoginClient = common_client.NewLoginClient(
		util.GetEnvDefault("POS_USER_HOST", "http://localhost:8080"),
		userTimeout,
	)
	bean.LoginRepository = common_repo.NewLoginRepository(bean.LoginClient)
	bean.APIRepository = persistence.NewAPIRepository(bean.APIClient, bean.LoginClient)

	bean.PageResourceFactory = resource.NewPageResourceFactory(util.GetEnvDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		util.GetEnvDefault("POS_DESIGN_URL", "http://localhost/pos/smarest-design"),
		util.GetEnvDefault("POS_IMAGE_URL", "http://localhost/pos/smarest-design"))

	bean.LoginService = common_application.NewLoginService(util.GetEnvDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		util.GetEnvDefault("POS_LOGIN_TOKEN", "pos_access_token"),
		bean.LoginRepository)

	bean.SecretService = service.NewSecretService(util.GetEnvDefault("SMAREST_RESTAURANT_TOKEN", "b3BlbnNzaC1rZXktdjEAAAAACm"))

	bean.CookieCheckService = NewCookieCheckService(bean)
	bean.RestaurantService = NewRestaurantService(bean)
	bean.CashierService = NewCashierService(bean)
	bean.KitchenService = NewKitchenService(bean)
	bean.OrderService = NewOrderService(bean)
	bean.CheckoutService = NewCheckoutService(bean)
	bean.PortalAPIService = NewPortalAPIService(bean)
	bean.ErrorService = NewErrorService(bean)
	bean.URLRepository = persistence.NewURLRepository(util.GetEnvDefault("SMAREST_LOGIN_URL", "http://localhost:8080/login"))
	// html layout
	bean.PageLayouts = make(map[string]string)
	bean.PageLayouts[PageCashier] = TemplatePath + "/page_cashier.html"
	bean.PageLayouts[PageKitchen] = TemplatePath + "/page_kitchen.html"
	bean.PageLayouts[PageOrder] = TemplatePath + "/page_order.html"
	bean.PageLayouts[PageCheckout] = TemplatePath + "/page_checkout.html"
	bean.PageLayouts[PageRestaurant] = TemplatePath + "/page_restaurant.html"
	bean.PageLayouts[PageError] = TemplatePath + "/page_error.html"

	//common
	bean.CommonLayouts = make(map[string]string)
	bean.CommonLayouts[CommonHeader] = TemplatePath + "/commons/header.html"
	bean.CommonLayouts[CommonFooter] = TemplatePath + "/commons/footer.html"
	bean.CommonLayouts[CommonTitle] = TemplatePath + "/commons/title.html"

	return bean, nil
}

func (bean *Bean) LoadTemplates() multitemplate.Renderer {
	render := multitemplate.NewRenderer()
	//cashierPage
	render.AddFromFiles(PageCashier,
		bean.PageLayouts[PageCashier],
		bean.CommonLayouts[CommonTitle],
		bean.CommonLayouts[CommonHeader],
		bean.CommonLayouts[CommonFooter])

	//orderPage
	render.AddFromFiles(PageOrder,
		bean.PageLayouts[PageOrder],
		bean.CommonLayouts[CommonTitle],
		bean.CommonLayouts[CommonHeader],
		bean.CommonLayouts[CommonFooter])

	//checkoutPage
	render.AddFromFiles(PageCheckout,
		bean.PageLayouts[PageCheckout],
		bean.CommonLayouts[CommonTitle],
		bean.CommonLayouts[CommonHeader],
		bean.CommonLayouts[CommonFooter])

	//kitchenPage
	render.AddFromFiles(PageKitchen,
		bean.PageLayouts[PageKitchen],
		bean.CommonLayouts[CommonTitle],
		bean.CommonLayouts[CommonHeader],
		bean.CommonLayouts[CommonFooter])

	//restaurantPage
	render.AddFromFiles(PageRestaurant, bean.PageLayouts[PageRestaurant])

	//errorPage
	render.AddFromFiles(PageError,
		bean.PageLayouts[PageError],
		bean.CommonLayouts[CommonTitle],
		bean.CommonLayouts[CommonHeader])
	return render
}
