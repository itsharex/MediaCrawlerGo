package platform

import (
	"MediaCrawlerGo/util"
	"bytes"
	"encoding/json"
	"github.com/playwright-community/playwright-go"
	"net/http"
)

// XhsHttpClient Packaged httpclient based on xhs
type XhsHttpClient struct {
	client         *http.Client
	headers        map[string]string
	timeout        int
	playwrightPage *playwright.Page
	cookiesMap     map[string]string
}

func (c *XhsHttpClient) PreHeaders(params ...any) map[string]string {
	// params []any
	headers := make(map[string]string)
	headers["X-s"] = ""
	headers["X-t"] = ""
	headers["X-s-common"] = ""
	headers["X-B3-Traceid"] = ""
	return headers
}

func (c *XhsHttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add pre headers method before send request
	for key, value := range c.PreHeaders(url, headers) {
		req.Header.Set(key, value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *XhsHttpClient) Post(url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// add pre headers method before send request
	for key, value := range c.PreHeaders(url, headers) {
		req.Header.Set(key, value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type SearchSortType string

const (
	GENERAL SearchSortType = "GENERAL" // replace with actual values
)

type SearchNoteType string

const (
	ALL SearchNoteType = "ALL" // replace with actual values
)

type NoteSearchResult struct {
	HasMore bool   `json:"has_more"`
	Items   []Note `json:"items"`
}

type Note struct {
	// Define your Note struct fields here
}

// XhsApiClient Packaged api client based on xhs
type XhsApiClient struct {
	httpClient *XhsHttpClient
}

// Ping xhs api pong before execute search note ...
func (api *XhsApiClient) Ping() bool {
	pong := true
	util.Log().Info("Begin to ping xhs ...")
	_, err := api.GetNoteByKeyword("小红书", 1, 10, "GENERAL", "2")
	if err != nil {
		pong = false
	}
	return pong
}

// UpdateCookies will call this method when invalid login
func (api *XhsApiClient) UpdateCookies(ctx playwright.BrowserContext) {
	cookies, err := ctx.Cookies()
	util.AssertErrorToNil("could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("convert cookie failed and error:", err)
	api.httpClient.headers["Cookie"] = convertResp.cookieStr
	api.httpClient.cookiesMap = convertResp.cookiesMap
}

// GetNoteByKeyword get note list by search keywords
func (api *XhsApiClient) GetNoteByKeyword(keyword string, page int, pageSize int, sort SearchSortType, noteType SearchNoteType) (NoteSearchResult, error) {
	uri := "/api/sns/web/v1/search/notes"
	data := map[string]interface{}{
		"keyword":   keyword,
		"page":      page,
		"page_size": pageSize,
		"search_id": GetSearchID(),
		"sort":      string(sort),     // convert SearchSortType to string
		"note_type": string(noteType), // convert SearchNoteType to string
	}

	body, err := json.Marshal(data)
	if err != nil {
		return NoteSearchResult{}, err
	}

	resp, err := api.httpClient.Post(uri, body, map[string]string{})
	if err != nil {
		return NoteSearchResult{}, err
	}

	var result NoteSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return NoteSearchResult{}, err
	}

	return result, nil
}
