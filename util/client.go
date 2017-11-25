package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// NewClient new http client
func NewClient() *Client {
	return &Client{
		client: http.Client{
			Timeout: time.Second * 10,
		},
		transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
}

// ClientPayload 请求内容
type ClientPayload struct {
	KeepAlive bool
	Method    string
	Data      interface{}
	Userinfo  *url.Userinfo
	Header    *http.Header
}

// Client http client
type Client struct {
	client    http.Client
	transport *http.Transport
}

// Read 从 URL 读取数据
func (c *Client) Read(url string, payload *ClientPayload) (*http.Response, error) {
	if nil == payload {
		payload = &ClientPayload{
			Method: "GET",
		}
	}
	if nil == payload.Header {
		payload.Header = &http.Header{
			"User-Agent": []string{"User-Agent:Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50"},
		}
	}

	var err error
	var req *http.Request
	var buf *bytes.Buffer
	if b, ok := payload.Data.(string); ok {
		buf = bytes.NewBufferString(b)
	} else if b, ok := payload.Data.([]byte); ok {
		buf = bytes.NewBuffer(b)
	} else if b, ok := payload.Data.(*bytes.Buffer); ok {
		buf = b
	} else if nil == payload.Data {

	} else {
		return nil, errors.New("[http client] unknown payload type")
	}

	if nil == buf {
		req, err = http.NewRequest(payload.Method, url, nil)
	} else {
		if "GET" == payload.Method {
			if strings.Index(url, "?") > 0 {
				url = url + "&" + buf.String()
			} else {
				url = url + "?" + buf.String()
			}

			req, err = http.NewRequest(payload.Method, url, nil)
		} else {
			req, err = http.NewRequest(payload.Method, url, buf)
		}
	}

	if err != nil {
		return nil, err
	}

	if payload.KeepAlive {
		c.transport.DisableKeepAlives = false
		payload.Header.Add("Connection", "keep-alive")
	}
	if "POST" == payload.Method && "" == payload.Header.Get("Content-Type") {
		payload.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if nil != payload.Userinfo {
		pwd, _ := payload.Userinfo.Password()
		req.SetBasicAuth(payload.Userinfo.Username(), pwd)
	}
	if nil != payload.Header {
		req.Header = *payload.Header
	}

	c.client.Transport = c.transport

	return c.client.Do(req)
}

// GetByte 从 URL 读取字节内容及状态码
func (c *Client) GetByte(url string, payload *ClientPayload) ([]byte, *http.Response, error) {
	var resp, err = c.Read(url, payload)
	if nil == err {
		var ret []byte
		ret, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, resp, err
		}

		return ret, resp, nil
	}

	return nil, nil, err
}

// GetDoc 从 URL 创建 goquery DOM 对象
func (c *Client) GetDoc(url string, payload *ClientPayload) (*goquery.Document, error) {
	var resp, err = c.Read(url, payload)
	if nil == err {
		return goquery.NewDocumentFromReader(resp.Body)
	}

	return nil, err
}

// GetCodec 从 URL 创建反序列化对象
func (c *Client) GetCodec(url string, payload *ClientPayload, codec string, out interface{}) error {
	var err error
	var data []byte
	var resp *http.Response
	if data, resp, err = c.GetByte(url, nil); nil == err && 200 == resp.StatusCode {
		switch codec {
		case "json":
			err = json.Unmarshal(data, out)
		default:
			err = errors.New("unknown codec")
		}
	}

	return err
}

// Download 下载文件
func (c *Client) Download(url string, filename string, auto bool) error {
	var err error
	var data []byte
	var resp *http.Response
	if data, resp, err = c.GetByte(url, nil); nil == err && 200 == resp.StatusCode {
		if auto && "" == c.GetURLExt(filename) {
			filename = filename + c.GetURLExt(resp.Request.URL.String())
		}

		filename = SafeFileName(filename)
		var output, err = os.Create(filename)
		if err != nil {
			return err
		}

		defer func() {
			output.Close()
			if nil != err {
				os.Remove(filename)
			}
		}()

		_, err = output.Write(data)
	}

	return err
}

// GetURLExt 获取 URL 中文件名扩展名
func (c *Client) GetURLExt(url string) string {
	var ext string

	url = c.GetURLFilename(url)
	if pos := strings.LastIndex(url, "."); pos > 0 {
		ext = url[pos:]
	}

	return ext
}

// GetURLFilename 获取 URL 中文件名名
func (c *Client) GetURLFilename(url string) string {
	var ext string

	if pos := strings.Index(url, "?"); pos > 0 {
		url = url[pos:]
	}

	url = strings.Trim(url, "/ ")
	if pos := strings.LastIndex(url, "/"); pos > 0 {
		ext = url[(pos + 1):]
	}

	return ext
}
