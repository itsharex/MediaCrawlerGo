package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/NanmiCoder/MediaCrawlerGo/util"
)

const (
	GENERAL SearchSortType = "general"
	POPULAR SearchSortType = "popularity_descending"
	LATEST  SearchSortType = "time_descending"
)

const (
	ALL SearchNoteType = iota
	VIDEO
	IMAGE
)

// XhsHttpClient Packaged httpclient based on xhs
type XhsHttpClient struct {
	client         *http.Client
	headers        map[string]interface{}
	timeout        int
	playwrightPage playwright.Page
	cookiesMap     map[string]string
	baseUrl        string
}
type SearchXhsNoteParams struct {
	Keyword  string         `json:"keyword"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Sort     SearchSortType `json:"sort"`
	NoteType SearchNoteType `json:"note_type"`
	SearchId string         `json:"search_id"`
}
type SearchSortType string
type SearchNoteType int

type NoteSearchResult struct {
	HasMore bool   `json:"has_more"`
	Items   []Note `json:"items"`
}

type Note struct {
	// Define your Note struct fields here
}

func (c *XhsHttpClient) PreHeaders(uri string, body []byte) map[string]interface{} {
	var fullHeaders map[string]interface{}
	var encryptData map[string]interface{}

	_ = json.Unmarshal(body, &encryptData)
	encryptParams, err := c.playwrightPage.Evaluate("(url, data) => window._webmsxyw(url,data)", uri, encryptData)
	if err != nil {
		util.Log().Panic("window._webmsxyw(url,data) err:%v", err)

	}
	if encryptParamsMap, ok := encryptParams.(map[string]interface{}); ok {
		headers := make(map[string]interface{})
		headers["X-s"] = encryptParamsMap["X-s"]
		headers["X-t"] = encryptParamsMap["X-t"]
		fullHeaders = util.MergeMap(headers, c.headers)
	} else {
		util.Log().Panic("encryptParams convert failed")
	}
	return fullHeaders
}

func (c *XhsHttpClient) Get(uri string, params map[string]string) (*http.Response, error) {
	url := c.baseUrl + uri
	if params != nil {
		paramSlice := make([]string, len(params))
		for k, v := range params {
			paramSlice = append(paramSlice, fmt.Sprintf("%s=%v", k, v))
		}
		url += "?" + strings.Join(paramSlice, "&")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add pre headers method before send request
	for key, value := range c.PreHeaders(uri, nil) {
		if stringValue, ok := value.(string); ok {
			req.Header.Set(key, stringValue)
		} else {
			continue
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *XhsHttpClient) Post(uri string, body []byte) (*http.Response, error) {
	util.Log().Info("[XhsHttpClient.Post] Begin execute post request, uri: %s, body: %v", uri, string(body))
	req, err := http.NewRequest("POST", c.baseUrl+uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// add pre headers method before send request
	for key, value := range c.PreHeaders(uri, body) {
		if stringValue, ok := value.(string); ok {
			req.Header.Set(key, stringValue)
		} else if intValue, ok := value.(int); ok {
			req.Header.Set(key, strconv.Itoa(intValue))
		} else {
			continue
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// XhsApiClient Packaged api client based on xhs
type XhsApiClient struct {
	httpClient *XhsHttpClient
}

// Ping xhs api pong before execute search note ...
func (api *XhsApiClient) Ping() bool {
	pong := true
	util.Log().Info("[XhsApiClient.Ping] Begin to ping xhs ...")
	_, err := api.GetNoteByKeyword(SearchXhsNoteParams{
		Keyword:  "小红书",
		Page:     1,
		PageSize: 20,
		Sort:     GENERAL,
		NoteType: ALL,
		SearchId: getSearchId(),
	})
	if err != nil {
		pong = false
		util.Log().Info("[XhsApiClient.Ping] Xhs ping failed and login again ...")
	}
	return pong
}

// UpdateCookies will call this method when invalid login
func (api *XhsApiClient) UpdateCookies(ctx playwright.BrowserContext) {
	cookies, err := ctx.Cookies()
	util.AssertErrorToNil("[XhsApiClient.UpdateCookies] could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("[XhsApiClient.UpdateCookies] convert cookie failed and error:", err)
	api.httpClient.headers["Cookie"] = convertResp.cookieStr
	api.httpClient.cookiesMap = convertResp.cookiesMap
}

// GetNoteByKeyword get note list by search keywords
func (api *XhsApiClient) GetNoteByKeyword(searchParams SearchXhsNoteParams) (NoteSearchResult, error) {
	uri := "/api/sns/web/v1/search/notes"
	searchBob, err := json.Marshal(searchParams)
	resp, err := api.httpClient.Post(uri, searchBob)
	if err != nil {
		return NoteSearchResult{}, err
	}

	if resp.StatusCode != 200 {
		var errorResp any
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)
		util.Log().Info("[XhsHttpClient.GetNoteByKeyword] errorResp: %v", errorResp)
		return NoteSearchResult{}, errors.New("[XhsHttpClient.GetNoteByKeyword] got note failed")
	}

	var result NoteSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return NoteSearchResult{}, err
	}

	return result, nil
}
