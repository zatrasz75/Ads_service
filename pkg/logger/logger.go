package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type LoggersInterface interface {
	Error(message string, err error)
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Fatal(message string, err error)
	Debug(message string, args ...interface{})
}
type MyLogger struct {
	logger        *log.Logger
	logFile       *os.File
	fileSizeMutex sync.Mutex
}

func NewLogger() LoggersInterface {
	filePath := "app.log"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия файла журнала:", err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	myLogger := &MyLogger{
		logger:  logger,
		logFile: file,
	}

	// Запускаем горутину для проверки размера файла и его обработки
	go myLogger.checkAndRotateLogFile()

	return myLogger
}

const maxLogFileSize = 5 * 1024 * 1024 // 10 MB

func logWithCallerInfo(file string, line int, level string, message string, l *log.Logger, args ...interface{}) {
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	messageWithCaller := fmt.Sprintf("[%s] %s %s %s", level, getFormattedTime(), caller, fmt.Sprintf(message, args...))
	l.Println(messageWithCaller)
	fmt.Println(messageWithCaller)
}

// Error записывает сообщение об ошибке в лог вместе с контекстом вызова функции.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Error(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "ERROR", "%s: %v", l.logger, message, err)
	} else {
		fmt.Println("недоступно.")
	}
}

// Info записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Info(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "INFO", message, l.logger, args...)
	} else {
		fmt.Println("No logger available.")
	}
}

func (l *MyLogger) Warn(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "WARN", message, l.logger, args...)
	} else {
		fmt.Println("No logger available.")
	}
}

// Fatal записывает фатальное сообщение в лог вместе с контекстом вызова функции
// и завершает приложение с кодом ошибки 1.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Fatal(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	logWithCallerInfo(file, line, "FATAL", "%s: %v", l.logger, message, err)
	os.Exit(1) // Завершаем приложение с кодом ошибки
}

// Debug записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Debug(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logWithCallerInfo(file, line, "DEBUG", message, l.logger, args...)
}

// getFormattedTime возвращает текущее время в заданном формате.
func getFormattedTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// checkAndRotateLogFile проверяет размер файла лога и при необходимости выполняет его ротацию.
func (l *MyLogger) checkAndRotateLogFile() {
	for {
		time.Sleep(time.Minute) // Пауза для проверки раз в минуту (настраивайте по желанию)
		l.fileSizeMutex.Lock()
		stat, err := l.logFile.Stat()
		if err != nil {
			l.logger.Println("Ошибка при получении информации о файле логов:", err)
			l.fileSizeMutex.Unlock()
			time.Sleep(time.Minute) // Делаем паузу если ошибка
			continue
		}

		if stat.Size() > maxLogFileSize {
			_ = l.logFile.Close()
			err = os.Rename("app.log", "app.old.log")
			if err != nil {
				l.logger.Println("Ошибка при переименовании файла логов:", err)
			}
			newFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				l.logger.Println("Ошибка при открытии нового файла журнала:", err)
			} else {
				l.logFile = newFile
				l.logger.SetOutput(l.logFile)
			}
		}
		l.fileSizeMutex.Unlock()
	}
}
