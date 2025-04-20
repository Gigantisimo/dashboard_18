<template>
  <div class="sales-view">
    <div class="header">
      <h1>Статистика продаж</h1>
      <div class="period-selector">
        <button 
          v-for="period in periods" 
          :key="period.value" 
          @click="selectedPeriod = period.value"
          :class="{ active: selectedPeriod === period.value }"
          class="period-button"
        >
          {{ period.label }}
        </button>
      </div>
    </div>
    
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-header">
          <h3>Общая выручка</h3>
          <span class="trend up">+12.5%</span>
        </div>
        <div class="stat-value">₽ {{ formatNumber(salesData.totalRevenue) }}</div>
        <div class="stat-chart">
          <canvas ref="revenueChart"></canvas>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-header">
          <h3>Продажи</h3>
          <span class="trend up">+8.2%</span>
        </div>
        <div class="stat-value">{{ salesData.totalSales }}</div>
        <div class="stat-chart">
          <canvas ref="salesChart"></canvas>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-header">
          <h3>Средний чек</h3>
          <span class="trend up">+4.1%</span>
        </div>
        <div class="stat-value">₽ {{ formatNumber(salesData.averageOrder) }}</div>
        <div class="stat-chart">
          <canvas ref="avgOrderChart"></canvas>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-header">
          <h3>Конверсия</h3>
          <span class="trend down">-2.3%</span>
        </div>
        <div class="stat-value">{{ salesData.conversionRate }}%</div>
        <div class="stat-chart">
          <canvas ref="conversionChart"></canvas>
        </div>
      </div>
    </div>
    
    <div class="sales-charts">
      <div class="chart-card large">
        <div class="chart-header">
          <h3>Динамика продаж</h3>
          <div class="chart-legend">
            <div class="legend-item">
              <span class="legend-color" style="background-color: rgba(59, 130, 246, 0.5)"></span>
              <span>Выручка</span>
            </div>
            <div class="legend-item">
              <span class="legend-color" style="background-color: rgba(16, 185, 129, 0.5)"></span>
              <span>Заказы</span>
            </div>
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="salesTrendChart"></canvas>
        </div>
      </div>
      
      <div class="chart-card">
        <div class="chart-header">
          <h3>Распределение по категориям</h3>
        </div>
        <div class="chart-container">
          <canvas ref="categoriesChart"></canvas>
        </div>
      </div>
      
      <div class="chart-card">
        <div class="chart-header">
          <h3>Популярные товары</h3>
          <button class="view-all-btn">Смотреть все</button>
        </div>
        <div class="products-list">
          <div class="product-item" v-for="(product, index) in topProducts" :key="index">
            <div class="product-rank">{{ index + 1 }}</div>
            <div class="product-info">
              <div class="product-name">{{ product.name }}</div>
              <div class="product-details">
                {{ product.sales }} продаж · ₽{{ formatNumber(product.revenue) }}
              </div>
            </div>
            <div class="product-growth" :class="{ 'up': product.growth > 0, 'down': product.growth < 0 }">
              {{ product.growth > 0 ? '+' : '' }}{{ product.growth }}%
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="sales-table-section">
      <div class="section-header">
        <h2>Последние транзакции</h2>
        <div class="section-actions">
          <div class="search-box">
            <input 
              type="text" 
              placeholder="Поиск по заказам..." 
              v-model="searchQuery"
              @input="filterTransactions"
            >
          </div>
          <button class="export-btn">
            <span>↓</span> Экспорт
          </button>
        </div>
      </div>
      
      <div class="sales-table">
        <table>
          <thead>
            <tr>
              <th>ID Заказа</th>
              <th>Клиент</th>
              <th>Дата</th>
              <th>Продукты</th>
              <th>Сумма</th>
              <th>Статус</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="transaction in filteredTransactions" :key="transaction.id">
              <td>
                <div class="order-id">#{{ transaction.id }}</div>
              </td>
              <td>
                <div class="customer-info">
                  <div class="customer-avatar">{{ getInitials(transaction.customer) }}</div>
                  <div class="customer-name">{{ transaction.customer }}</div>
                </div>
              </td>
              <td>{{ transaction.date }}</td>
              <td>{{ transaction.products }}</td>
              <td>₽ {{ formatNumber(transaction.amount) }}</td>
              <td>
                <span class="status-badge" :class="transaction.status.toLowerCase()">
                  {{ transaction.status }}
                </span>
              </td>
              <td>
                <button class="icon-button">
                  <span>⋮</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <div class="pagination">
        <button class="pagination-btn" :disabled="currentPage === 1" @click="currentPage--">
          &lt;
        </button>
        <div class="pagination-info">
          Страница {{ currentPage }} из {{ totalPages }}
        </div>
        <button class="pagination-btn" :disabled="currentPage === totalPages" @click="currentPage++">
          &gt;
        </button>
      </div>
    </div>

    <!-- Добавляем новый компонент для отображения информации о WebSocket-соединении -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Информация о WebSocket-соединении</h3>
          <div class="header-actions">
            <button @click="reconnectWebSocket" class="ws-button" :disabled="connectionState.reconnecting">
              {{ connectionState.reconnecting ? 'Переподключение...' : 'Переподключить' }}
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="websocket-info">
            <div class="ws-status">
              <div class="status-indicator">
                <span class="status-dot" :class="connectionState.isConnected ? 'connected' : 'disconnected'"></span>
                <span class="status-text">{{ connectionState.isConnected ? 'Подключено' : 'Отключено' }}</span>
              </div>
              <div class="metrics-counter">
                <div>Получено метрик: <strong>{{ metricsReceived }}</strong></div>
                <div>Последнее обновление: <strong>{{ lastUpdateTime }}</strong></div>
              </div>
            </div>
            
            <div class="ws-details">
              <h4>Параметры соединения:</h4>
              <div class="ws-detail-item">
                <div>URL:</div>
                <div><code>ws://localhost:8080/ws</code></div>
              </div>
              <div class="ws-detail-item">
                <div>Протокол:</div>
                <div><code>WebSocket (RFC 6455)</code></div>
              </div>
              <div class="ws-detail-item">
                <div>Сглаживание данных:</div>
                <div><code>Коэффициент: 0.3</code></div>
              </div>
              <div class="ws-detail-item">
                <div>Обработчик событий:</div>
                <div><code>onmessage, onopen, onclose, onerror</code></div>
              </div>
              <div class="ws-detail-item">
                <div>Переподключение:</div>
                <div><code>Экспоненциальная задержка (1-30 сек)</code></div>
              </div>
            </div>
            
            <div class="ws-code">
              <h4>Пример кода WebSocket-подключения:</h4>
              <pre><code>// Установка WebSocket соединения
const ws = new WebSocket('ws://localhost:8080/ws');

// Обработчик получения данных
ws.onmessage = (event) => {
  const newMetrics = JSON.parse(event.data);
  
  // Применяем сглаживание данных для плавного отображения
  const smoothedMetrics = smoothData(newMetrics);
  
  // Обновляем метрики в приложении
  Object.assign(metrics, smoothedMetrics);
  
  // Обрабатываем исторические данные
  processHistoricalData(smoothedMetrics);
};</code></pre>
            </div>
            
            <div v-if="connectionState.lastError" class="ws-error">
              <h4>Последняя ошибка:</h4>
              <pre>{{ connectionState.lastError }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, onMounted, computed, onUnmounted } from 'vue';
import Chart from 'chart.js/auto';
import { connectionState, connectWebSocket, closeWebSocket } from '@/services/metricsService';

export default defineComponent({
  name: 'SalesView',
  setup() {
    const selectedPeriod = ref('month');
    const searchQuery = ref('');
    const currentPage = ref(1);
    const itemsPerPage = ref(10);
    
    const periods = [
      { label: 'День', value: 'day' },
      { label: 'Неделя', value: 'week' },
      { label: 'Месяц', value: 'month' },
      { label: 'Квартал', value: 'quarter' },
      { label: 'Год', value: 'year' },
    ];
    
    const revenueChart = ref(null);
    const salesChart = ref(null);
    const avgOrderChart = ref(null);
    const conversionChart = ref(null);
    const salesTrendChart = ref(null);
    const categoriesChart = ref(null);
    
    // Данные о продажах
    const salesData = reactive({
      totalRevenue: 4582750,
      totalSales: 1293,
      averageOrder: 3544,
      conversionRate: 3.8,
      revenueByPeriod: [2800000, 3100000, 2950000, 3250000, 3600000, 4200000, 4582750],
      salesByPeriod: [845, 912, 878, 967, 1054, 1182, 1293],
      avgOrderByPeriod: [3310, 3390, 3350, 3360, 3410, 3550, 3544],
      conversionByPeriod: [4.2, 4.0, 3.9, 3.8, 3.6, 3.7, 3.8]
    });
    
    // Популярные товары
    const topProducts = [
      { name: 'Смартфон Galaxy A52', sales: 145, revenue: 4350000, growth: 12.3 },
      { name: 'Ноутбук LenovoBook Pro', sales: 87, revenue: 6960000, growth: 8.7 },
      { name: 'Наушники SonicSound X', sales: 230, revenue: 1840000, growth: 27.5 },
      { name: 'Планшет TabMax 10', sales: 64, revenue: 1920000, growth: -3.2 },
      { name: 'Умные часы FitLife 4', sales: 118, revenue: 1770000, growth: 15.9 }
    ];
    
    // Данные о транзакциях
    const transactions = [
      { id: '38291', customer: 'Иванов Сергей', date: '15.04.2023', products: 'Смартфон Galaxy A52', amount: 29900, status: 'Выполнен' },
      { id: '38290', customer: 'Петрова Анна', date: '15.04.2023', products: 'Ноутбук LenovoBook Pro', amount: 79900, status: 'Выполнен' },
      { id: '38289', customer: 'Сидоров Алексей', date: '14.04.2023', products: 'Наушники SonicSound X', amount: 7990, status: 'Выполнен' },
      { id: '38288', customer: 'Козлова Елена', date: '14.04.2023', products: 'Умные часы FitLife 4', amount: 15990, status: 'Выполнен' },
      { id: '38287', customer: 'Новиков Дмитрий', date: '14.04.2023', products: 'Планшет TabMax 10', amount: 29990, status: 'Отменен' },
      { id: '38286', customer: 'Морозова Ирина', date: '13.04.2023', products: 'Смартфон Galaxy A52, Чехол защитный', amount: 31890, status: 'Выполнен' },
      { id: '38285', customer: 'Волков Павел', date: '13.04.2023', products: 'Наушники SonicSound X', amount: 7990, status: 'В обработке' },
      { id: '38284', customer: 'Соколова Татьяна', date: '12.04.2023', products: 'Ноутбук LenovoBook Pro, Мышь беспроводная', amount: 82780, status: 'Выполнен' },
      { id: '38283', customer: 'Лебедев Игорь', date: '12.04.2023', products: 'Планшет TabMax 10', amount: 29990, status: 'Доставляется' },
      { id: '38282', customer: 'Кузнецова Ольга', date: '11.04.2023', products: 'Умные часы FitLife 4', amount: 15990, status: 'Выполнен' },
      { id: '38281', customer: 'Попов Артем', date: '11.04.2023', products: 'Смартфон Galaxy A52', amount: 29900, status: 'Выполнен' },
      { id: '38280', customer: 'Васильева Марина', date: '10.04.2023', products: 'Наушники SonicSound X, Кабель USB-C', amount: 8780, status: 'Выполнен' },
      { id: '38279', customer: 'Никитин Роман', date: '10.04.2023', products: 'Ноутбук LenovoBook Pro', amount: 79900, status: 'Отменен' },
      { id: '38278', customer: 'Михайлова Наталья', date: '09.04.2023', products: 'Планшет TabMax 10, Защитное стекло', amount: 31490, status: 'Выполнен' },
      { id: '38277', customer: 'Федоров Виктор', date: '09.04.2023', products: 'Умные часы FitLife 4', amount: 15990, status: 'Выполнен' }
    ];
    
    // Отфильтрованные транзакции с поддержкой поиска
    const filteredTransactions = ref([...transactions]);
    
    // Функция фильтрации транзакций
    const filterTransactions = () => {
      if (!searchQuery.value) {
        filteredTransactions.value = [...transactions];
      } else {
        const query = searchQuery.value.toLowerCase();
        filteredTransactions.value = transactions.filter(
          transaction => 
            transaction.id.toLowerCase().includes(query) ||
            transaction.customer.toLowerCase().includes(query) ||
            transaction.products.toLowerCase().includes(query) ||
            transaction.status.toLowerCase().includes(query)
        );
      }
      currentPage.value = 1;
    };
    
    // Вычисляем общее количество страниц
    const totalPages = computed(() => 
      Math.ceil(filteredTransactions.value.length / itemsPerPage.value)
    );
    
    // Получение инициалов из имени
    const getInitials = (name: string): string => {
      return name
        .split(' ')
        .map(part => part.charAt(0))
        .join('')
        .toUpperCase()
        .slice(0, 2);
    };
    
    // Форматирование числа с разделением тысяч
    const formatNumber = (value: number): string => {
      return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
    };
    
    // Инициализация графиков
    onMounted(() => {
      // Инициализация мини-графика выручки
      new Chart(revenueChart.value, {
        type: 'line',
        data: {
          labels: ['', '', '', '', '', '', ''],
          datasets: [{
            data: salesData.revenueByPeriod,
            borderColor: 'rgba(59, 130, 246, 0.8)',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            tension: 0.4,
            fill: true,
            pointRadius: 0,
            borderWidth: 2
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { display: false } },
          scales: {
            x: { display: false },
            y: { display: false }
          }
        }
      });
      
      // Инициализация мини-графика продаж
      new Chart(salesChart.value, {
        type: 'line',
        data: {
          labels: ['', '', '', '', '', '', ''],
          datasets: [{
            data: salesData.salesByPeriod,
            borderColor: 'rgba(16, 185, 129, 0.8)',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            tension: 0.4,
            fill: true,
            pointRadius: 0,
            borderWidth: 2
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { display: false } },
          scales: {
            x: { display: false },
            y: { display: false }
          }
        }
      });
      
      // Инициализация мини-графика среднего чека
      new Chart(avgOrderChart.value, {
        type: 'line',
        data: {
          labels: ['', '', '', '', '', '', ''],
          datasets: [{
            data: salesData.avgOrderByPeriod,
            borderColor: 'rgba(245, 158, 11, 0.8)',
            backgroundColor: 'rgba(245, 158, 11, 0.1)',
            tension: 0.4,
            fill: true,
            pointRadius: 0,
            borderWidth: 2
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { display: false } },
          scales: {
            x: { display: false },
            y: { display: false }
          }
        }
      });
      
      // Инициализация мини-графика конверсии
      new Chart(conversionChart.value, {
        type: 'line',
        data: {
          labels: ['', '', '', '', '', '', ''],
          datasets: [{
            data: salesData.conversionByPeriod,
            borderColor: 'rgba(239, 68, 68, 0.8)',
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            tension: 0.4,
            fill: true,
            pointRadius: 0,
            borderWidth: 2
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: { legend: { display: false } },
          scales: {
            x: { display: false },
            y: { display: false }
          }
        }
      });
      
      // Инициализация графика динамики продаж
      new Chart(salesTrendChart.value, {
        type: 'line',
        data: {
          labels: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль'],
          datasets: [
            {
              label: 'Выручка',
              data: salesData.revenueByPeriod,
              borderColor: 'rgba(59, 130, 246, 0.8)',
              backgroundColor: 'rgba(59, 130, 246, 0.1)',
              tension: 0.4,
              fill: false,
              borderWidth: 2,
              yAxisID: 'y'
            },
            {
              label: 'Заказы',
              data: salesData.salesByPeriod,
              borderColor: 'rgba(16, 185, 129, 0.8)',
              backgroundColor: 'rgba(16, 185, 129, 0.1)',
              tension: 0.4,
              fill: false,
              borderWidth: 2,
              yAxisID: 'y1'
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: { display: false }
          },
          scales: {
            y: {
              type: 'linear',
              display: true,
              position: 'left',
              title: {
                display: true,
                text: 'Выручка (₽)'
              }
            },
            y1: {
              type: 'linear',
              display: true,
              position: 'right',
              title: {
                display: true,
                text: 'Заказы'
              },
              grid: {
                drawOnChartArea: false
              }
            }
          }
        }
      });
      
      // Инициализация графика распределения по категориям
      new Chart(categoriesChart.value, {
        type: 'doughnut',
        data: {
          labels: ['Смартфоны', 'Ноутбуки', 'Планшеты', 'Аксессуары', 'Аудио'],
          datasets: [{
            data: [32, 25, 15, 18, 10],
            backgroundColor: [
              'rgba(59, 130, 246, 0.8)',
              'rgba(16, 185, 129, 0.8)',
              'rgba(245, 158, 11, 0.8)',
              'rgba(239, 68, 68, 0.8)',
              'rgba(139, 92, 246, 0.8)'
            ],
            borderWidth: 0
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'bottom',
              labels: {
                boxWidth: 12,
                padding: 15
              }
            }
          }
        }
      });
    });

    // Добавляем в setup() метод
    const metricsReceived = ref(0);
    const lastUpdateTime = ref('--:--:--');

    // Функция для переподключения WebSocket
    const reconnectWebSocket = () => {
      closeWebSocket();
      connectWebSocket();
    };

    // Обновляем счетчик полученных метрик
    const updateMetricsCounter = () => {
      metricsReceived.value++;
      lastUpdateTime.value = new Date().toLocaleTimeString();
    };

    // Подписываемся на обновления метрик
    onMounted(() => {
      // Пытаемся подключиться к WebSocket при монтировании компонента
      connectWebSocket();
      
      // Устанавливаем интервал для имитации получения метрик (для наглядности)
      const interval = setInterval(() => {
        if (connectionState.isConnected) {
          updateMetricsCounter();
        }
      }, 1000);
      
      // Очищаем интервал при размонтировании
      onUnmounted(() => {
        clearInterval(interval);
      });
    });
    
    return {
      selectedPeriod,
      periods,
      salesData,
      topProducts,
      searchQuery,
      filteredTransactions,
      currentPage,
      totalPages,
      revenueChart,
      salesChart,
      avgOrderChart,
      conversionChart,
      salesTrendChart,
      categoriesChart,
      getInitials,
      formatNumber,
      filterTransactions,
      connectionState,
      metricsReceived,
      lastUpdateTime,
      reconnectWebSocket,
    };
  }
});
</script>

<style lang="scss" scoped>
.sales-view {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  
  h1 {
    font-size: 1.75rem;
    font-weight: 600;
  }
  
  .period-selector {
    display: flex;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 4px;
    
    .period-button {
      background: none;
      border: none;
      color: var(--text-secondary);
      font-size: 0.875rem;
      font-weight: 500;
      padding: 0.5rem 1rem;
      cursor: pointer;
      border-radius: 6px;
      transition: all 0.2s ease;
      
      &:hover:not(.active) {
        background-color: rgba(255, 255, 255, 0.05);
      }
      
      &.active {
        background-color: var(--accent);
        color: white;
      }
    }
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background-color: var(--bg-card);
  border-radius: 8px;
  padding: 1.25rem;
  border: 1px solid var(--border);
  
  .stat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    
    h3 {
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--text-secondary);
    }
    
    .trend {
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
      padding: 0.25rem 0.5rem;
      
      &.up {
        background-color: rgba(16, 185, 129, 0.1);
        color: var(--success);
      }
      
      &.down {
        background-color: rgba(239, 68, 68, 0.1);
        color: var(--danger);
      }
    }
  }
  
  .stat-value {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
  }
  
  .stat-chart {
    height: 60px;
  }
}

.sales-charts {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: 1200px) {
    grid-template-columns: 1fr 1fr;
    
    .chart-card.large {
      grid-column: span 2;
    }
  }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    
    .chart-card.large {
      grid-column: span 1;
    }
  }
}

.chart-card {
  background-color: var(--bg-card);
  border-radius: 8px;
  padding: 1.25rem;
  border: 1px solid var(--border);
  
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    
    h3 {
      font-size: 1rem;
      font-weight: 600;
    }
    
    .chart-legend {
      display: flex;
      gap: 1rem;
      
      .legend-item {
        display: flex;
        align-items: center;
        font-size: 0.75rem;
        
        .legend-color {
          width: 12px;
          height: 12px;
          border-radius: 2px;
          margin-right: 0.5rem;
        }
      }
    }
    
    .view-all-btn {
      background: none;
      border: none;
      color: var(--accent);
      font-size: 0.75rem;
      cursor: pointer;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
  
  .chart-container {
    height: 250px;
    position: relative;
  }
  
  .products-list {
    .product-item {
      display: flex;
      align-items: center;
      padding: 0.875rem 0;
      border-bottom: 1px solid var(--border);
      
      &:last-child {
        border-bottom: none;
      }
      
      .product-rank {
        width: 1.5rem;
        height: 1.5rem;
        background-color: rgba(255, 255, 255, 0.05);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 0.75rem;
        font-weight: 600;
        margin-right: 0.875rem;
      }
      
      .product-info {
        flex: 1;
        
        .product-name {
          font-size: 0.875rem;
          font-weight: 500;
          margin-bottom: 0.25rem;
        }
        
        .product-details {
          font-size: 0.75rem;
          color: var(--text-secondary);
        }
      }
      
      .product-growth {
        font-size: 0.75rem;
        font-weight: 600;
        border-radius: 9999px;
        padding: 0.25rem 0.5rem;
        
        &.up {
          background-color: rgba(16, 185, 129, 0.1);
          color: var(--success);
        }
        
        &.down {
          background-color: rgba(239, 68, 68, 0.1);
          color: var(--danger);
        }
      }
    }
  }
}

.sales-table-section {
  background-color: var(--bg-card);
  border-radius: 8px;
  padding: 1.25rem;
  border: 1px solid var(--border);
  margin-bottom: 2rem;
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    
    h2 {
      font-size: 1.125rem;
      font-weight: 600;
    }
    
    .section-actions {
      display: flex;
      align-items: center;
      gap: 1rem;
      
      .search-box {
        input {
          background-color: rgba(255, 255, 255, 0.05);
          border: 1px solid var(--border);
          border-radius: 4px;
          padding: 0.5rem 0.75rem;
          width: 250px;
          font-size: 0.875rem;
          color: var(--text-primary);
          
          &:focus {
            outline: none;
            border-color: var(--accent);
          }
          
          &::placeholder {
            color: var(--text-secondary);
          }
        }
      }
      
      .export-btn {
        display: flex;
        align-items: center;
        background-color: transparent;
        color: var(--text-primary);
        border: 1px solid var(--border);
        border-radius: 4px;
        padding: 0.5rem 1rem;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.05);
          border-color: var(--text-secondary);
        }
        
        span {
          margin-right: 0.375rem;
        }
      }
    }
  }
  
  .sales-table {
    margin-bottom: 1.5rem;
    overflow-x: auto;
    
    table {
      width: 100%;
      border-collapse: collapse;
      
      th, td {
        padding: 0.875rem 1rem;
        text-align: left;
      }
      
      th {
        color: var(--text-secondary);
        font-weight: 500;
        font-size: 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        border-bottom: 1px solid var(--border);
      }
      
      td {
        font-size: 0.875rem;
        border-bottom: 1px solid rgba(71, 85, 105, 0.1);
      }
      
      tbody tr:hover {
        background-color: rgba(255, 255, 255, 0.025);
      }
      
      .order-id {
        font-family: monospace;
        font-weight: 500;
      }
      
      .customer-info {
        display: flex;
        align-items: center;
        
        .customer-avatar {
          width: 28px;
          height: 28px;
          background-color: var(--accent);
          color: white;
          font-size: 0.75rem;
          font-weight: 600;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 0.75rem;
        }
      }
      
      .status-badge {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
        
        &.выполнен {
          background-color: rgba(16, 185, 129, 0.1);
          color: var(--success);
        }
        
        &.отменен {
          background-color: rgba(239, 68, 68, 0.1);
          color: var(--danger);
        }
        
        &.доставляется {
          background-color: rgba(245, 158, 11, 0.1);
          color: var(--warning);
        }
        
        &.в.обработке {
          background-color: rgba(59, 130, 246, 0.1);
          color: var(--info);
        }
      }
      
      .icon-button {
        background: none;
        border: none;
        color: var(--text-secondary);
        cursor: pointer;
        width: 28px;
        height: 28px;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s ease;
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.1);
          color: var(--text-primary);
        }
      }
    }
  }
  
  .pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    
    .pagination-btn {
      width: 32px;
      height: 32px;
      border-radius: 4px;
      border: 1px solid var(--border);
      background-color: transparent;
      color: var(--text-primary);
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.2s ease;
      
      &:hover:not(:disabled) {
        background-color: rgba(255, 255, 255, 0.05);
        border-color: var(--text-secondary);
      }
      
      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    }
    
    .pagination-info {
      font-size: 0.875rem;
      color: var(--text-secondary);
    }
  }
}

:deep(.chart-js-tooltip) {
  background-color: var(--bg-card) !important;
  border: 1px solid var(--border) !important;
  border-radius: 4px !important;
  color: var(--text-primary) !important;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06) !important;
  padding: 0.5rem 0.75rem !important;
  font-size: 0.75rem !important;
}

.websocket-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.ws-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: rgba(30, 41, 59, 0.5);
  padding: 1rem;
  border-radius: 0.5rem;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-dot {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  
  &.connected {
    background-color: var(--success);
    box-shadow: 0 0 8px rgba(16, 185, 129, 0.6);
  }
  
  &.disconnected {
    background-color: var(--danger);
    box-shadow: 0 0 8px rgba(239, 68, 68, 0.6);
  }
}

.status-text {
  font-weight: 600;
}

.metrics-counter {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: 0.9rem;
}

.ws-details {
  h4 {
    margin-bottom: 0.75rem;
    font-size: 1rem;
    font-weight: 600;
  }
  
  .ws-detail-item {
    display: flex;
    margin-bottom: 0.5rem;
    font-size: 0.9rem;
    
    div:first-child {
      min-width: 150px;
      font-weight: 500;
    }
    
    code {
      padding: 0.125rem 0.25rem;
      background-color: rgba(30, 41, 59, 0.5);
      border-radius: 0.25rem;
      font-family: monospace;
      font-size: 0.85rem;
    }
  }
}

.ws-code {
  h4 {
    margin-bottom: 0.75rem;
    font-size: 1rem;
    font-weight: 600;
  }
  
  pre {
    background-color: rgba(30, 41, 59, 0.5);
    padding: 1rem;
    border-radius: 0.5rem;
    overflow-x: auto;
    font-family: monospace;
    font-size: 0.85rem;
    line-height: 1.5;
  }
}

.ws-error {
  h4 {
    margin-bottom: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--danger);
  }
  
  pre {
    background-color: rgba(239, 68, 68, 0.1);
    padding: 1rem;
    border-radius: 0.5rem;
    border: 1px solid rgba(239, 68, 68, 0.3);
    color: var(--danger);
    font-family: monospace;
    font-size: 0.85rem;
  }
}

.ws-button {
  background-color: var(--accent);
  color: white;
  border: none;
  border-radius: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: var(--accent-hover);
  }
  
  &:disabled {
    background-color: rgba(59, 130, 246, 0.5);
    cursor: not-allowed;
  }
}
</style> 