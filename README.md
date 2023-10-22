# smart-pb
智能数据传输协议，当pb变更后，无需通知使用方，可自动感知当前message的格式

主要原理是：
- 1: 传输的真正包是`Payload`
- 2: `Payload`有两个字段，data用于存储消息内容，pmd用于存储消息格式
- 3: 使用`Marshal`方法把整个`Payload`变为`Message`对象
- 4: 对`Message`再使用`Unmarshal`方法把包序列化发送给使用方
- 5: 使用方收到包后，可以使用`Sink`方法把整个包反序列化为`Message`对象
- 6: 使用方可以使用`GetXXX`系列方法，可以读取`Message`中的数据值

```protobuf
// 消息载体
message payload {
  bytes pmd = 1;  // proto message desc
  bytes data = 2; // message 内容
}
```

## 2、主要方法介绍

具体使用方式可以见 `message_test.go` 中单元测试

### 2.1、proto.Message -> Message

> message := Marshal(msg) // proto.Message -> Message
> buf := Unmarshal(message) // Message -> []byte
> message := Sink(buf) // []byte -> Message

### 2.2、Get方法

```golang
var api = &Api4{
    Users: []*Api2{
        {Id: 100, Info: &Api1{Name: "aa", Email: "aa@t.com"}},
        {Id: 101, Info: &Api1{Name: "bb", Email: "bb@t.com"}},
        {Id: 102, Info: &Api1{Name: "cc", Email: "cc@t.com"}},
    },
    Admin: &Api4Api6{
        V2: &Api2{
            Id: 12345,
            Info: &Api1{
                Name:  "tx",
                Email: "k@te.com",
            },
        },
    },
}

msg := Marshal(api)
Convey("Get\n", t, func() {
    var adminName string
    err := msg.Get("admin.v2.info.name", &adminName)
    So(err, ShouldBeNil)
    So(adminName, ShouldEqual, "tx")
})

Convey("Get_pb_array\n", t, func() {
    var users []*Api2
    err := msg.Get("users", &users)
    So(err, ShouldBeNil)
    So(users, ShouldHaveLength, 3)
    t.Logf("users:%+v", users)
})
``` 

### 2.3、GetXXX方法

```golang
var api = &Api4{
    Users: []*Api2{
        {Id: 100, Info: &Api1{Name: "aa", Email: "aa@t.com"}},
        {Id: 101, Info: &Api1{Name: "bb", Email: "bb@t.com"}},
        {Id: 102, Info: &Api1{Name: "cc", Email: "cc@t.com"}},
    },
    Admin: &Api4Api6{
        V2: &Api2{
            Id: 12345,
            Info: &Api1{
                Name:  "tx",
                Email: "k@te.com",
            },
        },
    },
}

msg := Marshal(api)
Convey("GetInt64\n", t, func() {
    id, err := msg.GetInt64("users[1].id")
    So(err, ShouldBeNil)
    So(id, ShouldEqual, 101)
})

Convey("GetInt64_nested\n", t, func() {
    adminId, err := msg.GetInt64("admin.v2.id")
    So(err, ShouldBeNil)
    So(adminId, ShouldEqual, 12345)
})

Convey("GetString_array\n", t, func() {
    email, err := msg.GetString("users[2].info.email")
    So(err, ShouldBeNil)
    So(email, ShouldEqual, "cc@t.com")
})
```