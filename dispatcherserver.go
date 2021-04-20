package main

import (
	"encoding/json"
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolAccounting/service-register/types"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	httpdaemon "github.com/NpoolRD/http-daemon"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	_ "strings"
	_ "time"
)

// etcd key
const accountingDomain = "accounting.npool.top"

// etcd value
type IpServerConfig struct {
	Ip   string
	Port string
}

type DispatcherConfig struct {
	Port int
}

type RegisterServer struct {
	config DispatcherConfig
}

func NewRegisterServer(configFile string) *RegisterServer {

	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot read file %v: %v", configFile, err)
		return nil
	}

	config := DispatcherConfig{}
	err = json.Unmarshal(buf, &config)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot parse file %v: %v", configFile, err)
		return nil
	}

	server := &RegisterServer{
		config: config,
	}

	log.Infof(log.Fields{}, "successful to create devops server")

	return server
}

func (s *RegisterServer) Run() error {

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.GetMinerPledgeAPI,
		Method:   "GET",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.GeMinerPledgeRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.GetMinerDailyRewardAPI,
		Method:   "GET",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.GetMinerDailyRewardRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.GetAccountInfoAPI,
		Method:   "GET",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.GetAccountInfoRequest(w, req)
		},
	})

	log.Infof(log.Fields{}, "start http daemon at %v", s.config.Port)
	httpdaemon.Run(s.config.Port)
	return nil
}

// minerPledgeInput
type MinerPledgeRequestInput struct {
	account string `gorm:"column:Account"`
}

// 获取抵押接口
func (s *RegisterServer) GeMinerPledgeRequest(writer http.ResponseWriter, request *http.Request) (interface{}, string, int) {
	// 接收请求参数
	query := request.URL.Query()
	// 账户
	id := string(query["account"][0])
	if id == "" {
		return nil, "account is must", -1
	}
	var host = ""
	// 获取Ip数组
	resp, err := etcdcli.Get(accountingDomain)
	if resp == nil {
		return nil, err.Error(), -1
	}
	var strs = ""
	// TODO 分发服务
	for i, v := range resp {
		if 0 < i {
			strs = fmt.Sprintf("%v,", strs)
		}
		var input = types.ServiceRegisterOutput{}
		fmt.Println("string(v):", string(v)[1:])
		err := json.Unmarshal([]byte(string(v)[1:]), &input)
		if err == nil {
			host = input.IP + ":" + input.Port
			break
		}
	}
	if host == "" {
		return nil, "host is null", -1
	}
	resps, err := httpdaemon.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("http://%v%v", host, types.GetMinerPledgeAPI) + "?account=" + id)
	if err != nil {
		log.Errorf(log.Fields{}, "heartbeat error: %v", err)
		return nil, err.Error(), -1
	}

	if resps.StatusCode() != 200 {
		return nil, xerrors.Errorf("NON-200 return").Error(), -1
	}

	apiResp, err := httpdaemon.ParseResponse(resps)
	if err != nil {
		return nil, err.Error(), -1
	}
	return apiResp.Body, "success", 0
}

// 某个时间段收益统计
func (s *RegisterServer) GetMinerDailyRewardRequest(writer http.ResponseWriter, request *http.Request) (interface{}, string, int) {
	// 接收请求参数
	query := request.URL.Query()
	// 账户
	account := string(query["account"][0])
	// 初始高度0 出块时间 1598306400 后面每一个加30s 北京时间戳
	startTime := string(query["startTime"][0])
	endTime := string(query["endTime"][0])

	if account == "" {
		return nil, "account is must", -1
	}
	var host string

	// 获取Ip数组
	resp, err := etcdcli.Get(accountingDomain)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot get %v: %v", accountingDomain, err)
		return nil, err.Error(), -1
	}
	var strs = ""
	// TODO 分发服务
	for i, v := range resp {
		if 0 < i {
			strs = fmt.Sprintf("%v,", strs)
		}
		var input = types.ServiceRegisterOutput{}
		fmt.Println("string(v):", string(v)[1:])
		err := json.Unmarshal([]byte(string(v)[1:]), &input)
		if err == nil {
			host = input.IP + ":" + input.Port
			break
		}
	}
	if host == "" {
		return nil, "host is null", -1
	}
	resps, err := httpdaemon.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("http://%v%v", host, types.GetMinerPledgeAPI) + "?account=" + account + "&startTime=" + startTime + "&endTime=" + endTime)
	if err != nil {
		log.Errorf(log.Fields{}, "heartbeat error: %v", err)
		return nil, err.Error(), -1
	}

	apiResp, err := httpdaemon.ParseResponse(resps)
	if err != nil {
		return nil, err.Error(), -1
	}
	return apiResp.Body, "success", 0
}

// accountInfo
func (s *RegisterServer) GetAccountInfoRequest(w http.ResponseWriter, request *http.Request) (interface{}, string, int) {
	// 接收请求参数
	query := request.URL.Query()
	// 账户
	account := string(query["account"][0])
	// 初始高度0 出块时间 1598306400 后面每一个加30s 北京时间戳
	startTime := string(query["startTime"][0])
	endTime := string(query["endTime"][0])
	pageSize := string(query["pageSize"][0]) // 分页大小

	currentPage := string(query["currentPage"][0]) // 当前第几页

	if account == "" {
		return nil, "account is must", -3
	}

	var host string

	// 获取Ip数组
	resp, err := etcdcli.Get(accountingDomain)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot get %v: %v", accountingDomain, err)
		return nil, err.Error(), -1
	}
	var strs = ""
	// TODO 分发服务
	for i, v := range resp {
		if 0 < i {
			strs = fmt.Sprintf("%v,", strs)
		}
		var input = types.ServiceRegisterOutput{}
		fmt.Println("string(v):", string(v)[1:])
		err := json.Unmarshal([]byte(string(v)[1:]), &input)
		if err == nil {
			host = input.IP + ":" + input.Port
			break
		}
	}
	if host == "" {
		return nil, "host is null", -1
	}
	resps, err := httpdaemon.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("http://%v%v", host, types.GetMinerPledgeAPI) + "?account=" + account + "&startTime=" + startTime + "&endTime=" + endTime + "&pageSize=" + pageSize + "&currentPage=" + currentPage)
	if err != nil {
		log.Errorf(log.Fields{}, "heartbeat error: %v", err)
		return nil, err.Error(), -1
	}

	if resps.StatusCode() != 200 {
		return nil, xerrors.Errorf("NON-200 return").Error(), -1
	}

	apiResp, err := httpdaemon.ParseResponse(resps)
	if err != nil {
		return nil, err.Error(), -1
	}
	return apiResp.Body, err.Error(), 0
}
