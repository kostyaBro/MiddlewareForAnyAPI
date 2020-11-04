package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	myhttp "middleware_for_any_api/common/http"
	"middleware_for_any_api/common/log"
	logcontext "middleware_for_any_api/common/log/context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var listenPort = 8080
var domain = "https://pro-api.coinmarketcap.com"
var mainContext = logcontext.WithLogger(
	context.Background(),
	logrus.NewEntry(log.ProvideLogrusLoggerUseFlags()),
)

func main() {
	muxRouter := mux.NewRouter()
	muxRouter.PathPrefix("/").HandlerFunc(apiMiddleware)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), muxRouter); err != nil {
		panic(err)
	}
}

func apiMiddleware(responseWriter http.ResponseWriter, request *http.Request) {
	print(fmt.Sprintf("%s%s?%s\n", domain, request.URL.Path, request.URL.RawQuery))
	response, err := myhttp.NewRequestSender().
		WithUrl(fmt.Sprintf("%s%s?%s", domain, request.URL.Path, request.URL.RawQuery)).
		WithMethod(request.Method).
		WithHeaders(request.Header).
		WithBody(request.Body).
		Do(mainContext)
	if err != nil {
		logcontext.FromContext(mainContext).
			WithError(err).
			Error("request was end with error")
		sendResponseWithResponseCode(
			responseWriter,
			http.StatusServiceUnavailable,
		)
		return
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logcontext.FromContext(mainContext).
				WithError(err).
				Error("Can not close response body")
		}
	}()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logcontext.FromContext(mainContext).
			WithError(err).
			Error("could not read body")
		sendResponseWithResponseCode(
			responseWriter,
			http.StatusServiceUnavailable,
		)
		return
	}
	sendResponseWithResponseCodeAndBody(
		responseWriter,
		response.StatusCode,
		response.Header,
		respBytes,
	)
}

func sendResponseWithResponseCode(
	responseWriter http.ResponseWriter,
	responseCode int,
) {
	sendResponseWithResponseCodeAndBody(
		responseWriter,
		responseCode,
		nil,
		[]byte(strconv.Itoa(responseCode)),
	)
}

func sendResponseWithResponseCodeAndBody(
	responseWriter http.ResponseWriter,
	responseCode int,
	header http.Header,
	body []byte,
) {
	for s, strings := range header {
		responseWriter.Header().Set(s, strings[0])
	}
	responseWriter.WriteHeader(responseCode)
	_, err := responseWriter.Write(body)
	if err != nil {
		logcontext.FromContext(mainContext).
			WithError(err).
			Error("response was sent with error")
	}
}
