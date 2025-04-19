package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// MetricsData структура для метрик
type MetricsData struct {
	Timestamp           int64   `json:"timestamp"`
	ActiveUsers         int     `json:"activeUsers"`
	RequestsPerSecond   float64 `json:"requestsPerSecond"`
	ResponseTimeMs      float64 `json:"responseTimeMs"`
	ConversionRate      float64 `json:"conversionRate"`
	Sales               int     `json:"sales"`
	ErrorRate           float64 `json:"errorRate"`
	ErrorsByType        map[string]int `json:"errorsByType"`
	ServerLoad          float64 `json:"serverLoad"`
	DatabaseConnections int     `json:"databaseConnections"`
}

// DataGenerator генератор данных
type DataGenerator struct {
	mu               sync.Mutex
	baseActiveUsers  int
	baseSales        int
	dayStartHour     int
	dayEndHour       int
	lastMetrics      MetricsData
	errorTypes       []string
	trends           map[string]float64
	seasonality      float64
	anomalyChance    float64
}

// Глобальные переменные
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Разрешаем любой источник для тестирования
		},
	}
	clients    = make(map[*websocket.Conn]bool)
	clientsMu  sync.Mutex
	generator  *DataGenerator
)

func NewDataGenerator() *DataGenerator {
	return &DataGenerator{
		baseActiveUsers: 1000,
		baseSales:       100,
		dayStartHour:    8,  // 8:00 AM
		dayEndHour:      20, // 8:00 PM
		errorTypes:      []string{"Server Error", "Client Error", "Network Error", "Database Error", "Validation Error"},
		trends: map[string]float64{
			"activeUsers": 0.2,       // Слабый рост
			"sales":       0.1,        // Слабый рост
			"errors":      -0.05,      // Слабое снижение
		},
		seasonality:   0.3,
		anomalyChance: 0.03,
		lastMetrics: MetricsData{
			Timestamp:           time.Now().Unix(),
			ActiveUsers:         1000,
			RequestsPerSecond:   50,
			ResponseTimeMs:      200,
			ConversionRate:      2.5,
			Sales:               100,
			ErrorRate:           1.0,
			ErrorsByType:        make(map[string]int),
			ServerLoad:          40,
			DatabaseConnections: 50,
		},
	}
}

// GenerateMetrics генерирует новые метрики
func (dg *DataGenerator) GenerateMetrics() MetricsData {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	
	now := time.Now()
	hour := now.Hour()
	
	// Коэффициент сезонности на основе времени дня
	timeOfDayFactor := 0.5
	if hour >= dg.dayStartHour && hour <= dg.dayEndHour {
		// Рабочие часы - больше активности
		timeOfDayFactor = 1.0 - 0.7*math.Abs(float64(hour-((dg.dayStartHour+dg.dayEndHour)/2)))/float64((dg.dayEndHour-dg.dayStartHour)/2)
	}
	
	// Применение долгосрочного тренда
	trendFactor := 1.0 + dg.trends["activeUsers"]*float64(now.Unix()-dg.lastMetrics.Timestamp)/86400
	
	// Случайные флуктуации
	randomFactor := 0.95 + 0.1*rand.Float64()
	
	// Аномальное поведение с малой вероятностью
	anomalyFactor := 1.0
	if rand.Float64() < dg.anomalyChance {
		// Резкий скачок или падение
		anomalyFactor = 0.5 + rand.Float64()
		if rand.Float64() < 0.5 {
			anomalyFactor = 1 / anomalyFactor // Иногда делаем падение вместо скачка
		}
	}
	 
	// Расчет активных пользователей
	activeUsers := int(float64(dg.lastMetrics.ActiveUsers) * timeOfDayFactor * trendFactor * randomFactor * anomalyFactor)
	if activeUsers < 10 {
		activeUsers = 10 // Минимум
	}
	
	// Расчет RPS на основе активных пользователей
	requestsPerSecond := float64(activeUsers) * (0.03 + 0.02*rand.Float64())
	
	// Время отклика зависит от RPS
	responseTimeFactor := 1.0 + 0.5*math.Log10(requestsPerSecond/50+0.1)
	responseTimeMs := dg.lastMetrics.ResponseTimeMs * responseTimeFactor * (0.95 + 0.1*rand.Float64())
	
	// Ограничения на время отклика
	if responseTimeMs < 100 {
		responseTimeMs = 100
	} else if responseTimeMs > 5000 {
		responseTimeMs = 5000
	}
	
	// Продажи на основе пользователей и конверсии
	conversionRate := dg.lastMetrics.ConversionRate * (0.98 + 0.04*rand.Float64())
	if responseTimeMs > 1000 {
		// При плохом времени отклика конверсия падает
		conversionRate *= 0.8
	}
	
	// Ограничения на конверсию
	if conversionRate < 0.5 {
		conversionRate = 0.5
	} else if conversionRate > 5 {
		conversionRate = 5
	}
	
	// Расчет продаж
	sales := int(float64(activeUsers) * conversionRate / 100)
	
	// Ошибки - обратно пропорциональны времени отклика
	errorFactor := 0.01 * (responseTimeMs / 200)
	if errorFactor < 0.001 {
		errorFactor = 0.001
	}
	
	errorRate := errorFactor * 100 // в процентах
	
	// Генерация ошибок по типам
	errorsByType := make(map[string]int)
	totalErrors := int(errorRate * requestsPerSecond / 100)
	remainingErrors := totalErrors
	
	for i, errType := range dg.errorTypes {
		var errCount int
		if i == len(dg.errorTypes)-1 {
			errCount = remainingErrors
		} else {
			errCount = int(float64(totalErrors) * (0.1 + 0.3*rand.Float64()))
			if errCount > remainingErrors {
				errCount = remainingErrors
			}
		}
		errorsByType[errType] = errCount
		remainingErrors -= errCount
		if remainingErrors <= 0 {
			break
		}
	}
	
	// Нагрузка сервера и подключения к БД
	serverLoad := 30 + 50*(requestsPerSecond/200) + 10*rand.Float64()
	if serverLoad > 100 {
		serverLoad = 100
	}
	
	dbConnections := int(20 + float64(activeUsers)/50 + rand.Float64()*30)
	
	metrics := MetricsData{
		Timestamp:           now.Unix(),
		ActiveUsers:         activeUsers,
		RequestsPerSecond:   requestsPerSecond,
		ResponseTimeMs:      responseTimeMs,
		ConversionRate:      conversionRate,
		Sales:               sales,
		ErrorRate:           errorRate,
		ErrorsByType:        errorsByType,
		ServerLoad:          serverLoad,
		DatabaseConnections: dbConnections,
	}
	
	dg.lastMetrics = metrics
	return metrics
}

// Обработчик WebSocket-соединений
func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer ws.Close()
	
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()
	
	for {
		// Читаем сообщения от клиента, чтобы обрабатывать отключения
		_, _, err := ws.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			break
		}
	}
}

// Отправка метрик всем подключенным клиентам
func broadcastMetrics() {
	for {
		metrics := generator.GenerateMetrics()
		
		// Сериализуем метрики в JSON
		data, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Error marshaling metrics: %v", err)
			continue
		}
		
		// Отправляем всем подключенным клиентам
		clientsMu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error sending metrics: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
		
		// Небольшая задержка между обновлениями
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Инициализация генератора данных
	gofakeit.Seed(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	generator = NewDataGenerator()

	// Настройка Gin
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Маршруты
	r.GET("/ws", handleConnections)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics/current", func(c *gin.Context) {
		c.JSON(http.StatusOK, generator.lastMetrics)
	})

	// Запуск широковещательной рассылки метрик
	go broadcastMetrics()

	// Запуск сервера
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 