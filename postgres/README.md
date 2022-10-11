## postgres

使用dbfactory包接口协议封装的实现

### 版本说明

### 定义模型

tag
1. postgres:"字段名"
2. pk:"" 主键
3. auto:"" 自增

#### 示例
```go
type test struct{
    ID string `postgres:"id" pk:"" auto:""`
    Name string `postgres:"name"`
}

// 联合主键未支持，但是日后实现方式如下
type test struct{
    ID string `postgres:"id" pk:"" other:"column:11,pk,auto"`
    Time time.time `postgres:"time" pk:""`
    Count int64 `postgres:"count"`
}
```

### 连接字符串
端口默认为 5432，有特殊配置自行个性
postgres://帐号:密码@host:5432/数据库

### 使用示例

```go

```

#### query

```go

```
