package persistence

import (
	"encoding/json"
	"fmt"
	"log"
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

	results, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return results, nil

}

func (r *APIRepository) GetRestaurantsByGroupID(groupID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantsByGroupID(groupID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	result, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return result, nil
}

func (r *APIRepository) GetRestaurantByCode(code string) (interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantByCode(code)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	result, err := r.GetResponse(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return result, nil

}

func (r *APIRepository) GetRestaurantByID(id int64) (interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantByID(id)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	result, err := r.GetResponse(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return result, nil

}

func (r *APIRepository) GetRestaurantOrdersByAreaID(restID int64, areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantOrdersByAreaID(restID, areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetRestaurantOrdersByOrderNumberID(restID int64, orderNumberID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantOrdersByOrderNumberID(restID, orderNumberID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	orders, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return orders, nil

}

func (r *APIRepository) GetRestaurantTablesByAreaID(restID int64, areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantTablesByAreaID(restID, areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetCategoriesByGroupID(groupID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetCategoriesByGroupID(groupID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetProductsByRestaurantIDAndCategoryID(areaID int64, categoryID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetProductsByRestaurantIDAndCategoryID(areaID, categoryID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetRestaurantCommentsByProductID(restID int64, productID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurantCommentsByProductID(restID, productID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) PutRestaurantOrders(restID int64, v interface{}) (interface{}, *exception.Error) {
	resp, err := r.Client.PutRestaurantOrders(restID, v)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	id, err := r.GetResponse(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return id, nil

}

func (s *APIRepository) GetResponses(resp *http.Response) ([]interface{}, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse exception.Error
		err := json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			return nil, fmt.Errorf("can not decode login response. [%w]", err)
		}
		if errResponse.ErrorCode == exception.CodeNotFound {
			return nil, fmt.Errorf("not found")
		}
		if errResponse.ErrorCode == exception.CodeValueInvalid {
			return nil, fmt.Errorf("request value invalid")
		}
		return nil, fmt.Errorf("api has error. [%s]", errResponse.ErrorMessage)
	}

	var items []interface{}
	err := json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		fmt.Printf("%v", resp.Body)
		return nil, fmt.Errorf("can not decode response. [%w]", err)
	}
	return items, nil
}

func (s *APIRepository) GetResponse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse exception.Error
		err := json.NewDecoder(resp.Body).Decode(&errResponse)

		log.Print(errResponse.ErrorMessage)
		if err != nil {
			return nil, fmt.Errorf("can not decode login response. [%w]", err)
		}
		if errResponse.ErrorCode == exception.CodeNotFound {
			return nil, fmt.Errorf("not found")
		}
		if errResponse.ErrorCode == exception.CodeValueInvalid {
			return nil, fmt.Errorf("request value invalid")
		}
		return nil, fmt.Errorf("api has error. [%s]", errResponse.ErrorMessage)
	}

	var items interface{}
	err := json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("can not decode response. [%w]", err)
	}
	return items, nil
}
