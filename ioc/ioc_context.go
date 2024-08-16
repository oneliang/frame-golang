package ioc

import (
	"errors"
	"fmt"
	"github.com/oneliang/frame-golang/context"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/logging"
	"reflect"
	"strings"
)

type IocContext struct {
	context.Context
	logger            logging.Logger
	instanceMap       map[string]any
	globalInstanceMap *map[string]any
}

func NewIocContext() *IocContext {
	return &IocContext{
		logger:      logging.LoggerManager.GetLoggerByPattern("IocContext"),
		instanceMap: make(map[string]any),
	}
}

func (this *IocContext) SetGlobalInstanceMap(globalInstanceMap *map[string]any) {
	this.globalInstanceMap = globalInstanceMap
}

func (this *IocContext) Initialize(param any) {
	if param == nil {
		this.logger.Error("Param can not be null, need to input array about []NewFunc", nil)
		return
	}
	if reflect.TypeOf(param).Kind() != reflect.Slice {
		this.logger.Error("Param must be array about []NewFunc", nil)
		return
	}
	err := this.newAllInstance(param.([]any))
	if err != nil {
		this.logger.Error("New all instance error", err)
		return
	}
	this.logger.Info("Ioc instance map:%v", this.instanceMap)
}

func (this *IocContext) newAllInstance(providerSet []any) error {
	for index, newInstanceFunction := range providerSet {
		newInstanceFunctionValue := reflect.ValueOf(newInstanceFunction)
		newInstanceFunctionType := reflect.TypeOf(newInstanceFunction)
		if newInstanceFunctionType.Kind() != reflect.Func {
			return errors.New(fmt.Sprintf("Only support new function, index:%d", index))
		}

		inputCount := newInstanceFunctionType.NumIn()
		if inputCount > 0 {
			return errors.New(fmt.Sprintf("Only support zero input about function, index:%d", index))
		}
		outputCount := newInstanceFunctionType.NumOut()
		if outputCount != 1 {
			return errors.New(fmt.Sprintf("Only support one output about function, index:%d", index))
		}

		outputArray := newInstanceFunctionValue.Call(nil)
		instanceValue := outputArray[0]
		instance := instanceValue.Interface()

		//instanceValue := reflect.ValueOf(instance)
		instanceType := reflect.TypeOf(instance)
		instanceTypeName := instanceType.Elem().Name()
		if common.StringIsBlank(instanceTypeName) {
			this.logger.Warning("Instance type name is empty, please check it, index:%d, type:%p", index, instanceType)
			continue
		}
		trimInstanceTypeName := strings.TrimSpace(instanceTypeName)
		var instanceKey = trimInstanceTypeName
		this.logger.Info("New initialize ioc instance, instance key:%s, instance:%p", instanceKey, instance)
		this.instanceMap[instanceKey] = instance
		(*this.globalInstanceMap)[instanceKey] = instance

	}
	return nil
}

func (this *IocContext) AutoInject() {
	this.logger.Info("Auto inject, global instance map:%v", this.globalInstanceMap)
	globalInstanceMap := *this.globalInstanceMap
	for instanceKey, instance := range globalInstanceMap {
		instanceType := reflect.TypeOf(instance).Elem()
		instanceValue := reflect.ValueOf(instance).Elem()

		for i := 0; i < instanceType.NumField(); i++ {
			field := instanceType.Field(i)
			fieldName := field.Name
			if !field.IsExported() {
				this.logger.Warning("Auto inject, field is unexported, can not auto inject, instance key:%s, field name:%s", instanceKey, fieldName)
				continue
			}

			fieldReferenceInstance, ok := globalInstanceMap[fieldName]
			this.logger.Info("Auto inject, instance key:%s, instance:%p, field name:%s, field reference instance:%p", instanceKey, instance, fieldName, fieldReferenceInstance)
			if ok { //exist
				fieldValue := instanceValue.FieldByName(fieldName)
				err := common.SetPointerValueByReflect(fieldValue, fieldReferenceInstance)
				if err != nil {
					this.logger.Error("Auto inject error", err)
				}
			} else {
				this.logger.Warning("Auto inject, field reference instance not exist, instance key:%s, field name:%s", instanceKey, fieldName)
			}
		}
	}
}

func (this *IocContext) Destroy() {
	clear(this.instanceMap)
}
