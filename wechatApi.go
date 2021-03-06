package wechatApi

import (
	"errors"
	"fmt"

	"github.com/duck-lab/wechat-api/accessToken"
	"github.com/duck-lab/wechat-api/menu"

	"time"
)

//API is the outlet of all APIs
type API struct {
	endpoint         string
	baseURL          string
	appID            string
	appSecret        string
	tokenStorageMode int
	currentToken     accessToken.Data
	calls            map[string]int
}

var api = API{
	endpoint:         "https://api.weixin.qq.com",
	baseURL:          "https://api.weixin.qq.com/cgi-bin",
	tokenStorageMode: accessToken.StoreInMemory,
	calls:            map[string]int{},
}

//New is the start point
//For the convience of get the accure calls, use single instance
func New(appID string, appSecret string) API {
	api.appID = appID
	api.appSecret = appSecret
	return api
}

//SetEndpoint let the user can chose which api server to use
func (api *API) SetEndpoint(endpoint string) {
	api.endpoint = endpoint
	api.baseURL = api.endpoint + "/cgi-bin"
}

//SetTokenStorageMode ...
func (api *API) SetTokenStorageMode(mode int, params interface{}) bool {
	if mode == accessToken.StoreInMemory {

	}
	if mode == accessToken.StoreInFile {
		file, ok := params.(accessToken.StoreInFileParams)
		fmt.Println(file, ok)
	}
	if mode == accessToken.StoreInRedis {
		redis, ok := params.(accessToken.StoreInRedisParams)
		fmt.Println(redis, ok)
	}
	return false
}

func (api *API) getCurrentToken() (accessToken.Data, error) {
	if api.tokenStorageMode == accessToken.StoreInMemory {
		return api.currentToken, nil
	}
	if api.tokenStorageMode == accessToken.StoreInFile {
		//TODO: read file
	}

	if api.tokenStorageMode == accessToken.StoreInRedis {
		//TODO: read redis
	}
	return accessToken.Data{}, errors.New("To be implemented")
}

func (api *API) setCurrentToken(token accessToken.Data) bool {
	if api.tokenStorageMode == accessToken.StoreInMemory {
		api.currentToken = token
		return true
	}
	return false
}

func getCallKey(apiName string) string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d-%s", now.Year(), now.Month(), now.Day(), apiName)
}

func (api *API) getCall(apiName string) int {
	return api.calls[getCallKey(apiName)]
}

func (api *API) callIncr(apiName string) int {
	//Todo: support multi routines
	key := getCallKey(apiName)
	api.calls[key] = api.calls[key] + 1
	fmt.Println(api.calls)
	return api.calls[key]
}

//GetAccessToken can let user fetch the access token
func (api *API) GetAccessToken() (accessToken.Data, error) {
	current, err := api.getCurrentToken()
	if err != nil {
		return accessToken.Data{}, err
	}
	if current.ExpiresTime.Unix() > time.Now().Unix() {
		return current, nil
	}
	if api.tokenStorageMode == accessToken.StoreInMemory {
		data, err := accessToken.Fetch(api.appID, api.appSecret, api.baseURL)
		api.callIncr(accessToken.APINameFetch)
		api.setCurrentToken(data)
		return data, err
	}

	return accessToken.Data{}, errors.New("To be implemented")
}

//CreateMenu is used to create Menu
func (api *API) CreateMenu(model menu.Model) error {
	token, err := api.GetAccessToken()
	if err != nil {
		return err
	}
	api.callIncr(menu.APINameCreate)
	return menu.Create(model, api.baseURL, token.AccessToken)
}

//CreateConditionalMenu is used to create customized Menu
func (api *API) CreateConditionalMenu(model menu.ConditionalMenu) (string, error) {
	token, err := api.GetAccessToken()
	if err != nil {
		return "", err
	}
	api.callIncr(menu.APINameCreateConditional)
	return menu.CreateConditional(model, api.baseURL, token.AccessToken)
}

//FetchMenuList is used to fetch all Menu
func (api *API) FetchMenuList() (menu.FetchedList, error) {
	token, err := api.GetAccessToken()
	if err != nil {
		return menu.FetchedList{}, err
	}
	api.callIncr(menu.APINameFetchedList)
	return menu.FetchList(api.baseURL, token.AccessToken)
}

//DeleteAllMenu ...
func (api *API) DeleteAllMenu() error {
	token, err := api.GetAccessToken()
	if err != nil {
		return err
	}
	api.callIncr(menu.APINameDeleteAll)
	return menu.DeleteAll(api.baseURL, token.AccessToken)
}
