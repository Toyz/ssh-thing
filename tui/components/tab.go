package components

import (
	"github.com/toyz/ssh-thing/ssh"
)

type TabContent struct {
	Client     *ssh.Client
	ScrollView *ScrollView
	HasError   bool
	ErrorMsg   string
	Name       string
}

func NewTabContent(name string) *TabContent {
	return &TabContent{
		Name:       name,
		ScrollView: NewScrollView(),
		HasError:   false,
	}
}

func (t *TabContent) HandleError(err error) {
	t.HasError = true
	if err != nil {
		t.ErrorMsg = err.Error()
	}
}

func (t *TabContent) SetClient(client *ssh.Client) {
	t.Client = client
}

func (t *TabContent) Close() {
	if t.Client != nil {
		t.Client.Close()
	}
}
