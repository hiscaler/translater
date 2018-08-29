# translater
文本翻译接口

# 使用方法
## 开启服务
```shell
./server
```

# 配置
## 配置参数示例
```json
{
  "accounts": [
    {
      "PID": "",
      "SecretKey": "",
      "Enabled": true,
      "YearMonth": 0
    },
    ...
  ],
  "debug": true,
  "languages": {
    "ar": "阿拉伯语",
    "et": "爱沙尼亚语",
    "bg": "保加利亚语",
    "pl": "波兰语",
    "ko": "韩语",
    "bs-Latn": "波斯尼亚语",
    "fa": "波斯语",
    "mww": "白苗文",
    "da": "丹麦语",
    "de": "德语",
    "ru": "俄语",
    "fr": "法语",
    "fi": "芬兰语",
    "tlh-Qaak": "克林贡语(piqaD)",
    "tlh": "克林贡语",
    "hr": "克罗地亚语",
    "otq": "克雷塔罗奥托米语",
    "ca": "加泰隆语",
    "cs": "捷克语",
    "ro": "罗马尼亚语",
    "lv": "拉脱维亚语",
    "ht": "海地克里奥尔语",
    "lt": "立陶宛语",
    "nl": "荷兰语",
    "ms": "马来语",
    "mt": "马耳他语",
    "pt": "葡萄牙语",
    "ja": "日语",
    "sl": "斯洛文尼亚语",
    "th": "泰语",
    "tr": "土耳其语",
    "sr-Latn": "塞尔维亚语(拉丁文)",
    "sr-Cyrl": "塞尔维亚语(西里尔文)",
    "sk": "斯洛伐克语",
    "sw": "斯瓦希里语",
    "af": "南非荷兰语",
    "no": "挪威语",
    "en": "英语",
    "es": "西班牙语",
    "uk": "乌克兰语",
    "ur": "乌尔都语",
    "el": "希腊语",
    "hu": "匈牙利语",
    "cy": "威尔士语",
    "yua": "尤卡坦玛雅语",
    "he": "希伯来语",
    "zh-CHS": "中文",
    "it": "意大利语",
    "hi": "印地语",
    "id": "印度尼西亚语",
    "zh-CHT": "中文繁体",
    "vi": "越南语",
    "sv": "瑞典语",
    "yue": "粤语(繁体)",
    "fj": "斐济",
    "fil": "菲律宾语",
    "sm": "萨摩亚语",
    "to": "汤加语",
    "ty": "塔希提语",
    "mg": "马尔加什语",
    "bn": "孟加拉语"
  },
  "listenport": "8080"
}
```
## 参数说明
参数 | 说明 | 默认值
---- | --- | ---
debug | 是否开启调试模式（可选值为：true, false） | true
listenport | 监听端口 | 8080
accounts | 翻译帐号 [申请地址](http://deepi.sogou.com/contact/fanyi)
languages | 支持的翻译语种 [语种简码](http://deepi.sogou.com/docs/fanyiDoc#lan)

## 接口

### 测试接口状态
GET /api/ping

正常状态下会返回 "OK" 字样

### 访问接口
#### 地址
POST /api/translate?from=en&to=zh-CHS

#### 参数说明
参数 | 说明 | 默认值
---- | --- | ----
from | 源语言[语种简码](http://deepi.sogou.com/docs/fanyiDoc#lan) | auto
to | 目标语言[语种简码](http://deepi.sogou.com/docs/fanyiDoc#lan) | 无
q | 需要翻译的文本（utf8编码待翻译文本） | 无

## 搜狗开发文档
http://deepi.sogou.com/docs/fanyiDoc
