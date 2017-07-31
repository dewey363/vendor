package tgo

//验证签名

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"fmt"
)

var (
	SignAppSecretKey string
)

/*
sha1的签名算法
   appsecret = "std::string"
   signature = sha1(appsecret+"babybirthday=1457578839&city=南京市&mobile=15324893018&province=江苏省
&signtimestamp=1457578839&username=张三"+appsecret)

1. 签名: signature, 时间戳: timestamp
2. 参数列表按参数Key字典序升序排列
3. 编码使用 UTF-8  //正式开始的时候SignSwitch配置=0
*/

func UtilSignCheckSign(c *gin.Context) bool {
	SignSwitch := ConfigAppGetString("SignSwitch", "1")
	if SignSwitch == "0" {
		return true
	}
	SignAppSecretKey = ConfigAppGetString("AppSecretKey", "")
	ps := UtilRequestGetAllParams(c)
	if !UtilSignCheckSignTimestamp(c.Request) {
		return false
	}
	signTimestamp, _ := c.Cookie("signtimestamp")
	fmt.Println("时间戳"+signTimestamp)
	ps["signtimestamp"] = []string{signTimestamp}
	sortedParams := UtilSignGetSortUpParamsString(ps)
	signString := SignAppSecretKey + sortedParams + SignAppSecretKey
	signature := UtilCryptoSha1(signString)
	//c.SetCookie("signature",signature,1000000,"/","192.168.1.106",false,false)
	signCookie, _ := c.Cookie("signature")
	fmt.Println("加密串"+signCookie)
	if signature == signCookie {
		return true
	}
	return false
}

//升序排序的参数拼接的字符串
func UtilSignGetSortUpParamsString(ps url.Values) string {
	psKey := []string{}
	for k, _ := range ps {
		psKey = append(psKey, k)
	}
	sort.Strings(psKey)
	ret := []string{}
	for _, v := range psKey {
		ret = append(ret, v+"="+ps.Get(v))
	}
	return strings.Join(ret, "&")
}

func UtilSignCheckSignTimestamp(req *http.Request) bool {
	appLimitTime, _ := strconv.Atoi(ConfigAppGetString("AppAccessLimitTime", ""))
	if appLimitTime == 0 {
		return true
	}
	ts, err := req.Cookie("signtimestamp")

	if err != nil {
		return false
	}
	signTimestamp, err := strconv.Atoi(ts.Value)

	if err != nil {
		return false
	}
	now := time.Now().Unix()
	if now < int64(appLimitTime+signTimestamp) {
		return true
	}

	return false
}
