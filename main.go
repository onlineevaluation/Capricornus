package main

import (
	"C"
	"bytes"
	"github.com/json-iterator/go"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/**
魔羯座
@author: 杨晓辉

测评机，提供c语言代码算法测试运行，可以扩展为其他语言

os:
	- windows 10

编译命令
编译为 so
------------------------------------------------------------------------------------------------------------------------
`go build -buildmode=c-shared -o .\out\libCapricornus.so .\main.go`
------------------------------------------------------------------------------------------------------------------------
code 说明
0 没有安装 gcc 环境
1 代码无法进行编译
2 运行超时
3 运行出错
4 未知错误
5 json格式错误

8 部分运行结果错误
9 运行全部通过
*/

type datas struct {
	Datas []data
}
type data struct {
	Input  string
	Output string
}

//export judgeCode
func judgeCode(filePath, outputPath, fileName string, data string, limitTime int64) *C.char {

	result := make(chan string)
	// 解析 json
	var d datas
	if err := jsoniter.Unmarshal([]byte(data), &d); err != nil {
		println("json 格式错误 " + data)
		return C.CString("code:5 json格式错误 " + err.Error())
	}

	// 系统是否安装有 gcc 环境
	if _, i := exec.LookPath("gcc"); i != nil {
		println("没有安装c语言环境", i.Error())
		return C.CString("code:0 没有安装 c 语言环境，请安装 gcc ")
	}

	// 获取系统信息
	osName := runtime.GOOS

	println("os is", osName)

	println("准备开始编译 C 语言")
	switch osName {
	case "windows":
		go runInWindows(filePath, outputPath, fileName, result, d.Datas, limitTime)
		break
	case "linux":
		go runInXnux(filePath, outputPath, fileName, result, d.Datas, limitTime)
		break
	case "macOs":
		go runInXnux(filePath, outputPath, fileName, result, d.Datas, limitTime)
		break
	default:
		result <- "还不支持该系统"
		break
	}

	select {
	case v := <-result:
		println("结果：", v)
		return C.CString(v)
	}

}

func main() {
	var filePath = "/mnt/f/testData/Add.cpp"
	var outPath = "/mnt/f/testData"
	var fileName = "add"
	// runCmdLine := outPath + "/" + fileName
	// {"datas":[{"input":"【0,4%$#","output":"【4%$#"},{"input":"【1,4%$#","output":"【5%$#"},{"input":"【2,4%$#","output":"【6%$#"},{"input":"【3,4%$#","output":"【7%$#"}]}
	judgeCode(filePath, outPath, fileName,
		`{"datas":[{"input":"#$%0,4%$#","output":"#$%4%$#"},{"input":"#$%1,4%$#","output":"#$%5%$#"},{"input":"#$%2,4%$#","output":"#$%6%$#"},{"input":"#$%3,4%$#","output":"#$%7%$#"}]}`, 5)
}

/**
 * 在 windows 下编译c语言
 */
func runInWindows(filePath, outputPath, fileName string, result chan string, data []data, limitTime int64) {
	// gcc -Wall e:/testData/Hello.cpp -o ollcode
	// 错误检查
	println("检查编译问题")
	cmdLine := "gcc -pedantic " + filePath + " -o " + outputPath + "/" + fileName + " "
	// 这里使用 powershell ,否者无法获取错误信息
	cmd := exec.Command("powershell", "/C ", cmdLine)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	_ = cmd.Run()
	// 代码错误
	if len(w.Bytes()) != 0 {
		result <- "code:1 " + string(w.Bytes())
	}
	println(string(w.Bytes()))
	// 编译 c 语言文件
	_, e := exec.Command("cmd", "/C", "gcc -g -o "+outputPath+"/"+fileName+" "+filePath).Output()
	// 异常处理
	if e != nil { // 无法编译
		result <- string("err:4 " + e.Error())
	}
	println("程序编译完成")
	// 运行 c 语言
	//开启协程
	go judge(outputPath, fileName, result, data, limitTime)

}

func runInXnux(filePath, outputPath, fileName string, result chan string, data []data, limitTime int64) {
	println("检查编译问题")
	cmdLine := "gcc -pedantic " + filePath + " -o " + outputPath + "/" + fileName
	cmd := exec.Command("/bin/bash","-c",cmdLine)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	_ = cmd.Run()
	if len(w.Bytes()) != 0 {
		result <- "code:1 " + string(w.Bytes())
	}
	println(string(w.Bytes()))
	// 编译 c/c++ 文件
	_, e := exec.Command("/bin/bash", "-c", "gcc -g -o "+outputPath+"/"+fileName+" "+filePath).Output()
	if e != nil {
		result <- string("err:4"+ e.Error())
	}
	println("程序编译完成")
	// 运行 c 语言
	//开启协程
	go judge(outputPath, fileName, result, data, limitTime)
}

/**
 * 代码运行
	将数据进行传输
*/
func judge(outputPath, fileName string, result chan string, data []data, limitTime int64) {

	// 程序开始时间
	start := time.Now().Unix()
	// go runCode(outputPath, fileName, result, data)
	// 判断系统
	osName := runtime.GOOS
	switch osName {
	case "windows":
		runCmdLine := outputPath + "/" + fileName
		go runCode(outputPath, fileName, result, data,runCmdLine)
		for {
			cur := time.Now().Unix()
			if cur-start >= limitTime {
				// 杀死进程
				// windows
				process := fileName + ".exe"
				_ = exec.Command("cmd", "/C", "taskkill /F /IM "+process).Run()
				result <- "code:2 运行超时"
			}
		}
		break;
	case	"linux":
		runCmdLine := outputPath + "/" + fileName
		println("run line "+runCmdLine)
		go runCode(outputPath, fileName, result, data,runCmdLine)
		for {
			cur := time.Now().Unix()
			if cur-start >= limitTime {
				// 杀死进程
				// linux
				process := fileName
				_ = exec.Command("cmd", "-c", "killall "+process).Run()
				result <- "code:2 运行超时"
			}
		}
		break;
	default:
		result <- "not support this os"
		break;
	}
	// 程序运行时间
	
}

/*
	代码运行
*/
func runCode(outputPath, fileName string, result chan string, data []data,runCmdLine string) {
	println("程序准备运行")
	var flag = 0
	for i := 0; i < len(data); i++ {
		sub := strings.Split(data[i].Input, "#$%")
		sub = strings.Split(sub[1], "%$#")
		sub = strings.Split(sub[0], ",")
		// 拼接参数
		var args string
		for i := 0; i < len(sub); i++ {
			args += sub[i] + " "
		}
		cmd := exec.Command(runCmdLine)
		cmd.Stdin = strings.NewReader(args)
		output, e := cmd.Output()
		if e != nil {
			result <- string("code:3 " + e.Error())
		}
		// 输出获取
		out := strings.Split(data[i].Output, "#$%")
		out = strings.Split(out[1], "%$#")

		println("第 ", i+1, "次答案", string(output))
		if string(output) == out[0] {
			flag++
		}
	}
	if flag == len(data) {
		result <- string("code:9 运行完美")
	} else if flag > 0 && flag < len(data) {
		result <- string("code:8 部分答案错误 正确数量为") + string(flag) + "/" + string(len(data))
	} else {
		result <- string("code:3 运行出错")
	}
}