package main

import (
	"fmt"
	"github.com/ixre/gof/log"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : upgrade.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-27 06:24
 * description :
 * history :
 */

var versionInfoURL = "https://github.com/ixre/tto/releases/"

type version struct{
	version string
	remark  string
	distURL string
}

func doUpdate(){
	v,err := checkVersion()
	if err != nil{
		log.Println(err.Error())
		os.Exit(1)
	}
	if v == nil{
		log.Println("已经是最新版本")
		return
	}
	log.Println("%#v",v)

}

func checkVersion()(*version,error){
	releases, err := getReleases()
	if err != nil{
		return nil,err
	}
	var ver *version
	verReg := regexp.MustCompile("\"(.+?/releases/download/v(.+?)/.+?)\"")
	submatch := verReg.FindAllStringSubmatch(releases, -1)
	if len(submatch)==0{
		return nil,nil
	}
	ver = &version{
		version: submatch[0][2],
		remark:  "",
		distURL: "https://github.com"+submatch[0][1],
	}
	blockR := regexp.MustCompile("release-entry[\\s\\S]+?releases/tag/v"+"[\\s\\S]+?markdown-body\">")
	blockMatches := blockR.FindAllStringSubmatch(releases,-1)
	log.Println(len(blockMatches))
	log.Println(fmt.Sprintf("%#v",ver))

	return nil,nil
}

func getReleases()(string,error) {
	cli := http.Client{}
	req, _ := http.NewRequest("GET", versionInfoURL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0")
	rsp, err := cli.Do(req)
	if err == nil {
		bytes, _ := ioutil.ReadAll(rsp.Body)
		return string(bytes), nil
	}
	return "",nil
}