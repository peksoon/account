package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Config 애플리케이션 설정 구조체
type Config struct {
	// 서버 설정
	Port string `env:"PORT"`

	// 데이터베이스 설정
	DBPath string `env:"DB_PATH"`

	// 로깅 설정
	LogLevel string `env:"LOG_LEVEL"`

	// 기타 설정
	MaxConnections int `env:"MAX_CONNECTIONS"`
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig 싱글톤 패턴으로 설정 인스턴스 반환
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			// 기본값 설정
			Port:           "8080",
			DBPath:         "./data/account_app.db",
			LogLevel:       "INFO",
			MaxConnections: 100,
		}
		instance.loadFromEnvFile()
		instance.loadFromEnvironment()
	})
	return instance
}

// loadFromEnvFile 환경별 .env 파일에서 설정 로드
func (c *Config) loadFromEnvFile() {
	// 환경별 설정 파일 경로 결정
	configFiles := []string{
		"config.env.production",  // 우선순위 1: 운영 설정
		"config.env.development", // 우선순위 2: 개발 설정
		"config.env",             // 우선순위 3: 기본 설정
	}

	for _, configFile := range configFiles {
		if err := c.loadEnvFile(configFile); err == nil {
			fmt.Printf("설정 파일 로드됨: %s\n", configFile)
			break
		}
	}
}

// loadEnvFile 특정 .env 파일에서 환경변수 로드
func (c *Config) loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 빈 줄이나 주석 무시
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// KEY=VALUE 형태로 파싱
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 환경변수로 설정 (기존 환경변수가 없는 경우만)
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}

	return scanner.Err()
}

// loadFromEnvironment 환경변수에서 설정값 로드
func (c *Config) loadFromEnvironment() {
	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		c.DBPath = dbPath
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}
}

// GetDBPath DB 파일 경로 반환 (디렉토리 자동 생성)
func (c *Config) GetDBPath() string {
	dbPath := c.DBPath

	// 디렉토리 경로 추출 및 생성
	if strings.Contains(dbPath, "/") {
		dir := dbPath[:strings.LastIndex(dbPath, "/")]
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("DB 디렉토리 생성 실패: %v\n", err)
		}
	}

	return dbPath
}

// Validate 설정값 유효성 검사
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT 설정이 필요합니다")
	}

	if c.DBPath == "" {
		return fmt.Errorf("DB_PATH 설정이 필요합니다")
	}

	// 로그 레벨 유효성 검사
	validLogLevels := []string{"DEBUG", "INFO", "WARNING", "ERROR"}
	logLevel := strings.ToUpper(c.LogLevel)
	found := false
	for _, level := range validLogLevels {
		if level == logLevel {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("유효하지 않은 LOG_LEVEL: %s (사용 가능: %v)", c.LogLevel, validLogLevels)
	}

	return nil
}

// PrintConfig 현재 설정값 출력 (디버깅용)
func (c *Config) PrintConfig() {
	fmt.Println("=== 애플리케이션 설정 ===")
	fmt.Printf("Port: %s\n", c.Port)
	fmt.Printf("DB Path: %s\n", c.DBPath)
	fmt.Printf("Log Level: %s\n", c.LogLevel)
	fmt.Printf("Max Connections: %d\n", c.MaxConnections)
	fmt.Println("========================")
}
