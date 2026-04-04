package api

import (
	"context"
	"os"
	"strconv"

	"golang.org/x/sync/semaphore"
)

// imageSem TASK-12: 控制图片生成最大并发数，防止 goroutine 无限累积
// 默认最多 10 个并发，可通过 IMAGE_GEN_CONCURRENCY 环境变量覆盖
var imageSem = newImageSemaphore()

func newImageSemaphore() *semaphore.Weighted {
	n := int64(10)
	if v := os.Getenv("IMAGE_GEN_CONCURRENCY"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed > 0 {
			n = parsed
		}
	}
	return semaphore.NewWeighted(n)
}

// acquireImageSlot TASK-12: 尝试获取图片生成并发槽位。
// 若槽位已满则返回 false，调用方应拒绝任务并退款。
func acquireImageSlot() bool {
	return imageSem.TryAcquire(1)
}

// releaseImageSlot TASK-12: 释放图片生成并发槽位。
func releaseImageSlot() {
	imageSem.Release(1)
}

// acquireImageSlotWait TASK-12: 阻塞等待获取槽位（用于高优先级任务）。
func acquireImageSlotWait(ctx context.Context) error {
	return imageSem.Acquire(ctx, 1)
}
