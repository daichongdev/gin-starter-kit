package logger

import (
	"context"
	"sync"
)

// SQLLog SQL执行日志
type SQLLog struct {
	SQL      string `json:"sql"`
	Duration string `json:"duration"`
	Rows     int64  `json:"rows"`
	Error    string `json:"error,omitempty"`
}

// TraceContext 链路追踪上下文
type TraceContext struct {
	TraceID string
	SQLLogs []SQLLog
	mutex   sync.RWMutex
}

type contextKey string

const (
	traceContextKey contextKey = "trace_context"
	traceIDKey      contextKey = "trace_id"
)

// SetTraceContext 设置链路追踪上下文到context
func SetTraceContext(ctx context.Context, traceID string) context.Context {
	traceCtx := &TraceContext{
		TraceID: traceID,
		SQLLogs: make([]SQLLog, 0),
	}
	// 同时写入 trace_context 和 trace_id，确保跨包读取一致
	ctx = context.WithValue(ctx, traceContextKey, traceCtx)
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	return ctx
}

// GetCurrentTraceID 从context获取链路追踪ID
func GetCurrentTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// AddSQLLog 添加SQL执行日志到context
func AddSQLLog(ctx context.Context, sql, duration string, rows int64, err error) {
	if traceCtx, ok := ctx.Value(traceContextKey).(*TraceContext); ok {
		sqlLog := SQLLog{
			SQL:      sql,
			Duration: duration,
			Rows:     rows,
		}
		if err != nil {
			sqlLog.Error = err.Error()
		}

		traceCtx.mutex.Lock()
		traceCtx.SQLLogs = append(traceCtx.SQLLogs, sqlLog)
		traceCtx.mutex.Unlock()
	}
}

// GetSQLLogs 从context获取SQL日志
func GetSQLLogs(ctx context.Context) []SQLLog {
	if traceCtx, ok := ctx.Value(traceContextKey).(*TraceContext); ok {
		traceCtx.mutex.RLock()
		defer traceCtx.mutex.RUnlock()

		// 返回副本避免并发问题
		logs := make([]SQLLog, len(traceCtx.SQLLogs))
		copy(logs, traceCtx.SQLLogs)
		return logs
	}
	return nil
}
