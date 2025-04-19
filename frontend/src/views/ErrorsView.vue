<template>
  <div class="errors-view">
    <!-- Сводка ошибок -->
    <div class="kpi-cards">
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Всего ошибок</h3>
            <div class="trend down">-12.3%</div>
          </div>
          <div class="card-value">{{ metrics.totalErrors }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot green"></span>
              За 24 часа
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Критические</h3>
            <div class="trend down">-5.2%</div>
          </div>
          <div class="card-value">{{ metrics.criticalErrors }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot red"></span>
              Требуют решения
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Среднее время решения</h3>
            <div class="trend up">+8.4%</div>
          </div>
          <div class="card-value">{{ metrics.avgResolutionTime }}</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot yellow"></span>
              Минут
            </div>
          </div>
        </div>
      </div>
      
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Успешность решения</h3>
            <div class="trend up">+3.1%</div>
          </div>
          <div class="card-value">{{ metrics.resolutionRate }}%</div>
          <div class="card-footer">
            <div class="card-badge">
              <span class="dot blue"></span>
              С первой попытки
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- График ошибок -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Динамика ошибок</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <LineChart 
              v-if="chartData.errors.labels.length > 0"
              :chartData="chartData.errors" 
              :options="chartOptions.errors"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Распределение ошибок по типам и службам -->
    <div class="chart-row">
      <div class="card chart-card">
        <div class="card-header">
          <h3>Типы ошибок</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <DoughnutChart 
              v-if="chartData.errorTypes.labels.length > 0"
              :chartData="chartData.errorTypes" 
              :options="chartOptions.pie"
            />
          </div>
        </div>
      </div>
      
      <div class="card chart-card">
        <div class="card-header">
          <h3>Распределение по службам</h3>
          <div class="header-actions">
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="chart-container">
            <BarChart 
              v-if="chartData.errorServices.labels.length > 0"
              :chartData="chartData.errorServices" 
              :options="chartOptions.bar"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Список ошибок -->
    <div class="chart-row">
      <div class="card chart-card wide">
        <div class="card-header">
          <h3>Последние ошибки</h3>
          <div class="header-actions">
            <div class="search-box">
              <input type="text" placeholder="Поиск ошибок..." v-model="searchQuery">
            </div>
            <button class="action-button">
              <span>⋮</span>
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="errors-table">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Тип</th>
                  <th>Сообщение</th>
                  <th>Служба</th>
                  <th>Статус</th>
                  <th>Время</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(error, index) in filteredErrors" :key="index">
                  <td>{{ error.id }}</td>
                  <td>
                    <span class="error-type" :class="error.type.toLowerCase()">
                      {{ error.type }}
                    </span>
                  </td>
                  <td class="error-message">{{ error.message }}</td>
                  <td>{{ error.service }}</td>
                  <td>
                    <span class="status-badge" :class="error.status.toLowerCase()">
                      {{ error.status }}
                    </span>
                  </td>
                  <td>{{ error.time }}</td>
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
import { defineComponent, reactive, ref, computed } from 'vue';
import LineChart from '@/components/charts/LineChart.vue';
import DoughnutChart from '@/components/charts/DoughnutChart.vue';
import BarChart from '@/components/charts/BarChart.vue';

export default defineComponent({
  name: 'ErrorsView',
  components: {
    LineChart,
    DoughnutChart,
    BarChart
  },
  setup() {
    const searchQuery = ref('');
    
    const metrics = reactive({
      totalErrors: 152,
      criticalErrors: 8,
      avgResolutionTime: 42,
      resolutionRate: 94.6
    });
    
    const errors = [
      {
        id: 'ERR-2305',
        type: 'Критическая',
        message: 'Ошибка подключения к базе данных',
        service: 'База данных',
        status: 'Активна',
        time: '10:45'
      },
      {
        id: 'ERR-2304',
        type: 'Высокая',
        message: 'Таймаут запроса к API платежного шлюза',
        service: 'Платежи',
        status: 'Решена',
        time: '10:32'
      },
      {
        id: 'ERR-2303',
        type: 'Средняя',
        message: 'Ошибка валидации данных пользователя',
        service: 'Аутентификация',
        status: 'Решена',
        time: '10:18'
      },
      {
        id: 'ERR-2302',
        type: 'Низкая',
        message: 'Предупреждение о высокой загрузке системы',
        service: 'Мониторинг',
        status: 'Активна',
        time: '09:55'
      },
      {
        id: 'ERR-2301',
        type: 'Критическая',
        message: 'Сбой в системе доставки уведомлений',
        service: 'Нотификации',
        status: 'Решена',
        time: '09:41'
      },
      {
        id: 'ERR-2300',
        type: 'Средняя',
        message: 'Ошибка загрузки кеша конфигурации',
        service: 'Конфигурация',
        status: 'Решена',
        time: '09:22'
      },
      {
        id: 'ERR-2299',
        type: 'Высокая',
        message: 'Сбой в работе поискового модуля',
        service: 'Поиск',
        status: 'Активна',
        time: '09:15'
      },
      {
        id: 'ERR-2298',
        type: 'Низкая',
        message: 'Ошибка отображения UI компонента',
        service: 'Фронтенд',
        status: 'Решена',
        time: '08:57'
      }
    ];
    
    const filteredErrors = computed(() => {
      if (!searchQuery.value) return errors;
      const query = searchQuery.value.toLowerCase();
      return errors.filter(error => 
        error.message.toLowerCase().includes(query) ||
        error.service.toLowerCase().includes(query) ||
        error.type.toLowerCase().includes(query) ||
        error.id.toLowerCase().includes(query)
      );
    });
    
    // Данные для графиков
    const chartData = reactive({
      errors: {
        labels: ['00:00', '03:00', '06:00', '09:00', '12:00', '15:00', '18:00', '21:00'],
        datasets: [
          {
            label: 'Критические',
            data: [5, 8, 3, 12, 8, 5, 2, 4],
            borderColor: '#ef4444',
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderWidth: 2,
            fill: true,
            tension: 0.4
          },
          {
            label: 'Высокие',
            data: [15, 23, 14, 31, 22, 14, 9, 11],
            borderColor: '#f59e0b',
            backgroundColor: 'rgba(245, 158, 11, 0.1)',
            borderWidth: 2,
            fill: true,
            tension: 0.4
          },
          {
            label: 'Средние',
            data: [25, 41, 30, 52, 39, 21, 18, 24],
            borderColor: '#3b82f6',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            borderWidth: 2,
            fill: true,
            tension: 0.4
          },
          {
            label: 'Низкие',
            data: [34, 52, 41, 65, 48, 35, 29, 38],
            borderColor: '#10b981',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            borderWidth: 2,
            fill: true,
            tension: 0.4
          }
        ]
      },
      errorTypes: {
        labels: ['Критические', 'Высокие', 'Средние', 'Низкие'],
        datasets: [
          {
            label: 'Типы ошибок',
            data: [8, 22, 84, 38],
            backgroundColor: [
              'rgba(239, 68, 68, 0.7)',
              'rgba(245, 158, 11, 0.7)',
              'rgba(59, 130, 246, 0.7)',
              'rgba(16, 185, 129, 0.7)'
            ],
            borderWidth: 1
          }
        ]
      },
      errorServices: {
        labels: ['База данных', 'Платежи', 'Аутентификация', 'Мониторинг', 'Нотификации', 'Поиск', 'Фронтенд'],
        datasets: [
          {
            label: 'Количество ошибок',
            data: [28, 22, 19, 31, 25, 14, 13],
            backgroundColor: 'rgba(59, 130, 246, 0.7)',
            borderColor: 'rgba(59, 130, 246, 1)',
            borderWidth: 1
          }
        ]
      }
    });
    
    // Опции для графиков
    const chartOptions = {
      errors: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'top',
            labels: {
              color: '#94a3b8',
              boxWidth: 12,
              padding: 15
            }
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
      pie: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'right',
            labels: {
              color: '#94a3b8',
              boxWidth: 12,
              padding: 10
            }
          }
        }
      },
      bar: {
        indexAxis: 'y',
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
      }
    };
    
    return {
      metrics,
      chartData,
      chartOptions,
      errors,
      searchQuery,
      filteredErrors
    };
  }
});
</script>

<style lang="scss" scoped>
.errors-view {
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
      align-items: center;
      
      .search-box {
        margin-right: 0.5rem;
        
        input {
          background-color: rgba(255, 255, 255, 0.05);
          border: 1px solid var(--border);
          border-radius: 4px;
          padding: 0.5rem 0.75rem;
          font-size: 0.875rem;
          color: var(--text-primary);
          width: 200px;
          
          &:focus {
            outline: none;
            border-color: var(--accent);
          }
          
          &::placeholder {
            color: var(--text-secondary);
          }
        }
      }
      
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
    height: 250px;
    width: 100%;
  }
}

.errors-table {
  width: 100%;
  overflow-x: auto;
  
  table {
    width: 100%;
    border-collapse: collapse;
    
    th, td {
      padding: 0.75rem 1rem;
      text-align: left;
    }
    
    th {
      color: var(--text-secondary);
      font-weight: 500;
      font-size: 0.875rem;
      border-bottom: 1px solid var(--border);
    }
    
    td {
      font-size: 0.875rem;
    }
    
    .error-message {
      max-width: 300px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    tbody tr {
      border-bottom: 1px solid rgba(71, 85, 105, 0.1);
      
      &:last-child {
        border-bottom: none;
      }
    }
    
    .error-type {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      
      &.критическая {
        color: var(--danger);
        background-color: rgba(239, 68, 68, 0.1);
      }
      
      &.высокая {
        color: var(--warning);
        background-color: rgba(245, 158, 11, 0.1);
      }
      
      &.средняя {
        color: var(--accent);
        background-color: rgba(59, 130, 246, 0.1);
      }
      
      &.низкая {
        color: var(--success);
        background-color: rgba(16, 185, 129, 0.1);
      }
    }
    
    .status-badge {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      
      &.активна {
        color: var(--danger);
        background-color: rgba(239, 68, 68, 0.1);
      }
      
      &.решена {
        color: var(--success);
        background-color: rgba(16, 185, 129, 0.1);
      }
    }
  }
}
</style> 