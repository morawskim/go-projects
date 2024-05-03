package main

import (
	"crypto/tls"
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func newCollyCollector(debugFlag bool) *colly.Collector {
	opt := []colly.CollectorOption{
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0"),
	}

	if debugFlag {
		opt = append(opt, colly.Debugger(&debug.LogDebugger{}))
	}

	c := colly.NewCollector(opt...)

	httpClient := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//InsecureSkipVerify: true, // Ignore SSL certificate errors
				// fix issue with amazon
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				},
			},
		},
	}
	c.SetClient(httpClient)
	//c.WithTransport(&http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore SSL certificate errors
	//})
	//proxyURL := "http://127.0.0.1:8080"
	//err := c.SetProxy(proxyURL)
	//if err != nil {
	//	panic(err)
	//}

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       5 * time.Second,
		RandomDelay: 4 * time.Second,
		Parallelism: 2,
	})

	return c
}

func collect(products []item2, selectors map[string]selector, pr map[string]string, ch chan metric) {
	data := make([]string, 0, len(products))

	c := newCollyCollector(false)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Host", r.URL.Host)
		r.Ctx.Put("product", pr[r.URL.String()])
	})

	for _, v := range selectors {
		c.OnHTML(v.Selector, func(e *colly.HTMLElement) {
			prodName, price := processHtmlElement(e, v.Expression)
			f, err := strconv.ParseFloat(price, 64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			ch <- metric{price: f, product: prodName}
		})
	}

	for _, product := range products {
		c.Visit(product.productUrl)
	}

	c.Wait()
	fmt.Println(data)
}

func processHtmlElement(e *colly.HTMLElement, expression string) (string, string) {
	prod := e.Response.Ctx.Get(productCtxKey)

	env := expressionEnv{Element: e}
	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	s := extractPrice(output.(string))

	return prod, s
}

func extractPrice(content string) string {
	content = onlyDigitsRegex.ReplaceAllString(content, "")
	content = strings.ReplaceAll(content, ",", ".")

	return content
}
