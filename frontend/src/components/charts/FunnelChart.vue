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
import { ConversionFunnel } from '../../services/metricsService';

const props = defineProps<{
  data: ConversionFunnel | undefined;
  title: string;
  isLoading?: boolean;
}>();

const chartContainer = ref<HTMLElement | null>(null);
const hasData = ref(false);
const isLoading = ref(props.isLoading || false);

const funnelStages = [
  { key: 'visitors', label: 'Посетители' },
  { key: 'productViews', label: 'Просмотр товаров' },
  { key: 'addedToCart', label: 'Добавление в корзину' },
  { key: 'beganCheckout', label: 'Начало оформления' },
  { key: 'purchasedItems', label: 'Покупка' }
];

// Функция для создания графика воронки
const drawChart = () => {
  if (!chartContainer.value || !props.data) return;
  
  // Проверяем наличие данных
  hasData.value = !!props.data;
  if (!hasData.value) return;
  
  // Очищаем предыдущий график
  d3.select(chartContainer.value).selectAll('*').remove();
  
  const container = chartContainer.value;
  const width = container.clientWidth;
  const height = container.clientHeight;
  
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
  
  // Создаем группу для воронки
  const g = svg.append('g')
    .attr('transform', `translate(${width / 2}, 40)`);
  
  // Подготавливаем данные
  const stages = funnelStages.map(stage => ({
    ...stage,
    value: props.data![stage.key as keyof ConversionFunnel] as number
  }));
  
  // Находим максимальное значение для масштабирования
  const maxValue = stages[0].value;
  
  // Определяем размеры и отступы
  const funnelHeight = height - 60;
  const stageHeight = funnelHeight / stages.length;
  const maxTrapezoidWidth = width * 0.8;
  
  // Создаем шкалу для ширины трапеций
  const widthScale = d3.scaleLinear()
    .domain([0, maxValue])
    .range([0, maxTrapezoidWidth]);
  
  // Определяем цветовую схему
  const colorScale = d3.scaleSequential()
    .domain([0, stages.length - 1])
    .interpolator(d3.interpolateBlues);
  
  // Рисуем каждую стадию воронки
  stages.forEach((stage, i) => {
    const y = i * stageHeight;
    const currentWidth = widthScale(stage.value);
    const nextWidth = i < stages.length - 1 ? widthScale(stages[i + 1].value) : 0;
    
    // Рисуем трапецию
    const trapezoid = [
      [-currentWidth / 2, y],
      [currentWidth / 2, y],
      [nextWidth / 2, y + stageHeight - 5],
      [-nextWidth / 2, y + stageHeight - 5]
    ];
    
    // Создаем строку path для трапеции
    const path = d3.line()(trapezoid);
    
    // Рисуем фигуру
    g.append('path')
      .attr('d', path)
      .attr('fill', colorScale(i))
      .attr('stroke', '#fff')
      .attr('stroke-width', 1);
    
    // Добавляем название стадии
    g.append('text')
      .attr('x', 0)
      .attr('y', y + stageHeight / 2)
      .attr('text-anchor', 'middle')
      .attr('fill', '#fff')
      .style('font-size', '12px')
      .style('font-weight', 'bold')
      .text(stage.label);
    
    // Добавляем значение
    g.append('text')
      .attr('x', 0)
      .attr('y', y + stageHeight / 2 + 15)
      .attr('text-anchor', 'middle')
      .attr('fill', '#fff')
      .style('font-size', '12px')
      .text(formatValue(stage.value));
    
    // Добавляем процент конверсии от предыдущей стадии
    if (i > 0) {
      const prevValue = stages[i - 1].value;
      const conversionRate = prevValue > 0 ? (stage.value / prevValue * 100).toFixed(1) : '0';
      
      g.append('text')
        .attr('x', currentWidth / 2 + 10)
        .attr('y', y + 10)
        .attr('text-anchor', 'start')
        .style('font-size', '10px')
        .text(`${conversionRate}%`);
    }
  });
  
  // Добавляем общую конверсию от первой до последней стадии
  const totalConversion = (stages[stages.length - 1].value / stages[0].value * 100).toFixed(2);
  
  svg.append('text')
    .attr('x', width / 2)
    .attr('y', height - 10)
    .attr('text-anchor', 'middle')
    .style('font-size', '14px')
    .style('font-weight', 'bold')
    .text(`Общая конверсия: ${totalConversion}%`);
};

// Форматирование значений
const formatValue = (value: number) => {
  return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
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