af# 学习笔记

## 作业

**在 work 目录下哦 😉**

## package sync

### Share Memory By Communicating

-   go 鼓励[用 chan 而不是锁][1]
-   Do not communicate by sharing memory, instead, share memory by communicating.

### Data Race Conditions With Go

-   data race 表示并发读写同一个资源导致的随机性 bug
    -   代码排查困难
    -   如何检测
        1. go build -race
        2. go test -race
-   通过分析并发[`count++`][2]了解 data race
-   [interface 的 data race][3]
    -   与 interface 的底层结构有关, 当两个结构体的底层结构一致时 interface 不会发生 data race
    -   有可能 interface 的 type 在 type1 data 在 type2
    -   Q: 如果是一个普通的指针、map、slice 可以安全赋值吗?  
        A: 普通指针、map 为 8B 可以安全赋值, slice 不行.
-   没有 safe data race. 程序要么没有 data race, 要么其操作未定义
-   锁设计原则: 最晚加锁、最早释放、锁内内容少、轻量、注意操作顺序避免死锁
-   案例: [for 循环 map 产生死锁][4]

### sync.atomic

-   案例代码中，cfg 为全局对象，同时被多个 goroutine 访问，存在 data race，使用 go 同步语义解决。
    -   [Mutex][5]
    -   [RWMutex][5]
    -   [Atomic][5]
-   **Copy-On-Write**
    -   **redis 的 BGSAVE**
    -   微服务中的定时更新: 风险-读写数据不一致，无法避免，但关系不大。
    -   微服务降级、local cache

### Mutex

-   [锁饥饿][6]
    1. g1 在获取锁后休眠 100ms，当 g2 试图获取锁时，将被添加到等待队列进行等待。
    2. 当 g1 完成工作时，释放锁，此时唤醒 g2，g2 等待运行，而此刻 g1 再次占有锁，g2 无奈又进入等待。
-   Mutex 锁的实现
    -   Barging: 提高了吞吐量，但不公平
    -   Hands-off: 吞吐量有所降低，但公平
    -   Spinning: 性能开销大
    -   Go 1.8 使用了 Barging 和 Spinning 结合实现，自旋几次后就会 park
    -   Go 1.9 添加了饥饿模式，如果等待锁 1ms, unlock 会 hands-off 把锁丢给第一个等待者,此时同样代码 g1:57 g2:10

### [errGroup][7]

-   核心原理
    -   利用 sync.WaitGroup 管理并执行 goroutine
-   主要功能
    -   并行工作流
    -   处理错误 或者 优雅降级
    -   context 传播与取消
    -   利用局部变量+闭包
-   设计缺陷 --- [改进][8]
    -   没有捕获 panic，导致程序异常退出 --- 改进 加 defer recover
    -   没有限制 goroutine 数量，存在大量创建 goroutine --- 改进 增加一个 channel 用来消费 func
    -   WithContext 返回的 context 可能被异常调用，当其在 errgroup 中被取消时，影响其它函数 --- 改进 代码内嵌 context

### sync.Pool

-   保存与复用临时对象
-   降低 GC 压力
-   不能放链接类型，有可能导致链接泄漏

## chan

### Channels

-   channels 是一种类型安全的消息队列，goroutine 之间的管道，创建 Go 同步机制
-   unbuffered Channels
    -   发送方在没有接收方时会阻塞
    -   接收方先结束发送方才结束
    -   好处：**100%保证收到**
    -   代价：延迟时间未知
-   buffer
    -   发送方在管道满时会阻塞
    -   发送方 happen before 接收方
    -   好处：延迟小
    -   代价：不保证数据到达、越大的 buffer，越小的保障到达。buffer = 1 时，给你一个延迟一个消息的保障

### Go Concurrency Patterns

-   [Timing out][10]
-   [Moving on][10]
-   [Pipeline][11]
-   [Fan-out,Fan-in][11]
-   [Cancellation][11]
-   [Context][12]
-   **一定要交给发送方 close chan**

### Design Philosophy

-   If any given Send on a channel CAN cause the sending goroutine to block:
    -   Not allowed to use a Buffered channel larger than 1.
        -   Buffers larger than 1 must have reason/measurements.
    -   Must know what happens when the sending goroutine blocks.
-   If any given Send on a channel WON’T cause the sending goroutine to block:
    -   You have the exact number of buffers for each send.
        -Fan Out pattern
    -   You have the buffer measured for max capacity.
        -Drop pattern
-   Less is more with buffers.
    -   Don’t think about performance when thinking about buffers.
    -   Buffers can help to reduce blocking latency between signaling.
        -   Reducing blocking latency towards zero does not necessarily mean better throughput.
        -   If a buffer of one is giving you good enough throughput then keep it.
        -   Question buffers that are larger than one and measure for size.
        -   Find the smallest buffer possible that provides good enough throughput.

## context

### Request-scoped context

-   实现传递数据，搞定超时控制，或者级联取消(显示传递)
-   context 集成到 API
    -   函数首参为 context
    -   创建对象时携带 context 对象: WithContext

### Don't store Contexts inside a struct type

-   不要把 context 放到结构体里，然后再把结构体当参数传输

### context.WithValue

-   从子向父递归查询 key-value
-   Background、TODO
-   **Debugging or tracing data is safe to pass in a Context**
-   context.WithValue 只读、安全 --- 染色、API 重要性、Trace
-   禁止在 context 中挂载与业务逻辑耦合的东西，不能放一些奇奇怪怪的东西进去
-   如果有必要修改 context 的内容，请使用 COW:
    1. 从源 ctx 获取到 v1
    2. 复制 v1 到 v2
    3. 修改 v2
    4. 将 v2 重新挂载到 ctx,产生 ctx2
    5. 将 ctx2 向下传递
-   ~~gin 的 context.Next 有缺陷，应~~参考 grpc 的 middleware
-   计算密集型耗时短，一般不处理超时。
-   go 标准网络库可被托管，~~吊打其它语言业务、中间件，~~不会因为超时导致 oom。[kratos 案例][9]
-   当一个 context 被 cancel 时，所有子 context 都会被 cancel
-   一定要 cancel 否者 context 会泄漏

### **Final Notes**

-   Incoming requests to a server should create a Context.
-   Outgoing calls to servers should accept a Context.
-   Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it.
-   The chain of function calls between them must propagate the Context.
-   Replace a Context using WithCancel, WithDeadline, WithTimeout, or WithValue.
-   When a Context is canceled, all Contexts derived from it are also canceled.
-   The same Context may be passed to functions running in different goroutines; Contexts are safe for simultaneous use by multiple goroutines.
-   Do not pass a nil Context, even if a function permits it. Pass a TODO context if you are unsure about which Context to use.
-   Use context values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
-   All blocking/long operations should be cancelable.
-   Context.Value obscures your program’s flow.
-   Context.Value should inform, not control.
-   Try not to use context.Value.

[1]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo1/main.go
[2]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo2/README.md
[3]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo3/main.go
[4]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo4/main.go
[5]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo5/README.md
[6]: https://github.com/XYZ0901/Go-000/blob/main/Week03/demo/demo6/main.go
[7]: https://pkg.go.dev/golang.org/x/sync/errgroup
[8]: https://github.com/go-kratos/kratos/blob/master/pkg/sync/errgroup/errgroup.go
[9]: https://github.com/go-kratos/kratos/blob/master/pkg/cache/redis/conn.go#L519
[10]: https://blog.golang.org/concurrency-timeouts
[11]: https://blog.golang.org/pipelines
[12]: https://blog.golang.org/context
