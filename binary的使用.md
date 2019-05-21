# binary的使用


## 几种形式

```go
buff:=make([]byte,20)
buff=([]byte)("woaini")
1：设置*Buffer
buf:=new(bytes.Buffer)
err:=binary.Write(buf,binary.BigEndian,buff)
2:直接写入的方式
length:=20
binary.BigEndian.PutUint16(buff[len(buff):],legnth)


```