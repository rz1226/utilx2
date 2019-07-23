package simplerpc

import (
	//"fmt"
	//"api.tagserv/errorlog"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Rpc struct {
	Host   string
	Client *http.Client
}

func NewRpc(host string) *Rpc {
	rpc := &Rpc{}
	rpc.Host = host
	transport := &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{Transport: transport}

	rpc.Client = client
	rpc.Client.Timeout = time.Second * 3
	return rpc
}

// name对应  fasthttp路由部分， params 是post参数
func (rpc *Rpc) Call(name string, params ...string) (string, error) {
	postdata := strings.Join(params, "&")
	return rpc.call(name, postdata, 1, nil)
}

func (rpc *Rpc) call(name string, params string, retry int, lasterr error) (string, error) {
	if retry > 3 {
		return "", errors.New("retry failed with error:" + lasterr.Error())
	}
	res, err := rpc.Client.Post("http://"+rpc.Host+"/"+name,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))
	if err != nil {
		//重试三次
		//errorlog.Log(err.Error())
		return rpc.call(name, params, retry+1, err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//errorlog.Log(err.Error())
		return "", err
	}
	return string(body), nil

}
