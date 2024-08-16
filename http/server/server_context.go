package server

import (
	"errors"
	"fmt"
	"github.com/oneliang/frame-golang/http/action"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/constants"
	"github.com/oneliang/util-golang/logging"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Interceptor func(requestUri string, method string, request *http.Request, writer http.ResponseWriter) bool
type ServerHandler struct {
	logger                      logging.Logger
	globalBeforeInterceptorList []Interceptor
	globalAfterInterceptorList  []Interceptor
	actionExecuteFunctionMap    map[string]action.ActionExecuteFunction
}

func NewServerHandler(actionExecuteFunctionMap map[string]action.ActionExecuteFunction) *ServerHandler {
	return &ServerHandler{
		logger:                      logging.LoggerManager.GetLoggerByPattern("ServerHandler"),
		globalBeforeInterceptorList: make([]Interceptor, 0),
		globalAfterInterceptorList:  make([]Interceptor, 0),
		actionExecuteFunctionMap:    actionExecuteFunctionMap,
	}
}

func (this *ServerHandler) generateRequestKey(requestUri string, httpMethod string) string {
	httpMethodUppercase := strings.ToUpper(httpMethod)
	return httpMethodUppercase + constants.SYMBOL_AT + requestUri
}

func (this *ServerHandler) AddGlobalBeforeInterceptor(interceptor Interceptor) {
	this.globalBeforeInterceptorList = append(this.globalBeforeInterceptorList, interceptor)
}

func (this *ServerHandler) AddGlobalAfterInterceptor(interceptor Interceptor) {
	this.globalAfterInterceptorList = append(this.globalAfterInterceptorList, interceptor)
}

func (this *ServerHandler) RegisterHandler(requestUri string, httpMethod string, actionExecuteFunction action.ActionExecuteFunction) error {
	httpMethodUppercase := strings.ToUpper(httpMethod)
	if httpMethodUppercase != constants.HTTP_REQUEST_METHOD_GET &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_POST &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_PUT &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_DELETE &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_HEAD &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_OPTIONS &&
		httpMethodUppercase != constants.HTTP_REQUEST_METHOD_TRACE {
		return errors.New(fmt.Sprintf("Not support the http method, input http method:%s, uppercase http method:%s", httpMethod, httpMethodUppercase))
	}
	requestKey := this.generateRequestKey(requestUri, httpMethod)
	this.actionExecuteFunctionMap[requestKey] = actionExecuteFunction
	return nil
}

// ServeHTTP can concurrent, no need to use go routine
func (this *ServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			errorMessage := fmt.Sprintf("Internal server error, has panic:%v", r)
			this.logger.Error(errorMessage, nil)
			http.Error(writer, errorMessage, 500)
		}
	}()

	requestUri := request.RequestURI
	httpMethod := request.Method
	for _, beforeInterceptor := range this.globalBeforeInterceptorList {
		result := beforeInterceptor(requestUri, httpMethod, request, writer)
		if result {
			continue
		} else {
			this.logger.Error("Global before intercept failure, %p", nil, beforeInterceptor)
			return
		}
	}

	httpHandlerKey := this.generateRequestKey(requestUri, httpMethod)
	this.logger.Info("Request http key:%s", httpHandlerKey)
	actionExecuteFunction, ok := this.actionExecuteFunctionMap[httpHandlerKey]
	if ok && actionExecuteFunction != nil {
		//func(request *http.Request, writer http.ResponseWriter) (error, []byte, int)
		err, data, statusCode := actionExecuteFunction(request, writer)
		//err, data, statusCode :=reflect.ValueOf(actionExecuteFunction).Call([]reflect.Value{reflect.ValueOf(request), reflect.ValueOf(writer)})
		dataString := ""
		if data != nil {
			dataString = string(data)
		}
		if err != nil {
			this.logger.Error("Execute http handler error, data:%s", nil, dataString)
			http.Error(writer, dataString, statusCode)
			return
		} else {
			writeResult, writeErr := writer.Write(data)
			if writeErr != nil {
				this.logger.Error("Response write error, write result:%d", nil, writeResult)
				http.Error(writer, "Response write error", 500)
				return
			} else {
				this.logger.Info("Response, data:%s, status code:%d", dataString, statusCode)
			}
		}
	} else {
		http.NotFound(writer, request)
	}
}

func GetRequestParameterInt(params url.Values, paramKey string, defaultValue int) int {
	var valueString = params.Get(paramKey)
	if common.StringIsNotBlank(valueString) {
		valueInt, err := strconv.Atoi(valueString)
		if err != nil {
			return defaultValue
		} else {
			return valueInt
		}
	} else {
		return defaultValue
	}
}

func GetRequestParameterString(params url.Values, paramKey string, defaultValue string) string {
	var valueString = params.Get(paramKey)
	if common.StringIsNotBlank(valueString) {
		return valueString
	} else {
		return defaultValue
	}
}
