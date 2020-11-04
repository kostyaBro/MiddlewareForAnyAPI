package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"middleware_for_any_api/common/errors"
	logcontext "middleware_for_any_api/common/log/context"
)

type RequestUpdater func(*http.Request) *http.Request

func GenerateUrl(host string, pathFormat string, pathParams ...interface{}) string {
	return fmt.Sprintf("%s%s", host, fmt.Sprintf(pathFormat, pathParams...))
}

type requestSender struct {
	httpClient *http.Client
	url        string
	method     string
	headers    http.Header
	body       io.Reader
}

func NewRequestSender() *requestSender {
	return &requestSender{
		httpClient: http.DefaultClient,
	}
}

func NewRequestSenderWithHttpClient(client *http.Client) *requestSender {
	return &requestSender{
		httpClient: client,
	}
}

func (rs *requestSender) WithUrl(url string) *requestSender {
	rs.url = url
	return rs
}

func (rs *requestSender) WithHeaders(headers http.Header) *requestSender {
	rs.headers = headers
	return rs
}

func (rs *requestSender) WithMethod(method string) *requestSender {
	rs.method = method
	return rs
}

func (rs *requestSender) WithBody(body io.Reader) *requestSender {
	rs.body = body
	return rs
}

func (rs *requestSender) Do(
	context context.Context,
) (*http.Response, error) {
	context = logcontext.WithLogger(
		context,
		logcontext.FromContext(context).
			WithField("req_url", rs.url).
			WithField("req_method", rs.method),
	)
	var response *http.Response
	var err error
	response, err = rs.doRequest(rs.url, rs.method, rs.body)
	if err != nil {
		logcontext.FromContext(context).WithError(err).Error("do request error")
		return nil, err
	}

	return response, nil
}

func (rs *requestSender) doRequest(
	url string, method string, body io.Reader,
) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept-Charset", "utf-8")
	request.Header = rs.headers
	return rs.httpClient.Do(request)
}

func (rs *requestSender) parseResponse(
	response *http.Response, parsedResponse interface{},
) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseBody, errors.Newf("сan not read response body: %s", err.Error())
	}
	err = json.Unmarshal(responseBody, parsedResponse)
	if err != nil {
		return responseBody, errors.Newf("can not parse response: %s, err: %s",
			string(responseBody), err.Error(),
		)
	}
	return responseBody, nil
}
