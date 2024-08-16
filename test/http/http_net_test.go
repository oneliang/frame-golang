package http

import (
	"fmt"
	"github.com/oneliang/frame-golang/context"
	httpAction "github.com/oneliang/frame-golang/http/action"
	"github.com/oneliang/frame-golang/http/server"
	"github.com/oneliang/frame-golang/ioc"
	"github.com/oneliang/frame-golang/test/http/action"
	"github.com/oneliang/frame-golang/test/http/service"
	"log"
	"net/http"
	"testing"
)

func TestHttpNet(t *testing.T) {
	globalContext := context.NewGlobalContext()
	actionContext := httpAction.NewActionContext()
	if putContextErr := globalContext.PutContext("ActionContext", actionContext, action.ProviderSet); putContextErr != nil {
		log.Fatalf("%+v", putContextErr)
		return
	}
	iocContext := ioc.NewIocContext()
	if putContextErr := globalContext.PutContext("IocContext", iocContext, service.ProviderSet); putContextErr != nil {
		log.Fatalf("%+v", putContextErr)
		return
	}

	globalContext.Initialize(nil)
	iocContext.AutoInject()

	a := actionContext.GetActionInstanceMap()
	for _, value := range a {
		fmt.Println(fmt.Sprintf("%p", value))
		userAction := value.(*action.UserAction)
		fmt.Println(fmt.Sprintf("%p", userAction.UserService))
		fmt.Println(fmt.Sprintf("%p", userAction.UserService.RoleService))
	}

	//return
	serverHandler := server.NewServerHandler(actionContext.GetActionExecuteFunctionMap())
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverHandler,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server.ListenAndServe error:%v", err)
	}
}
