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
	"net/url"
)

type GamesApi struct {
	Configuration Configuration
}

func NewGamesApi() *GamesApi {
	configuration := NewConfiguration()
	return &GamesApi{
		Configuration: *configuration,
	}
}

func NewGamesApiWithBasePath(basePath string) *GamesApi {
	configuration := NewConfiguration()
	configuration.BasePath = basePath

	return &GamesApi{
		Configuration: *configuration,
	}
}

/**
 * Games API
 * The games endpoint returns information about the supported games of this hosting system.
 *
 * @param gameId Game to look up.
 * @return []Game
 */
func (a GamesApi) GamesGet(gameId string) ([]Game, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/games"

	headerParams := a.Configuration.GenDefaultHeaders()
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte

	if gameId != "" {
		queryParams.Add("game_id", a.Configuration.APIClient.ParameterToString(gameId, ""))
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
	var successPayload = new([]Game)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	if httpResponse.StatusCode() != 200 {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), errors.New(string(httpResponse.Body()))
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}
