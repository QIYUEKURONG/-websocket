# log使用

## 设置flag

自定义的选项常量

```go
const(
Ldate         = 1 << iota     //日期示例： 2009/01/23
Ltime                         //时间示例: 01:23:23
Lmicroseconds                 //毫秒示例: 01:23:23.123123.
Llongfile                     //绝对路径和行号: /a/b/c/d.go:23
Lshortfile                    //文件和行号: d.go:23.
LUTC                          //日期时间转为0时区的
LstdFlags     = Ldate | Ltime //Go提供的标准抬头信息
)
```

例子：

```Go
log.SetPrefix("[user]") //输出数据的时候，就会在前面出现这几个字符。
log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

```
log除了Print系列的函数外，还有Fatal以及panic系列的函数，其中Fatal表示程序遇到了致命的错误，需要退出，这时候使用Fatal记录日志后，然后程序退出，也就是说Fatal相当于先调用Print打印日志，然后再调用os.Exit(1)退出程序。


