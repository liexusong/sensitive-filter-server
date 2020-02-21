敏感词检测服务
=============
s
> 原理: 使用双数组字典树作为敏感词查询结构

### 编译
```
$ go mod init
$ go mod vendor
$ go build .
```

### 与正则对比
```
go test
```

### 启动服务
```
./sensitive-filter-server -source ./keywords -host 127.0.0.1 -port 8080
```

### API列表

#### 添加敏感词
* __接口__：`/api/add_word`
* __方法__：`GET`
* __入参__：

|参数名称|参数类型|参数说明|
|---|---|---|
|keyword|string|敏感词|

* __出参__：

```json
// 成功返回
{
    "code": 200,
    "resp": "OK",
    "data": null
}

// 错误返回
{
    "code":  201,
    "resp": "Parameters Invalid",
    "data": null,
}
```

#### 删除敏感词
* __接口__：`/api/del_word`
* __方法__：`GET`
* __入参__：

|参数名称|参数类型|参数说明|
|---|---|---|
|keyword|string|敏感词|

* __出参__：

```json
// 成功返回
{
    "code": 200,
    "resp": "OK",
    "data": null
}

// 错误返回
{
    "code":  201,
    "resp": "Parameters Invalid",
    "data": null,
}
```

#### 匹配一个敏感词
* __接口__：`/api/match_first`
* __方法__：`POST`
* __入参__：通过body传输要检测的文本

* __出参__：

```json
// 成功返回
{
    "code": 200,
    "resp": "OK",
    "data": "敏感词"
}

// 错误返回
{
    "code":  201,
    "resp": "Parameters Invalid",
    "data": null,
}
```

#### 匹配所有敏感词
* __接口__：`/api/match_all`
* __方法__：`POST`
* __入参__：通过body传输要检测的文本

* __出参__：

```json
// 成功返回
{
    "code": 200,
    "resp": "OK",
    "data": ["敏感词1", "敏感词2", "敏感词3", "敏感词4", "敏感词5"]
}

// 错误返回
{
    "code":  201,
    "resp": "Parameters Invalid",
    "data": null,
}
```