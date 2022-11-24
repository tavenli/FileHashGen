# FileHashGen

本人发现市面上大多数文件指纹校验工具，都含有恶意代码或后门。所以自己写了这个小工具，方便固定电子数据证据使用，可根据电子数据文件批量生成HASH值。

- 支持MD5、SHA-1、SHA-256、SHA-512算法等。
- 支持超大文件计算Hash指纹，并且速度超快。
- 支持对整个目录下的文件，自动批量计算Hash指纹，一键输出报告。
- 支持指定单个文件计算Hash指纹。
- 可跨平台，支持Windows、Linux、MacOS等。
- 开源免费，欢迎安全工程师、司法机构电子取证人员使用。
- 该版本为命令行版本，图形化版本另外公布地址。


当然，本工具也可以用于验证网上下载的各种安装文件是否与官网一致，例如Windows安装镜像、Office等。


```
国际开源地址：https://github.com/tavenli/FileHashGen
国内开源地址：https://gitee.com/tavenli/FileHashGen

以上两个仓库，都会一起推送。
```

# 使用说明
```
以 Windows 系统下为例

自动生成目录下所有文件的指纹信息，不带任何参数，直接执行（推荐）：
FileHashCode.exe

指定单个文件：
FileHashCode.exe -f "d:\检材目录\检材1.docx"

指定文件夹下所有文件：
FileHashCode.exe -f "d:\检材目录\视频文件\"

生成目录下所有文件，只使用两种算法：
FileHashCode.exe -c "MD5,SHA-256"

指定一个文件，只使用两种算法：
FileHashCode.exe -c "MD5,SHA-256" -f "d:\检材目录\检材1.docx"

注意： -c 参数不指定时，默认所有支持的算法都会生成
```

