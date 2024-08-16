package action

import (
	"errors"
	"fmt"
	"github.com/oneliang/frame-golang/context"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/logging"
	"net/http"
	"reflect"
)

type ActionExecuteFunction func(request *http.Request, writer http.ResponseWriter) (error, []byte, int)

type Action interface {
	RouteMap() map[string]ActionExecuteFunction
}

type ActionContext struct {
	context.Context
	logger                   logging.Logger
	actionExecuteFunctionMap map[string]ActionExecuteFunction
	actionInstanceMap        map[string]any
	globalInstanceMap        *map[string]any
}

func NewActionContext() *ActionContext {
	return &ActionContext{
		logger:                   logging.LoggerManager.GetLoggerByPattern("ActionContext"),
		actionExecuteFunctionMap: make(map[string]ActionExecuteFunction),
		actionInstanceMap:        make(map[string]any),
	}
}

func (this *ActionContext) SetGlobalInstanceMap(globalInstanceMap *map[string]any) {
	this.globalInstanceMap = globalInstanceMap
}

func (this *ActionContext) Initialize(param any) {
	if param == nil {
		this.logger.Error("Param can not be null, need to input array about []NewFunc", nil)
		return
	}
	if reflect.TypeOf(param).Kind() != reflect.Slice {
		this.logger.Error("Param must be array about []NewFunc", nil)
		return
	}
	err := this.newAllActionInstance(param.([]any))
	if err != nil {
		this.logger.Error("New all action instance error", err)
		return
	}
	this.logger.Info("Action execute function map:%p", this.actionExecuteFunctionMap)
}

func (this *ActionContext) newAllActionInstance(providerSet []any) error {
	for index, newActionFunction := range providerSet {
		newActionFunctionValue := reflect.ValueOf(newActionFunction)
		newActionFunctionType := reflect.TypeOf(newActionFunction)

		if newActionFunctionType.Kind() != reflect.Func {
			return errors.New(fmt.Sprintf("Only support new function, index:%d", index))
		}

		inputCount := newActionFunctionType.NumIn()
		if inputCount > 0 {
			return errors.New(fmt.Sprintf("Only support zero input about function, index:%d", index))
		}
		outputCount := newActionFunctionType.NumOut()
		if outputCount != 1 {
			return errors.New(fmt.Sprintf("Only support one output about new function, index:%d", index))
		}

		outputType := newActionFunctionType.Out(0)
		if outputType.Kind() != reflect.Ptr {
			errorMessage := fmt.Sprintf("Only support one pointer type about output with new function, index:%d", index)
			this.logger.Error(errorMessage, nil)
			return errors.New(errorMessage)
		}

		outputArray := newActionFunctionValue.Call(nil)
		actionValue := outputArray[0]
		actionInstance := actionValue.Interface()

		err, actionInterface := common.ConvertType[Action](actionInstance, (*Action)(nil))
		if err != nil {
			this.logger.Error("Convert type error, can not convert to Action, index:%d", err, index)
			return err
		}
		actionInstanceType := reflect.TypeOf(actionInstance)
		actionInstanceName := actionInstanceType.Elem().Name()
		this.logger.Info("Initialize action, action name:%s, action:%p", actionInstanceName, actionInstance)
		//put it to map after check covert type
		this.actionInstanceMap[actionInstanceName] = actionInstance
		(*this.globalInstanceMap)[actionInstanceName] = actionInstance
		routeMap := actionInterface.RouteMap()
		for route, actionExecuteFunction := range routeMap {
			this.logger.Info("Initialize action execute function, route:%s, action execute function:%p", route, actionExecuteFunction)
			this.actionExecuteFunctionMap[route] = actionExecuteFunction
		}
	}
	return nil
}

func (this *ActionContext) GetActionExecuteFunctionMap() map[string]ActionExecuteFunction {
	return this.actionExecuteFunctionMap
}

func (this *ActionContext) GetActionInstanceMap() map[string]any {
	return this.actionInstanceMap
}

func (this *ActionContext) Destroy() {
	clear(this.actionExecuteFunctionMap)
	clear(this.actionInstanceMap)
}
