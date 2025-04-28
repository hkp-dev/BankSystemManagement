package utility

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveImage(base64Str, dir, name string) (string, error) {
	data := base64Str
	if idx := strings.Index(base64Str, ","); idx != -1 {
		data = base64Str[idx+1:]
	}

	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("lỗi giải mã Base64: %v", err)
	}

	ext := ".jpg"
	fileName := fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, fileName)

	// Lưu file
	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return "", fmt.Errorf("lỗi lưu file: %v", err)
	}

	return filePath, nil
}
