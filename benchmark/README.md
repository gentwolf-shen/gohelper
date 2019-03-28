### 基准测试方法使用说明

计算程序的运行时间

1. 使用说明
    ```
        benchmark.Start("start")
        
    
        // 计时开始
        // ...
        // 计时结束
        
    
        fmt.Println(benchmark.StopMilliseconds("start"), "s")
        fmt.Println(benchmark.StopMilliseconds("start"), "s")
    ```