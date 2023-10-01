// Package captureutil 验证码识别模块
package captureutil

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/otiai10/gosseract/v2"

	"GoCampusLogin/utils/httputil"
)

// GetCaptureCode 获取验证码
func GetCaptureCode() string {
	url := "https://wxjw.bistu.edu.cn/authserver/getCaptcha.htl?" + strconv.Itoa(int(time.Now().UnixMilli()))
	resp, err := httputil.Get(url, nil)
	if err != nil {
		log.Fatalf("无法下载图片: %v", err)
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("无法读取图片数据: %v", err)
	}

	// 图片去噪
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalf("image decode err: %v", err)
		return ""
	}
	img = effect.Grayscale(img)
	img = effect.Dilate(img, 0.3)
	img = effect.Erode(img, 0.2)

	// 测试：将图片写至本地
	//imaging.Save(img, "captureutil.png")

	client := gosseract.NewClient()
	defer client.Close()
	// 将 img 转回 []byte，提供 tesseract 识别
	var buf bytes.Buffer
	imgio.PNGEncoder()(&buf, img)
	client.SetImageFromBytes(buf.Bytes())
	code, err := client.Text()
	if err != nil {
		log.Fatalf("无法识别验证码: %v", err)
	}
	// 去除任何非单词字符
	return regexp.MustCompile("\\W").ReplaceAllString(code, "")
}

// IsNeedCaptcha 判断是否需要验证码
func IsNeedCaptcha(username string) bool {
	resp, err := httputil.Get("https://wxjw.bistu.edu.cn/authserver/checkNeedCaptcha.htl?username="+username, nil)
	if err != nil {
		log.Fatalf("get isNeedCapture err: %v", err)
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("readAll Body err: %v", err)
		return false
	}
	if strings.Contains(string(body), "true") {
		return true
	}
	return false
}
