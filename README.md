# Capricornuc 魔羯座

## 什么是魔羯座

魔羯座是用于测评系统的算法评测机

主要用于 基于 C/C++ 的数据结构使用。

项目采用 Go 语言编写,最后编译为 so , 使其他语言进行调用。


## API

### 接口

```c
judgeCode(GoString filePath,GoString outputPath,GoString fileName,GoString data):String
```

### 参数说明

- filePath 要编译的C语言文件路径
- outputPath 输出路径
- fileName 输出文件名称
- data 测评数据


**data** 类型是是一个 Json 结构,查看[data.json](./data.json)

```json
{
  "data": [
    {
      "input": "[1,2]",
      "output": "[3]"
    },
    {
      "input": "[3,6]",
      "output": "[9]"
    }
  ]
}
```

### 类型说明

`GoString` 是 Go 语言的 String 类型,直接用 String 可能无法使用，要手动实现该类型。 

### 返回值

返回 String 类型数据，数据中类型如下

```text
error:0 xxxxxxxxxxxxxxxx
```

error 表示错误码 
错误码类型目前如下

错误码说明

- 0 没有安装 gcc 环境
- 1 代码无法进行编译 包含错误信息
- 2 运行超时
- 3 运行出错
- 4 未知错误

## 编译

确保本地使用的 gcc 支持系统位数和Go语言位数一致

通过运行如下命令进行编译打包为 `*.so`

```shell
go build -buildmode=c-shared -o .\out\libCapricornus.so .\main.go
```

生成的 so 在 out 目录下，同时还会生成 .h 文件

还可以生成 dll 文件

```shell
go build -buildmode=c-shared -o .\out\libCapricornus.dll .\main.go
```

## Kotlin 使用案例

可以通过 JNI 或者 JNA 调用。如下使用 JNA 调用。

### 引入 JNA

JNA Github 地址 [https://github.com/java-native-access/jna](https://github.com/java-native-access/jna)

编写 JNA 接口

```kotlin
interface Capricornus : Library {

    fun judgeCode(
        filePath: GoString.ByValue,
        outputPath: GoString.ByValue,
        fileName: GoString.ByValue,
        data:GoString.ByValue
    ): String



    fun add(a: Int, b: Int): Int

    companion object {
        val INSTANCE =
            Native.load("C:\\Users\\young\\Desktop\\native\\cmder\\dll\\libCapricornus.so", Capricornus::class.java)!!
    }
}

```

对 GoString 进行实现

```java

public class GoString extends Structure {

    public String str;
    public long length;

    public GoString() {

    }

    public GoString(String str) {
        this.str = str;
        this.length = str.length();
    }

    @Override
    protected List<String> getFieldOrder() {
        List<String> files = new ArrayList<>();
        files.add("str");
        files.add("length");
        return files;
    }

    public static class ByValue extends GoString implements Structure.ByValue {
        public ByValue() {
        }

        public ByValue(String str) {
            super(str);
        }
    }

    public static class ByReference extends GoString implements Structure.ByReference {
        public ByReference() {
        }

        public ByReference(String str) {
            super(str);
        }
    }
}

```

调用该接口

```kotlin
fun main() {
    val filePath = GoString.ByValue("E:\\testData\\Hello.cpp")
    val outPath = GoString.ByValue("E:\\testData")
    val fileName = GoString.ByValue("HelloWorld")
    //language=JSON
    val data =
        GoString.ByValue("{\n  \"data\": [\n    {\n      \"input\": \"[1,2]\",\n      \"output\": \"[3]\"\n    },\n    {\n      \"input\": \"[3,6]\",\n      \"output\": \"[9]\"\n    }\n  ]\n}")
    val result = Capricornus.INSTANCE.judgeCode(filePath, outPath, fileName, data)
    val errorCode = result.substring(6, 7)
    val message = when (errorCode) {
        "0" -> "没有安装 Gcc 环境"
        "1" -> "代码语法错误，无法进行编译"
        "2" -> "代码运行超时"
        "3" -> "代码运行错误"
        else -> {
            "未知错误"
        }
    }
    println(message)
}

```