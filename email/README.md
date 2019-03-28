### 邮件发送

可以发送普通文本, HTML文本, 附件等.

1. 使用说明

    ```
        // 初始化配置
        config := &email.Config{}
        config.FromName = "SMG"
        config.FromAddress = "smg@sightp.com"
        config.Password = "********"
        config.Smtp = "smtp.exmail.qq.com"
        config.Port = "25"
     
        // 邮件标题与内容
        msg := email.NewHTMLMessage("test from golang", "test from golang HTML")
     
        // 将收件邮箱写在config中
        // config.To = []string{"787929691@qq.com"}
        // err := email.SendMessage(config, msg)
     
        // 增加附件
        msg.Attach("/path/to/file")
     
        // 发送到指定邮箱
        err := email.SendMessageTo(config, msg, []string{"787929691@qq.com"})
     
        fmt.Println(err)   
    ```