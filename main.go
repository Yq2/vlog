package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxFileBody = 10 * 1024 * 1024

func main() {
	// 实现读取文件的handler
	fileHandler := http.FileServer(http.Dir("./video"))
	// 注册金servermux 就是将不同的URL请求交给对应的handler处理
	http.HandleFunc("/sayHello", sayHello)
	http.Handle("/video/", http.StripPrefix("/video/", fileHandler))
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/list", getFileListHandler)
	// 启动web服务
	_ = http.ListenAndServe(":8090", nil)

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("vlog 视频网站"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1 限制客户端上传视频文件的大小
	r.Body = http.MaxBytesReader(w, r.Body, maxFileBody)
	// 对指定大小的文件进行截断读取，如果错误说明文件超过大小
	err := r.ParseMultipartForm(maxFileBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 获取上传的文件
	file, fileHeader, err := r.FormFile("uploadFile")

	// 检查文件类型
	ret := strings.HasSuffix(fileHeader.Filename, ".flv")
	if !ret {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取随机名称
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	md5Str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5Str + ".flv"
	dst, err := os.Create("./video/" + newFileName)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func getFileListHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("video/*")
	var ret []string
	for _, file := range files {
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(file))
	}
	retJson, _ := json.Marshal(ret)
	_, _ = w.Write(retJson)
	return
}
