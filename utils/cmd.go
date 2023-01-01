package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/404name/termui-demo/global"
)

func CallCommandRun(cmdName string, args []string) (string, error) {

	cmd := exec.Command(cmdName, args...)
	global.LOG.Infoln("CallCommand Run 参数=> ", args)
	global.LOG.Infoln("CallCommand Run 执行命令=> ", cmd)
	bytes, err := cmd.Output()
	if err != nil {
		global.LOG.Errorln("CallCommand Run 出错了.....", string(bytes), err.Error())
		return "", err
	}
	resp := string(bytes)
	global.LOG.Infoln("CallCommand Run 调用完成=> ", resp)
	return resp, nil
}

func ExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	re, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("The eacePath failed: %s\n", err.Error())
	}
	fmt.Println("The path is ", re)
	return filepath.Abs(file)
}
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
