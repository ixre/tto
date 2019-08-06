package tto

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// 生成Go仓储代码
func (s *Session) GenerateGoRepoCodes(tables []*Table, targetDir string) (err error) {
	wg := sync.WaitGroup{}
	for _, table := range tables {
		wg.Add(1)
		go func(wg *sync.WaitGroup, tb *Table) {
			defer wg.Done()
			//生成实体
			str, path := s.tableToGoStruct(tb)
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
			//生成仓储结构
			str, path = s.tableToGoRepo(tb, true, "")
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
			//生成仓储接口
			str, path = s.tableToGoIRepo(tb, true, "")
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
		}(&wg, table)
	}
	wg.Wait()
	// 生成仓储工厂
	code := s.GenerateCodeByTables(tables, GoRepoFactoryTemplate)
	path, _ := s.PredefineTargetPath(GoRepoFactoryTemplate, nil)
	return SaveFile(code, targetDir+"/"+path)
}

// 遍历模板文件夹, 并生成代码
func (s *Session) WalkGenerateCode(tables []*Table, tplDir string, outputDir string) error {
	tplMap := map[string]*CodeTemplate{}
	sliceSize := len(tplDir) - 1
	if tplDir[sliceSize] == '/' {
		tplDir = tplDir + "/"
		sliceSize += 1
	}
	err := filepath.Walk(tplDir, func(path string, info os.FileInfo, err error) error {
		// 如果模板名称以"_"开头，则忽略
		if info != nil && !info.IsDir() && info.Name()[0] != '_' {
			tp, err := s.ParseTemplate(path)
			if err != nil {
				return errors.New("template:" + info.Name() + "-" + err.Error())
			}
			s.Resolve(tp)
			tplMap[path[sliceSize:]] = tp
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(tplMap) == 0 {
		return errors.New("no any code template")
	}
	wg := sync.WaitGroup{}
	for _, tb := range tables {
		for path, tpl := range tplMap {
			wg.Add(1)
			go func(wg *sync.WaitGroup, tpl *CodeTemplate, tb *Table, path string) {
				defer wg.Done()
				str := s.GenerateCode(tb, tpl, "", true, "")
				dstPath, _ := s.PredefineTargetPath(tpl, tb)
				if dstPath == "" {
					dstPath = s.DefaultTargetPath(path, tb)
				}
				if err := SaveFile(str, outputDir+"/"+dstPath); err != nil {
					println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s ,template:%s",
						err.Error(), tpl.FilePath()))
				}
			}(&wg, tpl, tb, path)
		}
	}
	wg.Wait()
	return err
}
