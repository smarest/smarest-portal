package client

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/smarest/smarest-common/client"
)

type APIClient struct {
	LoginEntryPoint           string
	RestaurantGroupEntryPoint string
	RestaurantEntryPoint      string
	*client.BaseClient
}

func NewAPIClient(host string, timeout int) *APIClient {
	return &APIClient{"/v1/login", "/v1/portal/groups", "/v1/portal/restaurants", client.NewBaseClient(host, timeout)}
}

func (c *APIClient) GetAreasByRestaurantID(restID int64) (*http.Response, error) {
	url := c.Host + c.RestaurantEntryPoint + "/" + strconv.FormatInt(restID, 10) + "/areas?fields=name,id"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantsByGroupID(groupID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/restaurants?fields=id,restaurantGroupID,code,name,description,image,address,phone", c.Host, c.RestaurantGroupEntryPoint, groupID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantByCode(code string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s?code=%s&fields=id,restaurantGroupID,code,name,description,image,address,phone", c.Host, c.RestaurantEntryPoint, code)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantByID(id int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d?fields=id,restaurantGroupID,name,description,image,address,phone", c.Host, c.RestaurantEntryPoint, id)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantOrdersByAreaID(restID int64, areaID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/areas/%d/orders", c.Host, c.RestaurantEntryPoint, restID, areaID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantOrdersByOrderNumberID(restID int64, orderNumberID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/orders?orderNumberID=%d", c.Host, c.RestaurantEntryPoint, restID, orderNumberID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantTablesByAreaID(restID int64, areaID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/areas/%d/tables", c.Host, c.RestaurantEntryPoint, restID, areaID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetCategoriesByGroupID(groupID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/categories?fields=id,name", c.Host, c.RestaurantGroupEntryPoint, groupID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetProductsByRestaurantIDAndCategoryID(restID int64, cateID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/categories/%d/products?fields=id,name,price,quantityOnSingleOrder", c.Host, c.RestaurantEntryPoint, restID, cateID)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurantCommentsByProductID(restID int64, productID int64) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/products/%d/comments?fields=id,name", c.Host, c.RestaurantGroupEntryPoint, restID, productID)
	return c.DoGetRequest(url)
}

func (c *APIClient) PutRestaurantOrders(restID int64, request interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%d/orders", c.Host, c.RestaurantEntryPoint, restID)
	return c.DoPutRequest(url, request)
}
