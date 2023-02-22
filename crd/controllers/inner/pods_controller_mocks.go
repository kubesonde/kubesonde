package inner

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
	"kubesonde.io/controllers/probe_command"
)

type MockedCNIState struct {
	mock.Mock
}

func (mock *MockedCNIState) getClient() kubernetes.Interface {
	args := mock.Called()
	return args.Get(0).(kubernetes.Interface)
}
func (mock *MockedCNIState) logError(err error, message string) {
	mock.Called(err, message)

}
func (mock *MockedCNIState) logInfo(value string) {
	mock.Called(value)

}

func (mock *MockedCNIState) runCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	ret := mock.Called(client, namespace, command, checker)
	return ret.Get(0).(bool), ret.Error(1)

}

func (mock *MockedCNIState) runGenericCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand) (string, error) {
	ret := mock.Called(client, namespace, command)
	return ret.Get(0).(string), ret.Error(1)

}

func (mock *MockedCNIState) AddEdge(_ probe_command.KubesondeCommand, _ bool) {
	mock.Called()

}
