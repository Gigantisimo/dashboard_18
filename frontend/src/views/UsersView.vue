<template>
  <div class="users-view">
    <!-- KPI карточки -->
    <div class="kpi-cards">
      <div class="card kpi-card">
        <div class="card-content">
          <div class="card-header">
            <h3>Активные сейчас</h3>
            <div class="trend up">+5.2%</div>
          </div>
          <div class="card-value">842</div>
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
            <div class="trend up">+8.1%</div>
          </div>
          <div class="card-value">124</div>
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
            <div class="trend up">+3.5%</div>
          </div>
          <div class="card-value">3,842</div>
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
            <div class="trend up">+12.4%</div>
          </div>
          <div class="card-value">28,675</div>
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
import { defineComponent, reactive, onMounted } from 'vue';
import LineChart from '@/components/charts/LineChart.vue';
import DoughnutChart from '@/components/charts/DoughnutChart.vue';
import BarChart from '@/components/charts/BarChart.vue';
import { useMetrics, generateChartData } from '@/services/metricsService';

export default defineComponent({
  name: 'UsersView',
  components: {
    LineChart,
    DoughnutChart,
    BarChart
  },
  setup() {
    // Используем сервис метрик
    const { metrics } = useMetrics();
    
    // Данные для демографии не включены в общий сервис метрик,
    // поэтому создаем их локально
    const chartData = reactive({
      userActivity: generateChartData().userActivity,
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
      geoDistribution: generateChartData().geoDistribution
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
      demographics: {
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
      geoDistribution: {
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
    
    // Данные для таблицы удержания
    const retentionData = [
      {
        date: '2023-08',
        users: 5423,
        retention: [100, 45, 32, 24, 18]
      },
      {
        date: '2023-09',
        users: 6851,
        retention: [100, 48, 35, 26, 21]
      },
      {
        date: '2023-10',
        users: 7245,
        retention: [100, 52, 38, 29, 23]
      },
      {
        date: '2023-11',
        users: 8932,
        retention: [100, 54, 41, 32, 25]
      },
      {
        date: '2023-12',
        users: 10547,
        retention: [100, 56, 43, 35, 27]
      }
    ];
    
    const getRetentionClass = (rate: number) => {
      if (rate >= 40) return 'high';
      if (rate >= 25) return 'medium';
      return 'low';
    };
    
    return {
      metrics,
      chartData,
      chartOptions,
      retentionData,
      getRetentionClass
    };
  }
});
</script>

<style lang="scss" scoped>
.users-view {
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
        
        &.purple {
          background-color: #8b5cf6;
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

.geo-map {
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &.placeholder {
    .placeholder-text {
      color: var(--text-secondary);
      font-size: 0.875rem;
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
    }
    
    th {
      color: var(--text-secondary);
      font-weight: 500;
      font-size: 0.875rem;
      border-bottom: 1px solid var(--border);
    }
    
    td {
      font-size: 0.875rem;
      
      &.high {
        background-color: rgba(16, 185, 129, 0.1);
        color: var(--success);
        font-weight: 600;
      }
      
      &.medium {
        background-color: rgba(245, 158, 11, 0.1);
        color: var(--warning);
        font-weight: 600;
      }
      
      &.low {
        background-color: rgba(239, 68, 68, 0.1);
        color: var(--danger);
        font-weight: 600;
      }
    }
    
    tbody tr {
      border-bottom: 1px solid rgba(71, 85, 105, 0.1);
      
      &:hover {
        background-color: rgba(255, 255, 255, 0.02);
      }
    }
  }
}
</style> 