package inner

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
)

type MockedCNIState struct {
	mock.Mock
}

func (mock *MockedCNIState) getClient() *kubernetes.Clientset {
	mock.Called()
	return &kubernetes.Clientset{}
}
func (mock *MockedCNIState) logError(err error, message string) {
	mock.Called(err, message)
	return

}
func (mock *MockedCNIState) logInfo(value string) {
	mock.Called(value)
	return

}

func (mock *MockedCNIState) getKubesonde() v12.Kubesonde {
	mock.Called()
	return v12.Kubesonde{Spec: v12.KubesondeSpec{Actions: probingActions}}

}
func (mock *MockedCNIState) runCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	ret := mock.Called(client, namespace, command, checker)
	return ret.Get(0).(bool), ret.Error(1)

}

func (mock *MockedCNIState) runGenericCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand) (string, error) {
	ret := mock.Called(client, namespace, command)
	return ret.Get(0).(string), ret.Error(1)

}

func (mock *MockedCNIState) AddEdge(kubesondeCommand probe_command.KubesondeCommand, value bool) {
	mock.Called()

}
