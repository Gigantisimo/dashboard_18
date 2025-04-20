<template>
  <div class="users-view">
    <!-- KPI карточки -->
    <div class="kpi-cards">
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Активные сейчас</h3>
            <div class="trend up">+{{ trends.activeUsers }}%</div>
          </div>
          <div class="card-value">{{ formatNumber(kpiData.activeUsers) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot green"></span>
              Онлайн
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Новые сегодня</h3>
            <div class="trend up">+{{ trends.newUsers }}%</div>
          </div>
          <div class="card-value">{{ formatNumber(kpiData.newUsers) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot blue"></span>
              С начала дня
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>DAU</h3>
            <div class="trend up">+{{ trends.dau }}%</div>
          </div>
          <div class="card-value">{{ formatNumber(kpiData.dau) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot purple"></span>
              За 24 часа
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>MAU</h3>
            <div class="trend up">+{{ trends.mau }}%</div>
          </div>
          <div class="card-value">{{ formatNumber(kpiData.mau) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot yellow"></span>
              За 30 дней
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- График активности пользователей -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Активность пользователей</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <LineChart 
              v-if="chartData.userActivity.labels.length > 0"
              :chartData="chartData.userActivity" 
              :options="chartOptions.userActivity"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Географическая активность и демографические показатели -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Географическое распределение</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <BarChart 
              v-if="chartData.geoDistribution.labels.length > 0"
              :chartData="chartData.geoDistribution" 
              :options="chartOptions.geoDistribution"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Демографические показатели</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <DoughnutChart 
              v-if="chartData.demographics.labels.length > 0"
              :chartData="chartData.demographics" 
              :options="chartOptions.demographics"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Удержание пользователей -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Удержание пользователей</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="retention-table">
            <table>
              <thead>
                <tr>
                  <th>Когорта</th>
                  <th>Пользователи</th>
                  <th>День 1</th>
                  <th>День 3</th>
                  <th>День 7</th>
                  <th>День 14</th>
                  <th>День 30</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(cohort, index) in retentionData" :key="index">
                  <td>{{ cohort.date }}</td>
                  <td>{{ cohort.users }}</td>
                  <td v-for="(rate, i) in cohort.retention" :key="i" 
                      :class="getRetentionClass(rate)">
                    {{ rate }}%
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, onMounted, ref, watch, computed, onUnmounted } from 'vue';
import LineChart from '@/components/charts/LineChart.vue';
import DoughnutChart from '@/components/charts/DoughnutChart.vue';
import BarChart from '@/components/charts/BarChart.vue';
import { useMetrics, formatNumber, generateChartData } from '@/services/metricsService';
import { useTimePeriod, PeriodType, onPeriodChange } from '@/services/timePeriodService';
import dayjs from 'dayjs';

export default defineComponent({
  name: 'UsersView',
  components: {
    LineChart,
    DoughnutChart,
    BarChart
  },
  setup() {
    // Используем сервис метрик и периодов времени
    const { metrics } = useMetrics();
    const { activePeriod } = useTimePeriod();
    
    // KPI данные, которые будут меняться в зависимости от периода
    const kpiData = reactive({
      activeUsers: 842,
      newUsers: 124,
      dau: 3842,
      mau: 28675
    });
    
    // Тренды для показателей
    const trends = reactive({
      activeUsers: '5.2',
      newUsers: '8.1',
      dau: '3.5',
      mau: '12.4'
    });
    
    // Данные для графиков, которые будут меняться в зависимости от периода
    const chartData = reactive({
      userActivity: {
        labels: [] as string[],
        datasets: [
          {
            label: 'Активные пользователи',
            data: [] as number[],
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            borderColor: 'rgba(59, 130, 246, 1)',
            borderWidth: 2,
            tension: 0.4,
            fill: true
          }
        ]
      },
      demographics: {
        labels: ['18-24', '25-34', '35-44', '45-54', '55+'],
        datasets: [
          {
            label: 'Возрастные группы',
            data: [15, 32, 28, 17, 8],
            backgroundColor: [
              'rgba(59, 130, 246, 0.7)',
              'rgba(139, 92, 246, 0.7)',
              'rgba(16, 185, 129, 0.7)',
              'rgba(245, 158, 11, 0.7)',
              'rgba(239, 68, 68, 0.7)'
            ],
            borderWidth: 1
          }
        ]
      },
      geoDistribution: {
        labels: [] as string[],
        datasets: [
          {
            label: 'Активные пользователи',
            data: [] as number[],
            backgroundColor: 'rgba(59, 130, 246, 0.7)',
            borderColor: 'rgba(59, 130, 246, 1)',
            borderWidth: 1
          }
        ]
      }
    });
    
    // Опции для графиков
    const chartOptions = {
      userActivity: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            mode: 'index',
            intersect: false,
            backgroundColor: 'rgba(30, 41, 59, 0.9)',
            titleColor: '#f1f5f9',
            bodyColor: '#f1f5f9',
            borderColor: '#334155',
            borderWidth: 1
          }
        },
        scales: {
          x: {
            grid: {
              color: 'rgba(71, 85, 105, 0.1)',
              drawBorder: false
            },
            ticks: {
              color: '#94a3b8'
            }
          },
          y: {
            grid: {
              color: 'rgba(71, 85, 105, 0.1)',
              drawBorder: false
            },
            ticks: {
              color: '#94a3b8'
            }
          }
        }
      },
      geoDistribution: {
        responsive: true,
        maintainAspectRatio: false,
        indexAxis: 'y' as const,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            backgroundColor: 'rgba(30, 41, 59, 0.9)',
            titleColor: '#f1f5f9',
            bodyColor: '#f1f5f9',
            borderColor: '#334155',
            borderWidth: 1
          }
        },
        scales: {
          x: {
            grid: {
              color: 'rgba(71, 85, 105, 0.1)',
              drawBorder: false
            },
            ticks: {
              color: '#94a3b8'
            }
          },
          y: {
            grid: {
              color: 'rgba(71, 85, 105, 0.1)',
              drawBorder: false
            },
            ticks: {
              color: '#94a3b8'
            }
          }
        }
      },
      demographics: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'right' as const,
            labels: {
              color: '#94a3b8',
              boxWidth: 12,
              padding: 15
            }
          },
          tooltip: {
            backgroundColor: 'rgba(30, 41, 59, 0.9)',
            titleColor: '#f1f5f9',
            bodyColor: '#f1f5f9',
            borderColor: '#334155',
            borderWidth: 1
          }
        }
      }
    };
    
    // Тестовые данные удержания пользователей
    const retentionData = ref([
      {
        date: '1-7 Апр',
        users: '2,451',
        retention: [100, 64, 42, 35, 28, 22]
      },
      {
        date: '8-14 Апр',
        users: '3,125',
        retention: [100, 68, 45, 38, 30, 24]
      },
      {
        date: '15-21 Апр',
        users: '2,845',
        retention: [100, 65, 43, 36, 29, 23]
      },
      {
        date: '22-28 Апр',
        users: '3,512',
        retention: [100, 70, 48, 40, 32, 26]
      },
      {
        date: '29 Апр - 5 Мая',
        users: '3,842',
        retention: [100, 72, 50, 42, null, null]
      }
    ]);
    
    // Вычисление класса ячейки удержания
    const getRetentionClass = (rate: number | null) => {
      if (rate === null) return 'unavailable';
      if (rate >= 60) return 'high';
      if (rate >= 40) return 'medium';
      if (rate >= 25) return 'low';
      return 'very-low';
    };
    
    // Функция обновления данных в зависимости от периода
    const updateDataByPeriod = (period: PeriodType) => {
      console.log(`Обновление данных для периода: ${period}`);
      
      // Определяем значения метрик в зависимости от периода
      switch (period) {
        case 'realtime':
          // Для режима реального времени показываем текущие значения
          kpiData.activeUsers = 800 + Math.floor(Math.random() * 100); // Текущие активные пользователи (разброс для реализма)
          kpiData.newUsers = 120 + Math.floor(Math.random() * 30); // Новые пользователи сегодня
          kpiData.dau = 3800 + Math.floor(Math.random() * 200); // DAU (почти не меняется в пределах дня)
          kpiData.mau = 28000 + Math.floor(Math.random() * 1000); // MAU стабилен для всех периодов
          break;
          
        case 'hour':
          // Для часового периода данные отличаются от реального времени незначительно
          kpiData.activeUsers = 750 + Math.floor(Math.random() * 100); // Среднее за час
          kpiData.newUsers = 100 + Math.floor(Math.random() * 30); // Новые за последний час
          kpiData.dau = 3800 + Math.floor(Math.random() * 200); // DAU остается стабильным
          kpiData.mau = 28000 + Math.floor(Math.random() * 1000); // MAU стабилен
          break;
          
        case 'day':
          // Для дневного периода меняется только показ активных пользователей (усредненный)
          kpiData.activeUsers = 650 + Math.floor(Math.random() * 100); // Среднее за день
          kpiData.newUsers = 130 + Math.floor(Math.random() * 20); // Новые пользователи за день
          kpiData.dau = 3800 + Math.floor(Math.random() * 200); // Это и есть DAU
          kpiData.mau = 28000 + Math.floor(Math.random() * 1000); // MAU стабилен
          break;
          
        case 'week':
          // Для недельного периода показываем среднее значение активных пользователей за неделю
          kpiData.activeUsers = 580 + Math.floor(Math.random() * 50); // Среднее за неделю
          kpiData.newUsers = 800 + Math.floor(Math.random() * 100); // Новые за неделю (больше чем за день)
          kpiData.dau = 3800 + Math.floor(Math.random() * 200); // Среднее DAU за последнюю неделю
          kpiData.mau = 28000 + Math.floor(Math.random() * 1000); // MAU стабилен
          break;
      }
      
      // Трендовые показатели отличаются в зависимости от периода, но логически согласованы
      // Для коротких периодов тренды более волатильны
      switch (period) {
        case 'realtime':
          trends.activeUsers = (3 + Math.random() * 5).toFixed(1); // Более волатильный тренд для реального времени
          trends.newUsers = (5 + Math.random() * 7).toFixed(1);
          trends.dau = (2 + Math.random() * 3).toFixed(1); // DAU менее волатилен
          trends.mau = (8 + Math.random() * 4).toFixed(1); // MAU стабильный рост
          break;
          
        case 'hour':
          trends.activeUsers = (2.5 + Math.random() * 3).toFixed(1);
          trends.newUsers = (4 + Math.random() * 6).toFixed(1);
          trends.dau = (2 + Math.random() * 2.5).toFixed(1);
          trends.mau = (8 + Math.random() * 4).toFixed(1); // Неизменный для всех периодов
          break;
          
        case 'day':
          trends.activeUsers = (2 + Math.random() * 2).toFixed(1); // За день тренд более стабилен
          trends.newUsers = (3 + Math.random() * 4).toFixed(1);
          trends.dau = (2 + Math.random() * 2).toFixed(1);
          trends.mau = (8 + Math.random() * 4).toFixed(1);
          break;
          
        case 'week':
          trends.activeUsers = (1.5 + Math.random() * 1.5).toFixed(1); // Еще более стабильный тренд
          trends.newUsers = (2 + Math.random() * 3).toFixed(1);
          trends.dau = (2 + Math.random() * 1.5).toFixed(1);
          trends.mau = (8 + Math.random() * 4).toFixed(1);
          break;
      }
      
      // Обновляем данные графиков
      updateCharts(period);
      
      // Обновляем таблицу удержания пользователей в зависимости от периода
      updateRetentionData(period);
    };
    
    // Функция обновления данных таблицы удержания
    const updateRetentionData = (period: PeriodType) => {
      // Для разных периодов показываем разные данные когорт
      if (period === 'realtime' || period === 'hour') {
        // Для реального времени и часового периода показываем недельные когорты
        retentionData.value = [
          {
            date: '1-7 Апр',
            users: '2,451',
            retention: [100, 64, 42, 35, 28, 22]
          },
          {
            date: '8-14 Апр',
            users: '3,125',
            retention: [100, 68, 45, 38, 30, 24]
          },
          {
            date: '15-21 Апр',
            users: '2,845',
            retention: [100, 65, 43, 36, 29, 23]
          },
          {
            date: '22-28 Апр',
            users: '3,512',
            retention: [100, 70, 48, 40, 32, 26]
          },
          {
            date: '29 Апр - 5 Мая',
            users: '3,842',
            retention: [100, 72, 50, 42, null, null]
          }
        ];
      } else if (period === 'day') {
        // Для дневного периода показываем когорты по дням текущего месяца
        retentionData.value = [
          {
            date: '1-5 Мая',
            users: '4,215',
            retention: [100, 62, 40, 33, 26, 20]
          },
          {
            date: '6-10 Мая',
            users: '3,830',
            retention: [100, 65, 43, 36, 28, null]
          },
          {
            date: '11-15 Мая',
            users: '4,125',
            retention: [100, 67, 44, 37, null, null]
          },
          {
            date: '16-20 Мая',
            users: '3,975',
            retention: [100, 69, 47, null, null, null]
          },
          {
            date: '21-25 Мая',
            users: '4,320',
            retention: [100, 71, null, null, null, null]
          }
        ];
      } else if (period === 'week') {
        // Для недельного периода показываем месячные когорты
        retentionData.value = [
          {
            date: 'Январь',
            users: '18,542',
            retention: [100, 58, 39, 30, 24, 19]
          },
          {
            date: 'Февраль',
            users: '20,125',
            retention: [100, 59, 40, 32, 25, 20]
          },
          {
            date: 'Март',
            users: '23,678',
            retention: [100, 62, 43, 34, 27, 22]
          },
          {
            date: 'Апрель',
            users: '25,945',
            retention: [100, 65, 45, 37, 29, null]
          },
          {
            date: 'Май',
            users: '28,320',
            retention: [100, 69, null, null, null, null]
          }
        ];
      }
    };
    
    // Функция обновления графиков
    const updateCharts = (period: PeriodType) => {
      // Генерируем метки времени в зависимости от периода
      const labels: string[] = [];
      const now = dayjs();
      
      switch (period) {
        case 'realtime':
          // Последние 20 минут с интервалом в минуту
          for (let i = 19; i >= 0; i--) {
            labels.push(now.subtract(i, 'minute').format('HH:mm'));
          }
          break;
        case 'hour':
          // Последний час с интервалом в 5 минут
          for (let i = 11; i >= 0; i--) {
            labels.push(now.subtract(i * 5, 'minute').format('HH:mm'));
          }
          break;
        case 'day':
          // Последние 24 часа с интервалом в час
          for (let i = 23; i >= 0; i--) {
            labels.push(now.subtract(i, 'hour').format('HH:mm'));
          }
          break;
        case 'week':
          // Последние 7 дней
          for (let i = 6; i >= 0; i--) {
            labels.push(now.subtract(i, 'day').format('DD.MM'));
          }
          break;
      }
      
      // Обновляем метки графика активности пользователей
      chartData.userActivity.labels = labels;
      
      // Генерируем данные для графика активности
      const userActivityData: number[] = [];
      
      // Используем паттерны активности пользователей, которые наиболее приближены к реальности
      switch (period) {
        case 'realtime':
          // Для реального времени: тенденция с небольшими колебаниями
          for (let i = 0; i < labels.length; i++) {
            // Базовое значение равно текущему кол-ву активных пользователей
            const baseValue = kpiData.activeUsers * 0.95; // Немного меньше текущего значения
            
            // Небольшие колебания
            const noise = Math.random() * 50 - 25;
            
            // Небольшой тренд (увеличение к концу периода)
            const trend = (i / labels.length) * 70;
            
            userActivityData.push(Math.max(0, Math.round(baseValue + noise + trend)));
          }
          break;
          
        case 'hour':
          // Для часового периода: более заметные колебания
          for (let i = 0; i < labels.length; i++) {
            // Базовое значение со значительными колебаниями в течение часа
            const baseValue = kpiData.activeUsers * 0.9;
            
            // Более выраженные колебания в течение часа
            const variation = 100 * Math.sin((i / labels.length) * Math.PI * 1.5);
            
            // Случайный шум
            const noise = Math.random() * 30 - 15;
            
            // Тренд (увеличение к концу периода)
            const trend = (i / labels.length) * 80;
            
            userActivityData.push(Math.max(0, Math.round(baseValue + variation + noise + trend)));
          }
          break;
          
        case 'day':
          // Для дневного периода: суточные колебания активности
          // Учитываем время суток (наибольшая активность днем, спад ночью)
          for (let i = 0; i < labels.length; i++) {
            // Часовой паттерн с пиком активности днем
            const hourOfDay = (24 - i) % 24; // Индекс часа от 0 до 23
            
            // Активность пользователей по часам суток (ночью низкая, пик днем)
            // Коэффициент от 0.3 до 1.0
            let hourFactor = 0.3; // Минимальная активность ночью
            
            if (hourOfDay >= 8 && hourOfDay <= 23) {
              // Активность растет с 8 утра, пик в 14-15, затем спад
              hourFactor = 0.5 + 0.5 * Math.sin(((hourOfDay - 8) / 16) * Math.PI);
            }
            
            const baseValue = kpiData.activeUsers;
            const adjustedValue = baseValue * hourFactor;
            
            // Случайный шум
            const noise = Math.random() * 40 - 20;
            
            userActivityData.push(Math.max(0, Math.round(adjustedValue + noise)));
          }
          break;
          
        case 'week':
          // Для недельного периода: активность по дням недели
          // Обычно рабочие дни имеют более высокую активность, чем выходные
          for (let i = 0; i < labels.length; i++) {
            // Определяем день недели (0-воскресенье, 1-понедельник, ..., 6-суббота)
            const dayOfWeek = (now.subtract(6 - i, 'day').day() + 7) % 7;
            
            // Коэффициент для дня недели (рабочие дни выше, выходные ниже)
            let dayFactor = 1.0; // По умолчанию
            
            if (dayOfWeek === 0 || dayOfWeek === 6) {
              // Выходные дни (снижение активности)
              dayFactor = 0.7;
            } else if (dayOfWeek === 1 || dayOfWeek === 5) {
              // Понедельник и пятница (немного ниже середины недели)
              dayFactor = 0.9;
            } else {
              // Вторник-четверг (пиковая активность)
              dayFactor = 1.0;
            }
            
            const baseValue = kpiData.activeUsers * 1.1; // Немного выше, т.к. это пиковые значения
            const adjustedValue = baseValue * dayFactor;
            
            // Случайный шум
            const noise = Math.random() * 30 - 15;
            
            userActivityData.push(Math.max(0, Math.round(adjustedValue + noise)));
          }
          break;
      }
      
      chartData.userActivity.datasets[0].data = userActivityData;
      
      // Обновляем данные для географического распределения
      const regions = ['Москва', 'Санкт-Петербург', 'Новосибирск', 'Екатеринбург', 'Казань', 'Нижний Новгород'];
      const geoData: number[] = [];
      
      // Постоянные весовые коэффициенты для регионов
      const regionWeights = [0.35, 0.25, 0.15, 0.10, 0.08, 0.07];
      
      // Базовое значение зависит от общего количества активных пользователей
      const geoBaseValue = kpiData.activeUsers;
      
      // Распределяем пользователей по регионам согласно весам
      for (let i = 0; i < regions.length; i++) {
        const regionValue = Math.round(geoBaseValue * regionWeights[i] * (0.95 + Math.random() * 0.1));
        geoData.push(regionValue);
      }
      
      chartData.geoDistribution.labels = regions;
      chartData.geoDistribution.datasets[0].data = geoData;
      
      // Обновляем данные для демографических показателей
      // В демографии не должно быть существенных изменений при смене временных периодов
      // Небольшие колебания допустимы для реализма
      const demoBases = [15, 32, 28, 17, 8]; // Базовое распределение
      
      // Добавляем очень небольшие вариации для реализма (не более ±1%)
      for (let i = 0; i < demoBases.length; i++) {
        const variation = Math.random() * 2 - 1;
        chartData.demographics.datasets[0].data[i] = Math.max(5, Math.min(40, Math.round(demoBases[i] + variation)));
      }
    };
    
    // Подписываемся на изменения периода и обновляем данные
    const unsubscribe = onPeriodChange((period) => {
      updateDataByPeriod(period);
    });
    
    // Инициализируем данные при монтировании компонента
    onMounted(() => {
      updateDataByPeriod(activePeriod.value);
    });
    
    // Отписываемся при размонтировании компонента
    onUnmounted(() => {
      unsubscribe();
    });
    
    return {
      metrics,
      chartData,
      chartOptions,
      kpiData,
      trends,
      retentionData,
      getRetentionClass,
      formatNumber
    };
  }
});
</script>

<style lang="scss">
.users-view {
  width: 100%;
}

.kpi-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.chart-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
  
  .wide {
    grid-column: span 2;
  }
}

.card {
  background-color: var(--bg-card);
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--border);
    
    h3 {
      font-size: 1rem;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .trend {
      font-size: 0.75rem;
      font-weight: 600;
      padding: 0.125rem 0.375rem;
      border-radius: 9999px;
      
      &.up {
        color: var(--success);
        background-color: rgba(16, 185, 129, 0.1);
      }
      
      &.down {
        color: var(--danger);
        background-color: rgba(239, 68, 68, 0.1);
      }
    }
    
    .header-actions {
      .action-button {
        background: none;
        border: none;
        color: var(--text-secondary);
        cursor: pointer;
        font-size: 1.25rem;
        display: flex;
        align-items: center;
        justify-content: center;
        
        &:hover {
          color: var(--text-primary);
        }
      }
    }
  }
  
  .card-content {
    padding: 1.25rem;
  }
  
  &.kpi-card {
  .card-value {
    font-size: 2rem;
    font-weight: 700;
      margin: 0.75rem 0;
  }
  
  .card-footer {
    .card-badge {
      display: inline-flex;
      align-items: center;
      font-size: 0.75rem;
      color: var(--text-secondary);
      
      .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
          margin-right: 0.375rem;
        
        &.green {
          background-color: var(--success);
        }
        
          &.blue {
            background-color: var(--accent);
        }
        
        &.purple {
          background-color: #8b5cf6;
        }
        
          &.yellow {
            background-color: var(--warning);
          }
        }
      }
    }
  }
  
  &.chart-card {
  .chart-container {
      height: 300px;
    }
  }
}

.retention-table {
  width: 100%;
  overflow-x: auto;
  
  table {
    width: 100%;
    border-collapse: collapse;
    
    th, td {
      padding: 0.75rem 1rem;
      text-align: center;
      border-bottom: 1px solid var(--border);
    }
    
    th {
      color: var(--text-secondary);
      font-weight: 500;
      font-size: 0.9rem;
    }
    
    td {
      font-size: 0.9rem;
      
      &:first-child {
        text-align: left;
        font-weight: 500;
      }
      
      &.high {
        color: var(--success);
        background-color: rgba(16, 185, 129, 0.08);
      }
      
      &.medium {
        color: #10b981;
        background-color: rgba(16, 185, 129, 0.05);
      }
      
      &.low {
        color: var(--warning);
        background-color: rgba(245, 158, 11, 0.05);
      }
      
      &.very-low {
        color: var(--danger);
        background-color: rgba(239, 68, 68, 0.05);
      }
      
      &.unavailable {
        color: var(--text-secondary);
      }
    }
  }
}
</style> 