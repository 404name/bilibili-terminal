package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/404name/termui-demo/global"
)

func CallCommandRun(ctx context.Context, cmdName string, args []string) (string, error) {

	cmd := exec.CommandContext(ctx, cmdName, args...)
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

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// Copy 可中断的流复制
func Copy(ctx context.Context, dst io.Writer, src io.Reader) error {
	// Copy will call the Reader and Writer interface multiple time, in order
	// to copy by chunk (avoiding loading the whole file in memory).
	// I insert the ability to cancel before read time as it is the earliest
	// possible in the call process.
	_, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {

		// golang non-blocking channel: https://gobyexample.com/non-blocking-channel-operations
		select {
		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return src.Read(p)
		}
	}))
	return err
}
