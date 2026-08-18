# Bug 复现

- Bug 是什么：巡检任务状态机缺少 retrying 到 completed 的转换，重试成功仍回写旧状态且列表漏掉重试中任务。
- 如何触发：运行 `TestRetrySuccessCompletesAndIsVisible`，验证重试成功后的终态和进行中列表。
- 错误信息：`retrying task cannot complete`。
