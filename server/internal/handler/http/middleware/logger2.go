package middleware

import (
    "log/slog"
    "net/http"
    "time"
    
    "github.com/google/uuid"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    bytes      int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    size, err := rw.ResponseWriter.Write(b)
    rw.bytes += size
    return size, err
}

// RequestLogger логирует каждый HTTP запрос
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Генерируем request ID
            requestID := uuid.New().String()
            
            // Добавляем request ID в контекст
            ctx := r.Context()
            ctx = context.WithValue(ctx, "request_id", requestID)
            r = r.WithContext(ctx)
            
            // Оборачиваем response writer для захвата статуса
            wrapped := &responseWriter{
                ResponseWriter: w,
                statusCode:     http.StatusOK,
            }
            
            start := time.Now()
            
            // Выполняем запрос
            next.ServeHTTP(wrapped, r)
            
            // Логируем после выполнения
            duration := time.Since(start)
            
            logger.Info("HTTP request",
                slog.String("request_id", requestID),
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.String("remote_addr", r.RemoteAddr),
                slog.Int("status", wrapped.statusCode),
                slog.Int("bytes", wrapped.bytes),
                slog.Duration("duration", duration),
                slog.String("user_agent", r.UserAgent()),
            )
        })
    }
}
