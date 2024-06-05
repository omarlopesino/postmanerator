package themes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aubm/postmanerator/postman"
)

func curlSnippet(request postman.Request) string {
	var curlSnippet string
	payloadReady, _ := regexp.Compile("POST|PUT|PATCH|DELETE")
	curlSnippet += fmt.Sprintf("curl -X %v", request.Method)

	if payloadReady.MatchString(request.Method) {
		if request.PayloadType == "urlencoded" {
			curlSnippet += ` -H "Content-Type: application/x-www-form-urlencoded"`
		} else if request.PayloadType == "params" || request.PayloadType == "formdata" {
			curlSnippet += ` -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW"`
		} else if request.PayloadType == "graphql" {
			curlSnippet += ` -H "Content-Type: application/json"`
		}
	}

	for _, header := range request.Headers {
		curlSnippet += fmt.Sprintf(` -H "%v: %v"`, header.Name, header.Value)
	}

	if payloadReady.MatchString(request.Method) {
		if request.PayloadType == "raw" && request.PayloadRaw != "" {
			curlSnippet += fmt.Sprintf(` -d '%v'`, request.PayloadRaw)
		} else if len(request.PayloadParams) > 0 {
			if request.PayloadType == "urlencoded" {
				var dataList []string
				for _, data := range request.PayloadParams {
					dataList = append(dataList, fmt.Sprintf("%v=%v", data.Key, data.Value))
				}
				curlSnippet += fmt.Sprintf(` -d "%v"`, strings.Join(dataList, "&"))
			} else if request.PayloadType == "params" || request.PayloadType == "formdata" {
				for _, data := range request.PayloadParams {
					curlSnippet += fmt.Sprintf(` -F "%v=%v"`, data.Key, data.Value)
				}
			}
		} else if request.PayloadType == "graphql" {
			// Query and variables breaklines are removed
			// as curl may not interpret correctly the JSON.
			curlSnippet += fmt.Sprintf(" --data '{\"query\": \"%s\", \"variables\": %s}'", strings.Replace(request.PayloadGraphQL.Query, "\n", " ", -1), strings.Replace(request.PayloadGraphQL.Variables, "\n", " ", -1))
		}
	}

	curlSnippet += fmt.Sprintf(` "%v"`, request.URL)
	return curlSnippet
}
