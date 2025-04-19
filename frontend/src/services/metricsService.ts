import { ref, reactive, onUnmounted } from 'vue';
import dayjs from 'dayjs';

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

// Фактические метрики, которые будут обновляться
export const metrics = reactive<Metrics>({...defaultMetrics});

// История метрик для графиков
const historySize = 20;
export const historyMetrics = reactive<Metrics[]>([]);

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

// Соединение с WebSocket
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
      
      // Обновление текущих метрик
      Object.assign(metrics, newMetrics);
      
      // Добавление в историю
      historyMetrics.push({...newMetrics});
      if (historyMetrics.length > historySize) {
        historyMetrics.shift();
      }
      
      // Вызов колбэка
      if (onUpdateCallback) {
        onUpdateCallback();
      }
    };
    
    ws.onerror = (error) => {
      connectionState.lastError = 'Ошибка подключения';
      console.error('WebSocket ошибка:', error);
    };
    
    ws.onclose = () => {
      connectionState.isConnected = false;
      console.log('WebSocket соединение закрыто');
      
      // Пытаемся переподключиться
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
      }
      
      reconnectTimeout = window.setTimeout(connectWebSocket, 2000);
    };
  } catch (error) {
    connectionState.lastError = `Ошибка: ${error}`;
    connectionState.reconnecting = false;
    
    // Пытаемся переподключиться
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
    }
    
    reconnectTimeout = window.setTimeout(connectWebSocket, 2000);
  }
};

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

// Вспомогательный хук для компонентов
export const useMetrics = () => {
  // При монтировании компонента подключаемся к WebSocket если ещё не подключены
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    connectWebSocket();
  }
  
  // При размонтировании ничего не делаем, так как соединение должно поддерживаться
  // для всех компонентов
  
  return {
    metrics,
    historyMetrics,
    connectionState,
    formatNumber
  };
};

// Генерация данных для графиков на основе истории метрик
export const generateChartData = () => {
  const timeLabels = historyMetrics.map(m => 
    dayjs(m.timestamp).format('HH:mm:ss')
  );
  
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
          backgroundColor: 'rgba(139, 92, 246, 0.7)',
          borderColor: 'rgba(139, 92, 246, 1)',
          borderWidth: 1
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