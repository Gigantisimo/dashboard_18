<template>
  <div class="chart-wrapper">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch, PropType } from 'vue';
import Chart from 'chart.js/auto';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
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
  name: 'BarChart',
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
        type: 'bar',
        data: props.chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          animation: {
            duration: 800,
            easing: 'easeOutQuart'
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
            animation: {
              duration: 800,
              easing: 'easeOutQuart'
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
