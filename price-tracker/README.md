# price-tracker

Application designed to track prices from various websites.

Key features include:
* Web Scraping via Colly library
* Configuration via yaml file
* Prometheus' metrics (latest price, minimal price, errors)
* Language expressions
* Tests for colly and expressions

## Available expressions

* `GetInputValue` - function used to retrieve values from input elements on a webpage. Many websites store information, such as prices, in hidden input fields.

* `GetFirstChildTextContent` - function used to extract the text content from the first existing child element of a specified parent element in the HTML structure of a webpage.
For example expression `GetFirstChildTextContent(Element, "[color=\"success800\"]", "[color=\"neutral900\"]")`
will return text content of element found by selector `[color=\"success800\"]`
if exists otherwise the text content of element found by selector `[color=\"neutral900\"]`  
This expression is useful for stores which display special prices of product in another HTML element.
