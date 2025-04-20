package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
)

// Конфигурация OIDC
type OIDCConfig struct {
	ProviderURL  string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	JWTSecret    string
	TokenExpiry  time.Duration
}

// Конфигурация подключения к Redis для High Availability
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Пользователь системы
type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	Provider string   `json:"provider"`
}

// JWT Claims с расширенной информацией
type CustomClaims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	Provider string   `json:"provider"`
	jwt.RegisteredClaims
}

// Структура данных для метрик
type MetricsData struct {
	Timestamp           int64             `json:"timestamp"`
	ActiveUsers         int               `json:"activeUsers"`
	RequestsPerSecond   float64           `json:"requestsPerSecond"`
	ResponseTimeMs      float64           `json:"responseTimeMs"`
	ConversionRate      float64           `json:"conversionRate"`
	Sales               int               `json:"sales"`
	ErrorRate           float64           `json:"errorRate"`
	ErrorsByType        map[string]int    `json:"errorsByType"`
	ServerLoad          float64           `json:"serverLoad"`
	DatabaseConnections int               `json:"databaseConnections"`
	RegionalData        map[string]Region `json:"regionalData"`
	SourcesData         map[string]int    `json:"sourcesData"`
	ConversionFunnel    ConversionFunnel  `json:"conversionFunnel"`
	HistoricalData      HistoricalData    `json:"historicalData"`
}

// Структура для региональных данных
type Region struct {
	ActiveUsers    int     `json:"activeUsers"`
	Sales          int     `json:"sales"`
	ConversionRate float64 `json:"conversionRate"`
}

// Структура для воронки конверсии
type ConversionFunnel struct {
	Visitors       int `json:"visitors"`
	ProductViews   int `json:"productViews"`
	AddedToCart    int `json:"addedToCart"`
	BeganCheckout  int `json:"beganCheckout"`
	PurchasedItems int `json:"purchasedItems"`
}

// Структура для исторических данных
type HistoricalData struct {
	Hourly map[int64]HistoricalMetrics `json:"hourly"` // Почасовые данные (ключ - timestamp)
	Daily  map[int64]HistoricalMetrics `json:"daily"`  // Ежедневные данные (ключ - timestamp)
	Weekly map[int64]HistoricalMetrics `json:"weekly"` // Еженедельные данные (ключ - timestamp)
}

// Структура для исторических метрик
type HistoricalMetrics struct {
	ActiveUsers    int     `json:"activeUsers"`
	Sales          int     `json:"sales"`
	ConversionRate float64 `json:"conversionRate"`
	ResponseTimeMs float64 `json:"responseTimeMs"`
}

// Согласованный генератор данных
type CoherentDataGenerator struct {
	mu               sync.Mutex
	baseActiveUsers  int
	baseSales        int
	dayStartHour     int
	dayEndHour       int
	lastMetrics      MetricsData
	errorTypes       []string
	trafficSources   []string
	regions          []string
	regionWeights    map[string]float64
	sourceWeights    map[string]float64
	trends           map[string]float64
	seasonality      float64
	dayOfWeekFactors map[int]float64
	anomalyChance    float64
	baseDataTime     time.Time // Начальное время для расчета трендов
	currentDataTime  time.Time // Текущее время для генерации данных
	historicalHourly map[int64]HistoricalMetrics
	historicalDaily  map[int64]HistoricalMetrics
	historicalWeekly map[int64]HistoricalMetrics
}

// Менеджер OIDC авторизации
type OIDCManager struct {
	config       OIDCConfig
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	googleConfig oauth2.Config
	yandexConfig oauth2.Config
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
	clients      = make(map[*websocket.Conn]bool)
	clientsMu    sync.Mutex
	generator    *CoherentDataGenerator
	oidcManager  *OIDCManager
	redisClient  *redis.Client
	ctx          context.Context
	cancelFunc   context.CancelFunc
	shutdownChan = make(chan bool)
)

// Инициализация OIDC менеджера
func NewOIDCManager(config OIDCConfig) (*OIDCManager, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		return nil, err
	}

	// Создаем конфигурацию для Google
	googleConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	// Создаем конфигурацию для Яндекса
	yandexConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
		Scopes: config.Scopes,
	}

	return &OIDCManager{
		config:       config,
		provider:     provider,
		verifier:     provider.Verifier(&oidc.Config{ClientID: config.ClientID}),
		googleConfig: googleConfig,
		yandexConfig: yandexConfig,
	}, nil
}

// Авторизация через Google
func (om *OIDCManager) HandleGoogleLogin(c *gin.Context) {
	state := generateRandomState()

	// Сохраняем state в Redis или cookies для проверки при callback
	c.SetCookie("oidc_state", state, 60*15, "/", "", false, true)

	// Перенаправляем на страницу авторизации Google
	c.Redirect(http.StatusFound, om.googleConfig.AuthCodeURL(state))
}

// Авторизация через Яндекс
func (om *OIDCManager) HandleYandexLogin(c *gin.Context) {
	state := generateRandomState()

	// Сохраняем state в Redis или cookies для проверки при callback
	c.SetCookie("oidc_state", state, 60*15, "/", "", false, true)

	// Перенаправляем на страницу авторизации Яндекса
	c.Redirect(http.StatusFound, om.yandexConfig.AuthCodeURL(state))
}

// Обработка callback от провайдера OIDC
func (om *OIDCManager) HandleCallback(c *gin.Context) {
	// Получаем state и code из запроса
	state := c.Query("state")
	code := c.Query("code")

	// Проверяем state
	savedState, _ := c.Cookie("oidc_state")
	if state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Определяем провайдера по параметру запроса
	provider := c.Query("provider")
	var config oauth2.Config
	if provider == "google" {
		config = om.googleConfig
	} else if provider == "yandex" {
		config = om.yandexConfig
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown provider"})
		return
	}

	// Обмениваем code на токен
	oauth2Token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Получаем ID токен
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token in response"})
		return
	}

	// Верифицируем ID токен
	idToken, err := om.verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token"})
		return
	}

	// Извлекаем данные пользователя
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}

	// Создаем пользователя
	user := User{
		ID:       claims["sub"].(string),
		Email:    claims["email"].(string),
		Name:     claims["name"].(string),
		Roles:    []string{"user"}, // По умолчанию роль пользователя
		Provider: provider,
	}

	// Генерируем JWT токен
	token, err := om.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Устанавливаем токен в куки
	c.SetCookie("auth_token", token, int(om.config.TokenExpiry.Seconds()), "/", "", false, true)

	// Перенаправляем на главную страницу
	c.Redirect(http.StatusFound, "/")
}

// Генерация JWT токена
func (om *OIDCManager) GenerateJWT(user User) (string, error) {
	// Создаем claims для JWT
	claims := CustomClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Roles:    user.Roles,
		Provider: user.Provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(om.config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dashboard",
			Subject:   user.ID,
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	return token.SignedString([]byte(om.config.JWTSecret))
}

// Middleware для проверки авторизации
func (om *OIDCManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка или cookie
		authHeader := c.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			// Из заголовка Authorization: Bearer <token>
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		} else {
			// Из куки
			tokenString, _ = c.Cookie("auth_token")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Валидируем JWT
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(om.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Получаем claims из токена
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Сохраняем данные пользователя в контексте
		c.Set("user", User{
			ID:       claims.UserID,
			Email:    claims.Email,
			Name:     claims.Name,
			Roles:    claims.Roles,
			Provider: claims.Provider,
		})

		c.Next()
	}
}

// Middleware для проверки ролей
func (om *OIDCManager) RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		typedUser, ok := user.(User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}

		// Проверяем наличие хотя бы одной из требуемых ролей
		hasRole := false
		for _, role := range typedUser.Roles {
			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

// Инициализация Redis для High Availability
func initRedis(config RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Проверяем соединение
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Публикация метрик в Redis для других инстансов
func publishMetricsToRedis(metrics MetricsData) error {
	if redisClient == nil {
		return nil // Redis не инициализирован, пропускаем
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	// Публикуем данные в Redis канал
	return redisClient.Publish(context.Background(), "metrics_channel", data).Err()
}

// Подписка на метрики из Redis от других инстансов
func subscribeToMetricsFromRedis() {
	if redisClient == nil {
		return // Redis не инициализирован, пропускаем
	}

	pubsub := redisClient.Subscribe(context.Background(), "metrics_channel")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var metrics MetricsData
		if err := json.Unmarshal([]byte(msg.Payload), &metrics); err != nil {
			log.Printf("Error unmarshaling Redis metrics: %v", err)
			continue
		}

		// Отправляем метрики клиентам так же, как и локально сгенерированные
		data := []byte(msg.Payload)

		clientsMu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error sending Redis metrics: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

// Создание нового генератора согласованных данных
func NewCoherentDataGenerator() *CoherentDataGenerator {
	// Фиксируем базовое время для начала генерации данных
	baseTime := time.Now().Add(-24 * 7 * time.Hour) // Неделя назад для исторических данных

	// Регионы России с весами распределения трафика
	regions := []string{
		"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург",
		"Казань", "Нижний Новгород", "Челябинск", "Омск",
		"Самара", "Ростов-на-Дону", "Уфа", "Красноярск",
		"Пермь", "Воронеж", "Волгоград", "Краснодар",
	}

	regionWeights := map[string]float64{
		"Москва": 0.26, "Санкт-Петербург": 0.18, "Новосибирск": 0.06, "Екатеринбург": 0.05,
		"Казань": 0.045, "Нижний Новгород": 0.040, "Челябинск": 0.035, "Омск": 0.03,
		"Самара": 0.035, "Ростов-на-Дону": 0.035, "Уфа": 0.03, "Красноярск": 0.025,
		"Пермь": 0.02, "Воронеж": 0.025, "Волгоград": 0.02, "Краснодар": 0.03,
	}

	// Источники трафика с весами
	trafficSources := []string{
		"Органический поиск", "Прямые заходы", "Социальные сети",
		"Email-рассылки", "Реферальные ссылки", "Контекстная реклама",
		"Медийная реклама", "Партнерские программы",
	}

	sourceWeights := map[string]float64{
		"Органический поиск": 0.35, "Прямые заходы": 0.15, "Социальные сети": 0.2,
		"Email-рассылки": 0.1, "Реферальные ссылки": 0.05, "Контекстная реклама": 0.08,
		"Медийная реклама": 0.04, "Партнерские программы": 0.03,
	}

	// Факторы для дней недели (1 - Monday, 7 - Sunday)
	dayOfWeekFactors := map[int]float64{
		1: 0.85, // Понедельник
		2: 0.9,  // Вторник
		3: 1.0,  // Среда
		4: 1.05, // Четверг
		5: 1.2,  // Пятница
		6: 0.7,  // Суббота
		0: 0.6,  // Воскресенье
	}

	// Создание экземпляра генератора
	generator := &CoherentDataGenerator{
		baseActiveUsers: 1200,
		baseSales:       120,
		dayStartHour:    8,  // 8:00 AM
		dayEndHour:      20, // 8:00 PM
		errorTypes:      []string{"Server Error", "Client Error", "Network Error", "Database Error", "Validation Error"},
		trafficSources:  trafficSources,
		regions:         regions,
		regionWeights:   regionWeights,
		sourceWeights:   sourceWeights,
		trends: map[string]float64{
			"activeUsers": 0.15,  // Рост активных пользователей (% в день)
			"sales":       0.12,  // Рост продаж (% в день)
			"conversion":  0.02,  // Рост конверсии (% в день)
			"errors":      -0.05, // Снижение ошибок (% в день)
		},
		dayOfWeekFactors: dayOfWeekFactors,
		seasonality:      0.3,
		anomalyChance:    0.03,
		baseDataTime:     baseTime,
		currentDataTime:  baseTime,
		historicalHourly: make(map[int64]HistoricalMetrics),
		historicalDaily:  make(map[int64]HistoricalMetrics),
		historicalWeekly: make(map[int64]HistoricalMetrics),
	}

	// Генерируем начальные метрики
	initialMetrics := generator.generateInitialMetrics()
	generator.lastMetrics = initialMetrics

	// Генерируем исторические данные для заполнения графиков
	generator.generateHistoricalData()

	return generator
}

// Генерация начальных метрик
func (dg *CoherentDataGenerator) generateInitialMetrics() MetricsData {
	// Создаем начальные региональные данные
	regionalData := make(map[string]Region)
	sourcesData := make(map[string]int)
	totalUsers := 0

	// Создаем региональные данные на основе весов
	for _, region := range dg.regions {
		weight := dg.regionWeights[region]
		regionalUsers := int(float64(dg.baseActiveUsers) * weight)

		// Конверсия немного различается по регионам
		regionalConversion := 2.5 + (rand.Float64() - 0.5)
		regionalSales := int(float64(regionalUsers) * regionalConversion / 100)

		regionalData[region] = Region{
			ActiveUsers:    regionalUsers,
			Sales:          regionalSales,
			ConversionRate: regionalConversion,
		}

		totalUsers += regionalUsers
	}

	// Создаем данные по источникам трафика
	for _, source := range dg.trafficSources {
		weight := dg.sourceWeights[source]
		sourcesData[source] = int(float64(totalUsers) * weight)
	}

	// Создаем воронку конверсии
	visitors := totalUsers
	productViews := int(float64(visitors) * 0.7)
	addedToCart := int(float64(productViews) * 0.3)
	beganCheckout := int(float64(addedToCart) * 0.5)
	purchased := int(float64(beganCheckout) * 0.8)

	funnel := ConversionFunnel{
		Visitors:       visitors,
		ProductViews:   productViews,
		AddedToCart:    addedToCart,
		BeganCheckout:  beganCheckout,
		PurchasedItems: purchased,
	}

	// Создаем базовые метрики
	metrics := MetricsData{
		Timestamp:           time.Now().Unix(),
		ActiveUsers:         totalUsers,
		RequestsPerSecond:   float64(totalUsers) * 0.05,
		ResponseTimeMs:      200,
		ConversionRate:      2.5,
		Sales:               purchased,
		ErrorRate:           1.0,
		ErrorsByType:        make(map[string]int),
		ServerLoad:          40,
		DatabaseConnections: 50,
		RegionalData:        regionalData,
		SourcesData:         sourcesData,
		ConversionFunnel:    funnel,
		HistoricalData: HistoricalData{
			Hourly: make(map[int64]HistoricalMetrics),
			Daily:  make(map[int64]HistoricalMetrics),
			Weekly: make(map[int64]HistoricalMetrics),
		},
	}

	return metrics
}

// Генерация исторических данных
func (dg *CoherentDataGenerator) generateHistoricalData() {
	// Создаем исторические данные за последнюю неделю
	currentTime := dg.baseDataTime
	now := time.Now()

	// Генерируем почасовые данные за неделю
	for currentTime.Before(now) {
		hourTimestamp := time.Date(
			currentTime.Year(), currentTime.Month(), currentTime.Day(),
			currentTime.Hour(), 0, 0, 0, currentTime.Location(),
		).Unix()

		// Генерируем метрики для этого часа
		tempTime := dg.currentDataTime
		dg.currentDataTime = currentTime
		metrics := dg.GenerateMetrics()
		dg.currentDataTime = tempTime

		// Сохраняем почасовые метрики
		dg.historicalHourly[hourTimestamp] = HistoricalMetrics{
			ActiveUsers:    metrics.ActiveUsers,
			Sales:          metrics.Sales,
			ConversionRate: metrics.ConversionRate,
			ResponseTimeMs: metrics.ResponseTimeMs,
		}

		// Если это начало дня, создаем запись для дневных метрик
		if currentTime.Hour() == 0 {
			dayTimestamp := time.Date(
				currentTime.Year(), currentTime.Month(), currentTime.Day(),
				0, 0, 0, 0, currentTime.Location(),
			).Unix()

			// Вычисляем среднее значение за день (приближение)
			dailyActiveUsers := int(float64(metrics.ActiveUsers) * 16) // Примерно 16 часов активности
			dailySales := int(float64(metrics.Sales) * 16)

			dg.historicalDaily[dayTimestamp] = HistoricalMetrics{
				ActiveUsers:    dailyActiveUsers,
				Sales:          dailySales,
				ConversionRate: metrics.ConversionRate,
				ResponseTimeMs: metrics.ResponseTimeMs,
			}
		}

		// Если это начало недели (понедельник), создаем запись для недельных метрик
		if currentTime.Weekday() == time.Monday && currentTime.Hour() == 0 {
			weekTimestamp := time.Date(
				currentTime.Year(), currentTime.Month(), currentTime.Day(),
				0, 0, 0, 0, currentTime.Location(),
			).Unix()

			// Вычисляем примерное среднее значение за неделю
			weeklyActiveUsers := int(float64(metrics.ActiveUsers) * 16 * 5) // Учитываем рабочие дни
			weeklySales := int(float64(metrics.Sales) * 16 * 5)

			dg.historicalWeekly[weekTimestamp] = HistoricalMetrics{
				ActiveUsers:    weeklyActiveUsers,
				Sales:          weeklySales,
				ConversionRate: metrics.ConversionRate,
				ResponseTimeMs: metrics.ResponseTimeMs,
			}
		}

		// Переходим к следующему часу
		currentTime = currentTime.Add(1 * time.Hour)
	}
}

// GenerateMetrics генерирует новые согласованные метрики
func (dg *CoherentDataGenerator) GenerateMetrics() MetricsData {
	dg.mu.Lock()
	defer dg.mu.Unlock()

	// Используем текущее время для генерации или передвигаем виртуальное время для исторических данных
	var now time.Time
	if dg.currentDataTime.Equal(dg.baseDataTime) {
		now = time.Now()
	} else {
		now = dg.currentDataTime
	}

	// Получаем день недели (0 - Sunday, 1 - Monday, ...)
	dayOfWeek := int(now.Weekday())

	// Коэффициент дня недели
	dayOfWeekFactor := dg.dayOfWeekFactors[dayOfWeek]

	// Коэффициент времени дня
	hour := now.Hour()
	timeOfDayFactor := 0.5
	if hour >= dg.dayStartHour && hour <= dg.dayEndHour {
		// Рабочие часы - больше активности
		timeOfDayFactor = 1.0 - 0.7*math.Abs(float64(hour-((dg.dayStartHour+dg.dayEndHour)/2)))/float64((dg.dayEndHour-dg.dayStartHour)/2)
	}

	// Расчет долгосрочного тренда (дни с базового времени)
	daysSinceBase := now.Sub(dg.baseDataTime).Hours() / 24

	// Применяем трендовые факторы ко всем метрикам
	activeUsersTrendFactor := 1.0 + dg.trends["activeUsers"]*(daysSinceBase/30)
	salesTrendFactor := 1.0 + dg.trends["sales"]*(daysSinceBase/30)
	conversionTrendFactor := 1.0 + dg.trends["conversion"]*(daysSinceBase/30)
	errorTrendFactor := 1.0 + dg.trends["errors"]*(daysSinceBase/30)

	// Случайные флуктуации (меньше для более стабильных метрик)
	userRandomFactor := 0.97 + 0.06*rand.Float64()
	salesRandomFactor := 0.95 + 0.1*rand.Float64()
	conversionRandomFactor := 0.98 + 0.04*rand.Float64()

	// Аномальное поведение с малой вероятностью
	anomalyFactor := 1.0
	if rand.Float64() < dg.anomalyChance {
		// Резкий скачок или падение
		anomalyFactor = 0.7 + 0.6*rand.Float64()
		if rand.Float64() < 0.5 {
			anomalyFactor = 1 / anomalyFactor // Иногда делаем падение вместо скачка
		}
	}

	// Расчет общего количества активных пользователей
	baseActive := dg.lastMetrics.ActiveUsers
	if baseActive == 0 {
		baseActive = dg.baseActiveUsers
	}

	activeUsers := int(float64(baseActive) *
		timeOfDayFactor *
		dayOfWeekFactor *
		activeUsersTrendFactor *
		userRandomFactor *
		anomalyFactor)

	if activeUsers < 10 {
		activeUsers = 10 // Минимум
	}

	// Расчет общей конверсии с учетом факторов
	baseConversion := dg.lastMetrics.ConversionRate
	if baseConversion == 0 {
		baseConversion = 2.5
	}

	conversionRate := baseConversion *
		conversionTrendFactor *
		conversionRandomFactor

	// Ограничения на конверсию
	if conversionRate < 0.5 {
		conversionRate = 0.5
	} else if conversionRate > 7 {
		conversionRate = 7
	}

	// Расчет общих продаж на основе пользователей и конверсии
	sales := int(float64(activeUsers) *
		conversionRate / 100 *
		salesTrendFactor *
		salesRandomFactor)

	// Расчет RPS на основе активных пользователей
	requestsPerSecond := float64(activeUsers) * (0.03 + 0.02*rand.Float64())

	// Время отклика зависит от RPS
	baseResponseTime := dg.lastMetrics.ResponseTimeMs
	if baseResponseTime == 0 {
		baseResponseTime = 200
	}

	responseTimeFactor := 1.0 + 0.3*math.Log10(requestsPerSecond/50+0.1)
	responseTimeMs := baseResponseTime * responseTimeFactor * (0.95 + 0.1*rand.Float64())

	// Ограничения на время отклика
	if responseTimeMs < 100 {
		responseTimeMs = 100
	} else if responseTimeMs > 5000 {
		responseTimeMs = 5000
	}

	// Если время отклика плохое, конверсия падает
	if responseTimeMs > 1000 {
		conversionRate *= 0.9
		sales = int(float64(sales) * 0.9)
	}

	// Расчет ошибок
	errorFactor := 0.01 * (responseTimeMs / 200) * errorTrendFactor
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

	dbConnections := int(20 + float64(activeUsers)/40 + rand.Float64()*30)

	// Создаем региональные данные, распределяя пользователей и продажи по регионам
	regionalData := make(map[string]Region)
	totalRegionalUsers := 0
	totalRegionalSales := 0

	for _, region := range dg.regions {
		weight := dg.regionWeights[region]

		// Добавляем небольшую случайность в региональные веса
		adjustedWeight := weight * (0.9 + 0.2*rand.Float64())

		regionalUsers := int(float64(activeUsers) * adjustedWeight)

		// Конверсия может немного отличаться по регионам
		regionalConversion := conversionRate * (0.9 + 0.2*rand.Float64())
		regionalSales := int(float64(regionalUsers) * regionalConversion / 100)

		regionalData[region] = Region{
			ActiveUsers:    regionalUsers,
			Sales:          regionalSales,
			ConversionRate: regionalConversion,
		}

		totalRegionalUsers += regionalUsers
		totalRegionalSales += regionalSales
	}

	// Нормализуем региональные данные, чтобы суммы совпадали с общими значениями
	if totalRegionalUsers > 0 && totalRegionalSales > 0 {
		userFactor := float64(activeUsers) / float64(totalRegionalUsers)
		salesFactor := float64(sales) / float64(totalRegionalSales)

		for region, data := range regionalData {
			adjustedUsers := int(float64(data.ActiveUsers) * userFactor)
			adjustedSales := int(float64(data.Sales) * salesFactor)

			regionalData[region] = Region{
				ActiveUsers:    adjustedUsers,
				Sales:          adjustedSales,
				ConversionRate: data.ConversionRate,
			}
		}
	}

	// Создаем данные по источникам трафика
	sourcesData := make(map[string]int)
	totalSourceUsers := 0

	for _, source := range dg.trafficSources {
		weight := dg.sourceWeights[source]

		// Добавляем небольшую случайность в веса источников
		adjustedWeight := weight * (0.85 + 0.3*rand.Float64())

		sourceUsers := int(float64(activeUsers) * adjustedWeight)
		sourcesData[source] = sourceUsers

		totalSourceUsers += sourceUsers
	}

	// Нормализуем данные источников
	if totalSourceUsers > 0 {
		sourceFactor := float64(activeUsers) / float64(totalSourceUsers)

		for source, users := range sourcesData {
			sourcesData[source] = int(float64(users) * sourceFactor)
		}
	}

	// Создаем воронку конверсии
	visitors := activeUsers
	productViews := int(float64(visitors) * (0.65 + 0.1*rand.Float64()))
	addedToCart := int(float64(productViews) * (0.25 + 0.1*rand.Float64()))
	beganCheckout := int(float64(addedToCart) * (0.45 + 0.1*rand.Float64()))
	purchased := sales // Используем рассчитанные продажи для согласованности

	funnel := ConversionFunnel{
		Visitors:       visitors,
		ProductViews:   productViews,
		AddedToCart:    addedToCart,
		BeganCheckout:  beganCheckout,
		PurchasedItems: purchased,
	}

	// Генерируем или обновляем исторические данные
	hourTimestamp := time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(), 0, 0, 0, now.Location(),
	).Unix()

	dg.historicalHourly[hourTimestamp] = HistoricalMetrics{
		ActiveUsers:    activeUsers,
		Sales:          sales,
		ConversionRate: conversionRate,
		ResponseTimeMs: responseTimeMs,
	}

	// Создаем метрики с согласованными данными
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
		RegionalData:        regionalData,
		SourcesData:         sourcesData,
		ConversionFunnel:    funnel,
		HistoricalData: HistoricalData{
			Hourly: dg.historicalHourly,
			Daily:  dg.historicalDaily,
			Weekly: dg.historicalWeekly,
		},
	}

	// Публикуем метрики в Redis для других инстансов
	if err := publishMetricsToRedis(metrics); err != nil {
		log.Printf("Error publishing metrics to Redis: %v", err)
	}

	dg.lastMetrics = metrics
	return metrics
}

// Получение исторических данных для графиков
func (dg *CoherentDataGenerator) GetHistoricalData(period string, metric string) []map[string]interface{} {
	dg.mu.Lock()
	defer dg.mu.Unlock()

	var source map[int64]HistoricalMetrics

	switch period {
	case "hourly":
		source = dg.historicalHourly
	case "daily":
		source = dg.historicalDaily
	case "weekly":
		source = dg.historicalWeekly
	default:
		source = dg.historicalHourly
	}

	result := make([]map[string]interface{}, 0, len(source))

	// Сортируем ключи для хронологического порядка
	var timestamps []int64
	for ts := range source {
		timestamps = append(timestamps, ts)
	}
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	for _, ts := range timestamps {
		data := source[ts]
		item := map[string]interface{}{
			"timestamp": ts,
		}

		switch metric {
		case "activeUsers":
			item["value"] = data.ActiveUsers
		case "sales":
			item["value"] = data.Sales
		case "conversionRate":
			item["value"] = data.ConversionRate
		case "responseTime":
			item["value"] = data.ResponseTimeMs
		default:
			item["value"] = data.ActiveUsers
		}

		result = append(result, item)
	}

	return result
}

// Функция для получения метрик в реальном времени
func (dg *CoherentDataGenerator) GetCurrentMetrics() MetricsData {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	return dg.lastMetrics
}

// Генерация случайного state для OIDC
func generateRandomState() string {
	return gofakeit.UUID()
}

// Обработчик WebSocket-соединений с проверкой JWT
func handleConnections(c *gin.Context) {
	// Проверяем наличие JWT в запросе
	tokenString := c.Query("token")
	if tokenString == "" {
		// Пробуем получить из заголовка
		authHeader := c.GetHeader("Sec-WebSocket-Protocol")
		if authHeader != "" {
			tokenString = authHeader
		}
	}

	// Проверяем токен, если требуется аутентификация
	if oidcManager != nil && tokenString != "" {
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(oidcManager.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	}

	// Апгрейд соединения до WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	// Читаем сообщения для обработки отключений
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			break
		}

		// Проверяем состояние сервера
		select {
		case <-ctx.Done():
			// Сервер завершает работу, закрываем соединение
			ws.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"))
			return
		default:
			// Продолжаем работу
		}
	}
}

// Отправка метрик всем подключенным клиентам
func broadcastMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Генерируем новые метрики
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

		case <-ctx.Done():
			// Сигнал завершения работы
			log.Println("Stopping broadcast routine...")
			return
		}
	}
}

// Graceful shutdown
func setupGracefulShutdown(server *http.Server) {
	// Канал для получения сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Сначала отменяем контекст, чтобы остановить фоновые горутины
		cancelFunc()

		// Устанавливаем таймаут на завершение
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		// Останавливаем HTTP сервер
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		// Закрываем Redis, если используется
		if redisClient != nil {
			if err := redisClient.Close(); err != nil {
				log.Printf("Redis connection close error: %v", err)
			}
		}

		// Сигнализируем о завершении работы
		close(shutdownChan)
	}()
}

func main() {
	// Создаем контекст с возможностью отмены
	ctx, cancelFunc = context.WithCancel(context.Background())

	// Инициализация генератора данных
	gofakeit.Seed(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	generator = NewCoherentDataGenerator()

	// Попытка инициализации Redis для High Availability
	redisConfig := RedisConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	}

	var redisErr error
	redisClient, redisErr = initRedis(redisConfig)
	if redisErr != nil {
		log.Printf("Warning: Redis initialization failed: %v. Running in standalone mode.", redisErr)
	} else {
		log.Println("Redis connected. Running in high availability mode.")
		// Запускаем подписку на метрики от других инстансов
		go subscribeToMetricsFromRedis()
	}

	// Настройка OIDC авторизации, если указаны переменные окружения
	if getEnv("ENABLE_OIDC", "") == "true" {
		oidcConfig := OIDCConfig{
			ProviderURL:  getEnv("OIDC_PROVIDER_URL", "https://accounts.google.com"),
			ClientID:     getEnv("OIDC_CLIENT_ID", ""),
			ClientSecret: getEnv("OIDC_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback"),
			Scopes:       []string{"openid", "profile", "email"},
			JWTSecret:    getEnv("JWT_SECRET", "secret"),
			TokenExpiry:  24 * time.Hour,
		}

		var oidcErr error
		oidcManager, oidcErr = NewOIDCManager(oidcConfig)
		if oidcErr != nil {
			log.Printf("Warning: OIDC initialization failed: %v. Authentication disabled.", oidcErr)
		} else {
			log.Println("OIDC authentication enabled.")
		}
	}

	// Настройка Gin
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Общедоступные маршруты
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Маршруты для аутентификации
	if oidcManager != nil {
		auth := r.Group("/auth")
		{
			auth.GET("/google", oidcManager.HandleGoogleLogin)
			auth.GET("/yandex", oidcManager.HandleYandexLogin)
			auth.GET("/callback", oidcManager.HandleCallback)
		}

		// Защищенные маршруты
		protected := r.Group("/api")
		protected.Use(oidcManager.AuthMiddleware())
		{
			protected.GET("/metrics/current", func(c *gin.Context) {
				c.JSON(http.StatusOK, generator.GetCurrentMetrics())
			})

			// Добавляем маршрут для получения исторических данных
			protected.GET("/metrics/historical/:period/:metric", func(c *gin.Context) {
				period := c.Param("period") // hourly, daily, weekly
				metric := c.Param("metric") // activeUsers, sales, conversionRate, responseTime

				if period == "" || metric == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
					return
				}

				data := generator.GetHistoricalData(period, metric)
				c.JSON(http.StatusOK, data)
			})

			// Маршрут только для админов
			admin := protected.Group("/admin")
			admin.Use(oidcManager.RoleMiddleware("admin"))
			{
				admin.GET("/status", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"clients": len(clients),
						"uptime":  time.Since(time.Unix(0, 0)),
					})
				})
			}
		}

		// WebSocket с проверкой токена
		r.GET("/ws", handleConnections)
	} else {
		// Если OIDC не настроен, все маршруты без аутентификации
		r.GET("/ws", handleConnections)
		r.GET("/metrics/current", func(c *gin.Context) {
			c.JSON(http.StatusOK, generator.GetCurrentMetrics())
		})

		// Добавляем маршрут для получения исторических данных
		r.GET("/metrics/historical/:period/:metric", func(c *gin.Context) {
			period := c.Param("period") // hourly, daily, weekly
			metric := c.Param("metric") // activeUsers, sales, conversionRate, responseTime

			if period == "" || metric == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
				return
			}

			data := generator.GetHistoricalData(period, metric)
			c.JSON(http.StatusOK, data)
		})
	}

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    getEnv("SERVER_ADDR", ":8080"),
		Handler: r,
	}

	// Настраиваем graceful shutdown
	setupGracefulShutdown(server)

	// Запуск широковещательной рассылки метрик
	go broadcastMetrics()

	// Запуск сервера
	log.Printf("Server starting on %s...", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}

	// Ожидаем завершения graceful shutdown
	<-shutdownChan
	log.Println("Server gracefully stopped")
}

// Хелпер для получения переменных окружения
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
