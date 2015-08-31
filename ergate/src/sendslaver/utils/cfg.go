// utils
package utils

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	filePath   string
	configList map[string]string
}

func LoadConfig(FilePath string) *Config {
	myConfig := Config{
		filePath:   FilePath,
		configList: make(map[string]string),
	}

	f, err := os.Open(FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buff := bufio.NewReader(f)
	for {
		line, err := buff.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		lineStr := getUsefulLineStr(line)
		if !checkLineUseful(lineStr) {
			continue
		}

		key, value := getKeyValue(lineStr)
		myConfig.configList[key] = value
	}

	return &myConfig
}

func (c *Config) GetInt(key string) int {
	valueStr := c.configList[key]
	return stringToInt(valueStr)
}

func (c *Config) GetString(key string) string {
	value := c.configList[key]
	return value
}

func (c *Config) GetBool(key string) bool {
	value := c.configList[key]
	if strings.ToLower(value) == "false" {
		return false
	}
	return true
}

func (c *Config) GetStrArray(key string, spitStr string) []string {
	var result []string
	valueStr := c.configList[key]
	tmp := strings.Split(valueStr, spitStr)
	for i := 0; i < len(tmp); i++ {
		obj := strings.TrimSpace(tmp[i])
		if len(obj) > 0 {
			result = append(result, obj)
		}
	}
	return result
}

func (c *Config) GetIntArray(key string, spitStr string) []int {
	var resultIntArray []int
	valueStr := c.configList[key]
	strArray := strings.Split(valueStr, spitStr)
	for i := 0; i < len(strArray); i++ {
		tmp := stringToInt(strArray[i])
		resultIntArray = append(resultIntArray, tmp)
	}
	return resultIntArray
}

func stringToInt(str string) int {
	if len(str) == 0 {
		return 0
	}
	tmp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(tmp)
}

func getUsefulLineStr(line string) (lineStr string) {
	lineStr = strings.TrimSpace(line)
	return
}

func checkLineUseful(lineStr string) bool {
	if len(lineStr) == 0 {
		return false
	}

	if lineStr[0:1] == "#" {
		return false
	}

	if !strings.Contains(lineStr, " = ") {
		return false
	}

	return true
}

func getKeyValue(lineStr string) (key string, value string) {
	resArray := strings.SplitN(lineStr, " = ", 2)
	if len(resArray) == 2 {
		key = strings.TrimSpace(resArray[0])
		value = strings.TrimSpace(resArray[1])
	}
	return
}
