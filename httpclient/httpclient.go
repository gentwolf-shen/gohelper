package httpclient

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func UrlValuesToMap(urlValues url.Values) map[string]string {
	params := make(map[string]string, len(urlValues))
	for k, v := range urlValues {
		params[k] = v[0]
	}
	return params
}

func MapToUrlValues(data map[string]string) url.Values {
	params := url.Values{}

	for k, v := range data {
		params.Add(k, v)
	}
	return params
}

func Get(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return doRequest("GET", buildQuery(hostUrl, params), headers, nil)
}

func Delete(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return doRequest("DELETE", buildQuery(hostUrl, params), headers, nil)
}

func Put(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return doRequest("PUT", hostUrl, formHeader(headers), strings.NewReader(params.Encode()))
}

func Post(hostUrl string, params url.Values, headers map[string]string) ([]byte, error) {
	return doRequest("POST", hostUrl, formHeader(headers), strings.NewReader(params.Encode()))
}

func PutToBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest("PUT", hostUrl, jsonHeader(headers), bytes.NewReader(body))
}

func PostToBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest("POST", hostUrl, jsonHeader(headers), bytes.NewReader(body))
}

func DeleteBody(hostUrl string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest("DELETE", hostUrl, jsonHeader(headers), bytes.NewReader(body))
}

func doRequest(method string, hostUrl string, headers map[string]string, reader io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, hostUrl, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var responseBody []byte
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		responseBody, _ = ioutil.ReadAll(reader)
		reader.Close()
	default:
		responseBody, _ = ioutil.ReadAll(response.Body)
	}

	response.Body.Close()

	if response.StatusCode == 200 {
		return responseBody, nil
	}

	return responseBody, errors.New(strconv.Itoa(response.StatusCode))
}

func buildQuery(hostUrl string, params url.Values) string {
	if !strings.Contains(hostUrl, "?") {
		hostUrl += "?"
	} else {
		hostUrl += "&"
	}

	return hostUrl + params.Encode()
}

func formHeader(headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return headers
}

func jsonHeader(headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json;charset=utf-8"
	return headers
}

func PostFile(url string, param url.Values, files map[string]string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range param {
		for _, item := range v {
			_ = writer.WriteField(k, item)
		}
	}

	for k, v := range files {
		file, err := os.Open(v)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile(k, v)
		if err == nil {
			_, err = io.Copy(part, file)
		}

		file.Close()
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if nil != err {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func Download(filename string, url string) error {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}
