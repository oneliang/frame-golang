package context

type Context interface {
	SetGlobalInstanceMap(globalInstanceMap *map[string]any)
	Initialize(param any)
	Destroy()
}

type contextBean struct {
	id           string
	context      Context
	contextParam any
}

func newContextBean(id string, context Context, contextParam any) contextBean {
	return contextBean{
		id:           id,
		context:      context,
		contextParam: contextParam,
	}
}
