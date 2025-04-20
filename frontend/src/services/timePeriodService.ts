import { ref, watch } from 'vue';

// Типы периодов
export type PeriodType = 'realtime' | 'hour' | 'day' | 'week';

// Интерфейс для периодов
export interface TimePeriod {
  id: PeriodType;
  label: string;
}

// Доступные периоды времени
export const timePeriods: TimePeriod[] = [
  { id: 'realtime', label: 'Реальное время' },
  { id: 'hour', label: 'Час' },
  { id: 'day', label: 'День' },
  { id: 'week', label: 'Неделя' }
];

// Текущий выбранный период (по умолчанию - реальное время)
const activePeriod = ref<PeriodType>('realtime');

// Массив колбэков для оповещения о изменениях
const changeCallbacks: Function[] = [];

// Функция для смены активного периода
export const setActivePeriod = (period: PeriodType): void => {
  activePeriod.value = period;
  
  // Уведомляем всех подписчиков об изменении
  changeCallbacks.forEach(callback => callback(period));
};

// Подписка на изменения периода
export const onPeriodChange = (callback: (period: PeriodType) => void): () => void => {
  changeCallbacks.push(callback);
  
  // Сразу вызываем колбэк с текущим значением
  callback(activePeriod.value);
  
  // Возвращаем функцию отписки
  return () => {
    const index = changeCallbacks.indexOf(callback);
    if (index !== -1) {
      changeCallbacks.splice(index, 1);
    }
  };
};

// Получение текущего активного периода
export const getActivePeriod = (): PeriodType => {
  return activePeriod.value;
};

// Хук для компонентов
export const useTimePeriod = () => {
  return {
    activePeriod,
    timePeriods,
    setActivePeriod,
    onPeriodChange
  };
}; 