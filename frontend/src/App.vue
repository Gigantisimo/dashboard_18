<template>
  <div class="app-container">
    <div class="sidebar">
      <div class="logo">
        <div class="logo-icon"></div>
        <h1 class="logo-text">MetricsFlow</h1>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/" class="nav-item" active-class="active">
          <span class="nav-icon">📊</span>
          <span class="nav-text">Панель</span>
        </router-link>
        <router-link to="/users" class="nav-item" active-class="active">
          <span class="nav-icon">👥</span>
          <span class="nav-text">Пользователи</span>
        </router-link>
        <router-link to="/sales" class="nav-item" active-class="active">
          <span class="nav-icon">💼</span>
          <span class="nav-text">Продажи</span>
        </router-link>
        <router-link to="/errors" class="nav-item" active-class="active">
          <span class="nav-icon">⚠️</span>
          <span class="nav-text">Ошибки</span>
        </router-link>
        <router-link to="/settings" class="nav-item" active-class="active">
          <span class="nav-icon">⚙️</span>
          <span class="nav-text">Настройки</span>
        </router-link>
      </nav>
    </div>
    
    <div class="main-content">
      <header class="app-header">
        <div class="header-title">
          <h2>{{ pageTitle }}</h2>
          <span class="status online">● Онлайн</span>
        </div>
        <div class="time-controls">
          <button 
            v-for="period in timePeriods" 
            :key="period.id"
            :class="['time-button', { active: activePeriod === period.id }]"
            @click="setActivePeriod(period.id)"
          >
            {{ period.label }}
          </button>
        </div>
      </header>
      
      <div class="dashboard-container">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';
import { useRoute } from 'vue-router';

export default defineComponent({
  name: 'App',
  setup() {
    const route = useRoute();
    const activePeriod = ref('realtime');
    
    const timePeriods = [
      { id: 'realtime', label: 'Реальное время' },
      { id: 'hour', label: 'Час' },
      { id: 'day', label: 'День' },
      { id: 'week', label: 'Неделя' }
    ];
    
    const setActivePeriod = (period: string) => {
      activePeriod.value = period;
    };
    
    const pageTitle = computed(() => {
      switch (route.path) {
        case '/':
          return 'Бизнес-метрики';
        case '/users':
          return 'Пользователи';
        case '/sales':
          return 'Продажи';
        case '/errors':
          return 'Ошибки';
        case '/settings':
          return 'Настройки';
        default:
          return 'Бизнес-метрики';
      }
    });
    
    return {
      activePeriod,
      timePeriods,
      setActivePeriod,
      pageTitle
    };
  }
});
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

:root {
  --bg-dark: #0f172a; 
  --bg-card: #1e293b;
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --accent: #3b82f6;
  --accent-hover: #2563eb;
  --success: #10b981;
  --warning: #f59e0b;
  --danger: #ef4444;
  --border: #334155;
  --sidebar-width: 240px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', sans-serif;
  background-color: var(--bg-dark);
  color: var(--text-primary);
  line-height: 1.5;
}

.app-container {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: var(--sidebar-width);
  background-color: var(--bg-card);
  border-right: 1px solid var(--border);
  padding: 1.5rem 0;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  padding: 0 1.5rem;
  margin-bottom: 2rem;
  
  .logo-icon {
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
    border-radius: 8px;
    margin-right: 0.75rem;
  }
  
  .logo-text {
    font-size: 1.25rem;
    font-weight: 600;
    background: linear-gradient(90deg, #3b82f6, #8b5cf6);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  
  .nav-item {
    display: flex;
    align-items: center;
    padding: 0.75rem 1.5rem;
    color: var(--text-secondary);
    text-decoration: none;
    font-weight: 500;
    transition: all 0.2s ease;
    margin-bottom: 0.25rem;
    border-left: 3px solid transparent;
    
    &:hover {
      color: var(--text-primary);
      background-color: rgba(255, 255, 255, 0.05);
    }
    
    &.active {
      color: var(--accent);
      background-color: rgba(59, 130, 246, 0.08);
      border-left: 3px solid var(--accent);
    }
    
    .nav-icon {
      margin-right: 0.75rem;
      font-size: 1.25rem;
    }
    
    .nav-text {
      font-size: 0.9rem;
    }
  }
}

.main-content {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  
  .header-title {
    display: flex;
    align-items: center;
    
    h2 {
      font-size: 1.5rem;
      font-weight: 600;
      margin-right: 1rem;
    }
    
    .status {
      font-size: 0.8rem;
      padding: 0.25rem 0.5rem;
      border-radius: 9999px;
      
      &.online {
        color: var(--success);
      }
    }
  }
  
  .time-controls {
    display: flex;
    
    .time-button {
      background: none;
      border: 1px solid var(--border);
      color: var(--text-secondary);
      padding: 0.5rem 1rem;
      margin-left: 0.5rem;
      border-radius: 0.25rem;
      cursor: pointer;
      font-size: 0.9rem;
      transition: all 0.2s ease;
      
      &:hover {
        background-color: rgba(255, 255, 255, 0.05);
        color: var(--text-primary);
      }
      
      &.active {
        background-color: var(--accent);
        border-color: var(--accent);
        color: white;
      }
    }
  }
}

.dashboard-container {
  display: grid;
  gap: 1.5rem;
}

// Анимации
@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}
</style> 