import { ref, reactive, onUnmounted } from 'vue';
import dayjs from 'dayjs';
import { PeriodType } from './timePeriodService';

// Интерфейс для региональных данных
export interface RegionData {
  activeUsers: number;
  sales: number;
  conversionRate: number;
}

// Интерфейс для воронки конверсии
export interface ConversionFunnel {
  visitors: number;
  productViews: number;
  addedToCart: number;
  beganCheckout: number;
  purchasedItems: number;
}

// Интерфейс для исторических метрик
export interface HistoricalMetrics {
  activeUsers: number;
  sales: number;
  conversionRate: number;
  responseTimeMs: number;
}

// Интерфейс для исторических данных
export interface HistoricalData {
  hourly: Record<string, HistoricalMetrics>;
  daily: Record<string, HistoricalMetrics>;
  weekly: Record<string, HistoricalMetrics>;
}

// Интерфейс для метрик
export interface Metrics {
  timestamp: number;
  activeUsers: number;
  requestsPerSecond: number;
  responseTimeMs: number;
  conversionRate: number;
  sales: number;
  errorRate: number;
  errorsByType: Record<string, number>;
  serverLoad: number;
  databaseConnections: number;
  regionalData?: Record<string, RegionData>;
  sourcesData?: Record<string, number>;
  conversionFunnel?: ConversionFunnel;
  historicalData?: HistoricalData;
}

// Начальные данные для метрик
const defaultMetrics: Metrics = {
  timestamp: Date.now(),
  activeUsers: 842,
  requestsPerSecond: 0,
  responseTimeMs: 0,
  conversionRate: 0,
  sales: 0,
  errorRate: 0,
  errorsByType: {},
  serverLoad: 0,
  databaseConnections: 0
};

// Буфер для сглаживания данных
const smoothingFactor = 0.3; // Фактор сглаживания (0-1): 0 - нет сглаживания, 1 - максимальное сглаживание
let previousMetrics: Metrics | null = null;

// Время последнего обновления для каждого периода
const lastUpdateTime: Record<PeriodType, number> = {
  hour: 0,
  day: 0,
  week: 0,
  realtime: 0
};

// Частота обновления для каждого периода (в миллисекундах)
const updateFrequency: Record<PeriodType, number> = {
  realtime: 1000,     // Реальное время - каждую секунду
  hour: 5000,         // Час - каждые 5 секунд
  day: 15000,         // День - каждые 15 секунд
  week: 60000         // Неделя - каждую минуту
};

// Конфигурация истории для каждого периода
const historyConfigs: Record<PeriodType, { maxItems: number, interval: number }> = {
  realtime: { maxItems: 20, interval: 1000 },
  hour: { maxItems: 60, interval: 60000 },  // 60 точек по минуте
  day: { maxItems: 24, interval: 3600000 },   // 24 точки по часу
  week: { maxItems: 7, interval: 86400000 }    // 7 точек по дню
};

// История метрик для каждого периода
const metricsHistory: Record<PeriodType, Metrics[]> = {
  realtime: [],
  hour: [],
  day: [],
  week: []
};

// История метрик для графиков - сохраняем истории для разных периодов
export const historyMetrics = reactive<Metrics[]>([]);

// Текущий активный период
let currentPeriod: PeriodType = 'realtime';

// Функция сглаживания данных
const smoothData = (newData: Metrics): Metrics => {
  if (!previousMetrics) {
    previousMetrics = {...newData};
    return newData;
  }
  
  const smoothed = {...newData};
  
  // Обрабатываем каждый числовой параметр
  if (typeof smoothed.activeUsers === 'number' && typeof previousMetrics.activeUsers === 'number') {
    smoothed.activeUsers = previousMetrics.activeUsers * smoothingFactor + newData.activeUsers * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.requestsPerSecond === 'number' && typeof previousMetrics.requestsPerSecond === 'number') {
    smoothed.requestsPerSecond = previousMetrics.requestsPerSecond * smoothingFactor + newData.requestsPerSecond * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.responseTimeMs === 'number' && typeof previousMetrics.responseTimeMs === 'number') {
    smoothed.responseTimeMs = previousMetrics.responseTimeMs * smoothingFactor + newData.responseTimeMs * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.conversionRate === 'number' && typeof previousMetrics.conversionRate === 'number') {
    smoothed.conversionRate = previousMetrics.conversionRate * smoothingFactor + newData.conversionRate * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.sales === 'number' && typeof previousMetrics.sales === 'number') {
    smoothed.sales = previousMetrics.sales * smoothingFactor + newData.sales * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.errorRate === 'number' && typeof previousMetrics.errorRate === 'number') {
    smoothed.errorRate = previousMetrics.errorRate * smoothingFactor + newData.errorRate * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.serverLoad === 'number' && typeof previousMetrics.serverLoad === 'number') {
    smoothed.serverLoad = previousMetrics.serverLoad * smoothingFactor + newData.serverLoad * (1 - smoothingFactor);
  }
  
  if (typeof smoothed.databaseConnections === 'number' && typeof previousMetrics.databaseConnections === 'number') {
    smoothed.databaseConnections = previousMetrics.databaseConnections * smoothingFactor + newData.databaseConnections * (1 - smoothingFactor);
  }
  
  // Для errorsByType нужна отдельная обработка, так как это объект
  smoothed.errorsByType = {...newData.errorsByType};
  
  // Сохраняем текущие сглаженные данные для следующего обновления
  previousMetrics = {...smoothed};
  
  return smoothed;
};

// Фактические метрики, которые будут обновляться
export const metrics = reactive<Metrics>({...defaultMetrics});

// WebSocket подключение
let ws: WebSocket | null = null;
let reconnectTimeout: number | null = null;

// Состояние подключения
export const connectionState = reactive({
  isConnected: false,
  lastError: '',
  reconnecting: false
});

// Функция для форматирования чисел
export const formatNumber = (value: number, decimals = 0) => {
  return Number(value).toLocaleString('ru-RU', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals
  });
};

// Колбэк, вызываемый при обновлении метрик
let onUpdateCallback: (() => void) | null = null;

// Установка колбэка на обновление
export const setOnUpdateCallback = (callback: () => void) => {
  onUpdateCallback = callback;
};

// Установка активного периода для получения соответствующих данных
export const setActivePeriod = (period: PeriodType) => {
  currentPeriod = period;
  
  // Обновляем текущую историю на основе выбранного периода
  // Важно: для реактивности Vue нужно заменить содержимое массива, а не сам массив
  historyMetrics.splice(0, historyMetrics.length, ...metricsHistory[period]);
  
  // Если данных для этого периода еще нет, создадим тестовые данные
  if (historyMetrics.length === 0) {
    // Создаем тестовые данные для демонстрации
    const now = Date.now();
    const config = historyConfigs[period];
    
    // Генерируем тестовые данные для текущего периода
    for (let i = 0; i < config.maxItems; i++) {
      const timestamp = now - (config.maxItems - i - 1) * config.interval;
      
      // Создаем данные с некоторой случайностью
      const randomFactor = 0.5 + Math.random();
      
      // Добавляем в историю периода
      metricsHistory[period].push({
        timestamp,
        activeUsers: Math.round(800 * randomFactor),
        requestsPerSecond: (0.2 + Math.random() * 0.5) * (period === 'realtime' ? 1 : (period === 'hour' ? 5 : (period === 'day' ? 20 : 50))),
        responseTimeMs: 50 + Math.random() * 150,
        conversionRate: 1.5 + Math.random() * 2,
        sales: Math.round(100 * randomFactor) * (period === 'realtime' ? 1 : (period === 'hour' ? 10 : (period === 'day' ? 100 : 1000))),
        errorRate: 0.5 + Math.random() * 1.5,
        errorsByType: {
          'Сетевые': Math.round(10 * randomFactor),
          'База данных': Math.round(5 * randomFactor),
          'Авторизация': Math.round(3 * randomFactor),
          'API': Math.round(8 * randomFactor),
          'Прочие': Math.round(2 * randomFactor),
        },
        serverLoad: 20 + Math.random() * 60,
        databaseConnections: Math.round(10 + Math.random() * 40)
      });
    }
    
    // Обновляем массив historyMetrics
    historyMetrics.splice(0, historyMetrics.length, ...metricsHistory[period]);
  }
  
  // Вызываем колбэк для обновления графиков
  if (onUpdateCallback) {
    onUpdateCallback();
  }
  
  console.log(`Период изменен на: ${period}, количество точек: ${historyMetrics.length}`);
};

// Функция для получения исторических данных от API
export const fetchHistoricalData = async (period: string, metric: string): Promise<any[]> => {
  try {
    const response = await fetch(`http://localhost:8080/metrics/historical/${period}/${metric}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Ошибка при получении исторических данных:', error);
    return [];
  }
};

// Получение текущих метрик по REST API (для аутентифицированных пользователей)
export const fetchCurrentMetrics = async (): Promise<Metrics | null> => {
  try {
    const response = await fetch('http://localhost:8080/metrics/current');
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Ошибка при получении текущих метрик:', error);
    return null;
  }
};

// Соединение с WebSocket с учётом новых данных
export const connectWebSocket = () => {
  if (ws && (ws.readyState === WebSocket.CONNECTING || ws.readyState === WebSocket.OPEN)) {
    return;
  }
  
  connectionState.reconnecting = true;
  
  try {
    ws = new WebSocket('ws://localhost:8080/ws');
    
    ws.onopen = () => {
      connectionState.isConnected = true;
      connectionState.lastError = '';
      connectionState.reconnecting = false;
      console.log('WebSocket соединение установлено');
    };
    
    ws.onmessage = (event) => {
      const newMetrics: Metrics = JSON.parse(event.data);
      
      // Применяем сглаживание данных только для числовых полей верхнего уровня
      const smoothedMetrics = smoothData(newMetrics);
      
      // Обновление текущих метрик (включая новые поля)
      Object.assign(metrics, smoothedMetrics);
      
      // Добавление в историю реального времени
      metricsHistory.realtime.push({...smoothedMetrics});
      if (metricsHistory.realtime.length > historyConfigs.realtime.maxItems) {
        metricsHistory.realtime.shift();
      }
      
      // Обрабатываем исторические данные для других периодов
      processHistoricalData(smoothedMetrics);
      
      // Обновляем текущую историю на основе выбранного периода
      // Используем splice для правильного обновления реактивного массива
      historyMetrics.splice(0, historyMetrics.length, ...metricsHistory[currentPeriod]);
      
      // Вызов колбэка
      if (onUpdateCallback) {
        onUpdateCallback();
      }
    };
    
    ws.onclose = (event) => {
      connectionState.isConnected = false;
      connectionState.reconnecting = true;
      
      if (event.code !== 1000) {
        connectionState.lastError = `Соединение закрыто с кодом ${event.code}`;
        console.log(`WebSocket соединение закрыто с кодом ${event.code}, причина: ${event.reason}`);
        
        // Попытка переподключения через некоторое время
        if (reconnectTimeout) {
          clearTimeout(reconnectTimeout);
        }
        
        // Экспоненциальная задержка: от 1 до 30 секунд
        const delay = Math.min(30000, (Math.random() + 1) * 1000 * Math.pow(1.5, Math.min(10, 1)));
        console.log(`Попытка переподключения через ${delay}ms...`);
        
        reconnectTimeout = window.setTimeout(() => {
          connectWebSocket();
        }, delay);
      }
    };
    
    ws.onerror = (error) => {
      connectionState.isConnected = false;
      connectionState.lastError = 'Ошибка WebSocket';
      console.error('WebSocket ошибка:', error);
    };
  } catch (error) {
    connectionState.lastError = 'Ошибка при создании WebSocket';
    console.error('Ошибка при создании WebSocket:', error);
    
    // Также пробуем переподключиться при ошибке
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
    }
    
    reconnectTimeout = window.setTimeout(() => {
      connectWebSocket();
    }, 1000);
  }
};

// Интервал фейковых данных
let fakeDataInterval: number | null = null;

// Генерация фейковых метрик для тестирования
const generateFakeMetrics = (): Metrics => {
  const baseActiveUsers = 800 + Math.round(Math.sin(Date.now() / 10000) * 100);
  const baseSales = 80 + Math.round(Math.cos(Date.now() / 15000) * 20);
  
  return {
    timestamp: Date.now(),
    activeUsers: baseActiveUsers + Math.round(Math.random() * 50),
    requestsPerSecond: 0.2 + Math.random() * 0.4,
    responseTimeMs: 50 + Math.random() * 150,
    conversionRate: 1.5 + Math.random() * 1.5,
    sales: baseSales + Math.round(Math.random() * 30),
    errorRate: 0.5 + Math.random() * 1.0,
    errorsByType: {
      'Сетевые': Math.round(5 + Math.random() * 5),
      'База данных': Math.round(2 + Math.random() * 5),
      'Авторизация': Math.round(1 + Math.random() * 3),
      'API': Math.round(3 + Math.random() * 5),
      'Прочие': Math.round(1 + Math.random() * 2),
    },
    serverLoad: 20 + Math.random() * 60,
    databaseConnections: Math.round(10 + Math.random() * 40)
  };
};

// Запуск генерации фейковых данных
const startFakeDataGeneration = () => {
  // Если интервал уже запущен, не создаем новый
  if (fakeDataInterval !== null) {
    return;
  }
  
  console.log('Запущена генерация тестовых данных');
  connectionState.isConnected = true;
  
  // Инициализируем начальные значения времени обновления
  const now = Date.now();
  lastUpdateTime.realtime = now;
  lastUpdateTime.hour = now;
  lastUpdateTime.day = now;
  lastUpdateTime.week = now;
  
  // Устанавливаем интервал для генерации данных
  fakeDataInterval = window.setInterval(() => {
    const fakeMetrics = generateFakeMetrics();
    
    // Применяем сглаживание данных
    const smoothedMetrics = smoothData(fakeMetrics);
    
    // Обновление текущих метрик
    Object.assign(metrics, smoothedMetrics);
    
    // Добавление в историю реального времени
    metricsHistory.realtime.push({...smoothedMetrics});
    if (metricsHistory.realtime.length > historyConfigs.realtime.maxItems) {
      metricsHistory.realtime.shift();
    }
    
    // Обрабатываем исторические данные для других периодов
    processHistoricalData(smoothedMetrics);
    
    // Всегда уведомляем UI о новых данных для выбранного периода
    historyMetrics.splice(0, historyMetrics.length, ...metricsHistory[currentPeriod]);
    
    // Вызов колбэка для обновления UI
    if (onUpdateCallback) {
      onUpdateCallback();
    }
  }, 1000);
};

// Остановка генерации фейковых данных
const stopFakeDataGeneration = () => {
  if (fakeDataInterval !== null) {
    clearInterval(fakeDataInterval);
    fakeDataInterval = null;
  }
};

// Обработка исторических данных для разных периодов
function processHistoricalData(newMetrics: Metrics) {
  const now = Date.now();
  
  // Периоды и их частота обновления (в миллисекундах)
  const periods: PeriodType[] = ['hour', 'day', 'week'];
  
  for (const period of periods) {
    // Проверяем, нужно ли обновлять данные для этого периода
    if (now - lastUpdateTime[period] >= updateFrequency[period]) {
      // Для каждого периода данные меняются с разной интенсивностью
      let variation = 0.05; // Базовая вариация для часа
      
      if (period === 'day') {
        variation = 0.03; // Меньшая вариация для дня
      } else if (period === 'week') {
        variation = 0.01; // Минимальная вариация для недели
      }
      
      // Генерируем новые метрики с уменьшенной вариацией для более длительных периодов
      const historicalMetrics = generateHistoricalMetricsForPeriod(
        {...metricsHistory[period][metricsHistory[period].length - 1] || newMetrics},
        variation
      );
      
      // Обновляем историю для данного периода
      metricsHistory[period].push(historicalMetrics);
      if (metricsHistory[period].length > historyConfigs[period].maxItems) {
        metricsHistory[period].shift();
      }
      
      // Обновляем время последнего обновления
      lastUpdateTime[period] = now;
    }
  }
}

// Генерация исторических метрик с заданной вариацией
function generateHistoricalMetricsForPeriod(baseMetrics: Metrics, variationFactor: number): Metrics {
  return {
    timestamp: baseMetrics.timestamp,
    activeUsers: applyVariation(baseMetrics.activeUsers, variationFactor),
    requestsPerSecond: applyVariation(baseMetrics.requestsPerSecond, variationFactor),
    responseTimeMs: applyVariation(baseMetrics.responseTimeMs, variationFactor),
    errorRate: applyVariation(baseMetrics.errorRate, variationFactor),
    conversionRate: applyVariation(baseMetrics.conversionRate, variationFactor),
    sales: applyVariation(baseMetrics.sales, variationFactor),
    errorsByType: {
      'Сетевые': baseMetrics.errorsByType['Сетевые'] ? applyVariation(baseMetrics.errorsByType['Сетевые'], variationFactor) : 0,
      'База данных': baseMetrics.errorsByType['База данных'] ? applyVariation(baseMetrics.errorsByType['База данных'], variationFactor) : 0,
      'Авторизация': baseMetrics.errorsByType['Авторизация'] ? applyVariation(baseMetrics.errorsByType['Авторизация'], variationFactor) : 0,
      'API': baseMetrics.errorsByType['API'] ? applyVariation(baseMetrics.errorsByType['API'], variationFactor) : 0,
      'Прочие': baseMetrics.errorsByType['Прочие'] ? applyVariation(baseMetrics.errorsByType['Прочие'], variationFactor) : 0
    },
    serverLoad: applyVariation(baseMetrics.serverLoad, variationFactor),
    databaseConnections: applyVariation(baseMetrics.databaseConnections, variationFactor)
  };
}

// Применяет случайную вариацию к числовому значению
function applyVariation(value: number, factor: number): number {
  const variation = value * factor * (Math.random() * 2 - 1);
  return Math.max(0, value + variation);
}

// Закрытие соединения
export const closeWebSocket = () => {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }
  
  if (ws) {
    ws.close();
    ws = null;
  }
  
  connectionState.isConnected = false;
};

// Вызов соединения при импорте сервиса
connectWebSocket();

// Сразу запускаем генерацию фейковых данных, чтобы не ждать ошибки соединения
startFakeDataGeneration();

// Вспомогательный хук для компонентов
export const useMetrics = () => {
  // При монтировании компонента подключаемся к WebSocket если ещё не подключены
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    connectWebSocket();
  }
  
  return {
    metrics,
    historyMetrics,
    connectionState,
    formatNumber,
    setActivePeriod
  };
};

// Метод очистки ресурсов при размонтировании приложения
export const cleanupMetricsService = () => {
  closeWebSocket();
  stopFakeDataGeneration();
};

// Генерация данных для графиков на основе истории метрик
export const generateChartData = () => {
  const timeLabels = historyMetrics.map(m => {
    // Форматирование метки времени в зависимости от текущего периода
    switch (currentPeriod) {
      case 'realtime':
        return dayjs(m.timestamp).format('HH:mm:ss');
      case 'hour':
        return dayjs(m.timestamp).format('HH:mm');
      case 'day':
        return dayjs(m.timestamp).format('HH:00');
      case 'week':
        return dayjs(m.timestamp).format('DD.MM');
      default:
        return dayjs(m.timestamp).format('HH:mm:ss');
    }
  });
  
  // Общие настройки анимации для всех графиков
  const animation = {
    duration: 800,
    easing: 'easeOutQuad'
  };
  
  return {
    userActivity: {
      labels: timeLabels,
      datasets: [
        {
          label: 'Активные пользователи',
          data: historyMetrics.map(m => m.activeUsers),
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        }
      ]
    },
    requests: {
      labels: timeLabels,
      datasets: [
        {
          label: 'Запросы',
          data: historyMetrics.map(m => m.requestsPerSecond),
          borderColor: '#10b981',
          backgroundColor: 'rgba(16, 185, 129, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        }
      ]
    },
    responseTime: {
      labels: timeLabels,
      datasets: [
        {
          label: 'мс',
          data: historyMetrics.map(m => m.responseTimeMs),
          borderColor: '#f59e0b',
          backgroundColor: 'rgba(245, 158, 11, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        }
      ]
    },
    sales: {
      labels: timeLabels,
      datasets: [
        {
          label: 'Продажи',
          data: historyMetrics.map(m => m.sales),
          borderColor: '#8b5cf6',
          backgroundColor: 'rgba(139, 92, 246, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4,
          animation
        }
      ]
    },
    conversion: {
      labels: timeLabels,
      datasets: [
        {
          label: '%',
          data: historyMetrics.map(m => m.conversionRate),
          borderColor: '#8b5cf6',
          backgroundColor: 'rgba(139, 92, 246, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        }
      ]
    },
    errorRate: {
      labels: timeLabels,
      datasets: [
        {
          label: '%',
          data: historyMetrics.map(m => m.errorRate),
          borderColor: '#ef4444',
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        }
      ]
    },
    errorTypes: {
      labels: Object.keys(metrics.errorsByType),
      datasets: [
        {
          label: 'Ошибки',
          data: Object.values(metrics.errorsByType),
          backgroundColor: [
            'rgba(239, 68, 68, 0.7)',
            'rgba(245, 158, 11, 0.7)',
            'rgba(16, 185, 129, 0.7)',
            'rgba(59, 130, 246, 0.7)',
            'rgba(139, 92, 246, 0.7)'
          ],
          borderWidth: 1
        }
      ]
    },
    // Добавляем географическое распределение
    geoDistribution: {
      labels: ['Москва', 'Санкт-Петербург', 'Екатеринбург', 'Новосибирск', 'Казань', 'Краснодар', 'Владивосток', 'Калининград'],
      datasets: [
        {
          label: 'Пользователи',
          data: [284, 152, 98, 76, 64, 53, 42, 37],
          backgroundColor: 'rgba(59, 130, 246, 0.7)',
          borderColor: 'rgba(59, 130, 246, 1)',
          borderWidth: 1
        }
      ]
    }
  };
}; 