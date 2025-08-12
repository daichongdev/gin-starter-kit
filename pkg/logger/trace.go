package logger

import (
	"sync"
)

// TraceContext 链路追踪上下文
type TraceContext struct {
	TraceID string
	SQLLogs []SQLLog
	mutex   sync.RWMutex
}

// SQLLog SQL执行日志
type SQLLog struct {
	SQL      string `json:"sql"`
	Duration string `json:"duration"`
	Rows     int64  `json:"rows"`
	Error    string `json:"error,omitempty"`
}

var (
	// 使用context来存储链路追踪信息
	traceContextKey = "trace_context"
	// 当前goroutine的链路追踪上下文
	currentTraceContext *TraceContext
	traceContextMutex   sync.RWMutex
)

// SetTraceContext 设置当前goroutine的链路追踪上下文
func SetTraceContext(traceID string) {
	traceContextMutex.Lock()
	defer traceContextMutex.Unlock()

	currentTraceContext = &TraceContext{
		TraceID: traceID,
		SQLLogs: make([]SQLLog, 0),
	}
}

// GetCurrentTraceID 获取当前goroutine的链路追踪ID
func GetCurrentTraceID() string {
	traceContextMutex.RLock()
	defer traceContextMutex.RUnlock()

	if currentTraceContext != nil {
		return currentTraceContext.TraceID
	}
	return ""
}

// AddSQLLog 添加SQL执行日志到当前链路追踪
func AddSQLLog(sql, duration string, rows int64, err error) {
	traceContextMutex.Lock()
	defer traceContextMutex.Unlock()

	if currentTraceContext != nil {
		sqlLog := SQLLog{
			SQL:      sql,
			Duration: duration,
			Rows:     rows,
		}
		if err != nil {
			sqlLog.Error = err.Error()
		}

		currentTraceContext.mutex.Lock()
		currentTraceContext.SQLLogs = append(currentTraceContext.SQLLogs, sqlLog)
		currentTraceContext.mutex.Unlock()
	}
}

// GetSQLLogs 获取当前链路追踪的SQL日志
func GetSQLLogs() []SQLLog {
	traceContextMutex.RLock()
	defer traceContextMutex.RUnlock()

	if currentTraceContext != nil {
		currentTraceContext.mutex.RLock()
		defer currentTraceContext.mutex.RUnlock()

		// 返回副本避免并发问题
		logs := make([]SQLLog, len(currentTraceContext.SQLLogs))
		copy(logs, currentTraceContext.SQLLogs)
		return logs
	}
	return nil
}

// ClearTraceContext 清除当前goroutine的链路追踪上下文
func ClearTraceContext() {
	traceContextMutex.Lock()
	defer traceContextMutex.Unlock()

	currentTraceContext = nil
}
