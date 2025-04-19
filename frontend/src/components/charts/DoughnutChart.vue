<template>
  <div class="chart-wrapper">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch, PropType } from 'vue';
import Chart from 'chart.js/auto';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';

ChartJS.register(
  ArcElement,
  Tooltip,
  Legend
);

interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor?: string | string[];
  borderColor?: string | string[];
  borderWidth?: number;
}

interface ChartDataInterface {
  labels: string[];
  datasets: ChartDataset[];
}

export default defineComponent({
  name: 'DoughnutChart',
  props: {
    chartData: {
      type: Object as PropType<ChartDataInterface>,
      required: true
    },
    options: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const chartCanvas = ref<HTMLCanvasElement | null>(null);
    let chart: Chart | null = null;

    const createChart = () => {
      if (!chartCanvas.value) return;

      // Уничтожаем предыдущий график, если есть
      if (chart) {
        chart.destroy();
      }

      // Создаем новый график
      chart = new Chart(chartCanvas.value, {
        type: 'doughnut',
        data: props.chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          cutout: '65%',
          animation: {
            animateRotate: true,
            animateScale: true,
            duration: 1000,
            easing: 'easeOutCubic'
          },
          plugins: {
            legend: {
              position: 'right',
              labels: {
                boxWidth: 12,
                padding: 15
              }
            },
            tooltip: {
              backgroundColor: 'rgba(30, 41, 59, 0.9)',
              titleColor: '#f1f5f9',
              bodyColor: '#f1f5f9',
              borderColor: '#334155',
              borderWidth: 1,
              padding: 10,
              cornerRadius: 6,
              boxPadding: 6
            }
          },
          ...props.options
        }
      });
    };

    onMounted(() => {
      createChart();
    });

    // Обновляем график при изменении данных
    watch(
      () => props.chartData,
      () => {
        if (chart) {
          chart.data = props.chartData;
          chart.update();
        } else {
          createChart();
        }
      },
      { deep: true }
    );

    // Обновляем параметры при их изменении
    watch(
      () => props.options,
      () => {
        if (chart) {
          chart.options = {
            responsive: true,
            maintainAspectRatio: false,
            cutout: '65%',
            animation: {
              animateRotate: true,
              animateScale: true,
              duration: 1000,
              easing: 'easeOutCubic'
            },
            plugins: {
              legend: {
                position: 'right',
                labels: {
                  boxWidth: 12,
                  padding: 15
                }
              },
              tooltip: {
                backgroundColor: 'rgba(30, 41, 59, 0.9)',
                titleColor: '#f1f5f9',
                bodyColor: '#f1f5f9',
                borderColor: '#334155',
                borderWidth: 1,
                padding: 10,
                cornerRadius: 6,
                boxPadding: 6
              }
            },
            ...props.options
          };
          chart.update();
        }
      },
      { deep: true }
    );

    return {
      chartCanvas
    };
  }
});
</script>

<style scoped>
.chart-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}
</style> 