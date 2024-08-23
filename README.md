# youtube music bot

## 功能

下载youtube music中的音乐连接

！gemini接口！

与群友交互
chatgpt待续

## 使用教程

修改config copy.toml为config.toml

token为bot的token

address为本地监听的ip：端口

domain为tg服务端向你发送请求的域名，注意不要带https前缀，tg会使用https访问，所以这里不要填ip

连接tg后直接发送音乐的分享链接就可以

## youtube music

向机器人发送音乐的分享链接
bot会自动下载

## chat 使用

私聊chat 会自动触发chat功能，bot会使用gemini（目前不支持openai）进行答复

在群组中 使用 "/chat " 开头的消息会被bot识别，注意后面有个空格

## 交互功能

需要满足两个条件

- 引用了群组中的某个留言
- 包含以下关键字

注：左侧为关键字，右侧为处理方案，触发时会对比左侧

```go
var quotationsKey = map[string]string{
	"骂她": 骂人,
	"骂他": 骂人,
	"骂它": 骂人,
	"咬他": 骂人,
	"咬它": 骂人,
	"咬她": 骂人,

	"舔ta":  舔狗,
	"舔":    舔狗,
	"t":    舔狗,
	"tian": 舔狗,

	"有病":   神经,
	"神经":   神经,
	"发神经":  神经,
	"神经病":  神经,
	"有病吧":  神经,
	"你有病吧": 神经,

	"爱你":   情侣,
	"mua":  情侣,
	"mua~": 情侣,
	"宝":    情侣,
	"宝儿":   情侣,
	"宝儿~":  情侣,
	"摸摸":   情侣,
	"抱抱":   情侣,
	"亲亲":   情侣,
	"摸摸~":  情侣,
	"抱抱~":  情侣,
	"亲亲~":  情侣,
}
```
