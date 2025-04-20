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

const props = defineProps<{
  data: Record<string, number> | undefined;
  title: string;
  isLoading?: boolean;
}>();

const chartContainer = ref<HTMLElement | null>(null);
const hasData = ref(false);
const isLoading = ref(props.isLoading || false);

// Функция для создания графика источников трафика
const drawChart = () => {
  if (!chartContainer.value || !props.data) return;
  
  // Проверяем наличие данных
  hasData.value = Object.keys(props.data).length > 0;
  if (!hasData.value) return;
  
  // Очищаем предыдущий график
  d3.select(chartContainer.value).selectAll('*').remove();
  
  const container = chartContainer.value;
  const width = container.clientWidth;
  const height = container.clientHeight;
  const radius = Math.min(width, height) / 2 - 40;
  
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
  
  // Преобразуем данные в массив
  const data = Object.entries(props.data).map(([name, value]) => ({ name, value }));
  
  // Вычисляем общую сумму для процентов
  const total = d3.sum(data, d => d.value);
  
  // Создаем цветовую шкалу
  const colorScale = d3.scaleOrdinal(d3.schemeCategory10);
  
  // Создаем функцию для генерации секторов круговой диаграммы
  const pie = d3.pie<{name: string, value: number}>()
    .sort(null)
    .value(d => d.value);
  
  // Создаем арки для каждого сектора
  const arc = d3.arc<d3.PieArcDatum<{name: string, value: number}>>()
    .innerRadius(0)
    .outerRadius(radius);
  
  // Создаем арку для текста
  const labelArc = d3.arc<d3.PieArcDatum<{name: string, value: number}>>()
    .innerRadius(radius * 0.6)
    .outerRadius(radius * 0.6);
  
  // Создаем группу для диаграммы
  const g = svg.append('g')
    .attr('transform', `translate(${width / 2}, ${height / 2})`);
  
  // Создаем секторы
  const arcs = g.selectAll('.arc')
    .data(pie(data))
    .enter()
    .append('g')
    .attr('class', 'arc');
  
  // Рисуем секторы
  arcs.append('path')
    .attr('d', arc)
    .attr('fill', (d, i) => colorScale(i.toString()))
    .attr('stroke', '#fff')
    .style('stroke-width', '2px');
  
  // Добавляем метки
  arcs.append('text')
    .attr('transform', d => `translate(${labelArc.centroid(d)})`)
    .attr('dy', '.35em')
    .style('text-anchor', 'middle')
    .style('font-size', '12px')
    .style('fill', '#fff')
    .style('font-weight', 'bold')
    .text(d => {
      const percent = (d.data.value / total * 100).toFixed(1);
      return percent + '%';
    });
  
  // Создаем легенду
  const legend = svg.append('g')
    .attr('transform', `translate(${width - 150}, 50)`);
  
  data.forEach((d, i) => {
    const legendRow = legend.append('g')
      .attr('transform', `translate(0, ${i * 20})`);
    
    legendRow.append('rect')
      .attr('width', 10)
      .attr('height', 10)
      .attr('fill', colorScale(i.toString()));
    
    legendRow.append('text')
      .attr('x', 15)
      .attr('y', 10)
      .attr('text-anchor', 'start')
      .style('font-size', '12px')
      .text(d.name);
  });
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