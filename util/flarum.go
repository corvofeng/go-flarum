package util

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func isValidFilename(filename string) bool {
	return strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml")
}

func recursiveRead(localeData map[string]interface{}, prefix string, localeDataArr *map[string]string) {
	logger := GetLogger()
	// 递归的获取信息
	for key, element := range localeData {

		switch v := element.(type) {
		// case int:
		// 	// 	// v is an int here, so e.g. v + 1 is possible.
		// 	fmt.Printf("Integer: %v", v)
		// case float64:
		// 	// v is a float64 here, so e.g. v + 1.0 is possible.
		// 	fmt.Printf("Float64: %v", v)
		case string:
			// v is a string here, so e.g. v + " Yeah!" is possible.
			// logger.Debugf("%s: %s\n", key, v)
			(*localeDataArr)[fmt.Sprintf("%s.%s", prefix[1:], key)] = v
		case map[string]interface{}:
			recursiveRead(v, fmt.Sprintf("%s.%s", prefix, key), localeDataArr)
		case map[interface{}]interface{}:
			newDict := make(map[string]interface{})
			for _k, _v := range v {
				_key := ""
				switch _kdata := _k.(type) {
				case int:
					_key = fmt.Sprintf("%d", _kdata)
				case string:
					_key = _kdata
				default:
					logger.Errorf("Can't parse type %+v, %s", _k, reflect.TypeOf(_k))
				}
				newDict[_key] = _v
			}
			recursiveRead(newDict, fmt.Sprintf("%s.%s", prefix, key), localeDataArr)
		default:
			logger.Errorf("Can't parse type %+v, %s", v, reflect.TypeOf(v))
		}
	}
}

func readLocaleData(localeData map[string]interface{}, localeDataArr *map[string]string) {
	// localeDataArr := make(map[string]string)
	recursiveRead(localeData, "", localeDataArr)

	for key, value := range *localeDataArr {
		if strings.HasPrefix(value, "=> ") {
			if val, ok := (*localeDataArr)[value[3:]]; ok {
				//do something here
				// fmt.Println(key, value)
				(*localeDataArr)[key] = val
			}
		}

	}
}

// FlarumReadLocale 读取Flarum的语言包
func FlarumReadLocale(flarumDir string, extDirs []string, localeDir, locale string) map[string]string {
	logger := GetLogger()
	localeDataArr := make(map[string]string)

	doParseFile := func(fn string) {
		yamlFile, err := os.ReadFile(fn)
		if err != nil {
			logger.Errorf("yamlFile.Get err %v ", err)
		}
		localeData := make(map[string]interface{})

		err = yaml.Unmarshal(yamlFile, &localeData)
		if err != nil {
			logger.Fatalf("Unmarshal: %v", err)
			return
		}
		readLocaleData(localeData, &localeDataArr)
	}

	// flarum/locale
	// 首先解析flarum中自带的语言包
	flarumDatas, err := os.ReadDir(path.Join(flarumDir, "locale"))
	if err == nil {
		for _, fi := range flarumDatas {
			// 过滤指定格式
			if ok := isValidFilename(fi.Name()); ok {
				doParseFile(path.Join(flarumDir, "locale", fi.Name()))
			}
		}
	}

	// 解析flarum-lang中携带的语言包
	// locale/en/*.yml, locale/zh/*.yml
	dirPath := path.Join(localeDir, locale, "locale")
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Errorf("Can't read '%s' with err '%v'", dirPath, err)
		return localeDataArr
	}
	for _, fi := range dir {
		// 过滤指定格式
		if ok := isValidFilename(fi.Name()); ok {
			doParseFile(path.Join(dirPath, fi.Name()))
		}
	}

	// 解析各个插件的语言包
	// flarum-tags/en/xxx.yml
	for _, extDir := range extDirs {
		extDirDatas, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}
		for _, fi := range extDirDatas {
			if !fi.IsDir() {
				continue
			}

			extLocaleDir := path.Join(extDir, fi.Name(), locale)
			dir, err := os.ReadDir(extLocaleDir)

			if err != nil {
				continue
			}
			for _, fi := range dir {
				// 过滤指定格式
				if ok := isValidFilename(fi.Name()); ok {
					doParseFile(path.Join(extDir, fi.Name()))
				}
			}

		}
	}

	// flarum-tags/locale/en.yml
	for _, extDir := range extDirs {
		extDirDatas, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}
		for _, fi := range extDirDatas {
			if !fi.IsDir() {
				continue
			}

			extLocaleDir := path.Join(extDir, fi.Name(), "locale")
			dir, err := os.ReadDir(extLocaleDir)
			if err != nil {
				continue
			}

			for _, fi := range dir {
				// 过滤指定格式
				if ok := isValidFilename(fi.Name()); ok {
					doParseFile(path.Join(extLocaleDir, fi.Name()))
				}
			}

		}
	}

	// fof-upload/resources/locale/en.yml
	for _, extDir := range extDirs {
		extDirDatas, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}
		for _, fi := range extDirDatas {
			if !fi.IsDir() {
				continue
			}

			extLocaleDir := path.Join(extDir, fi.Name(), "resources", "locale")
			dir, err := os.ReadDir(extLocaleDir)
			if err != nil {
				continue
			}

			for _, fi := range dir {
				// 过滤指定格式
				if ok := isValidFilename(fi.Name()); ok {
					doParseFile(path.Join(extLocaleDir, fi.Name()))
				}
			}

		}
	}

	return localeDataArr
}
