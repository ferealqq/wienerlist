package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type ApiService struct {
	BaseUrl string
	Token   string
}

type JsonResponse struct {
	data map[string]interface{}
	err  error
}

func (j JsonResponse) BindModel(m interface{}) error {
	if j.data != nil && j.err == nil {
		return mapstructure.Decode(j.data, &m)
	}
	return j.err
}

func (j JsonResponse) RawData() (map[string]interface{}, error) {
	if j.data != nil {
		return j.data, nil
	}
	return nil, j.err
}

func (a *ApiService) Get(path string) JsonResponse {
	return get(a.BaseUrl + path)
}

type ApiGetRequest struct {
	*ApiService

	Params map[string]interface{}
}

// FIXME figure out how to comment properly in golang
// name => parameter name, value => paramter value, should be used only ones
func (a *ApiService) Params(name string, value ...interface{}) *ApiGetRequest {
	params := make(map[string]interface{})
	params[name] = value
	return &ApiGetRequest{
		a,
		params,
	}
}

func (a *ApiGetRequest) AddParam(name string, value ...interface{}) *ApiGetRequest {
	a.Params[name] = value
	return a
}

func (a *ApiGetRequest) Get(path string) JsonResponse {
	fullpath := a.BaseUrl + path
	i := 0
	for key, val := range a.Params {
		var prm string
		if i == 0 {
			prm = "?" + key + "="
		} else {
			prm = "&" + key + "="
		}

		switch t := val.(type) {
		case []int:
			for j, v := range t {
				if j == 0 {
					prm = prm + strconv.Itoa(v)
				} else {
					prm = prm + "&" + key + "=" + strconv.Itoa(v)
				}
			}
		case []string:
			for j, v := range t {
				if j == 0 {
					prm = prm + v
				} else {
					prm = prm + "&" + key + "=" + v
				}
			}
		case string:
			prm = prm + t
		case int:
			prm = prm + strconv.Itoa(t)
		case []interface{}:
			for j, v := range t {
				switch d := v.(type) {
				case string:
					if j == 0 {
						prm = prm + d
					} else {
						prm = prm + "&" + key + "=" + d
					}
				case int:
					if j == 0 {
						prm = prm + strconv.Itoa(d)
					} else {
						prm = prm + "&" + key + "=" + strconv.Itoa(d)
					}
				default:
					return JsonResponse{nil, errors.New("param values should only be type of int | int[] | string | string []")}
				}
			}
		default:
			return JsonResponse{nil, errors.New("param values should only be type of int | int[] | string | string []")}
		}

		fullpath = fullpath + prm
		i++
	}

	return get(fullpath)
}

func get(fullpath string) JsonResponse {
	resp, err := http.Get(fullpath)
	if err == nil {
		defer resp.Body.Close()
		defer println("body closed")
		var data map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err == nil {
			return JsonResponse{
				data,
				nil,
			}
		}
		return JsonResponse{
			nil,
			err,
		}
	}
	return JsonResponse{
		nil,
		err,
	}
}

func NewApi(u string) *ApiService {
	return &ApiService{
		BaseUrl: u,
	}
}
