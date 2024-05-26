package main

import (
	"fmt"
	"log"
	"plugin"
	"plugins/common"
)

func main() {
	pluginInstance, err := plugin.Open("./plugin-example.so")

	if err != nil {
		panic(fmt.Errorf("error opening plugin: %v", err))
	}

	symbol, err := pluginInstance.Lookup("GetPluginInfo")
	if err != nil {
		panic(fmt.Errorf("error looking up symbol: %v", err))
	}

	f, ok := symbol.(func() *common.PluginInfo)
	if !ok {
		log.Fatal("unexpected type from module symbol")
	}

	pluginInfo := f()
	fmt.Printf("plugin info: %+v\n", pluginInfo)

	pluginInfo.OnInit()
	err = pluginInfo.OnSomething(&common.SomeEvent{Foo: "Message passed from app to plugin"})

	if err != nil {
		panic(fmt.Errorf("error calling OnSomething: %v", err))
	}
}
