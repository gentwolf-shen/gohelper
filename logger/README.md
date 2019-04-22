### 日志方法使用说明

按不同的日志级别(TRACE < DEBUG < INFO < WARN < ERROR), 使用不同的颜色输出到控制台.
并可以设置将日志保存到文件的级别, 每天日志文件.

1. 配置文件说明

    ```
    {
      "level": "DEBUG", // 日志级别(TRACE < DEBUG < INFO < WARN < ERROR),在这个级别以下的,不输出
      "file": { // 定义需要输出到文件的日志
        "logPath": "/path/to/log/",    // 日志的保存目录, 或使用 ${application.path} (你的bin程序所在目录)
        "level": "INFO" // 在这个级别及以上的日志,输出到文件
      }
    }
    ```
    
2. 使用说明

    ```
        // 从配置文件初始化
        // logger.InitFromJson("/path/to/logger.json")
        
        // 使用默认初始化
        logger.InitDefault()
        
        // 输出日志信息
        logger.Trace("msg from trace")
        logger.Debug("msg from debug")
        logger.Info("msg from info")
        logger.Warn("msg from warn")
        logger.Error("msg from error")
        
        logger.Log(logger.LEVEL_ERROR, "ERROR MSG")
        
        // 在初始化后, 还可以动态的设置输出
        logger.SetLevel(logger.LEVEL_WARN)
    ```

    