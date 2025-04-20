<template>
  <div class="h-full">
    <div v-if="isLoading" class="flex items-center justify-center h-full">
      <div class="text-gray-400">Загрузка данных...</div>
    </div>
    <div v-else-if="!hasData" class="flex items-center justify-center h-full">
      <div class="text-gray-400">Нет данных для отображения</div>
    </div>
    <div v-else ref="chartContainer" class="h-full"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import * as d3 from 'd3';
import { RegionData } from '../../services/metricsService';

const props = defineProps<{
  data: Record<string, RegionData> | undefined;
  metricKey: 'activeUsers' | 'sales' | 'conversionRate';
  title: string;
  isLoading?: boolean;
}>();

const chartContainer = ref<HTMLElement | null>(null);
const hasData = ref(false);
const isLoading = ref(props.isLoading || false);

// Создаем карту России
const drawChart = () => {
  if (!chartContainer.value || !props.data) return;
  
  hasData.value = Object.keys(props.data).length > 0;
  if (!hasData.value) return;
  
  // Очищаем предыдущий график
  d3.select(chartContainer.value).selectAll('*').remove();
  
  const container = chartContainer.value;
  const width = container.clientWidth;
  const height = container.clientHeight;
  
  // Подготавливаем данные для визуализации
  const regions = Object.keys(props.data);
  const metricValues = regions.map(r => props.data![r][props.metricKey]);
  
  // Находим максимальное значение для нормализации
  const maxValue = Math.max(...metricValues);
  
  // Создаем цветовую шкалу
  const colorScale = d3.scaleSequential()
    .domain([0, maxValue])
    .interpolator(d3.interpolateBlues);
  
  // Создаем SVG
  const svg = d3.select(container)
    .append('svg')
    .attr('width', width)
    .attr('height', height)
    .attr('viewBox', `0 0 ${width} ${height}`)
    .attr('preserveAspectRatio', 'xMidYMid meet');
  
  // Добавляем заголовок
  svg.append('text')
    .attr('x', width / 2)
    .attr('y', 20)
    .attr('text-anchor', 'middle')
    .style('font-size', '16px')
    .style('font-weight', 'bold')
    .text(props.title);
  
  // Создаем группу для визуализации в виде горизонтальных баров
  const g = svg.append('g')
    .attr('transform', `translate(10, 40)`);
  
  // Сортируем регионы по значению метрики (по убыванию)
  const sortedRegions = [...regions].sort((a, b) => 
    props.data![b][props.metricKey] - props.data![a][props.metricKey]
  );
  
  // Задаем высоту строки и отступы
  const barHeight = 25;
  const barPadding = 5;
  const maxBarWidth = width - 200; // Оставляем место для названий и значений
  
  // Создаем шкалу для ширины баров
  const xScale = d3.scaleLinear()
    .domain([0, maxValue])
    .range([0, maxBarWidth]);
  
  // Рисуем бары для каждого региона
  sortedRegions.forEach((region, i) => {
    const value = props.data![region][props.metricKey];
    const y = i * (barHeight + barPadding);
    
    // Фон бара
    g.append('rect')
      .attr('x', 150)
      .attr('y', y)
      .attr('width', maxBarWidth)
      .attr('height', barHeight)
      .attr('fill', '#f0f0f0')
      .attr('rx', 3)
      .attr('ry', 3);
    
    // Бар со значением
    g.append('rect')
      .attr('x', 150)
      .attr('y', y)
      .attr('width', xScale(value))
      .attr('height', barHeight)
      .attr('fill', colorScale(value))
      .attr('rx', 3)
      .attr('ry', 3);
    
    // Название региона
    g.append('text')
      .attr('x', 145)
      .attr('y', y + barHeight / 2 + 5)
      .attr('text-anchor', 'end')
      .style('font-size', '12px')
      .text(region);
    
    // Значение метрики
    g.append('text')
      .attr('x', 150 + xScale(value) + 5)
      .attr('y', y + barHeight / 2 + 5)
      .style('font-size', '12px')
      .text(formatValue(value));
  });
  
  // Если регионов слишком много, добавляем скролл
  if (sortedRegions.length * (barHeight + barPadding) > height - 50) {
    container.style.overflowY = 'auto';
    container.style.overflowX = 'hidden';
    svg.attr('height', sortedRegions.length * (barHeight + barPadding) + 50);
  }
};

// Форматирование значений в зависимости от метрики
const formatValue = (value: number) => {
  if (props.metricKey === 'conversionRate') {
    return `${value.toFixed(2)}%`;
  } else if (props.metricKey === 'sales') {
    return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
  } else {
    return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
  }
};

// Обновляем график при изменении данных
watch(() => props.data, () => {
  if (props.data) {
    drawChart();
  }
}, { deep: true });

// Обновляем график при изменении размера окна
const handleResize = () => {
  drawChart();
};

onMounted(() => {
  window.addEventListener('resize', handleResize);
  if (props.data) {
    drawChart();
  }
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script> 