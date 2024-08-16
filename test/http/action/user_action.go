package action

import (
	"fmt"
	httpAction "github.com/oneliang/frame-golang/http/action"
	"github.com/oneliang/frame-golang/test/http/service"
	"github.com/oneliang/util-golang/constants"
	"net/http"
)

type UserAction struct {
	UserService *service.UserService
}

func NewUserAction() *UserAction {
	return &UserAction{}
}

func (this *UserAction) RouteMap() map[string]httpAction.ActionExecuteFunction {
	return map[string]httpAction.ActionExecuteFunction{
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/create": this.Create,
		constants.HTTP_REQUEST_METHOD_GET + constants.SYMBOL_AT + "/read":    this.Read,
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/update": this.Update,
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/delete": this.Delete,
	}
}

func (this *UserAction) Create(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Create")
	this.UserService.Create()
	return nil, nil, 200
}

func (this *UserAction) Read(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Read")
	this.UserService.Read()
	return nil, nil, 200
}

func (this *UserAction) Update(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Update")
	this.UserService.Update()
	return nil, nil, 200
}

func (this *UserAction) Delete(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Delete")
	this.UserService.Delete()
	return nil, nil, 200
}
