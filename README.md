### 1. 项目结构规划
```
kgm2flac-backend/
├── cmd/
│   └── server/
│       └── main.go          # 程序入口点
├── internal/
│   ├── config/
│   │   └── config.go        # 配置处理
│   ├── handler/
│   │   ├── convert.go       # 文件转换处理
│   │   └── middleware.go    # 中间件
│   ├── utils/
│   │   └── utils.go         # 工具函数
│   └── service/
│       └── decrypt.go       # 解密服务
├── pkg/
│   └── types/
│       └── types.go         # 类型定义
├── go.mod                   # Go模块定义
├── go.sum                   # 依赖校验
└── config.example.yaml      # 示例配置文件
```

### 2. 构建二进制

#### 1.1 初始化Go模块
```
cd kgm2flac-backend
go mod init kgm2flac-backend
go mod tidy
go mod download
```

#### 1.2 交叉编译

```
go mod tidy

# 构建 Linux amd64 二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kgm2flac-linux-amd64 ./cmd/server

# 构建 Linux arm64 二进制
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o kgm2flac-linux-arm64 ./cmd/server

# 构建 macOS amd64 二进制（Intel）
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o kgm2flac-darwin-amd64 ./cmd/server

# 构建 macOS arm64 二进制（Apple Silicon M1/M2）
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o kgm2flac-darwin-arm64 ./cmd/server

# 构建 Windows 64位 可执行 .exe（Windows powershell）
$env:CGO_ENABLED="0"; $env:GOOS="windows"; $env:GOARCH="amd64"; go build -o kgm2flac-windows-amd64.exe ./cmd/server

# 构建 Windows 64位 可执行 .exe（Linux shell）
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o kgm2flac-windows-amd64.exe ./cmd/server
```

### 2. 运行程序

```
# 使用默认配置
./kgm2flac-linux-amd64

# 指定配置文件
./kgm2flac-linux-amd64 --config config.yaml

# 指定监听地址
./kgm2flac-linux-amd64 --addr :9090

# 指定ffmpeg路径
./kgm2flac-linux-amd64 --ffmpeg /usr/local/bin/ffmpeg

# 显示帮助
./kgm2flac-linux-amd64 --help
```
