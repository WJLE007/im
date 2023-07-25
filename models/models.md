# 集合列表

## 用户表

```text

    "account": "账号",
    "password": "密码",
    "nickname": "昵称",
    "sex": 1, //0是未知1是男2是女
    "email": "邮箱",
    "avatar": "头像",
    "creat_at": 1, //创建时间
    "updated_at": 1 // 更新时间
```

## 消息列表

```text
"user_identify":"用户唯一标识",
"room_identify":"房间唯一标识",
"data":"发送的数据",
"created_at":1,
"updated_at":1,

```

## 房间列表

```text
"number":"房间号",
"name":"房间名称",
"info":"房间简介",
"user_identify":"房主标识",
"created_at":1,
"updated_at":1,
```

## 用户房间关联列表

```text
{
  "user_identity": "用户的唯一标识",
  "room_identity": "房间的唯一标识",
  "message_identity": "消息的唯一标识",
  "created_at": 1, // 创建时间
  "updated_at": 1, // 更新时间
}
```