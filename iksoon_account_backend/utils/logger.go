package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel 로그 레벨 타입
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// 로그 레벨 문자열 매핑
var logLevelStrings = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

// Logger 구조체
type Logger struct {
	level   LogLevel
	logger  *log.Logger
	prefix  string
	enabled bool
}

var (
	// 전역 로거 인스턴스
	defaultLogger *Logger
)

// init 함수에서 로거 초기화
func init() {
	// 환경변수에서 로그 레벨 설정 (기본값: INFO)
	level := getLogLevelFromEnv()

	defaultLogger = &Logger{
		level:   level,
		logger:  log.New(os.Stdout, "", 0), // 커스텀 포맷을 위해 기본 플래그 제거
		prefix:  "[ACCOUNT_API]",
		enabled: true,
	}
}

// getLogLevelFromEnv 환경변수에서 로그 레벨 가져오기
func getLogLevelFromEnv() LogLevel {
	levelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING", "WARN":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		return INFO // 기본값
	}
}

// formatLog 로그 메시지 포맷팅
func (l *Logger) formatLog(level LogLevel, message string) string {
	now := time.Now().Format("2006-01-02 15:04:05")

	// 호출자 정보 가져오기
	_, file, line, ok := runtime.Caller(3)
	caller := ""
	if ok {
		// 파일 경로에서 파일명만 추출
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			caller = fmt.Sprintf("%s:%d", parts[len(parts)-1], line)
		}
	}

	levelStr := logLevelStrings[level]

	if caller != "" {
		return fmt.Sprintf("%s [%s] %s [%s] %s",
			now, levelStr, l.prefix, caller, message)
	} else {
		return fmt.Sprintf("%s [%s] %s %s",
			now, levelStr, l.prefix, message)
	}
}

// shouldLog 로그를 출력할지 판단
func (l *Logger) shouldLog(level LogLevel) bool {
	return l.enabled && level >= l.level
}

// Debug 디버그 로그
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.shouldLog(DEBUG) {
		message := fmt.Sprintf(format, args...)
		l.logger.Println(l.formatLog(DEBUG, message))
	}
}

// Info 정보 로그
func (l *Logger) Info(format string, args ...interface{}) {
	if l.shouldLog(INFO) {
		message := fmt.Sprintf(format, args...)
		l.logger.Println(l.formatLog(INFO, message))
	}
}

// Warning 경고 로그
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.shouldLog(WARNING) {
		message := fmt.Sprintf(format, args...)
		l.logger.Println(l.formatLog(WARNING, message))
	}
}

// Error 에러 로그
func (l *Logger) Error(format string, args ...interface{}) {
	if l.shouldLog(ERROR) {
		message := fmt.Sprintf(format, args...)
		l.logger.Println(l.formatLog(ERROR, message))
	}
}

// SetLevel 로그 레벨 설정
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetEnabled 로거 활성화/비활성화
func (l *Logger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

// 전역 함수들 - 기본 로거 사용
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warning(format string, args ...interface{}) {
	defaultLogger.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// SetLogLevel 전역 로그 레벨 설정
func SetLogLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetLogEnabled 전역 로거 활성화/비활성화
func SetLogEnabled(enabled bool) {
	defaultLogger.SetEnabled(enabled)
}

// GetLogLevel 현재 로그 레벨 반환
func GetLogLevel() LogLevel {
	return defaultLogger.level
}

// IsDebugEnabled 디버그 로그가 활성화되어 있는지 확인
func IsDebugEnabled() bool {
	return defaultLogger.shouldLog(DEBUG)
}

// IsInfoEnabled 정보 로그가 활성화되어 있는지 확인
func IsInfoEnabled() bool {
	return defaultLogger.shouldLog(INFO)
}

// LogHTTPRequest HTTP 요청 로그 (개발환경에서만)
func LogHTTPRequest(method, path, remoteAddr string) {
	if IsInfoEnabled() {
		Info("HTTP %s %s from %s", method, path, remoteAddr)
	}
}

// LogHTTPResponse HTTP 응답 로그 (개발환경에서만)
func LogHTTPResponse(method, path string, statusCode int, duration time.Duration) {
	if IsInfoEnabled() {
		Info("HTTP %s %s responded with %d in %v", method, path, statusCode, duration)
	}
}

// LogDatabaseQuery 데이터베이스 쿼리 로그 (디버그 레벨)
func LogDatabaseQuery(query string, args ...interface{}) {
	if IsDebugEnabled() {
		if len(args) > 0 {
			Debug("DB Query: %s, Args: %v", query, args)
		} else {
			Debug("DB Query: %s", query)
		}
	}
}

// LogDatabaseError 데이터베이스 에러 로그
func LogDatabaseError(operation string, err error) {
	Error("Database %s failed: %v", operation, err)
}

// LogError 에러 로그 (스택 트레이스 포함)
func LogError(operation string, err error) {
	Error("%s failed: %v", operation, err)
}

// LogStartup 시작 로그
func LogStartup(port string) {
	Info("서버가 %s 포트에서 실행 중입니다...", port)
	Info("새로운 구조의 API가 적용되었습니다.")
	Info("로그 레벨: %s", logLevelStrings[defaultLogger.level])
}
