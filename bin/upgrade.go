package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/ixre/gof/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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

func doUpdate(force bool){
	v,err := checkVersion()
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if v == nil{
		fmt.Println("已经是最新版本")
		return
	}
	printVersion(v)
	if !force {
		fmt.Printf("\n是否现在立即更新?[Y/N] ")
		input := bufio.NewScanner(os.Stdin) //初始化一个扫表对象
		input.Scan()
		if strings.ToLower(input.Text()) != "y"{
			os.Exit(0)
		}
	}
	tmpFile := fmt.Sprintf("%s/tto-release-%s.tar.gz",os.TempDir(),v.version)
	err = prepareFiles(v.distURL,tmpFile)
	if err != nil{
		fmt.Println("下载失败:",err.Error())
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout,"安装中...\r")
	err = install(tmpFile)
	if err != nil{
		fmt.Fprint(os.Stdout,"安装失败:",err.Error())
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout,"恭喜! 安装完成！\r")
}

func install(file string) error {
	tmpDir := filepath.Dir(file)+"/tto"
	err := decompressTarFile(file,tmpDir)
	dstPath := getProgramPath()
	switch runtime.GOOS{
	case "linux":
		 err = overwriteFile(tmpDir+"/tto",dstPath+"/tto")
	case "windows":
		err = overwriteFile(tmpDir+"/tto.exe",dstPath+"/tto.exe")
	case "darwin":
		err = overwriteFile(tmpDir+"/tto-mac",dstPath+"/tto")
	}
	if err != nil{
		log.Println(err.Error())
	}
	return err
}

// 删除原程序,替换为新程序
func overwriteFile(src string, dst string)error {
	sf, _ := os.Open(src)
	_ = os.Remove(dst)
	df, err := os.Create(dst)
	df.Chmod(755)
	_, err = io.Copy(df, sf)
	return err
}

// 解压文件
func decompressTarFile(file string,dstDir string) error {
	srcFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dstDir +"/"+ hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		if file != nil{
			io.Copy(file, tr)
		}
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	dirPath := filepath.Dir(name)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}
	if dirPath != name[:len(name)-1] {
		return os.Create(name)
	}
	return nil,nil
}

func prepareFiles(distURL string,file string)error {
	//lastLen := 0
	return down(distURL,file,func(total int,reads int,time int){
		prg := int(float32(reads)/float32(total)*100)
		if prg > 0 {
			line := "已下载 "+strconv.Itoa(prg)+"% 速度："+
				strconv.Itoa(reads/time/1000)+"kb/s \r"
			fmt.Fprint(os.Stdout,line)
		}
	},-1)
}

func printVersion(v *version) {
	line := fmt.Sprintf("**  new version v%s  **",v.version)
	lineFill := strings.Repeat("-",len(line))
	fmt.Println(lineFill)
	fmt.Println(line)
	fmt.Println(lineFill)
	fmt.Println(v.remark)
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
	// 获取版本更新信息
	blockR := regexp.MustCompile("release-entry[\\s\\S]+?releases/tag/v"+
		ver.version+"[\\s\\S]+?markdown-body\">([\\s\\S]+?)</div>")
	blockMatches := blockR.FindStringSubmatch(releases)
	remark := regexp.MustCompile("<[^>]+>").ReplaceAllString(blockMatches[1],"")
	ver.remark = strings.TrimSpace(remark)
	return ver,nil
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


func down(ur, target string,onProgress func(total int,reads int,seconds int), timeout int)error {
	// 创建目录
	dir := filepath.Dir(target)
	if _,err := os.Stat(dir);os.IsNotExist(err){
		if err := os.MkdirAll(dir, os.ModePerm);err != nil{
			return err
		}
	}
	var size int64
	file, err := os.OpenFile(target, os.O_RDWR, os.ModePerm)
	if err == nil {
		// 计算断点位置
		stat, _ := file.Stat()
		size = stat.Size()
		sk, err := file.Seek(size, 0)
		if err != nil {
			_ = file.Close()
			return err
		}
		if sk != size {
			_ = file.Close()
			return errors.New(fmt.Sprintf("%s seek length not equal file size,"+
				"seek=%d,size=%d\n", target, sk, size))
		}
	 }
	 if os.IsNotExist(err){
		file, err = os.Create(target)
		if err != nil {
			return err
		}
	}
	client := &http.Client{}
	if timeout >0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	request := http.Request{}
	request.Method = http.MethodGet
	if size != 0 {
		header := http.Header{}
		header.Set("Range", "bytes="+strconv.FormatInt(size, 10)+"-")
		request.Header = header
	}
	parse, err := url.Parse(ur)
	if err != nil {
		return err
	}
	request.URL = parse
	get, err := client.Do(&request)
	if err != nil {
		return err
	}
	defer func() {
		err := get.Body.Close()
		if err != nil {
			fmt.Println(target, "body close:", err.Error())
		}
		err = file.Close()
		if err != nil {
			fmt.Println(target, "file close:", err.Error())
		}
	}()
	if get.ContentLength == 0 {
		return errors.New("already downloaded")
	}
	reads := 0
	total := int(get.ContentLength)
	begin := time.Now().Unix()

	body := get.Body
	writer := bufio.NewWriter(file)
	bs := make([]byte, 10*1024*1024) //每次读取的最大字节数，不可为0

	for {
		var read int
		read, err = body.Read(bs)
		reads += read
		onProgress(total,reads,int(time.Now().Unix()-begin+1))
		if err != nil {
			if err != io.EOF {
				fmt.Println(target, "read err:"+err.Error())
			} else {
				err = nil
			}
			break
		}
		_, err = writer.Write(bs[:read])
		if err != nil {
			fmt.Println(target, "write err:"+err.Error())
			break
		}
	}
	if err == nil {
		err = writer.Flush()
	}
	return err
}

func getProgramPath()string{
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return path.Join(path.Dir(filename), "") // the the main function file directory
	}
	return "./"
}