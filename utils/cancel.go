package utils

import "sync"

var (
	// IsCanceled 标记任务是否被取消
	IsCanceled bool
	// CancelLock 保护取消标记的锁
	CancelLock sync.Mutex
)

// ResetCancel 重置取消标记
func ResetCancel() {
	CancelLock.Lock()
	defer CancelLock.Unlock()
	IsCanceled = false
}

// SetCancel 设置取消标记
func SetCancel() {
	CancelLock.Lock()
	defer CancelLock.Unlock()
	IsCanceled = true
}

// CheckCanceled 检查是否取消
func CheckCanceled() bool {
	CancelLock.Lock()
	defer CancelLock.Unlock()
	return IsCanceled
}
