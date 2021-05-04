package entity

type Restaurant struct {
	ID                int64
	RestaurantGroupID int64
	Name              string
	Description       string
	Image             string
	Address           string
	Phone             string
}

func CreateRestaurantFromSlice(restaurantRaw interface{}) *Restaurant {
	result := &Restaurant{
		ID:                int64(restaurantRaw.(map[string]interface{})["id"].(float64)),
		RestaurantGroupID: int64(restaurantRaw.(map[string]interface{})["restaurantGroupID"].(float64)),
		Name:              restaurantRaw.(map[string]interface{})["name"].(string),
	}
	if description, ok := restaurantRaw.(map[string]interface{})["description"].(string); ok {
		result.Description = description
	}

	if image, ok := restaurantRaw.(map[string]interface{})["image"].(string); ok {
		result.Image = image
	}

	if address, ok := restaurantRaw.(map[string]interface{})["address"].(string); ok {
		result.Address = address
	}

	if phone, ok := restaurantRaw.(map[string]interface{})["phone"].(string); ok {
		result.Phone = phone
	}
	return result
}
