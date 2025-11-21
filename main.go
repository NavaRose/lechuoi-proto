package main

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func main() {
	generateProto("auth/auth.proto")
}

func generateProto(path string) {
	var out bytes.Buffer
	cmd := exec.Command(
		"protoc",
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		path,
	)

	// Lấy đầu ra (stdout/stderr) của FFmpeg để debug
	cmd.Stdout = &out
	//cmd.Stderr = &stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("lỗi tạo StderrPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		//log.Printf("FFmpeg stderr output:\n%s", stderr.String())
		log.Fatalf("FFmpeg thất bại: %v", err)
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("[PROTOC]", line) // In log ra console để debug

		// Kiểm tra dấu hiệu READY
		//if strings.Contains(line, *publishUrl) || strings.Contains(line, "frame=") {
		//	select {
		//	case readyChan <- true:
		//	default:
		//	}
		//	// Bạn có thể dừng phân tích log nếu muốn, nhưng để lại để xem các lỗi sau đó
		//}
		// Nếu phát hiện lỗi nghiêm trọng ngay lập tức
		if strings.Contains(line, "Could not write header") || strings.Contains(line, "Invalid data") {
			// Đóng pipe để báo hiệu lỗi
			stderr.Close()
		}
	}
}
