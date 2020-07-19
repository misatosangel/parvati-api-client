package parvatigo

import (
	"github.com/misatosangel/parvati-api-client/pkg/swagger"
	"net/http"
)

type ApiError struct {
	ResponseError error
	RawResponse   *swagger.APIResponse
}

func ApiErr(resp *swagger.APIResponse, err error) *ApiError {
	if err == nil {
		return nil
	}
	return &ApiError{
		ResponseError: err,
		RawResponse:   resp,
	}
}

func HttpErr(resp *http.Response, err error) *ApiError {
	return ApiErr(swagger.NewAPIResponse(resp), err)
}

func (self *ApiError) HasResponse() bool {
	if self.RawResponse == nil || self.RawResponse.Response == nil {
		return false
	}
	return true
}

func (self *ApiError) Error() string {
	if self.RawResponse == nil {
		if self.ResponseError == nil {
			return "No error"
		}
		return self.ResponseError.Error()
	}
	var err string
	var is404 bool
	if self.RawResponse.Response != nil {
		resp := self.RawResponse.Response
		req := resp.Request
		if req != nil {
			err += "Request:  " + req.Method + " on " + req.URL.String() + "\n"
		}
		err += "Response: " + resp.Proto + " - " + resp.Status
		is404 = resp.StatusCode == 404
		if self.RawResponse.Message != "" {
			err += " - " + self.RawResponse.Message
		}
		err += "\n"
	} else {
		err += "No response object; message never sent.\n"
		if self.RawResponse.Message != "" {
			err += self.RawResponse.Message
		}
	}
	if !is404 && self.ResponseError != nil {
		err += self.ResponseError.Error()
	}
	return err
}
