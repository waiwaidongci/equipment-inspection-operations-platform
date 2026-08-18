# Bug 复现

- Bug 是什么：设备概览快照和缓存直接持有仓库内部 map，后续更新或缓存调用方写入会改变历史快照。
- 如何触发：运行 `TestOverviewSnapshotsAreStable`，先读取概览，再更新设备并修改缓存返回值。
- 错误信息：`old snapshot changed after update` 或 `cache caller changed service snapshot`。
