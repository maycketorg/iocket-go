package iocketsdk

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

//======================================================================
// LOGGER CUSTOMIZADO (Integrado diretamente no pacote)
//======================================================================

const (
	green  = "\033[97;42m"
	yellow = "\033[97;43m"
	red    = "\033[97;41m"
	reset  = "\033[0m"
)

// Logger é uma struct que encapsula os métodos de log.
type Logger struct{}

// NewLogger cria uma instância do seu logger.
func NewLogger() *Logger {
	return &Logger{}
}

// Error loga mensagens de erro (vermelho).
func (l *Logger) Error(args ...any) {
	_, file, line, _ := runtime.Caller(2) // Pula para pegar o local da chamada
	prefix := fmt.Sprintf("%s[IOCKET]%s %s %s:%d |", red, reset, time.Now().Format("15:04:05"), simplifyPath(file), line)
	fmt.Println(append([]any{prefix}, args...)...)
}

// Warn loga mensagens de aviso (amarelo).
func (l *Logger) Warn(args ...any) {
	_, file, line, _ := runtime.Caller(2)
	prefix := fmt.Sprintf("%s[IOCKET]%s %s %s:%d |", yellow, reset, time.Now().Format("15:04:05"), simplifyPath(file), line)
	fmt.Println(append([]any{prefix}, args...)...)
}

// Info loga mensagens informativas (verde).
func (l *Logger) Info(args ...any) {
	_, file, line, _ := runtime.Caller(2)
	prefix := fmt.Sprintf("%s[IOCKET]%s %s %s:%d |", green, reset, time.Now().Format("15:04:05"), simplifyPath(file), line)
	fmt.Println(append([]any{prefix}, args...)...)
}

func simplifyPath(path string) string {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return path
	}
	return path[i+1:]
}
