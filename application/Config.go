package application

import (
	"strconv"

	"github.com/gin-contrib/multitemplate"
	"github.com/smarest/smarest-common/application"
	"github.com/smarest/smarest-common/util"
	"github.com/smarest/smarest-portal/application/resource"
	"github.com/smarest/smarest-portal/infrastructure/client"
	"github.com/smarest/smarest-portal/infrastructure/persistence"

	cClient "github.com/smarest/smarest-common/client"
	cRepo "github.com/smarest/smarest-common/infrastructure/persistence"
)

const (
	PageCashier    = "cashier"
	PageOrder      = "order"
	PageCheckout   = "checkout"
	PageRestaurant = "restaurant"
	CommonTitle    = "title"
	CommonHeader   = "header"
	CommonFooter   = "footer"
	TemplatePath   = "./templates"
)

type Bean struct {
	PageLayouts       map[string]string
	CommonLayouts     map[string]string
	CashierService    *CashierService
	OrderService      *OrderService
	CheckoutService   *CheckoutService
	RestaurantService *RestaurantService
	PortalAPIService  *PortalAPIService
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

	apiClient := client.NewAPIClient(
		util.GetEnvDefault("SMAREST_API_HOST", "http://localhost:8080"),
		apiTimeout,
	)

	loginClient := cClient.NewLoginClient(
		util.GetEnvDefault("POS_USER_HOST", "http://localhost:8080"),
		userTimeout,
	)
	loginRepository := cRepo.NewLoginRepository(loginClient)
	apiRepository := persistence.NewAPIRepository(apiClient, loginClient)

	pageResourceFactory := resource.NewPageResourceFactory(util.GetEnvDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		util.GetEnvDefault("POS_DESIGN_URL", "http://localhost/pos/pos-lib"),
		util.GetEnvDefault("POS_IMAGE_URL", "http://localhost"))

	loginService := application.NewLoginService(util.GetEnvDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		util.GetEnvDefault("POS_LOGIN_TOKEN", "pos_access_token"),
		loginRepository)

	cashierService := NewCashierService(loginService, apiRepository, pageResourceFactory)
	orderService := NewOrderService(loginService, apiRepository, pageResourceFactory)
	checkoutService := NewCheckoutService(loginService, apiRepository, pageResourceFactory)
	restaurantService := NewRestaurantService(apiRepository, pageResourceFactory)
	portalAPIService := NewPortalAPIService(loginService, apiRepository)
	// html layout
	pageLayouts := make(map[string]string)
	pageLayouts[PageCashier] = TemplatePath + "/page_cashier.html"
	pageLayouts[PageOrder] = TemplatePath + "/page_order.html"
	pageLayouts[PageCheckout] = TemplatePath + "/page_checkout.html"
	pageLayouts[PageRestaurant] = TemplatePath + "/page_restaurant.html"

	//common
	commonLayouts := make(map[string]string)
	commonLayouts[CommonHeader] = TemplatePath + "/commons/header.html"
	commonLayouts[CommonFooter] = TemplatePath + "/commons/footer.html"
	commonLayouts[CommonTitle] = TemplatePath + "/commons/title.html"

	return &Bean{
		PageLayouts:       pageLayouts,
		CommonLayouts:     commonLayouts,
		CashierService:    cashierService,
		CheckoutService:   checkoutService,
		RestaurantService: restaurantService,
		PortalAPIService:  portalAPIService,
		OrderService:      orderService}, nil
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

	//restaurantPage
	render.AddFromFiles(PageRestaurant, bean.PageLayouts[PageRestaurant])
	return render
}
