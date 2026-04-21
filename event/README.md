# EventBus 使用指南

## 概述

EventBus 是一个基于 Channel 驱动的事件总线，支持泛型和可切换的底层驱动。

## 架构

```
contracts.EventDriver (驱动接口)
    ↓
event.ChannelDriver (默认 Channel 实现)
    ↓
event.EventBus[T] (泛型事件总线)
    ↓
event.Get[T]("driver") (获取事件总线实例)
```

## 特性

- **泛型支持**: 支持任意类型的事件
- **Channel 驱动**: 使用 Go Channel 实现高性能事件传递
- **可切换驱动**: 通过 `contracts.EventDriver` 接口支持不同驱动实现
- **主题订阅**: 支持按主题(topic)订阅事件
- **多订阅者**: 同一主题可有多个订阅者
- **线程安全**: 内置互斥锁保证并发安全

## 基本使用

### 1. 注册驱动 (在 config 中)

```go
import (
    "github.com/xframe-go/x/event"
)

func registerEvent() {
    event.Register("default", event.NewChannelDriver())
}
```

### 2. 定义事件类型

```go
type UserCreatedEvent struct {
    UserID   uint64
    Username string
    Email    string
}
```

### 3. 获取事件总线

```go
import "github.com/xframe-go/x/event"

bus := event.Get[UserCreatedEvent]("default")
```

### 4. 订阅事件

```go
bus.Subscribe("user.created", func(event UserCreatedEvent) {
    fmt.Printf("User created: %s (%s)\n", event.Username, event.Email)
})
```

### 5. 发布事件

```go
bus.Publish("user.created", UserCreatedEvent{
    UserID:   1,
    Username: "testuser",
    Email:    "test@example.com",
})
```

## 完整示例

### 注册和初始化

```go
// internal/config/event.go
package config

import (
    "github.com/xframe-go/x/event"
)

func registerEvent() {
    event.Register("default", event.NewChannelDriver())
}

// internal/config/app.go
package config

func Register() {
    registerDB()
    registerStorage()
    registerAuth()
    registerEvent()  // 注册事件驱动
}
```

### 在业务代码中使用

```go
package api

import (
    "github.com/xframe-go/x/event"
)

type UserCreatedEvent struct {
    UserID   uint64 `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (h *Traveler) Register(c echo.Context) error {
    // ... 注册逻辑 ...

    // 发布用户创建事件
    bus := event.Get[UserCreatedEvent]("default")
    bus.Publish("user.created", UserCreatedEvent{
        UserID:   traveler.ID,
        Username: traveler.Username,
        Email:    traveler.Email,
    })

    return h.Success(c, detail)
}
```

### 订阅事件处理

```go
package listeners

import (
    "github.com/xframe-go/x/event"
)

func RegisterUserListeners() {
    bus := event.Get[UserCreatedEvent]("default")

    bus.Subscribe("user.created", func(event UserCreatedEvent) {
        // 发送欢迎邮件
        sendWelcomeEmail(event.Email)

        // 记录审计日志
        logAuditEvent("user.created", event.UserID)

        // 更新统计数据
        updateUserStats()
    })

    bus.Subscribe("user.created", func(event UserCreatedEvent) {
        // 另一个监听器：创建用户个人空间
        createUserSpace(event.UserID)
    })
}
```

## 高级用法

### 自定义驱动

```go
package custom_driver

import (
    "github.com/xframe-go/x/contracts"
)

type RedisDriver struct {
    redisClient *redis.Client
}

func (d *RedisDriver) Publish(topic string, data interface{}) error {
    jsonData, _ := json.Marshal(data)
    return d.redisClient.Publish(ctx, topic, jsonData).Err()
}

func (d *RedisDriver) Subscribe(topic string) (<-chan interface{}, error) {
    ch := make(chan interface{}, 100)

    // 实现订阅逻辑
    go func() {
        pubsub := d.redisClient.Subscribe(ctx, topic)
        for msg := range pubsub.Channel() {
            ch <- msg.Payload
        }
    }()

    return ch, nil
}

func (d *RedisDriver) Close() error {
    return d.redisClient.Close()
}
```

### 注册自定义驱动

```go
func registerEvent() {
    redisDriver := &RedisDriver{
        redisClient: redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
    }
    event.Register("redis", redisDriver)
}
```

### 使用自定义驱动

```go
bus := event.Get[UserCreatedEvent]("redis")
bus.Subscribe("user.created", handler)
```

## API 参考

### EventDriver 接口

```go
type EventDriver interface {
    Publish(topic string, data interface{}) error
    Subscribe(topic string) (<-chan interface{}, error)
    Close() error
}
```

### EventBus[T] 方法

```go
// 订阅主题
Subscribe(topic string, handler func(T)) error

// 发布事件
Publish(topic string, data T) error

// 关闭事件总线
Close() error
```

### 全局函数

```go
// 注册驱动
event.Register(name string, driver contracts.EventDriver)

// 获取事件总线实例
event.Get[T any](name string) *EventBus[T]
```

## 最佳实践

1. **事件命名**: 使用 `resource.action` 格式，如 `user.created`, `post.published`
2. **事件结构**: 使用结构体而非基本类型，方便扩展
3. **错误处理**: 在订阅者函数内部处理错误，不要影响其他订阅者
4. **性能考虑**: 避免在订阅者中执行耗时操作
5. **解耦合**: 事件发布者不需要知道订阅者的存在
6. **测试友好**: 使用 MockDriver 进行单元测试

## 注意事项

- Channel 默认缓冲大小为 100，可通过修改 `ChannelDriver` 调整
- 订阅者函数中发生 panic 不会影响其他订阅者
- 建议在应用启动时注册所有事件监听器
- 关闭 EventBus 会关闭所有底层 channel
