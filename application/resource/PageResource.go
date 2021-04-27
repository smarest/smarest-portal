package resource

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type PageResource struct {
	ErrorMessage  string
	Message       string
	IsCashier     bool
	IsOrder       bool
	OrderNumberID string
	Restaurants   []interface{}
	Areas         []interface{}
	Tables        []interface{}
	Orders        []interface{}
	Categories    []interface{}
	Products      []interface{}
	Comments      []interface{}
	AreaID        string
	CategoryID    string
	LoginUrl      string
	DesignUrl     string
	ImageUrl      string
	PageTitle     string
	FromURL       string
}

func (s *PageResource) IsEqual(value1 string, value2 interface{}) bool {
	return value1 == fmt.Sprintf("%v", value2)
}

func (s *PageResource) SumOrdersPrice() int64 {
	var result = int64(0)
	for _, item := range s.Orders {
		result += int64(item.(map[string]interface{})["price"].(float64))
	}
	return result
}

func (s *PageResource) SumOrdersCount() int64 {
	var result = int64(0)
	for _, item := range s.Orders {
		result += int64(item.(map[string]interface{})["count"].(float64))
	}
	return result
}

func (s *PageResource) Int64(value interface{}) int64 {
	switch value := value.(type) {
	case string:
		result, err := strconv.ParseInt(fmt.Sprintf("%v", value), 0, 64)
		if err == nil {
			return result
		}
	case float64:
		return int64(value)
	}
	return 0

}

func (s *PageResource) Comma(v int64) string {
	sign := ""

	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,37Ò2,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}
