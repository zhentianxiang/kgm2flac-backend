package main

import (
	"flag"
	"fmt"
	"kgm2flac-backend/internal/config"
	"kgm2flac-backend/internal/handler"
	"log"
	"runtime"
)

// 编译时通过 -ldflags 注入
var (
	version    = "dev"     // 版本号
	buildDate  = "unknown" // 编译时间
	commitHash = "unknown" // Git 提交哈希
	appEnv     = "unknown" // 运行环境 (dev/pre/prod)
)

func main() {
	// 命令行参数解析
	configPath := flag.String("config", "", "配置文件路径")
	showHelp := flag.Bool("help", false, "显示帮助信息")
	showVersion := flag.Bool("version", false, "显示版本信息")
	showEnv := flag.Bool("env", false, "显示当前运行环境")
	addr := flag.String("addr", ":8080", "服务器监听地址")
	ffmpegBin := flag.String("ffmpeg", "ffmpeg", "ffmpeg可执行文件路径")

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	if *showVersion {
		printVersion()
		return
	}

	if *showEnv {
		printEnv()
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(*configPath, *addr, *ffmpegBin)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 启动服务器
	if err := handler.StartServer(cfg); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func printHelp() {
	fmt.Println("KGM to FLAC 转换服务")
	fmt.Println("用法: server [选项]")
	fmt.Println()
	fmt.Println("选项:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  server --config config.yaml --addr :8080")
	fmt.Println("  server --ffmpeg /usr/local/bin/ffmpeg")
	fmt.Println("  server --version")
	fmt.Println("  server --env")
}

func printVersion() {
	fmt.Printf("KGM to FLAC 转换服务\n")
	fmt.Printf("版本: %s\n", version)
	//fmt.Printf("环境: %s\n", appEnv)
	fmt.Printf("编译日期: %s\n", buildDate)
	fmt.Printf("Git提交: %s\n", commitHash)
	fmt.Printf("Go 版本: %s\n", runtime.Version())
	fmt.Printf("操作系统/架构: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func printEnv() {
	fmt.Printf("当前运行环境: %s\n", appEnv)
}
