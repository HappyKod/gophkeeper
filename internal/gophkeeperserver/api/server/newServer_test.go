package server

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewServer(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	addressService := "localhost:5000"
	go NewServer(r, addressService)
	time.Sleep(2 * time.Second)
	// Проверяем, что сервер запущен
	resp, err := http.Get("http://" + addressService)
	if err != nil {
		t.Errorf("Ошибка: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Неверный статус код: %d", resp.StatusCode)
	}
}

func TestNewServer_Shutdown(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
	addressService := "localhost:8080"
	go NewServer(r, addressService)
	// Отправляем сигнал остановки сервера
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	// Проверяем, что сервер выключен
	resp, err := http.Get("http://" + addressService)
	if err == nil {
		t.Errorf("Сервер не выключен")
	}
	if resp != nil && resp.StatusCode != http.StatusNotFound {
		t.Errorf("Неверный статус код: %d", resp.StatusCode)
	}
}
