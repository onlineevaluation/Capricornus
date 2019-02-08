package main

import (
	"C"
	"bytes"
	"os/exec"
	"runtime"
	"time"
)

/**
魔羯座
测评机
编译命令
编译为 so
`go build -buildmode=c-shared -o .\out\libCapricornus.so .\main.go`

错误码说明
0 没有安装 gcc 环境
1 代码无法进行编译
2 运行超时
3 运行出错
4 未知错误
 */

//export judgeCode
func judgeCode(filePath, outputPath, fileName string, data string) *C.char {

	println("文件路径为 "+filePath, "文件输出路径为"+outputPath, "文件名为", fileName, "data ", data)

	result := make(chan string)
	// 系统是否安装有 gcc 环境
	_, i := exec.LookPath("gcc")
	if i != nil {
		println("没有安装c语言环境", i.Error())
		return C.CString("error:0 没有安装 c 语言环境，请安装 gcc")
	}
	// 获取系统信息
	osName := runtime.GOOS
	println("os is", osName)
	println("准备开始编译 C 语言")
	switch osName {
	case "windows":
		go runInWindows(filePath, outputPath, fileName, result)
		break
	case "linux":
		runInLinux()
		break
	default:
		runInMacOs()
		break
	}

	select {
	case v := <-result:
		println("结果：", v)
		return C.CString(v)
	}

}

func main() {
	//var filePath = "e:/testData/Hello.cpp"
	//var outPath = "e:/testData"
	//var fileName = "hello"
	//
	//judgeCode(filePath, outPath, fileName, "")
}

/**
 * 在 windows 下编译c语言
 */
func runInWindows(filePath, outputPath, fileName string, result chan string) {
	// gcc -Wall e:/testData/Hello.cpp -o ollcode
	// 错误检查
	println("检查编译问题")
	cmdLine := "gcc -pedantic " + filePath + " -o " + fileName
	cmd := exec.Command("powershell", "/C ", cmdLine)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	_ = cmd.Run()
	if len(w.Bytes()) != 0 {
		result <- "error:1 " + string(w.Bytes())
	}
	println(string(w.Bytes()))
	// 编译 c 语言文件
	_, e := exec.Command("cmd", "/C", "gcc -g -o "+outputPath+"\\"+fileName+" "+filePath).Output()
	// 异常处理

	if e != nil {
		// 无法编译
		result <- string("err:4 " + e.Error())
	}
	println("程序编译完成")
	// 运行 c 语言
	//开启协程
	go judge(outputPath, fileName, result)

}

func runInLinux() {

}

func runInMacOs() {

}

/**
 * 代码运行
 */
func judge(outputPath, fileName string, result chan string) {
	//var i1 = ""

	// 程序开始时间
	start := time.Now().Unix()

	go runCode(outputPath, fileName, result)
	// 程序运行时间
	for ; ; {
		cur := time.Now().Unix()
		if cur-start >= 2 {
			// 杀死进程
			// windows
			process := fileName + ".exe"
			_ = exec.Command("cmd", "/C", "taskkill /F /IM "+process).Run()
			result <- "error:2 运行超时"
		}
	}
}

/*
	代码运行
 */
func runCode(outputPath, fileName string, result chan string) {
	println("程序准备运行")
	bytes, e := exec.Command(outputPath + "/" + fileName).Output()
	if e != nil {
		println("error:3 " + e.Error())
	}
	result <- string(bytes)
}
