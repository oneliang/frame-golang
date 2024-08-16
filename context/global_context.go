package context

import (
	"errors"
	"github.com/oneliang/util-golang/logging"
)

type GlobalContext struct {
	Context
	logger            logging.Logger
	contextMap        map[string]contextBean
	contextList       []contextBean
	globalInstanceMap map[string]any
}

func NewGlobalContext() *GlobalContext {
	return &GlobalContext{
		logger:            logging.LoggerManager.GetLoggerByPattern("GlobalContext"),
		contextMap:        make(map[string]contextBean),
		contextList:       make([]contextBean, 0),
		globalInstanceMap: make(map[string]any),
	}
}
func (this *GlobalContext) SetGlobalInstanceMap(globalInstanceMap *map[string]any) {}

func (this *GlobalContext) Initialize(param any) {
	for _, contextBeanInstance := range this.contextList {
		context := contextBeanInstance.context
		contextParam := contextBeanInstance.contextParam
		this.logger.Info("Initialize context:%p, context param:%p", context, contextParam)
		context.SetGlobalInstanceMap(&this.globalInstanceMap)
		this.logger.Info("global instance map size:%d", len(this.globalInstanceMap))
		context.Initialize(contextParam)
	}
}

func (this *GlobalContext) PutContext(id string, context Context, param any) error {
	if context == nil {
		this.logger.Error("context can not be nil", nil)
		return errors.New("context can not be nil")
	}
	contextBeanInstance := newContextBean(id, context, param)
	this.contextMap[id] = contextBeanInstance
	this.contextList = append(this.contextList, contextBeanInstance)
	return nil
}

func (this *GlobalContext) GetContext(id string) Context {
	return this.contextMap[id].context
}

func (this *GlobalContext) Destroy() {

}
