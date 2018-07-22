package cli

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	//"strconv"
	"crypto/tls"
	"math"
	"strings"
	"time"
)

const (
	Retry = 5
	UA    = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

// 包裹
type Package struct {
	GlobalExpressCode string
	CNExpressCode     string
	CNPostCode        string
	CNStatus          bool
	NLStatus          bool
	CompanyCode       string
	Weight            float64
}

func (p *Package) TraceNL() (err error) {

	cacheFile := fmt.Sprintf("%s/%s/%s", Conf.CacheDir, "postnl", p.GlobalExpressCode)
	cacheDir := fmt.Sprintf("%s/%s", Conf.CacheDir, "postnl")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.Mkdir(cacheDir, os.ModeDir); err != nil {
			return err
		}
	}

	getBody := func() (body []byte, err error) {

		if _, err = os.Stat(cacheFile); err == nil {
			body, _ = ioutil.ReadFile(cacheFile)
			log.Info("read cache " + cacheFile)
			return
		}

		BaseUrl := "https://www.internationalparceltracking.com/api/shipment"

		// 接口URL
		url := fmt.Sprintf("%s?barcode=%s&country=CN&language=zh&postalCode=%s", BaseUrl, p.GlobalExpressCode, p.CNPostCode)
		//log.Info(url)

		// 初始化
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)

		// 设置UA
		req.Header.Set("User-Agent", UA)

		// 访问
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("access failed")
			return
		}

		// 获取内容
		ioutil.WriteFile(cacheFile, body, 0644)
		log.Info("write cache " + cacheFile)
		return ioutil.ReadAll(resp.Body)
	}

	body, err := getBody()
	if err != nil {
		return
	}

	weight := gjson.GetBytes(body, "measurements.weight")
	p.Weight = weight.Float()
	p.NLStatus = true
	log.Info(weight)
	return
}

func (p *Package) TraceCN() (err error) {
	endpoint := "https://poll.kuaidi100.com/poll/query.do"
	config := Conf.Kuaidi100Config
	cacheFile := fmt.Sprintf("%s/%s/%s", Conf.CacheDir, "kuaidi100", p.CNExpressCode)
	cacheDir := fmt.Sprintf("%s/%s", Conf.CacheDir, "kuaidi100")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.Mkdir(cacheDir, os.ModeDir); err != nil {
			return err
		}
	}

	getSign := func(data string) (sign string) {
		has := md5.Sum([]byte(data))
		sign = fmt.Sprintf("%x", has)
		return
	}

	getBody := func(cacheFile string) (body []byte, err error) {

		if _, err = os.Stat(cacheFile); err == nil {
			body, _ = ioutil.ReadFile(cacheFile)
			log.Info("read cache " + cacheFile)
			return
		}

		type Param struct {
			Com      string      `json:"com"`
			Num      string      `json:"num"`
			From     string      `json:"from"`
			To       interface{} `json:"to"`
			Resultv2 int         `json:"resultv2"`
		}
		orderParam := Param{
			Com:      "postnl",
			Num:      p.CNExpressCode,
			From:     "",
			To:       "",
			Resultv2: 0,
		}

		paramJson, err := json.Marshal(orderParam)

		u, _ := url.Parse(endpoint)
		param := url.Values{}
		param.Set("customer", config.Customer)
		//param.Set("param", string(paramJson))
		param.Set("sign", strings.ToUpper(getSign(string(paramJson)+config.Key+config.Customer)))

		u.RawQuery = param.Encode()

		tr := &http.Transport{ //解决x509: certificate signed by unknown authority
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		nt := &http.Client{
			Timeout:   time.Duration(10 * time.Second),
			Transport: tr,
		}
		req, _ := http.NewRequest("GET", u.String()+"&param="+string(paramJson), nil)
		resp, err := nt.Do(req)
		if err != nil {
			return
		}

		// status code check
		log.Info(resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("http code: %d", resp.StatusCode)
			return
		}

		body, _ = ioutil.ReadAll(resp.Body)
		err = ioutil.WriteFile(cacheFile, body, 0644)
		log.Info("write cache " + cacheFile)
		return
	}

	retry := Retry
	var body []byte
	for retry > 0 {
		body, err = getBody(cacheFile)
		if err != nil {
			retry--
			log.Info(fmt.Sprintf("%s retry %d", cacheFile, Retry-retry))
			time.Sleep(time.Second * time.Duration(math.Pow(5, float64(Retry-retry))))
			if retry == 0 {
				return
			}
			continue
		}
		break
	}

	//log.Info(string(body))
	message := gjson.GetBytes(body, "message")
	//log.Info(message.String())
	if message.String() == "ok" {
		p.CNStatus = true
	}
	return
}
