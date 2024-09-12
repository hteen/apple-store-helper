# Apple Store 预约助手

## 更新 iPhone 16 系列

## 重要提示
* *这不是外挂，不能全自动一劳永逸*
* *提前登录*
* *提前将需要购买的型号加入购物车，检测有货会打开购物车页面，需要在购物车页面手动选择门店*


## 关于开发
* 代码不优雅, 注释不完善, review须谨慎
* GUI框架 [fyne](https://github.com/fyne-io/fyne)

### 运行
```shell script
go run main.go
```

### 打包
```
# macos下打包其他平台
go get github.com/fyne-io/fyne-cross
fyne-cross darwin -arch=amd64,arm64
fyne-cross darwin -arch=amd64,arm64 -app-id=apple-store-helper
```

## 使用方法


1. 前往 `release` 页面下载，启动 
2. 提前将需要购买的型号加入购物车，检测有货会打开购物车页面，需要在购物车页面手动选择门店
3. 选择门店和型号，点击 `添加` 到监控列表
4. 点击 `开始` 即可

匹配到之后会暂停监听，直到再次点击 `开始`

## 一杯卡布奇诺 ☕️

<img src='https://tva1.sinaimg.cn/large/0081Kckwly1gls6d2nnicj30i00pcq9i.jpg' width='200px'/>
