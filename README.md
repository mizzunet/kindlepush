KindlePush [![GitHub release](https://img.shields.io/github/release/zhengchun/kindlepush.svg)](https://GitHub.com/zhengchun/kindlepush/releases/)
===
Subscribe the RSS/Atom feeds and delivered to your Kindle.

Download [KindlePush](https://github.com/zhengchun/kindlepush/releases) latest version.

### How to Use

- Download [KindleGen](https://www.amazon.com/gp/feature.html?docId=1000765211) software and install into your computer.

- Edit `config.yaml` file.

- Run `kindlepush` or `kindlepush -config=config.yaml` command.

### config.yaml

- Required
    * kindleAddress - Kindle email address.
    * smtp - About SMTP configuration used to send email. 
    * feeds - List of Atom/RSS feed URLs.
- Optional
    * cacheDir - Output directory.
    * resizeImage - Resize image with specifies width and height.
    * maxFileSize - Max mobi file size.
    * kindlegen - kindlegen.exe path.
    * proxy - HTTP proxy.


`config.yaml` example

```yaml
kindleAddress: test@kindle.cn
smtp:
    # domain:port,(etc. gmail.com:25), port number is 25 in default.
    hostAndPort: smtp.gmail.com:465
    # Display your email address.
    # Default `account` value will used if no set this value.
    senderAddress: xxx@gmail.com
    # Use SSL to connect to the email server.
    ssl: true
    # Enable SMTP Authentication.
    account: xxx@gmail.com
    password: $password$
feeds:
  - https://feedocean.com/feeds/sry69h
  - https://feedocean.com/feeds/d519ydb
```

NOTES: Make sure `senderAddress` add to your kindle approved e-mail address list. See [Add an Email Address to Receive Documents](https://www.amazon.com/gp/help/customer/display.html?nodeId=201974240).

中文说明
===

`KindlePush`是Kindle电子书推送软件，能够将你订阅多个RSS/Atom源生成MOBI格式的电子杂志，并推送到你的Kindle设备上阅读。

### 使用说明

- 下载[KindleGen](https://www.amazon.com/gp/feature.html?docId=1000765211)并安装到本地机器。

- 编辑 `config.yaml` 配置文件。

- 运行 `kindlepush` 或 `kindlepush -config=config.yaml` 命令。

请注意：正常接收投递电子杂志，你必须确保你的`senderAddress`邮件地址已经加入到kindle的邮箱列中。[添加用于接收文档的电子邮箱](https://www.amazon.cn/gp/help/customer/display.html?nodeId=201974240)。

### config.yaml说明

- 必选
    * kindleAddress - Kindle的邮箱地址。
    * smtp - SMTP配置，用于邮箱投递。
    * feeds - Atom/RSS源列表。
- 可选
    * cacheDir - 生成的输出目录。
    * resizeImage - 自定义略缩图像大小，可以减少电子书的大小。
    * maxFileSize - Mobi文件最大的允许大小。
    * kindlegen - kindlegen.exe文件路径。
    * proxy - HTTP代理。

A list of RSS/Atom feeds
===

|Name |URL |
|--------------------------|----------------|
|36氪 | [http://36kr.com/feed](http://36kr.com/feed)|
|Cnbeta IT News | [https://feedocean.com/feeds/gjn0yf](https://feedocean.com/feeds/gjn0yf)|
|Engadget | [https://feedocean.com/feeds/3soj5w](https://feedocean.com/feeds/3soj5w)|
|Engadget 中国版 | [https://feedocean.com/feeds/6381j2](https://feedocean.com/feeds/6381j2)|
|以太坊中文区 | [https://ethfans.org/feed](https://ethfans.org/feed) |
|FT中文版 | [https://feedocean.com/feeds/15qqtz](https://feedocean.com/feeds/15qqtz)|
|The New York Times | [https://feedocean.com/feeds/5yoe4g](https://feedocean.com/feeds/5yoe4g)|
|纽约时报中文网 | [https://feedocean.com/feeds/fpdymq](https://feedocean.com/feeds/fpdymq)|
|RFA | [https://www.rfa.org/english/rss2.xml](https://www.rfa.org/english/rss2.xml)|
|知乎日报 | [https://feedocean.com/feeds/sry69h](https://feedocean.com/feeds/sry69h)|

[FeedOcean.com](https://feedocean.com): Provides convert HTML to RSS and Atom feeds with full-text, Its free.

**You're welcome to contribute.**

![alt text](https://github.com/zhengchun/kindlepush/blob/master/docs/001.png)

![alt text](https://github.com/zhengchun/kindlepush/blob/master/docs/002.png)