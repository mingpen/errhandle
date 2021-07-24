# 错误处理demo

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

# 解

应该用Wrap。
sql.ErrNoRows 直接抛给调用方，缺少错误的上下文信息，不便于调试。

1. 定义特殊的错误变量 ErrNotFound；
2. 调用sql接口遇到 sql.ErrNoRows，fmt.Errorf %w 做 Wrap an error；
3. 定义 QueryErr 对错误信息进行更多的描述，定义Error、Unwrap、As方法;
4. 以及增加 Is 方法，
当 参数值为 ErrNoRows ，且err本身Wrap 了ErrNotFound，就返回true；
增加该方法是防止，之前调用者使用了 Sentinel Error（ErrNoRows）进行判断。
