## 方法1
```
bee generate docs
bee run watchall true

#自动化生成文档
bee run -gendoc=true -downdoc=true
```
## govendor
```
govendor init
govendor add +external
```
## 添加包
```
govendor add +outside
```
## 同步包
```
govendor sync
```