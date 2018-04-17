KindlePush
===
Subscribe the RSS/Atom feeds and delivered to your Kindle.

[FeedOcean.com](https://feedocean.com): Convert HTML to RSS and Atom feeds with full-text, Its free.

### How to Use

- Download a [KindleGen](https://www.amazon.com/gp/feature.html?docId=1000765211) software and install it into your computer.

- Edit `config.yaml` file.

- Run `kindlepush` command.

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

Tips: Make sure your `senderAddress` already add to the approved e-mail address list. See [Add an Email Address to Receive Documents](https://www.amazon.com/gp/help/customer/display.html?nodeId=201974240).


KindlePush
===
KindlePush能将多个订阅的RSS/Atom源生成杂志模式的MOBI文件，并推送到你的Kindle设备。

### 使用说明

- 下载[KindleGen](https://www.amazon.com/gp/feature.html?docId=1000765211)并安装到本地机器。
- 编辑 `config.yaml` 配置文件。
- 运行 `kindlepush` 命令。

请注意：正常接收投递电子杂志，你必须确保你的`senderAddress`邮件地址已经加入到kindle的邮箱列中。[Add an Email Address to Receive Documents](https://www.amazon.com/gp/help/customer/display.html?nodeId=201974240)。

### config.yaml说明

- 必选
    * kindleAddress - Kindle email address.
    * smtp - About SMTP configuration used to send email. 
    * feeds - List of Atom/RSS feed URLs.
- 可选
    * cacheDir - Output directory.
    * resizeImage - Resize image with specifies width and height.
    * maxFileSize - Max mobi file size.
    * kindlegen - kindlegen.exe path.

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


Welcome to contribute and make is too better.


![alt text](https://github.com/zhengchun/kindlepush/blob/master/docs/001.png)

![alt text](https://github.com/zhengchun/kindlepush/blob/master/docs/002.png)