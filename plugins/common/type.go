package common

type onInitFunc func()
type onSomethingHappen func(*SomeEvent) error

type SomeEvent struct {
	Foo string
}

type PluginInfo struct {
	Name        string
	Version     string
	OnInit      onInitFunc
	OnSomething onSomethingHappen
}
