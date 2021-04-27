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

func NewAPIRepository(client *client.APIClient, loginClient *commonClient.LoginClient) APIRepository {
	return APIRepository{Client: client, LoginClient: loginClient}
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

func (r *APIRepository) GetRestaurants() ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetRestaurants()
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	result, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return result, nil

}

func (r *APIRepository) GetOrdersByAreaID(areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetOrdersByAreaID(areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetOrdersByOrderNumberID(orderNumberID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetOrdersByOrderNumberID(orderNumberID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	orders, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return orders, nil

}

func (r *APIRepository) GetTablesByAreaID(areaID int64) ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetTablesByAreaID(areaID)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) GetCategories() ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetCategories()
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

func (r *APIRepository) GetComments() ([]interface{}, *exception.Error) {

	resp, err := r.Client.GetComments()
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponses(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}

func (r *APIRepository) PutOrders(v interface{}) (interface{}, *exception.Error) {
	resp, err := r.Client.PutOrders(v)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	id, err := r.GetResponse(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return id, nil

}

func (r *APIRepository) PostRestaurant(v interface{}) (interface{}, *exception.Error) {
	resp, err := r.LoginClient.PostRestaurant(v)
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
			return nil, fmt.Errorf("Not found.")
		}
		if errResponse.ErrorCode == exception.CodeValueInvalid {
			return nil, fmt.Errorf("request value invalid.")
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
			return nil, fmt.Errorf("Not found.")
		}
		if errResponse.ErrorCode == exception.CodeValueInvalid {
			return nil, fmt.Errorf("request value invalid.")
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
