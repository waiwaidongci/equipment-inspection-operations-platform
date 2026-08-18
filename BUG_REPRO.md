# Bug 复现

- Bug 是什么：巡检队列过滤和报表拼接直接复用了切片底层数组，调用方的设备列表和已生成队列会互相改写。
- 如何触发：运行 `TestQueuePlanningKeepsSourceAndQueueIndependent`，先规划活动设备，再追加报表延后项并修改报表首项。
- 错误信息：`planning changed source candidates` 或 `report mutation leaked into queue`。
