package gohttp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

type GoHttpRequest struct {
	uri    string
	method string

	queryValues url.Values
	formValues  url.Values
	headers     map[string]string
	cookies     map[int]*http.Cookie

	requestBody    []byte
	responseResult interface{}
}

type GoHttpResponse struct {
	Body       []byte            `json:"body"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
}

func new(method, uri string) *GoHttpRequest {
	request := &GoHttpRequest{}
	request.uri = uri
	request.method = method

	request.queryValues = url.Values{}
	request.formValues = url.Values{}
	request.headers = make(map[string]string)

	return request
}

func (this *GoHttpRequest) BindResponseJson(v interface{}) *GoHttpRequest {
	this.responseResult = v
	return this
}

func (this *GoHttpRequest) SetQueryValue(key, value string) *GoHttpRequest {
	this.queryValues.Add(key, value)
	return this
}

func (this *GoHttpRequest) SetQueryValues(values map[string]string) *GoHttpRequest {
	for key, value := range values {
		this.queryValues.Add(key, value)
	}
	return this
}

func (this *GoHttpRequest) SetFormValue(key, value string) *GoHttpRequest {
	this.formValues.Add(key, value)
	return this
}

func (this *GoHttpRequest) SetFormValues(values map[string]string) *GoHttpRequest {
	for key, value := range values {
		this.formValues.Add(key, value)
	}
	return this
}

func (this *GoHttpRequest) SetBodyJson(v interface{}) *GoHttpRequest {
	b, _ := json.Marshal(v)
	return this.SetBodyByte(b)
}

func (this *GoHttpRequest) SetBodyByte(body []byte) *GoHttpRequest {
	this.requestBody = body
	return this
}

func (this *GoHttpRequest) SetHeader(key, value string) *GoHttpRequest {
	this.headers[key] = value
	return this
}

func (this *GoHttpRequest) SetHeaders(values map[string]string) *GoHttpRequest {
	for key, value := range values {
		this.headers[key] = value
	}
	return this
}

func (this *GoHttpRequest) SetCookie(cookie *http.Cookie) *GoHttpRequest {
	if this.cookies == nil {
		this.cookies = make(map[int]*http.Cookie)
	}

	this.cookies[len(this.cookies)] = cookie
	return this
}

func (this *GoHttpRequest) buildQuery() string {
	query := ""

	if len(this.queryValues) > 0 {
		if strings.Contains(this.uri, "?") {
			query = "&" + this.queryValues.Encode()
		} else {
			query = "?" + this.queryValues.Encode()
		}
	}

	return query
}

func (this *GoHttpRequest) Do() (*GoHttpResponse, error) {
	if this.method == MethodHead {
		return this.header()
	}

	r := &GoHttpResponse{}
	request, err := http.NewRequest(this.method, this.uri+this.buildQuery(), this.getReader())
	if err != nil {
		return r, err
	}

	for key, value := range this.headers {
		request.Header.Add(key, value)
	}

	for _, c := range this.cookies {
		request.AddCookie(c)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return r, err
	}

	r.StatusCode = response.StatusCode
	this.parseHeader(r, response.Header)

	if this.method == MethodOptions {
		return r, nil
	}

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		r.Body, _ = ioutil.ReadAll(reader)
		_ = reader.Close()
	default:
		r.Body, _ = ioutil.ReadAll(response.Body)
	}

	_ = response.Body.Close()

	if this.responseResult != nil {
		if err := json.Unmarshal(r.Body, this.responseResult); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (this *GoHttpRequest) header() (*GoHttpResponse, error) {
	r := &GoHttpResponse{}

	response, err := http.Head(this.uri + this.buildQuery())
	if err != nil {
		return r, err
	}

	r.StatusCode = response.StatusCode
	this.parseHeader(r, response.Header)

	return r, nil
}

func (this *GoHttpRequest) parseHeader(r *GoHttpResponse, headers http.Header) {
	r.Headers = make(map[string]string, len(headers))
	for key, value := range headers {
		r.Headers[key] = value[0]
	}
}

func (this *GoHttpRequest) getReader() io.Reader {
	var reader io.Reader

	if this.method == MethodGet || this.method == MethodOptions {
		return reader
	}

	if this.requestBody != nil {
		reader = bytes.NewReader(this.requestBody)
	} else {
		if len(this.formValues.Encode()) > 0 {
			if _, ok := this.headers["Content-Type"]; !ok {
				this.SetHeader("Content-Type", "application/x-www-form-urlencoded")
			}

			reader = strings.NewReader(this.formValues.Encode())
		}
	}

	return reader
}
