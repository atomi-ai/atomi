# atomi

## 启动服务器
```shell
go install github.com/atomi-ai/atomi && ~/work/bin/atomi
```

## TODO
- Move the package to github.com/atomi-ai/server

## Coding Tips
### 如果你在goland底下有些包显示虹色，感觉是import错误很annoying的时候：
这时候往往这个包其实是已经正确的download了的，你的go.mod也没有问题，你可以用
```go mod vendor```
试试看。虽然vendor是已经deprecated的方式，但是对于IDE而言，每个project有自己dedicated的实现是一个非常方便的事情。所以它常常可以识别vendor的库，但是却不能识别mod的库。以后的IDE应该会优化掉，我现在用的是2023.1版本的goland。

### wire_gen.go中的初始设置一般是这样的：
```go
package main

//go:generate wire
```
### wire的引入，导致golang测试需要放到另外的packages里面，否则会有循环依赖，譬如：
- 创建一个controllers/login_controller_test.go
- 它会用一个testing application，假设这个在app/ package里面。
- 那么application必然需要controller来构建LoginController for tests
- 这样，我们必然有这样的循环依赖： controller => app => controller。理论上，只要test cases在controller package下，我们就没有办法跳出这个循环。
- 解决这样的循环，我目前想到的一个办法（理论上应该是唯一的），就是test cases放到test/controllers package里面去，这样就变成 test/controllers => app => controllers，就没有循环了。

## tax-rates-csv的生成
网站下载的csv文件比较大，使用以下命令将其处理保留State,ZipCode,EstimatedCombinedRate三列:
```
for file in TAXRATES_ZIP5/*.csv; do
  filename=$(basename "$file" .csv)
  cut -d ',' -f 1,2,4 "$file" > "tax-rates-csv/$filename.csv"
done
```
