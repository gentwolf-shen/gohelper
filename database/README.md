### 数据库操作类使用说明

当前支持MySQL，SQLite，PostgreSQL，Oracle。

1. 驱动及配置

    1. MySQL： github.com/Go-SQL-Driver/MySQL
       
       DSN配置： dbuser:dbpassword@tcp(127.0.0.1:3306)/dbname?charset=utf8
    
    2. SQLite： github.com/mattn/go-sqlite3
    
       DSN配置： dbuser:dbpassword@tcp(127.0.0.1:3306)/dbname?charset=utf8
    
    3. PostgreSQL： github.com/lib/pq
    
       DSN配置： postgres://dbuser:dbpassword@127.0.0.1:5432/log_show?sslmode=disable
       
    4. Oracle： github.com/mattn/go-oci8
    
        DSN配置： system/oracle@127.0.0.1:1521/xe

2. 初始化
    1. 直接使用，初始化数据库操作对象
    ```
    db := database.New()    
    err := db.Open("mysql", "dbuser:dbpassword@tcp(127.0.0.1:3306)/dbname?charset=utf8", 1, 1)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    ```
    
    2. 使用配置文件，可以初始化多个数据库操作对象
    ```    
    // 从字符串中加载配置
    str := `
            {
                "db": {
                    "default": {
                        "type": "mysql",
                        "dsn": "dbuser:dbpassword@tcp(127.0.0.1:3306)/dbname?charset=utf8",
                        "maxOpenConnections": 1,
                        "maxIdleConnections": 1
                    },
                    "db1": {
                        "type": "sqlite3",
                        "dsn": "/path/to/app.db",
                        "maxOpenConnections": 1,
                        "maxIdleConnections": 1
                    },
                    "db2": {
                        "type": "postgres",
                        "dsn": "postgres://dbuser:dbpassword@127.0.0.1:5432/dbname?sslmode=disable",
                        "maxOpenConnections": 1,
                        "maxIdleConnections": 1
                    },
                     "db2": {
                        "type": "oci8",
                        "dsn": "system/oracle@127.0.0.1:1521/xe",
                        "maxOpenConnections": 1,
                        "maxIdleConnections": 1
                    }
                }
            }
            `
    cfg, err := config.LoadFromStr(str)
    if err != nil {
        logger.Error(err)
        return
    }
    
    // 可以使用config.Load从文件中加载配置
    // cfg , err := config.Load("path/to/application.json")

    if err := database.LoadFromConfig(cfg.Db); err != nil {
        logger.Error(err)
        return
    }

    db := database.Driver("default")
    ```

3. 查询数据
    ```
    /**
        数据查询
        Query: 返回[]map[string]string
        QueryRow: 返回 map[string]string
        QueryScalar: 返回 string
    **/
    rows, err := db.Query("SELECT id,username FROM user WHERE id<=?", 10)
    fmt.Println(err)
    fmt.Println(rows)
    ```
    
4. 添加数据
    ```
    // 返回LastInsertId，需要数据库驱动支持
    LastInsertId, err := db.Insert("INSERT user (username,email) VALUES(?,?)", "USERNAME", "EMAIL")
    fmt.Println(err)
    fmt.Println(LastInsertId)
    ```

5. 更新数据
    ```
    // 返回更新影响的记录数
    n, err := db.Update("UPDATE user SET status=? WHERE id=?", 1, 2)
    fmt.Println(err)
    fmt.Println(n)
    ```
    
6. 删除数据
    ```
    // 返回被删除的记录籹
    n, err := db.Delete("DELETE FROM user WHERE id>=?", 2)
    fmt.Println(err)
    fmt.Println(n)
    ```
    
7. 使用预处理stmt操作数据库
    ```
    // db.CreateStmt("SQL语句", "stmt的名称，如果不传入，则使用md5(SQL语句)作为键名，缓存在内存中")
    stmt, err := db.CreateStmt("SELECT id,username FROM user WHERE id<=?", "queryStmt")
    stmt, err := db.CreateStmt("INSERT user (username,email) VALUES(?,?)" "insertStmt")
    stmt, err := db.CreateStmt("UPDATE user SET status=? WHERE id=?", "updateStmt")
   
    if err != nil {
    	fmt.Println(err)
    } else {
    
    
    //rows, err := stmt.Query(55432)
    //lastId, err := stmt.Insert("USERNAME", "EMAIL")
    n, err := stmt.Update(1, 2)
    
    //fmt.Println(err)
    fmt.Println(n)
    
    // 关闭stmt， 也可以不显示关闭，在db.Close()或database.CloseAll()中会自动关闭
    db.CloseStmt("stmt的名称)
    ```    
   
7. 返回struct数据

    ```
    type User struct {
        Id       int64  `db:"id"`
        Username string `db:"username"`
        Email    string `db:"email"`
    }
   
    user := User{}
   
    err = db.QueryObject(&user, "SELECT id,username,email FROM user where id = ?", 1)
    fmt.Println(err)
    fmt.Println(user)
    
    all, err := db.QueryObjects(&user, "SELECT id,username,email  FROM id < ?", 10)
    fmt.Println(err)
    fmt.Println(all)
    ```