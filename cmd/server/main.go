package main

import (
	"flag"
	"fmt"
	"kgm2flac-backend/internal/config"
	"kgm2flac-backend/internal/handler"
	"log"
)

func main() {
	// 命令行参数解析
	configPath := flag.String("config", "", "配置文件路径")
	showHelp := flag.Bool("help", false, "显示帮助信息")
	addr := flag.String("addr", ":8080", "服务器监听地址")
	ffmpegBin := flag.String("ffmpeg", "ffmpeg", "ffmpeg可执行文件路径")

	flag.Parse()

	if *showHelp {
		printHelp()
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
}
