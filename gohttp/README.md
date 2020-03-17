### gohttp 使用说明

使用类似stream的方式处理http的请求，支持 GET, POST, PUT, PATH, HEAD, OPTIONS。

1. GET

    ```
    var v interface{}
    url := "http://test.com/test.php"
    response, err := gohttp.Get(url).SetQueryValue("ver", "1.0").BindResponseJson(&v).Do()
    fmt.Println(err)
    fmt.Println(string(response.Body))
    if err == nil {
        fmt.Println(v)
    }
    ```

2. POST (form)

    ```
    url := "http://test.com/test.php"
    response, err := gohttp.Post(url).
        SetQueryValue("ver", "1.0").
        SetFormValue("name", "Tom").
        Do()
    fmt.Println(err)
    fmt.Println(string(response.Body))
    ```

3. POST (application/json)

    ```
   p := map[string]interface{}{
   		"name": "Tom",
   		"age":  3,
   	}
    url := "http://test.com/test.php"
    response, err := gohttp.Post(url).
        SetBodyJson(p).
        SetHeader("Content-Type", "application/json").
        Do()
    fmt.Println(err)
    fmt.Println(string(response.Body))
    ```

4. PUT, PATCH

5. DELETE
