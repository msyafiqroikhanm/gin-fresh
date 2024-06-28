package helpers

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// type Log struct {
// 	StartTime time.Time
// 	Location  string
// }

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func GetFunctionAndStructName(i interface{}) (string, string, string) {
	pc, file, _, ok := runtime.Caller(2) // 2 untuk mengambil caller dari fungsi pemanggil
	if !ok {
		return "", "", ""
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := fullFuncName[strings.LastIndex(fullFuncName, ".")+1:]

	contrName := reflect.TypeOf(i).Elem().Name()

	projectRoot := getProjectRoot()

	relativePath, err := filepath.Rel(projectRoot, file)
	if err != nil {
		relativePath = file
	}

	packagePath := strings.ReplaceAll(relativePath, string(filepath.Separator), "/")
	packagePath = strings.TrimSuffix(packagePath, filepath.Ext(packagePath))

	projectName := filepath.Base(projectRoot)
	if idx := strings.Index(packagePath, projectName); idx != -1 {
		packagePath = packagePath[idx+len(projectName)+1:]
	}

	return funcName, contrName, packagePath
}

func CreateLog(c *gin.Context, i interface{}) handlers.Log {
	funcName, contrName, packagePath := GetFunctionAndStructName(i)
	return handlers.Log{
		StartTime: time.Now(),
		Location:  fmt.Sprintf("%s/%s.%s", packagePath, contrName, funcName),
	}
}
