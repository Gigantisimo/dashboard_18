<template>
  <div class="chart-wrapper">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch, PropType } from 'vue';
import Chart from 'chart.js/auto';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler } from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor?: string;
  borderColor?: string;
  borderWidth?: number;
  fill?: boolean;
  tension?: number;
}

interface ChartDataInterface {
  labels: string[];
  datasets: ChartDataset[];
}

export default defineComponent({
  name: 'LineChart',
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
        type: 'line',
        data: props.chartData,
        options: props.options
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
          chart.options = props.options;
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