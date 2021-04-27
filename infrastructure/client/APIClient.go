package client

import (
	"net/http"
	"strconv"

	"github.com/smarest/smarest-common/client"
)

type APIClient struct {
	LoginEntryPoint      string
	RestaurantEntryPoint string
	AreaEntryPoint       string
	CategoryEntryPoint   string
	CommentEntryPoint    string
	OrderEntryPoint      string
	*client.POSClient
}

func NewAPIClient(host string, timeout int) *APIClient {
	return &APIClient{"/v1/login", "/v1/restaurants", "/v1/areas", "/v1/categories", "/v1/comments", "/v1/orders", client.NewPOSClient(host, timeout)}
}

func (c *APIClient) GetAreasByRestaurantID(restID int64) (*http.Response, error) {
	url := c.Host + c.RestaurantEntryPoint + "/" + strconv.FormatInt(restID, 10) + "/areas?fields=name,id"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetRestaurants() (*http.Response, error) {
	url := c.Host + c.RestaurantEntryPoint + "?fields=name,id"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetOrdersByAreaID(areaID int64) (*http.Response, error) {
	url := c.Host + c.AreaEntryPoint + "/" + strconv.FormatInt(areaID, 10) + "/orders"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetOrdersByOrderNumberID(orderNumberID int64) (*http.Response, error) {
	url := c.Host + c.OrderEntryPoint + "?orderNumberID=" + strconv.FormatInt(orderNumberID, 10)
	return c.DoGetRequest(url)
}

func (c *APIClient) GetTablesByAreaID(areaID int64) (*http.Response, error) {
	url := c.Host + c.AreaEntryPoint + "/" + strconv.FormatInt(areaID, 10) + "/tables"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetCategories() (*http.Response, error) {
	url := c.Host + c.CategoryEntryPoint + "?type=P&available=true&fields=id,name"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetProductsByRestaurantIDAndCategoryID(restID int64, cateID int64) (*http.Response, error) {
	url := c.Host + c.RestaurantEntryPoint + "/" + strconv.FormatInt(restID, 10) + "/products?categoryID=" + strconv.FormatInt(cateID, 10) + "&fields=id,name,price,quantityOnSingleOrder"
	return c.DoGetRequest(url)
}

func (c *APIClient) GetComments() (*http.Response, error) {
	url := c.Host + c.CommentEntryPoint + "?fields=id,name"
	return c.DoGetRequest(url)
}

func (c *APIClient) PutOrders(request interface{}) (*http.Response, error) {
	url := c.Host + c.OrderEntryPoint
	return c.DoPutRequest(url, request)
}
