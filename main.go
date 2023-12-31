package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var cdnUrl = "cf.0sm.com"
var v2u = "aop.ssfree.ru"
var pst = "isun cdn 节点"

func main() {
	port := "0.0.0.0:59399"
	r := gin.Default()
	r.GET("/", helloWord)
	r.Run(port)
}

func helloWord(c *gin.Context) {
	//writer.Write([]byte(GetV2rayString()))
	c.Writer.Write([]byte(GetV2rayString(c)))
}

type V2rayTime struct {
	v2      string
	timeInt int
	vmess   []string
}

var v2rayTime = V2rayTime{
	v2:      "",
	timeInt: 0,
	vmess:   []string{},
}

type V2rayJson struct {
	Host string `json:"host"`
	Path string `json:"path"`
	Port string `json:"port"`
	Tls  string `json:"tls"`
	Ps   string `json:"ps"`
	Id   string `json:"id"`
	Add  string `json:"add"`
	V    string `json:"v"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
}

// 获取当前时间戳
func GetTime() int {
	return int(time.Now().Unix())
}

// 是否从内存读取
func IsReadMemory() bool {
	// 时间戳转换为时间
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	timeStr := time.Unix(int64(v2rayTime.timeInt), 0).Format(timeLayout)
	timeStr1 := time.Unix(int64(GetTime()), 0).Format(timeLayout)

	fmt.Printf("当前时间：%s,上次时间：%s\r\n", timeStr1, timeStr)
	//fmt.Printf("v2ary:%s", v2rayTime.v2)
	if v2rayTime.timeInt == 0 {
		fmt.Println("超过120秒，重新获取")
		return false
	}
	if GetTime()-v2rayTime.timeInt > 120 {
		fmt.Println("超过120秒，重新获取")
		return false
	}
	fmt.Println("从内存读取")
	return true
}

func GenerateV2(one []string, c *gin.Context) string {
	bytes := []byte(one[1])
	// 替换vmess://为""
	bytes = bytes[8:]
	// base64解码
	bytes, err := base64.StdEncoding.DecodeString(string(bytes))
	println(string(bytes))
	// 字符串转json
	var v2rayJson V2rayJson
	err = json.Unmarshal(bytes, &v2rayJson)
	if err != nil {
		fmt.Println(err)
	}
	// 接收get请求参数string
	cdn := c.DefaultQuery("c", cdnUrl)
	// 定义个数组
	var arr []string
	// 追加
	arr = append(arr, cdn)
	u := c.DefaultQuery("u", v2u)
	arr = append(arr, u)
	ps := c.DefaultQuery("p", pst)
	arr = append(arr, ps)
	// 替换
	v2rayJson.Host = arr[1]
	v2rayJson.Ps = arr[2]
	v2rayJson.Add = arr[0]
	// v2rayJson 转 string
	bytes, err = json.Marshal(v2rayJson)
	// 拼接字符串
	bytes = []byte("vmess://" + base64.StdEncoding.EncodeToString(bytes))
	//println(string(bytes))
	str := base64.StdEncoding.EncodeToString(bytes)
	v2rayTime.v2 = str
	v2rayTime.timeInt = GetTime()
	return str
}

func GetV2rayString(c *gin.Context) string {
	if IsReadMemory() {
		//return v2rayTime.v2
		return GenerateV2(v2rayTime.vmess, c)
	}
	url := "https://view.ssfree.ru/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
	}
	// 正则匹配
	regxpstr := `var vmess = \"(.*?)\"`
	compile, err := regexp.Compile(regxpstr)
	if err != nil {
		fmt.Println(err)
	}
	// 查找符合正则的第一个
	one := compile.FindStringSubmatch(string(body))
	v2rayTime.vmess = one
	return GenerateV2(one, c)

}
