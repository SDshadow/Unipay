# Unipay
Payment center  支付中台 微信支付 支付宝
DDD风格
Gin+GORM+gRPC+rocketmq+nacos支付中台项目
数据库：MySQL+redis+redlock（redsync）
鉴权：APP key
日志：logrus（优化性能用zap）
监控：Prometheus
熔断：Sentinel
计算库：shopspring/decimal


pay-center
├── api
│   ├── http
│   │   ├── payment_handler.go
│   │   └── callback_handler.go
│   └── grpc
│       └── payment.proto
│
├── application                # 核心
│   ├── payment
│   │   ├── create_order.go    # 下单用例
│   │   ├── pay.go             # 发起支付
│   │   ├── callback.go        # 支付回调
│   │   └── query.go           # 查询用例（读模型）
│   │
│   └── dto                    # 入参 / 出参
│       └── payment_dto.go
│
├── domain
│   └── payment
│       ├── order.go           # 聚合根
│       ├── status.go          # 状态机
│       ├── money.go           # 值对象
│       └── repository.go      # 接口
│
├── channel                    # 支付通道能力层
│   ├── channel.go             # 接口定义
│   ├── wechat.go
│   ├── alipay.go
│   └── router.go              # 简单路由
│
├── infrastructure
│   ├── persistence
│   │   └── payment_order_repo.go
│   ├── mq
│   │   └── callback_consumer.go
│   ├── cache
│   └── config
│
├── middleware                 # Gin / gRPC 中间件
│   ├── auth.go
│   ├── idempotent.go
│   └── limiter.go
│
├── common
│   ├── errors
│   ├── logger
│   └── utils
│
├── cmd
│   └── pay-center
│       └── main.go
│
└── docs
