import { check } from 'k6';
import { WebSocket } from 'k6/experimental/websockets';
import { sleep } from 'k6';
import http from 'k6/http';

// Конфигурация теста
export let options = {
  // Основной сценарий: постепенное увеличение нагрузки
  stages: [
    { duration: '30s', target: 100 }, // Разогрев: увеличение до 100 пользователей за 30 секунд
    { duration: '1m', target: 500 },  // Постепенное увеличение до 500 пользователей за 1 минуту
    { duration: '2m', target: 1000 }, // Увеличение до 1000 пользователей за 2 минуты
    { duration: '1m', target: 3000 }, // Пиковая нагрузка: 3000 пользователей за 1 минуту
    { duration: '1m', target: 3000 }, // Поддержание пиковой нагрузки в течение 1 минуты
    { duration: '2m', target: 0 },    // Постепенное снижение до 0 за 2 минуты
  ],
  thresholds: {
    // Пороговые значения для метрик
    'http_req_duration': ['p(95)<500'], // 95% запросов должны завершаться менее чем за 500 мс
    'ws_connecting_duration': ['p(95)<1000'], // 95% WebSocket соединений должны устанавливаться менее чем за 1 секунду
    'ws_session_duration': ['p(95)>30000'], // 95% WebSocket сессий должны длиться более 30 секунд
  }
};

// Переменные для подсчета метрик
let successfulConnections = 0;
let failedConnections = 0;
let messageReceived = 0;
let connectionErrors = {};

// Основная функция теста
export default function() {
  // Сначала выполняем обычный HTTP запрос к эндпоинту healthcheck
  let healthcheck = http.get('http://localhost:8080/health');
  
  // Проверяем успешность HTTP запроса
  check(healthcheck, {
    'health endpoint status is 200': (r) => r.status === 200,
    'health endpoint response contains status:ok': (r) => r.body.includes('ok'),
  });
  
  // Если healthcheck не прошел, не пытаемся устанавливать WebSocket соединение
  if (healthcheck.status !== 200) {
    console.log('Healthcheck failed, skipping WebSocket connection');
    failedConnections++;
    return;
  }
  
  // Устанавливаем WebSocket соединение
  const url = 'ws://localhost:8080/ws';
  let ws = new WebSocket(url);
  let messagesCount = 0;
  let isConnected = false;
  
  // Обработчик события при успешном соединении
  ws.onopen = () => {
    isConnected = true;
    successfulConnections++;
  };
  
  // Обработчик получения сообщения
  ws.onmessage = (message) => {
    messageReceived++;
    messagesCount++;
    
    try {
      // Пробуем парсить JSON, чтобы убедиться, что получаем валидные метрики
      const data = JSON.parse(message);
      check(data, {
        'message contains timestamp': (d) => d.hasOwnProperty('timestamp'),
        'message contains activeUsers': (d) => d.hasOwnProperty('activeUsers'),
        'message contains requestsPerSecond': (d) => d.hasOwnProperty('requestsPerSecond'),
      });
    } catch (e) {
      console.log(`Error parsing message: ${e.message}`);
    }
  };
  
  // Обработчик ошибок
  ws.onerror = (e) => {
    failedConnections++;
    const errorMessage = e.error || 'unknown error';
    connectionErrors[errorMessage] = (connectionErrors[errorMessage] || 0) + 1;
  };
  
  // Обработчик закрытия соединения
  ws.onclose = () => {
    // Ничего не делаем при закрытии соединения
  };
  
  // Имитируем работающего клиента в течение некоторого времени
  // Время варьируется, чтобы имитировать разных пользователей
  const sessionDuration = Math.floor(Math.random() * 20) + 10; // от 10 до 30 секунд
  sleep(sessionDuration);
  
  // Закрываем соединение если оно все еще открыто
  if (isConnected) {
    ws.close();
  }
  
  // Проверка: получили ли мы хотя бы одно сообщение
  check(messagesCount, {
    'received at least one message': (count) => count > 0,
  });
}

// Функция, которая выполняется после завершения всех итераций
export function handleSummary(data) {
  console.log('Successful connections:', successfulConnections);
  console.log('Failed connections:', failedConnections);
  console.log('Total messages received:', messageReceived);
  console.log('Connection errors by type:', connectionErrors);
  
  return {
    'stdout': JSON.stringify(data, null, 2), // Выводим данные в стандартный вывод
    './load_test_results.json': JSON.stringify(data, null, 2), // Сохраняем в файл
  };
} 