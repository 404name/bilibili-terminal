package utils

import (
	"os/exec"
)

func CallCommandRun(cmdName string, args []string) (string, error) {

	cmd := exec.Command(cmdName, args...)
	Log.Infoln("CallCommand Run 参数=> ", args)
	Log.Infoln("CallCommand Run 执行命令=> ", cmd)
	bytes, err := cmd.Output()
	if err != nil {
		Log.Errorln("CallCommand Run 出错了.....", string(bytes), err.Error())
		return "", err
	}
	resp := string(bytes)
	Log.Infoln("CallCommand Run 调用完成=> ", resp)
	return resp, nil
}
