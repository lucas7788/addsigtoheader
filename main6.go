package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"os"
	"os/signal"
	"syscall"
	core2 "github.com/CandyDrop/core"
	restful2 "github.com/CandyDrop/restful/restful"
	"sync"
)

func main() {

	restStudy()
	//waitToExit()

}

type RestfulReq struct {
	Action  string
	Version string
	Type    int
	Data    string
}
type Args struct {
	Ontid   string
	Address string
}

type Num struct {
	num int
	lock sync.Locker
}

func (self *Num) Add1() {
	self.lock.Lock()
	self.num+=1
	self.lock.Unlock()
}
func (self *Num) GetNum() int{
	return self.num
}
func restStudy() {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   5,
			DisableKeepAlives:     true, //enable keepalive
			IdleConnTimeout:       time.Second * 300,
			ResponseHeaderTimeout: time.Second * 300,
		},
	}
	args := core2.Param{
		Ontid:   "did:ont:TX4hBHqE4FG94B4B9RQCSp1jLDvdmto8a3",
		Address: "ALz86Xe1FR6oJ4xJW6PKfNLj88cMEPKDuV",
	}
	argsBytes,_ := json.Marshal(args)
	restReq := &RestfulReq{
		Action:  restful2.POST_AIRDROP,
		Version: "1",
		Data:    hex.EncodeToString(argsBytes),
	}
	reqData, err := json.Marshal(restReq)
	if err != nil {
		return
	}
	reqUrl, err := new(url.URL).Parse("http://127.0.0.1:30334")
	if err != nil {
		fmt.Println("err:", err)
	}
	reqUrl.Path = restful2.POST_AIRDROP
	num := &Num{
		num:0,
		lock:new(sync.Mutex),
	}
	var w *sync.WaitGroup
	w = new(sync.WaitGroup)
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			resp, err := client.Post(reqUrl.String(), "application/json", bytes.NewReader(reqData))
			if err != nil {
				return
			}
			defer resp.Body.Close()
			fmt.Println(i)
			num.Add1()
			//data, err := ioutil.ReadAll(resp.Body)
			//if err != nil {
			//	return
			//}
			//restRsp := &RestfulResp{}
			//err = json.Unmarshal(data, restRsp)
			//fmt.Println("resp:", restRsp)
			//fmt.Println("resp:", string(restRsp.Result))
		}(i)
	}
	w.Wait()
	fmt.Println("num:",num.GetNum())
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			fmt.Println("Ontology received exit signal:", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}


func time_study() {
	t1 := "2019-01-08 13:50:30"
	timeTemplate1 := "2006-01-02 15:04:05"                          //常规类型
	stamp, _ := time.ParseInLocation(timeTemplate1, t1, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	fmt.Println(stamp.Unix())                                       //输出：1546926630

	fmt.Println(time.Now().Unix())
	fmt.Println()
}
