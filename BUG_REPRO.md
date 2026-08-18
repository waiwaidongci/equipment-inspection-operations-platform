# Bug 复现

- Bug 是什么：设备仓储包装错误时丢失未找到错误链，服务把不存在设备当成内部错误并触发重试。
- 如何触发：运行 `TestMissingDeviceReturnsNotFoundWithoutRetry`，让后端返回 `domain.ErrNotFound`。
- 错误信息：`status = 500, want 404`。
