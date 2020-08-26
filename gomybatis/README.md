### gomybatis 使用说明

模仿MyBatis开发一个简单的go版本，实现了一些基本功能。

Query, Update, Delete, Insert的传入参数类型为 map[string]interface{},

Query的返回类型为 []map[string]string

xml文件请查看 sample.xml

1. 初始化

    ```
        import (
            "database/sql"
            _ "github.com/go-sql-driver/mysql"
            "gohelper/gomybatis"
            "gohelper/logger"
        )
        
        // ...
        
        // 初始化日志
        logger.InitDefault()
        
        // 初始化MySQL
        dbConn, err := sql.Open("mysql", "username:password@tcp(host:3306)dbname?charset=utf8")        
        defer dbConn.Close()
        
        // 传入mapper的XML文件所在目录
        gomybatis.SetMapperPath(dbConn, "/path/to/mapper/xml/")
        defer gomybatis.Close()
    ```

2. 查询
    
    ```
        args := make(map[string]interface{})            
        args["status"] = 1
        args["limit"] = 20
    
        // gomybatis.Query, gomybatis.QueryRow, gomybatis.QueryScalar
        if rows, err := gomybatis.Query("sample.query", args); err != nil {
            logger.Error(err)
        } else {
            logger.Debug(rows)
        }
    ```

3. 更新

    ```
        args := make(map[string]interface{})        
        args["id"] = "10"
        args["username"] = "test-username"
        args["email"] = "test@email.com"
        
        if n, err := gomybatis.Update("sample.update", args); err != nil {
            logger.Error(err)
        } else {
            logger.Debug(n)
        }
    ```

4. 删除

    ```
        args := make(map[string]interface{})        
        args["id"] = 11
        
        if n, err := gomybatis.Delete("sample.delete", args); err != nil {
            logger.Error(err)
        } else {
            logger.Debug(n)
        }
    ```

5. 添加

    ```
        args := make(map[string]interface{})        
        args["username"] = "test-username"
        args["email"] = "test@email.com"
        
        if lastId, err := gomybatis.Insert("sample.insert", args); err != nil {
            logger.Error(err)
        } else {
            logger.Debug(lastId)
        }
    ```