### memcache 使用说明

在 "github.com/bradfitz/gomemcache" 上扩展了几个方法


1. 初始化

    ```
        memcache.Init("key_", 60, []string{"127.0.0.1:11211"})
    ```

2. 添加/查询 byte 类型
    
    ```
        str := "abcd"
        if err := memcache.SetByte("testa", []byte(str), 60); err != nil {
            logger.Error(err)
        }

        if b, err := memcache.GetByte("testa"); err == nil {
            logger.Debug(string(b))
        } else {
            logger.Error(err)
        }
    ```

3. 添加/查询任意类型 (内部使用 encoding/json 转换, 性能较好)

    ```
        person := &Person{"Jerry", 3}
        if err := memcache.Set("person", person, 60); err != nil {
            logger.Error(err)
        }
    
        if err := memcache.Get("person", person); err == nil {
            logger.Debug(person)
        } else {
            logger.Error(err)
        }
    ```
    
4. 添加/查询任意类型 (内部使用 encoding/gob 转换)

    ```
        person := &Person{"Jerry", 3}
        if err := memcache.SetObject("person", person, 60); err != nil {
            logger.Error(err)
        }

        for i := 0; i < 1000; i++ {
            if err := memcache.GetObject("person", person); err == nil {
                logger.Debug(person)
            } else {
                logger.Error(err)
            }
        }
    ```