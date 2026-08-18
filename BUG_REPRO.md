# Bug 复现

- Bug 是什么：巡检批量派发在协程内部增加 WaitGroup 计数，结果通道又可能在 worker 结束前关闭，错误通道也没有可靠接收者。
- 如何触发：运行 `TestRunClosesAfterError`，同时提交一条正常任务和一条返回错误的任务。
- 错误信息：`panic: send on closed channel` 或测试超时。
