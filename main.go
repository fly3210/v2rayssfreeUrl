package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

/*
*
cf.0sm.com, *.cloudflare.182682.xyz
*/
var cdnUrl = "api.cloudflare.182682.xyz"

// var v2u = "aop.ssfree.ru"
var pst = "view.ssfree.ru提供免费节点"

func main() {
	port := "0.0.0.0:59399"
	//r := gin.Default()
	r := gin.New()
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

func GenerateV2(one []string, c *gin.Context, b bool) string {
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
	u := c.DefaultQuery("u", v2rayJson.Add)
	arr = append(arr, u)
	ps := c.DefaultQuery("p", pst)
	arr = append(arr, ps)
	// 替换
	v2rayJson.Host = arr[1]
	v2rayJson.Ps = arr[2]
	v2rayJson.Add = arr[0]
	// v2rayJson 转 string
	bytes, err = json.Marshal(v2rayJson)

	// 接收一个参数 是否是clash
	clash := c.DefaultQuery("is_clash", "0")
	if clash == "1" {
		return GenerateClash(v2rayJson)
	}

	// 拼接字符串
	bytes = []byte("vmess://" + base64.StdEncoding.EncodeToString(bytes))
	//println(string(bytes))
	str := base64.StdEncoding.EncodeToString(bytes)
	v2rayTime.v2 = str
	if b {
		v2rayTime.timeInt = GetTime()
	}
	return str
}

// 生成clash的节点
func GenerateClash(rayJson V2rayJson) string {
	// 读取配置文件 clash.yml
	// 打开文件
	file, err := os.Open("./clash.yml")
	if err != nil {
		return "Error opening file: " + err.Error()
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件失败")
		}
	}(file) // 确保文件被关闭

	// 读取文件的全部内容
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "Error reading file: " + err.Error()
	}

	rayJson.Type = "vmess"
	// 替换字符串 %s 根据v2rayjson 去生成
	// {name: 自用纽约, server: 新服务器地址, port: 新端口, type: vmess, uuid: 新UUID, alterId: 0, cipher: auto, tls: false, network: ws, ws-opts: {path: /新路径, headers: {Host: 新主机名}}}
	v2 := fmt.Sprintf("{name: %s, server: %s, port: %s, type: %s, uuid: %s, alterId: %s, cipher: auto, tls: true, network: %s, ws-opts: {path: %s, headers: {Host: %s}}}",
		rayJson.Ps,
		rayJson.Add,
		rayJson.Port,
		rayJson.Type,
		rayJson.Id,
		rayJson.Aid,
		rayJson.Net,
		rayJson.Path,
		rayJson.Host,
	)
	// 替换
	return fmt.Sprintf(string(contents), v2, rayJson.Ps, rayJson.Ps)
}

func GetV2rayString(c *gin.Context) string {
	if IsReadMemory() {
		//return v2rayTime.v2
		return GenerateV2(v2rayTime.vmess, c, false)
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
	return GenerateV2(one, c, true)

}
