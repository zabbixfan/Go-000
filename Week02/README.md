学习笔记# 学习笔记

## Error vs Exception

### Exception In Other Language

-   Error In C
    -   单返回值，通过指针入参获取数据，返回值 int 表示成功或失败
-   Error In C++
    -   引入 exception 但无法知道被调用方抛出的是什么类型的异常
-   Error In Java
    -   引入 checked exception 但不同的使用者会有不同处理方法，变得太司空见惯，严重程度只能人为区分，
        并且容易被使用者滥用，如经常 catch (e Exception) { // ignore }

### Error In Go

-   Go 中的 error 只是一个普通的 interface 包含一个 Error() string 方法
-   使用 errors.New() 创建一个 error 对象，返回的是 errorString 结构体的指针
-   利用指针 在基础库内部预定义了大量的 error, 用于返回及上层 err 与预定义的 err 做对比，
    预防 err 文本内容一致但实际意义及环境不同的两个 err 对比成功。
-   go 支持多参数返回，一般最后一个参数是 err，必须先判断 err 才使用 value，除非你不关心 value，即可忽略 err.
-   go 的 panic 与别的语言的 exception 不一样，需谨慎或不使用，一般在 api 中第一个 middleware 就是 recover panic.
-   野生 goroutine 如果 panic 无法被 recover， 需要构造一个 func Go(x func()) 在其内部 recover
-   强依赖、配置文件: panic , 弱依赖: 不需要 panic
    -   Q1: 案例: 如果数据库连不上但 redis 连得上,是否需要 panic.
    -   A1: 取决于业务，如果读多写少，可以先不 panic，等待数据库重连。
    -   Q2: 案例: 服务更新中导致 gRPC 初始化的 client 连不上
    -   A2: 也是看业务，如果 gRPC 是 Blocking(阻塞):等待重连、nonBlocking(非阻塞):立刻返回一个 default、
        nonBlocking+timeout(非阻塞+超时/推荐):先尝试重连如果超时返回 default
-   只有真正意外、不可恢复的程序错误才会使用 panic , 如 索引越界、不可恢复的环境问题、栈溢出，才使用 panic。除此之外都是 error。
-   go error 特点:
    -   简单
    -   Plan for failure not success
    -   没有隐藏的控制流
    -   完全交给你来控制 error
    -   Error are values

## Error Type

### Sentinel Error

-   sentinel error: 预定义错误,特定的不可能进行进一步处理的做法
-   if err == ErrSomething { ... } 类似的 sentinel error 比如: io.EOF、syscall.ENOENT
-   最不灵活，必须利用==判断，无法提供上下文。只能利用 error.Error()查看错误输出。
-   会变成 API 公共部分
    -   增加 API 表面积
    -   所有的接口都会被限制为只能返回该类型的错误，即使可以提供更具描述性的错误
-   在两个包中间产生依赖关系：无法二次修改现在包所返回的 error，存在高耦合、无法重构
-   **总结:尽可能避免 sentinel errors**

### Error types

-   Error type 是实现了 error 接口的自定义类型，可以自定义需要的上下文及各种信息
-   Error type 是一个 type 所以可以被**断言**用来获取更多的上下文信息
-   VS Sentinel Error
    -   Error type 可以提供更多的上下文
    -   一样会 public，与调用者产生强耦合，导致 API 变得脆弱。
    -   也需要尽量避免 Error types

### Opaque errors (最标准、建议的方法)

-   只知道对或错 只能 err != nil
-   但无法携带上下文信息
-   Assert errors for behaviour, not type
    -   通过定一个 interface ，然后暴露相关的 Is 方法去判断 err，调用库内部断言。
-   **具体选择哪种还是得需要看场景**

## Handing Error

### Indented flow is for errors

-   err != nil 而不是 err == nil

### Eliminate error handling by eliminating errors

-   代码编写时可以直接返回 err 的 别用 err != nil
-   利用已经封装好的方法去消除代码中的 err
