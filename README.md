# tg bot

## 功能

下载youtube music中的音乐连接

与gemini对话

与群友交互

chatgpt待续...

## 使用教程（tg用户）

群互动机器人

用法：

    cp语录:
    
        1. 引用其他人的消息
        
        2. 回复关键词 [mua,mua~,摸摸,抱抱]
        
        3. 有60%概率触发，摘星会引用你引用的话，并发🍬
        
    骂人：
    
        1. 引用其他人的消息
        
        2. 回复 [骂他，咬他]，其中 他 可以替换为 她 它 ta
        
        3. 100%触发,摘星会引用你引用而话，并骂他
        
    chatgpt：
    
        1. 在群聊中使用 `/chat msg` 可以与摘星聊天，MSG可以是任意内容
        
        2. 在群聊里应用摘星的话，摘星会以为你在和他聊天，@则不会
        
        3. 私聊摘星，摘星会与你对话
        
    youtubeMusic下载：
    
        私聊或者群聊里发送youtubeMusic链接，摘星会下载音乐并唱给你听
        
> 摘星是bot的名字：@ytbmusicPlaerBot

另外你们的start是作者最大的动力😀


## 使用教程（服务端）

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


