### dict 使用说明

将json文件解析为map[string]string类型的值， 支持从环境变量中替换value。

dict.json示例：
```
{
  "version": "1.0",
  "name": "${ENV.name:Tom}",
  "age": "${ENV.age:0}"
}
```

使用说明：

```go
package main

import (
    "fmt"
    "github.com/gentwolf-shen/gohelper/dict"
)

func main() {
    //开启环境变量支持
    dict.EnableEnv = true
    _ = dict.LoadDefault()

    fmt.Println(dict.Get("name"))
    fmt.Println(dict.Get("age"))
}
```
