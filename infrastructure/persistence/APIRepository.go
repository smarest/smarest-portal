package persistence

import (
	"encoding/json"
	"io"
	"net/http"

	commonClient "github.com/smarest/smarest-common/client"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-portal/infrastructure/client"
)

type APIRepository struct {
	Client      *client.APIClient
	LoginClient *commonClient.LoginClient
}

func NewAPIRepository(client *client.APIClient, loginClient *commonClient.LoginClient) *APIRepository {
	return &APIRepository{Client: client, LoginClient: loginClient}
}

func (r *APIRepository) GetAreasByRestaurantID(restID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetAreasByRestaurantID(restID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) GetRestaurantsByGroupID(groupID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantsByGroupID(groupID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) GetRestaurantByCode(code string) (interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantByCode(code)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponse(resp)
}

func (r *APIRepository) GetRestaurantByID(id int64) (interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantByID(id)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponse(resp)
}

func (r *APIRepository) GetRestaurantOrdersByAreaID(restID int64, areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantOrdersByAreaID(restID, areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) GetRestaurantOrders(restID int64, orderBy string, groupBy string) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantOrders(restID, orderBy, groupBy)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)

}

func (r *APIRepository) GetRestaurantOrdersByOrderNumberID(restID int64, orderNumberID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantOrdersByOrderNumberID(restID, orderNumberID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)

}

func (r *APIRepository) GetRestaurantTablesByAreaID(restID int64, areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantTablesByAreaID(restID, areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)

}

func (r *APIRepository) GetCategoriesByGroupID(groupID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetCategoriesByGroupID(groupID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) GetProductsByRestaurantIDAndCategoryID(areaID int64, categoryID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetProductsByRestaurantIDAndCategoryID(areaID, categoryID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) GetRestaurantCommentsByProductID(restID int64, productID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantCommentsByProductID(restID, productID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return r.GetResponses(resp)
}

func (r *APIRepository) PutRestaurantOrders(restID int64, v interface{}) (interface{}, *exception.Error) {
	resp, err := r.Client.PutRestaurantOrders(restID, v)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}
	return r.GetResponse(resp)

}

func (s *APIRepository) GetResponses(resp *http.Response) ([]interface{}, *exception.Error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.GetErrorResponse(resp.Body)
	}

	var items []interface{}
	err := json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, exception.CreateErrorWithRootCause(exception.CodeUnknown, "can not decode responses.", err)
	}
	return items, nil
}

func (s *APIRepository) GetResponse(resp *http.Response) (interface{}, *exception.Error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.GetErrorResponse(resp.Body)
	}

	var items interface{}
	err := json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, exception.CreateErrorWithRootCause(exception.CodeUnknown, "can not decode response.", err)
	}
	return items, nil
}

func (s *APIRepository) GetErrorResponse(body io.ReadCloser) *exception.Error {
	var errResponse exception.Error
	err := json.NewDecoder(body).Decode(&errResponse)
	if err != nil {
		return exception.CreateErrorWithRootCause(exception.CodeUnknown, "Can not decode Error Message", err)
	}
	return &errResponse
}
