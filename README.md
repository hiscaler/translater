# translater
文本翻译接口

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
to |  	目标语言[语种简码](http://deepi.sogou.com/docs/fanyiDoc#lan) | 无
q |  需要翻译的文本（utf8编码待翻译文本）| 无

## 搜狗开发文档
http://deepi.sogou.com/docs/fanyiDoc
