# logger
Yapo custom logger 

# Usage
### Package import
```go
import "github.schibsted.io/Yapo/logger"
```

### Init module
```go
conf := logger.LogConfig{
  logger.SyslogConfig{
    true,             // ON | OFF syslog
    "exampleLogger"   // syslog tag
  }, 
  logger.StdlogConfig{
    true              // ON | OFF stdout log
  }
}
logger.Init(conf)
``` 

### Set Log Level
```go
logger.SetLogLevel(logger.[DEBUG|INFO|WARN|ERROR|CRIT])
```

### Do log
```go
logger.[Debug|Info|Warn|Error|Crit]("This is log")
```

### Close logger
user should call this before ending execution so that the logger frees syslog connection 
```go
logger.CloseSyslog()
```

# Example file
for further reference check the example file in the example folder. at example directory, do `go build` and run `./example`
