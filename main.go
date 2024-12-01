package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

/*
*
cf.0sm.com, *.cloudflare.182682.xyz
*/
//var cdnUrl = "api.cloudflare.182682.xyz"
var nodeUrl = "https://1208.fly321.pw/"

func main() {
	port := "0.0.0.0:59399"
	//r := gin.Default()
	r := gin.New()
	r.GET("/", GetYaml)
	r.Run(port)
}

type CacheYaml struct {
	yamlStr map[string]string
	// 加一个时间戳
	timeStamp int64
}

var cacheYamlMap = CacheYaml{yamlStr: make(map[string]string), timeStamp: time.Now().Unix()}

func GetYaml(c *gin.Context) {
	var isOld = c.Query("old")
	var year, month, day = GetYmd(isOld)
	if (cacheYamlMap.timeStamp < time.Now().Unix()-600) || cacheYamlMap.yamlStr[isOld] == "" {
		var url = fmt.Sprintf("%s/%d/%02d/%d%02d%02d.yaml", nodeUrl, year, month, year, month, day)
		curl, err := Curl(url)
		if err != nil {
			c.String(http.StatusInternalServerError, "获取yaml失败")
		}
		fmt.Println("curl请求")
		cacheYamlMap.yamlStr[isOld] = curl
		c.String(http.StatusOK, curl)
	}
	// 如果缓存中有，直接返回
	fmt.Println("缓存中有")
	c.String(http.StatusOK, cacheYamlMap.yamlStr[isOld])
}

// Curl 发送GET请求到指定的URL并返回响应体
func Curl(url string) (string, error) {
	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err // 如果请求失败，返回错误
	}
	defer resp.Body.Close() // 确保在函数返回时关闭响应体

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err // 如果读取失败，返回错误
	}

	return string(body), nil // 返回响应体的字符串表示
}

func GetYmd(old string) (year int, month time.Month, day int) {
	if old == "1" {
		// 昨天的Ymd
		t := time.Now().AddDate(0, 0, -1)
		year, month, day = t.Year(), t.Month(), t.Day()
	} else {
		// 今天的Ymd
		t := time.Now()
		year, month, day = t.Year(), t.Month(), t.Day()
	}
	return
}
