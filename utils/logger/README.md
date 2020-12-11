如果使用常规的`Log = zap.New()`,输出日志需可能需要`Field`类型：
```go
zap.Int("line", 47) //得到的是Field类型
zap.String("level", `{"a":"4","b":"5"}`)
Log.info(msg string, fields ...Field)
```

如果使用`Log = zap.New().Sugar()`则不需要。


