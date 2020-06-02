package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

func recursiveRead(localeData map[interface{}]interface{}, prefix string, localeDataArr *map[string]string) {
	// 递归的获取信息
	for key, element := range localeData {

		switch v := element.(type) {
		// case int:
		// 	// v is an int here, so e.g. v + 1 is possible.
		// 	fmt.Printf("Integer: %v", v)
		// case float64:
		// 	// v is a float64 here, so e.g. v + 1.0 is possible.
		// 	fmt.Printf("Float64: %v", v)
		case string:
			// v is a string here, so e.g. v + " Yeah!" is possible.
			(*localeDataArr)[fmt.Sprintf("%s.%s", prefix[1:], key)] = v
			// fmt.Printf("%s: %s\n", , v)
		case map[interface{}]interface{}:
			recursiveRead(v, fmt.Sprintf("%s.%s", prefix, key), localeDataArr)
		default:
			// And here I'm feeling dumb. ;)
			fmt.Printf("I don't know, ask stackoverflow. %+v", v)
		}
	}
}

func readLocaleData(localeData map[interface{}]interface{}, localeDataArr *map[string]string) {
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
func FlarumReadLocale(localeDir, locale string) map[string]string {
	localeDataArr := make(map[string]string)

	dirPath := path.Join(localeDir, locale, "locale")
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return localeDataArr
	}

	doParseFile := func(fn string) {

		yamlFile, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		localeData := make(map[interface{}]interface{})

		err = yaml.Unmarshal(yamlFile, &localeData)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			return
		}
		readLocaleData(localeData, &localeDataArr)
	}

	for _, fi := range dir {
		// 过滤指定格式
		if ok := strings.HasSuffix(fi.Name(), ".yml"); ok {
			doParseFile(path.Join(dirPath, fi.Name()))
		}
	}

	// fmt.Println(localeData)
	// fmt.Println(reflect.TypeOf(localeData))

	return localeDataArr
}
