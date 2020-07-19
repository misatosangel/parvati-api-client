/*
 * Host API
 *
 * API for posting hosts and waiting for hosts
 *
 * OpenAPI spec version: 1.0.0
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the MIT License
 */

package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type UsersApi struct {
	Configuration Configuration
}

func NewUsersApi() *UsersApi {
	configuration := NewConfiguration()
	return &UsersApi{
		Configuration: *configuration,
	}
}

func NewUsersApiWithBasePath(basePath string) *UsersApi {
	configuration := NewConfiguration()
	configuration.BasePath = basePath

	return &UsersApi{
		Configuration: *configuration,
	}
}

/**
 * User Activity
 * The User Activity endpoint returns data about a user&#39;s hosts.
 *
 * @param offset Offset the list of returned results by this amount. Default is zero.
 * @param limit Number of items to retrieve. Default is 5, maximum is 100.
 * @return Array of found Host objects.
 */
func (a UsersApi) HistoryGet(userId string, since, before *time.Time, seePrivate bool, limit int32) ([]Host, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	userId = strings.Replace(userId, "/", "%2F", -1)
	path := a.Configuration.BasePath + "/users/" + userId + "/history"

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte

	//    queryParams.Add("offset", a.Configuration.APIClient.ParameterToString(offset, ""))
	queryParams.Add("limit", a.Configuration.APIClient.ParameterToString(limit, ""))
	if since != nil {
		queryParams.Add("since", since.Format(time.RFC3339Nano))
	}
	if before != nil {
		queryParams.Add("before", before.Format(time.RFC3339Nano))
	}
	if seePrivate {
		queryParams.Add("unmask-private", "1")
	}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload []Host
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * User Profile
 * The User Profile endpoint returns information about the user that has authorized with the application.
 *
 * @return *User
 */
func (a UsersApi) MeGet() (*User, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/me"

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new(User)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
func (a UsersApi) UsersGetWithChallonge() ([]User, *APIResponse, error) {
 * Users API
 * The Users endpoint returns information about the Users of the system. The response includes the display name and other details about each user. Not all information is public for a given user.
 *
 * @return []User
*/
func (a UsersApi) UsersGetWithChallonge() ([]User, *APIResponse, error) {
	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/users"
	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	queryParams.Add("challonge", "yes")

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]User)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Users API
 * The Users endpoint returns information about the Users of the system. The response includes the display name and other details about each user. Not all information is public for a given user.
 *
 * @param userId A specific user-id or username to query.
 * @param name Name or part of of a name of the user.
 * @param country Country of the user.
 * @return []User
 */
func (a UsersApi) UsersGet(userId string, name string, country string) ([]User, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/users"

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte

	if userId != "" {
		queryParams.Add("user_id", a.Configuration.APIClient.ParameterToString(userId, ""))
	}
	if name != "" {
		queryParams.Add("name", a.Configuration.APIClient.ParameterToString(name, ""))
	}
	if country != "" {
		queryParams.Add("country", a.Configuration.APIClient.ParameterToString(country, ""))
	}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]User)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	if err != nil {
		fmt.Printf("Failure to unmarshall '" + string(httpResponse.Body()) + "'\n")
	}
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Users API
 * Grab information on a specific user.
 *
 * @param userId A specific user-id or username to query.
 * @return User
 */
func (a UsersApi) UserGet(userId string) (User, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/users"
	//	path = strings.Replace(path, "{"+"user_id"+"}", fmt.Sprintf("%v", userId), -1)

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	queryParams["user_id"] = []string{userId}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload User
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	if httpResponse.StatusCode() != 200 {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), errors.New(string(httpResponse.Body()))
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

func (a UsersApi) UpdateUser(userId string, postBody map[string]string) (UserDelta, *APIResponse, error) {
	var httpMethod = "Post"
	// create path and map variables
	path := a.Configuration.BasePath + "/users/" + userId

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)

	var fileName string
	var fileBytes []byte

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload UserDelta
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		if httpResponse != nil {
			return successPayload, NewAPIResponse(httpResponse.RawResponse), err
		}
		return successPayload, nil, err
	}
	if httpResponse.StatusCode() != 200 {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), errors.New(string(httpResponse.Body()))
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

func (a UsersApi) UserCreate(credential, backendNick string, ip net.IP, port int) (UserDelta, *APIResponse, error) {
	var httpMethod = "Post"
	// create path and map variables
	path := a.Configuration.BasePath + "/users"

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody = make(map[string]string)
	postBody["nick"] = backendNick
	postBody["credential"] = credential
	if ip != nil {
		if ip.To4() == nil {
			postBody["ipv6"] = ip.String()
		} else {
			postBody["ip"] = ip.String()
		}
	}
	if port != 0 {
		postBody["port"] = fmt.Sprintf("%d", port)
	}

	var fileName string
	var fileBytes []byte

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload UserDelta
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	if httpResponse.StatusCode() != 200 {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), errors.New(string(httpResponse.Body()))
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}