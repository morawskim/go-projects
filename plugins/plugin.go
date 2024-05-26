package main

import (
	"fmt"
	"plugins/common"
)

func GetPluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:    "Simple plugin",
		Version: "0.10",
		OnInit: func() {
			fmt.Println("Plugin on init has been called")
		},
		OnSomething: func(event *common.SomeEvent) error {
			fmt.Println("Something happen! Foo:", event.Foo)

			return nil
		},
	}
}
