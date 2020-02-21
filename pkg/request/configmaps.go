package request

import (
	"00pf00/https-apiserver/pkg/util"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetConfigMap() {
	cert, err := tls.LoadX509KeyPair(util.CERT, util.KEY)
	if err != nil {
		fmt.Printf("load error fail !")
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}
	httpclient := &http.Client{
		Transport: tr,
	}
	request, err := http.NewRequest("GET", "https://"+util.IP+":"+util.PORT+util.EDIT_CONFIGMAP, nil)
	if err != nil {
		fmt.Printf("Get request error !")
		return
	}
	resp, err := httpclient.Do(request)
	if err != nil {
		fmt.Printf("Get resp error!")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body error,err= %v", err)
		return
	}
	fmt.Printf("body = %s\n", string(body))
	cm, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("get configmap json err\n")
		return
	}
	core, err := cm.Get("data").Get("Corefile").String()
	if err != nil {
		fmt.Printf("get corefile error err = %v", err)
		return
	}

	start := strings.Index(core, "}")
	fmt.Printf("start = %s", core[start])
	fmt.Printf("Corefile = %s\n", core)
	var hosts string
	if strings.Contains(core, "hosts") {
		hosts += core[:start+2]
		d := core[start+2:]
		e := strings.Index(core[start+2:], "}")
		hosts +=d[e+2:]
	} else {
		hosts += core[:start+2]
		//hosts+="    hosts {\n    \t127.0.0.1\tlocalhost\n    \tfallthrough\n    }\n"
		hosts += "    hosts {\n"
		hosts += "        127.0.0.1     localhost\n"
		hosts += "        fallthrough\n"
		hosts += "    }\n"
		hosts += core[start+2:]
	}

	fmt.Printf("hosts = %s\n", hosts)
	//cm.Get("data").Set("Corefile", hosts)
	cm.Get("data").Set("Corefile", hosts)
	cmbody, err := cm.MarshalJSON()
	if err != nil {
		fmt.Printf("host2json fail !")
		return
	}
	putb := bytes.NewReader(cmbody)
	put, err := http.NewRequest("PUT", "https://"+util.IP+":"+util.PORT+util.EDIT_CONFIGMAP, putb)
	presp, err := httpclient.Do(put)
	if err != nil {
		fmt.Printf("Put resp error!")
		return
	}
	p, err := ioutil.ReadAll(presp.Body)
	if err != nil {
		fmt.Printf("put error err= %v\n", err)
		return
	}
	fmt.Printf("put resp.body = %s\n", string(p))
}
