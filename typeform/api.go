package typeform

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Api struct {
	client *http.Client
	token  string
}

func NewTypeformApi(token string) *Api {
	api := new(Api)
	api.client = &http.Client{}
	api.token = token
	return api
}

func (t *Api) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	var correctedPath string
	if strings.HasPrefix(path, "/") && strings.HasSuffix(ApiUrl, "/") {
		correctedPath = strings.TrimPrefix(path, "/")
	} else if !strings.HasPrefix(path, "/") && !strings.HasSuffix(ApiUrl, "/") {
		correctedPath = "/" + path
	} else {
		correctedPath = path
	}

	// Create Request
	request, err := http.NewRequest(method, ApiUrl+correctedPath, body)
	if err != nil {
		return nil, err
	}

	// Set Auth Header
	request.Header.Set("Authorization", "Bearer "+t.token)

	return request, nil
}

func (t *Api) GetResponses(params GetResponsesParams) (*GetResponsesResponse, error) {
	// Form ID for in path
	if params.FormId == "" {
		log.Fatal("INVALID_FORM_ID")
	}
	path := "/forms/" + params.FormId + "/responses"
	request, err := t.newRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create Query
	q := request.URL.Query()

	// Page Size
	if params.PageSize > 1000 {
		log.Fatal("EXCEEDED_MAX_PAGE_SIZE")
	} else if params.PageSize == 0 {
		params.PageSize = 25
	}

	q.Add("page_size", strconv.Itoa(params.PageSize))

	// Since
	if params.Since.Valid {
		q.Add("since", params.Since.Time.Format(time.RFC3339))
	}

	// Until
	if params.Until.Valid {
		q.Add("until", params.Until.Time.Format(time.RFC3339))
	}

	// After
	if params.After != "" {
		q.Add("after", params.After)
	}

	// Before
	if params.Before != "" {
		q.Add("before", params.Before)
	}

	// Included Response IDs
	if len(params.IncludedResponseIds) > 0 {
		q.Add("included_response_ids", strings.Join(params.IncludedResponseIds, ","))
	}

	// Completed
	if params.Completed != nil {
		value := "false"
		if *params.Completed {
			value = "true"
		}
		q.Add("completed", value)
	}

	// Query
	if params.Query != "" {
		q.Add("query", params.Query)
	}

	// Fields
	if len(params.Fields) > 0 {
		q.Add("fields", strings.Join(params.Fields, ","))
	}

	// Encode Query String
	request.URL.RawQuery = q.Encode()

	// Get Response
	response, err := t.client.Do(request)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse json
	var data GetResponsesResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	// Return Response
	return &data, nil
}
