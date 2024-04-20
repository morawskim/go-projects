package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><body>
<div id="buybox">
	<p color="success800">99,99 zł</p>
	<p color="neutral900">199,99 zł</p>
</div>
</body></html>`)
	})
	mux.HandleFunc("/html2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><body>
<div id="buybox">
	<input id="price" type="hidden" value="19.99">
</div>
</body></html>`)
	})
	ts := httptest.NewUnstartedServer(mux)

	return ts
}

func TestGetFirstChildTextContent(t *testing.T) {
	ts := createTestServer()
	ts.Start()
	defer ts.Close()

	expectedProductName := "Foo"
	expectedPrice := "99.99"
	hasBeenCalled := false
	c := colly.NewCollector()

	c.OnHTML("#buybox", func(e *colly.HTMLElement) {
		hasBeenCalled = true
		e.Request.Ctx.Put(productCtxKey, expectedProductName)
		prod, price := processHtmlElement(e, `GetFirstChildTextContent(Element, "[color=\"success800\"]", "[color=\"neutral900\"]")`)

		if prod != expectedProductName {
			t.Fatalf("Expected produt %#v. Got %#v", expectedProductName, prod)
		}

		if price != expectedPrice {
			t.Fatalf("Expected produt %#v. Got %#v", expectedPrice, price)
		}
	})

	c.Visit(ts.URL + "/html")
	c.Wait()

	if !hasBeenCalled {
		t.Fatalf("The callback has not been called")
	}
}

func TestGetInputValue(t *testing.T) {
	ts := createTestServer()
	ts.Start()
	defer ts.Close()

	expectedProductName := "Foo"
	expectedPrice := "19.99"
	hasBeenCalled := false
	c := colly.NewCollector()

	c.OnHTML("input#price", func(e *colly.HTMLElement) {
		hasBeenCalled = true
		e.Request.Ctx.Put(productCtxKey, expectedProductName)
		prod, price := processHtmlElement(e, `GetInputValue(Element)`)

		if prod != expectedProductName {
			t.Fatalf("Expected produt %#v. Got %#v", expectedProductName, prod)
		}

		if price != expectedPrice {
			t.Fatalf("Expected produt %#v. Got %#v", expectedPrice, price)
		}
	})

	c.Visit(ts.URL + "/html2")
	c.Wait()

	if !hasBeenCalled {
		t.Fatalf("The callback has not been called")
	}
}
