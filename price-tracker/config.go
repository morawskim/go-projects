package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type selector struct {
	Expression string `yaml:"expression"`
	Selector   string `yaml:"selector"`
}

type configProductItem struct {
	Name        string `yaml:"name"`
	Url         string `yaml:"url"`
	SelectorRef string `yaml:"selectorRef"`
	Expression  string `yaml:"expression"`
}

type config struct {
	Selectors map[string]selector `yaml:"selectors"`
	Products  []configProductItem `yaml:"products"`
}

func isValidFile(fileName string) bool {
	// Use filepath.Abs to resolve any relative paths
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// Check if the file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func processConfig(c *config) ([]item2, map[string]string) {
	productsToObserve := make([]item2, 0, len(c.Products))
	prodUrlToProdNameMap := make(map[string]string, len(c.Products))

	for _, i := range c.Products {
		prodUrlToProdNameMap[i.Url] = i.Name
		productsToObserve = append(productsToObserve, item2{
			productName: i.Name,
			productUrl:  i.Url,
		})
	}

	return productsToObserve, prodUrlToProdNameMap
}
