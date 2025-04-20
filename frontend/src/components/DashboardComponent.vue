<template>
  <div class="dashboard">
    <!-- Индикатор выбранного периода -->
    <div class="period-indicator">
      <span>Данные {{ getPeriodLabel() }}</span>
      <div class="status-indicator" :class="{'updating': isUpdating}"></div>
    </div>
    
    <!-- KPI карточки -->
    <div class="kpi-cards">
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Активные пользователи</h3>
            <div class="trend up">+5.2%</div>
          </div>
          <div class="card-value">{{ formatNumber(metrics.activeUsers) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot green"></span>
              <span v-if="activePeriod === 'realtime'">Онлайн сейчас</span>
              <span v-else-if="activePeriod === 'hour'">За последний час</span>
              <span v-else-if="activePeriod === 'day'">За последний день</span>
              <span v-else-if="activePeriod === 'week'">За последнюю неделю</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Запросов в секунду</h3>
            <div class="trend up">+3.8%</div>
          </div>
          <div class="card-value">{{ formatNumber(metrics.requestsPerSecond, 1) }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot blue"></span>
              <span v-if="activePeriod === 'realtime'">В реальном времени</span>
              <span v-else-if="activePeriod === 'hour'">Среднее за час</span>
              <span v-else-if="activePeriod === 'day'">Среднее за день</span>
              <span v-else-if="activePeriod === 'week'">Среднее за неделю</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Конверсия</h3>
            <div class="trend" :class="metrics.conversionRate > 2.0 ? 'up' : 'down'">
              {{ metrics.conversionRate > 2.0 ? '+' : '-' }}{{ Math.abs((metrics.conversionRate - 2.0) / 2.0 * 100).toFixed(1) }}%
            </div>
          </div>
          <div class="card-value">{{ metrics.conversionRate.toFixed(1) }}%</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot purple"></span>
              <span v-if="activePeriod === 'realtime'">Текущая</span>
              <span v-else-if="activePeriod === 'hour'">Средняя за час</span>
              <span v-else-if="activePeriod === 'day'">Средняя за день</span>
              <span v-else-if="activePeriod === 'week'">Средняя за неделю</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Ошибки</h3>
            <div class="trend" :class="metrics.errorRate < 1.5 ? 'up' : 'down'">
              {{ metrics.errorRate < 1.5 ? '-' : '+' }}{{ Math.abs((metrics.errorRate - 1.5) / 1.5 * 100).toFixed(1) }}%
            </div>
          </div>
          <div class="card-value">{{ metrics.errorRate.toFixed(2) }}%</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot" :class="metrics.errorRate < 1.0 ? 'green' : metrics.errorRate < 2.0 ? 'yellow' : 'red'"></span>
              <span v-if="metrics.errorRate < 1.0">Нормально</span>
              <span v-else-if="metrics.errorRate < 2.0">Внимание</span>
              <span v-else>Критично</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Графики активности пользователей -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Активные пользователи (DAU)</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="userActivityChart">
            <!-- График в компоненте LineChart -->
            <LineChart 
              v-if="chartData.userActivity.labels.length > 0"
              :chartData="chartData.userActivity" 
              :options="chartOptions.userActivity"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Запросы и время отклика -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Запросы в секунду</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="requestsChart">
            <LineChart 
              v-if="chartData.requests.labels.length > 0"
              :chartData="chartData.requests" 
              :options="chartOptions.requests"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Время отклика</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="responseTimeChart">
            <LineChart 
              v-if="chartData.responseTime.labels.length > 0"
              :chartData="chartData.responseTime" 
              :options="chartOptions.responseTime"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Продажи и конверсия -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Продажи</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="salesChart">
            <LineChart 
              v-if="chartData.sales.labels.length > 0"
              :chartData="chartData.sales" 
              :options="chartOptions.sales"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Конверсия</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="conversionChart">
            <LineChart 
              v-if="chartData.conversion.labels.length > 0"
              :chartData="chartData.conversion" 
              :options="chartOptions.conversion"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Региональные данные и источники трафика -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Региональная активность</h3>
          <div class="header-buttons">
            <button 
              v-for="metric in ['activeUsers', 'sales', 'conversionRate']" 
              :key="metric"
              @click="activeRegionalMetric = metric"
              :class="{ active: activeRegionalMetric === metric }"
              class="metric-button"
            >
              {{ getRegionalMetricLabel(metric) }}
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <RegionalChart 
              :data="metrics.regionalData" 
              :metricKey="activeRegionalMetric"
              :title="getRegionalMetricLabel(activeRegionalMetric)"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Источники трафика</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <SourcesChart 
              :data="metrics.sourcesData" 
              title="Распределение трафика по источникам"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Воронка конверсии -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Воронка конверсии</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <FunnelChart 
              :data="metrics.conversionFunnel" 
              title="Путь пользователя от просмотра до покупки"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Ошибки -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Ошибки по типам</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="errorTypesChart">
            <DoughnutChart 
              v-if="chartData.errorTypes.labels.length > 0"
              :chartData="chartData.errorTypes" 
              :options="chartOptions.errorTypes"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Процент ошибок</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container" ref="errorRateChart">
            <LineChart 
              v-if="chartData.errorRate.labels.length > 0"
              :chartData="chartData.errorRate" 
              :options="chartOptions.errorRate"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Системные метрики -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Системные метрики</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="system-metrics">
            <div class="gauge-container">
              <div class="gauge">
                <div class="gauge-inner" :style="`--percentage: ${metrics.serverLoad}%;`">
                  <span class="gauge-value">{{ metrics.serverLoad.toFixed(1) }}%</span>
                </div>
                <div class="gauge-label">Нагрузка сервера</div>
              </div>
            </div>
            
            <div class="gauge-container">
              <div class="gauge">
                <div class="gauge-inner" :style="`--percentage: ${metrics.databaseConnections / 2}%;`">
                  <span class="gauge-value">{{ metrics.databaseConnections }}</span>
                </div>
                <div class="gauge-label">БД соединений</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref, computed, onMounted, onUnmounted, watch } from 'vue';
import LineChart from '@/components/charts/LineChart.vue';
import BarChart from '@/components/charts/BarChart.vue';
import DoughnutChart from '@/components/charts/DoughnutChart.vue';
import { useMetrics, generateChartData, formatNumber, setOnUpdateCallback } from '@/services/metricsService';
import { useTimePeriod, PeriodType, timePeriods } from '@/services/timePeriodService';
import type { Metrics } from '@/services/metricsService';
import RegionalChart from './charts/RegionalChart.vue';
import SourcesChart from './charts/SourcesChart.vue';
import FunnelChart from './charts/FunnelChart.vue';

export default defineComponent({
  name: 'DashboardComponent',
  components: {
    LineChart,
    BarChart,
    DoughnutChart,
    RegionalChart,
    SourcesChart,
    FunnelChart
  },
  setup() {
    const { metrics, setActivePeriod: setMetricsActivePeriod } = useMetrics();
    const { activePeriod, onPeriodChange } = useTimePeriod();
    
    // Состояние обновления данных
    const isUpdating = ref(false);
    
    // Данные для графиков
    const chartData = reactive(generateChartData());
    
    // Переменная для отслеживания запроса обновления
    let updateRequested = false;
    // Интервал обновления в миллисекундах
    let updateInterval = 1000;
    // Таймер для дебаунсинга обновлений
    let updateTimer: number | null = null;
    
    // Активная метрика для региональных данных
    const activeRegionalMetric = ref<'activeUsers' | 'sales' | 'conversionRate'>('activeUsers');
    
    // Функция для получения человеко-читаемого названия периода
    const getPeriodLabel = () => {
      const period = timePeriods.find(p => p.id === activePeriod.value);
      return period ? period.label.toLowerCase() : 'в реальном времени';
    };
    
    // Функция получения метки для региональной метрики
    const getRegionalMetricLabel = (metric: string) => {
      switch (metric) {
        case 'activeUsers': return 'Пользователи';
        case 'sales': return 'Продажи';
        case 'conversionRate': return 'Конверсия';
        default: return '';
      }
    };
    
    // Функция обновления графиков с ограничением частоты
    const updateCharts = () => {
      // Если обновление уже запрошено, просто отмечаем запрос
      if (updateRequested) {
        return;
      }
      
      // Помечаем, что запрос на обновление получен
      updateRequested = true;
      
      // Если таймер уже запущен, не устанавливаем новый
      if (updateTimer !== null) {
        return;
      }
      
      // Показываем анимацию обновления
      isUpdating.value = true;
      
      // Устанавливаем таймер для обновления
      updateTimer = window.setTimeout(() => {
        // Обновляем данные графиков
        Object.assign(chartData, generateChartData());
        
        // Сбрасываем флаги и таймер
        updateRequested = false;
        updateTimer = null;
        
        // Скрываем индикатор обновления после небольшой задержки
        setTimeout(() => {
          isUpdating.value = false;
        }, 300);
        
        // Если были запросы на обновление во время ожидания, запускаем новое обновление
        if (updateRequested) {
          updateCharts();
        }
      }, updateInterval);
    };
    
    // Функция для получения данных в зависимости от выбранного периода
    const updateDataBasedOnPeriod = (period: PeriodType) => {
      // Уведомляем сервис метрик об изменении периода
      setMetricsActivePeriod(period);
      
      console.log(`Дашборд: переключение на период ${period}`);
      
      // Обновляем частоту обновления в зависимости от периода
      switch(period) {
        case 'realtime':
          // Данные в реальном времени - обновляем часто
          setUpdateInterval(1000);
          break;
        case 'hour':
          // Данные за час - обновляем реже
          setUpdateInterval(5000);
          break;
        case 'day':
          // Данные за день - еще реже
          setUpdateInterval(10000);
          break;
        case 'week':
          // Данные за неделю - очень редко
          setUpdateInterval(30000);
          break;
      }
      
      // Запрашиваем обновление графиков
      updateCharts();
      
      // Также обновляем информацию о периоде в KPI карточках
      updateKpiInfo(period);
    };
    
    // Обновление информации в KPI карточках в зависимости от периода
    const updateKpiInfo = (period: PeriodType) => {
      // Здесь можно обновить дополнительную информацию в KPI карточках,
      // если это необходимо в зависимости от выбранного периода
      console.log(`Обновление KPI информации для периода: ${period}`);
    };
    
    // Функция изменения интервала обновления
    const setUpdateInterval = (interval: number) => {
      // Обновляем только если интервал изменился
      if (updateInterval !== interval) {
        // Очищаем существующий таймер
        if (updateTimer !== null) {
          clearTimeout(updateTimer);
          updateTimer = null;
        }
        
        // Устанавливаем новый интервал
        updateInterval = interval;
        
        // Запускаем обновление с новым интервалом
        updateRequested = false;
        updateCharts();
      }
    };
    
    // Подписываемся на изменения периода
    const unsubscribe = onPeriodChange(updateDataBasedOnPeriod);
    
    onMounted(() => {
      // Устанавливаем колбэк для обновления данных графиков
      setOnUpdateCallback(updateCharts);
      // Вызываем начальное обновление
      updateCharts();
    });
    
    // Очищаем таймеры и подписки при размонтировании компонента
    onUnmounted(() => {
      if (updateTimer !== null) {
        clearTimeout(updateTimer);
      }
      // Отписываемся от изменений периода
      unsubscribe();
    });
    
    // Опции для графиков
    const chartOptions = {
      userActivity: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
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
              color: '#94a3b8',
              maxTicksLimit: 5 // Ограничиваем количество делений для лучшей читаемости
            },
            beginAtZero: true
          }
        }
      },
      requests: { 
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            display: false
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
              color: '#94a3b8',
              maxTicksLimit: 5
            },
            beginAtZero: true
          }
        }
      },
      responseTime: { 
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            display: false
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
              color: '#94a3b8',
              maxTicksLimit: 5
            },
            beginAtZero: true
          }
        }
      },
      sales: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            display: false
          }
        },
        scales: {
          x: {
            grid: {
              display: false
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
              color: '#94a3b8',
              maxTicksLimit: 5
            },
            beginAtZero: true
          }
        }
      },
      conversion: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            display: false
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
              color: '#94a3b8',
              maxTicksLimit: 5
            },
            beginAtZero: true
          }
        }
      },
      errorRate: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            display: false
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
              color: '#94a3b8',
              maxTicksLimit: 5
            },
            beginAtZero: true
          }
        }
      },
      errorTypes: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 800,
          easing: 'easeOutQuad'
        },
        plugins: {
          legend: {
            position: 'right',
            labels: {
              color: '#94a3b8',
              boxWidth: 12,
              padding: 15
            }
          }
        }
      }
    };
    
    return {
      metrics,
      chartData,
      chartOptions,
      formatNumber,
      activePeriod,
      getPeriodLabel,
      isUpdating,
      activeRegionalMetric,
      getRegionalMetricLabel
    };
  }
});
</script>

<style lang="scss" scoped>
.dashboard {
  width: 100%;
}

.kpi-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.chart-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.card {
  background-color: var(--bg-card);
  border-radius: 0.75rem;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: transform 0.2s, box-shadow 0.2s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  }
}

.kpi-card {
  padding: 1.5rem;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
    
    h3 {
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--text-secondary);
    }
    
    .trend {
      font-size: 0.75rem;
      font-weight: 600;
      padding: 0.25rem 0.5rem;
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
  }
  
  .card-value {
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 0.75rem;
  }
  
  .card-footer {
    .card-badge {
      display: inline-flex;
      align-items: center;
      font-size: 0.75rem;
      color: var(--text-secondary);
      
      .dot {
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        margin-right: 0.5rem;
        
        &.green {
          background-color: var(--success);
        }
        
        &.yellow {
          background-color: var(--warning);
        }
        
        &.red {
          background-color: var(--danger);
        }
        
        &.blue {
          background-color: var(--accent);
        }
        
        &.purple {
          background-color: #8b5cf6;
        }
      }
    }
  }
}

.chart-card {
  display: flex;
  flex-direction: column;
  
  &.wide {
    grid-column: span 2;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border);
    
    h3 {
      font-size: 1rem;
      font-weight: 600;
    }
    
    .header-actions {
      display: flex;
      
      .action-button {
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
  
  .card-content {
    padding: 1.5rem;
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  
  .chart-container {
    height: 200px;
    width: 100%;
  }
}

.system-metrics {
  display: flex;
  justify-content: space-around;
  width: 100%;
  
  .gauge-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    .gauge {
      width: 160px;
      height: 160px;
      position: relative;
      padding: 10px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      
      .gauge-inner {
        width: 140px;
        height: 140px;
        border-radius: 50%;
        background: conic-gradient(
          var(--accent) 0% calc(var(--percentage) * 1%),
          rgba(255, 255, 255, 0.1) calc(var(--percentage) * 1%) 100%
        );
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
        
        &::before {
          content: "";
          position: absolute;
          width: 110px;
          height: 110px;
          border-radius: 50%;
          background-color: var(--bg-card);
        }
        
        .gauge-value {
          position: relative;
          z-index: 2;
          font-size: 1.5rem;
          font-weight: 700;
        }
      }
      
      .gauge-label {
        margin-top: 0.75rem;
        font-size: 0.875rem;
        color: var(--text-secondary);
      }
    }
  }
}

.period-indicator {
  padding: 0.5rem 1rem;
  background-color: rgba(59, 130, 246, 0.1);
  border-radius: 0.25rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
  display: inline-flex;
  align-items: center;
  
  span {
    font-weight: 500;
    margin-right: 0.5rem;
  }
  
  .status-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--success);
    opacity: 0.5;
    transition: opacity 0.2s ease-in-out;
    
    &.updating {
      opacity: 1;
      animation: pulse 0.8s infinite;
    }
  }
}

@keyframes pulse {
  0% {
    opacity: 0.2;
    transform: scale(0.8);
  }
  50% {
    opacity: 1;
    transform: scale(1.2);
  }
  100% {
    opacity: 0.2;
    transform: scale(0.8);
  }
}

.metric-button {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  background-color: #f3f4f6;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  margin-left: 5px;
}

.metric-button.active {
  background-color: #3b82f6;
  color: white;
  border-color: #2563eb;
}

.header-buttons {
  display: flex;
  gap: 5px;
}
</style> 