<template>
  <div class="settings-view">
    <div class="tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id" 
        class="tab" 
        :class="{ active: activeTab === tab.id }"
        @click="activeTab = tab.id"
      >
        {{ tab.name }}
      </button>
    </div>
    
    <!-- Основные настройки -->
    <div class="settings-panel" v-if="activeTab === 'general'">
      <div class="panel-header">
        <h2>Основные настройки</h2>
        <p>Настройте основные параметры приложения</p>
      </div>
      
      <div class="settings-group">
        <h3>Отображение</h3>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Тема</span>
            <span class="description">Выберите тему оформления</span>
          </div>
          <div class="settings-item-control">
            <select v-model="settings.theme" class="select-input">
              <option value="system">Системная</option>
              <option value="light">Светлая</option>
              <option value="dark">Тёмная</option>
            </select>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Язык</span>
            <span class="description">Язык интерфейса</span>
          </div>
          <div class="settings-item-control">
            <select v-model="settings.language" class="select-input">
              <option value="ru">Русский</option>
              <option value="en">English</option>
              <option value="es">Español</option>
            </select>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Сжатый режим</span>
            <span class="description">Уменьшить размер элементов интерфейса</span>
          </div>
          <div class="settings-item-control">
            <label class="toggle">
              <input type="checkbox" v-model="settings.compactMode">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Уведомления</h3>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Системные уведомления</span>
            <span class="description">Отображать системные уведомления</span>
          </div>
          <div class="settings-item-control">
            <label class="toggle">
              <input type="checkbox" v-model="settings.systemNotifications">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Email-уведомления</span>
            <span class="description">Получать оповещения на почту</span>
          </div>
          <div class="settings-item-control">
            <label class="toggle">
              <input type="checkbox" v-model="settings.emailNotifications">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Интервал проверки</span>
            <span class="description">Как часто проверять наличие новых данных</span>
          </div>
          <div class="settings-item-control">
            <select v-model="settings.updateInterval" class="select-input">
              <option value="1">1 минута</option>
              <option value="5">5 минут</option>
              <option value="15">15 минут</option>
              <option value="30">30 минут</option>
              <option value="60">1 час</option>
            </select>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Настройки профиля -->
    <div class="settings-panel" v-if="activeTab === 'profile'">
      <div class="panel-header">
        <h2>Профиль пользователя</h2>
        <p>Управление профилем и учетными данными</p>
      </div>
      
      <div class="profile-top">
        <div class="avatar-container">
          <div class="avatar">{{ getInitials(user.name) }}</div>
          <button class="avatar-change-btn">Изменить</button>
        </div>
        
        <div class="user-info">
          <h3>{{ user.name }}</h3>
          <p>{{ user.email }}</p>
          <span class="role-badge">{{ user.role }}</span>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Персональные данные</h3>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Имя</span>
          </div>
          <div class="settings-item-control">
            <input type="text" v-model="user.name" class="text-input" />
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Email</span>
          </div>
          <div class="settings-item-control">
            <input type="email" v-model="user.email" class="text-input" />
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Телефон</span>
          </div>
          <div class="settings-item-control">
            <input type="tel" v-model="user.phone" class="text-input" />
          </div>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Безопасность</h3>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Двухфакторная аутентификация</span>
            <span class="description">Повышенная защита аккаунта</span>
          </div>
          <div class="settings-item-control">
            <label class="toggle">
              <input type="checkbox" v-model="user.twoFactorEnabled">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Изменить пароль</span>
          </div>
          <div class="settings-item-control">
            <button class="secondary-button">Изменить</button>
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Сессии</span>
            <span class="description">Управление активными сессиями</span>
          </div>
          <div class="settings-item-control">
            <button class="secondary-button">Управление</button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Настройки API -->
    <div class="settings-panel" v-if="activeTab === 'api'">
      <div class="panel-header">
        <h2>Настройки API</h2>
        <p>Ключи API и интеграции</p>
      </div>
      
      <div class="settings-group">
        <h3>API ключи</h3>
        
        <div class="settings-item api-key-item">
          <div class="settings-item-label">
            <span>Основной ключ API</span>
            <span class="description">Используется для запросов к API</span>
          </div>
          <div class="settings-item-control api-key-control">
            <div class="api-key-input">
              <input 
                :type="showApiKey ? 'text' : 'password'" 
                :value="apiKeys.main" 
                readonly 
                class="text-input api-key"
              />
              <button class="icon-button" @click="showApiKey = !showApiKey">
                <span>👁️</span>
              </button>
              <button class="icon-button">
                <span>📋</span>
              </button>
            </div>
            <button class="secondary-button">Пересоздать</button>
          </div>
        </div>
        
        <div class="settings-item api-key-item">
          <div class="settings-item-label">
            <span>Webhook Secret</span>
            <span class="description">Для верификации webhook запросов</span>
          </div>
          <div class="settings-item-control api-key-control">
            <div class="api-key-input">
              <input 
                :type="showWebhookSecret ? 'text' : 'password'" 
                :value="apiKeys.webhook" 
                readonly 
                class="text-input api-key"
              />
              <button class="icon-button" @click="showWebhookSecret = !showWebhookSecret">
                <span>👁️</span>
              </button>
              <button class="icon-button">
                <span>📋</span>
              </button>
            </div>
            <button class="secondary-button">Пересоздать</button>
          </div>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Webhook URL</h3>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>URL для webhook</span>
            <span class="description">Адрес для получения событий</span>
          </div>
          <div class="settings-item-control">
            <input type="url" v-model="webhookUrl" class="text-input" />
          </div>
        </div>
        
        <div class="settings-item">
          <div class="settings-item-label">
            <span>Тестировать webhook</span>
          </div>
          <div class="settings-item-control">
            <button class="secondary-button">Отправить тестовое событие</button>
          </div>
        </div>
      </div>
      
      <div class="settings-group">
        <h3>Интеграции</h3>
        
        <div class="settings-item" v-for="(integration, index) in integrations" :key="index">
          <div class="settings-item-label">
            <span>{{ integration.name }}</span>
            <span class="description">{{ integration.description }}</span>
          </div>
          <div class="settings-item-control">
            <label class="toggle">
              <input type="checkbox" v-model="integration.enabled">
              <span class="toggle-slider"></span>
            </label>
            <button class="secondary-button ml-2" v-if="integration.enabled">Настроить</button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Настройки команды -->
    <div class="settings-panel" v-if="activeTab === 'team'">
      <div class="panel-header">
        <h2>Команда</h2>
        <p>Управление участниками и разрешениями</p>
      </div>
      
      <div class="panel-actions">
        <div class="search-box">
          <input type="text" placeholder="Поиск сотрудников..." v-model="teamSearch">
        </div>
        <button class="primary-button">
          <span>+</span> Добавить сотрудника
        </button>
      </div>
      
      <div class="team-members">
        <table>
          <thead>
            <tr>
              <th>Сотрудник</th>
              <th>Роль</th>
              <th>Отдел</th>
              <th>Статус</th>
              <th>Последний вход</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(member, index) in filteredTeamMembers" :key="index">
              <td>
                <div class="member-info">
                  <div class="member-avatar">{{ getInitials(member.name) }}</div>
                  <div>
                    <div class="member-name">{{ member.name }}</div>
                    <div class="member-email">{{ member.email }}</div>
                  </div>
                </div>
              </td>
              <td>{{ member.role }}</td>
              <td>{{ member.department }}</td>
              <td>
                <span class="status-badge" :class="member.status.toLowerCase()">
                  {{ member.status }}
                </span>
              </td>
              <td>{{ member.lastLogin }}</td>
              <td>
                <button class="icon-button">
                  <span>⋮</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <div class="settings-group">
        <h3>Отделы</h3>
        
        <div class="departments-list">
          <div 
            v-for="(dept, index) in departments" 
            :key="index"
            class="department-item"
          >
            <div class="department-info">
              <h4>{{ dept.name }}</h4>
              <span class="department-count">{{ dept.memberCount }} сотрудников</span>
            </div>
            <div class="department-actions">
              <button class="secondary-button small">Редактировать</button>
            </div>
          </div>
        </div>
        
        <button class="secondary-button mt-4">
          <span>+</span> Добавить отдел
        </button>
      </div>
    </div>
    
    <div class="settings-footer">
      <button class="primary-button" @click="saveSettings">Сохранить изменения</button>
      <button class="secondary-button" @click="resetSettings">Отменить</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, computed } from 'vue';

export default defineComponent({
  name: 'SettingsView',
  setup() {
    const activeTab = ref('general');
    const showApiKey = ref(false);
    const showWebhookSecret = ref(false);
    const teamSearch = ref('');
    const webhookUrl = ref('https://api.example.com/webhooks/incoming');
    
    const tabs = [
      { id: 'general', name: 'Основные' },
      { id: 'profile', name: 'Профиль' },
      { id: 'api', name: 'API' },
      { id: 'team', name: 'Команда' }
    ];
    
    const settings = reactive({
      theme: 'system',
      language: 'ru',
      compactMode: false,
      systemNotifications: true,
      emailNotifications: true,
      updateInterval: '5'
    });
    
    const user = reactive({
      name: 'Иван Петров',
      email: 'ivan.petrov@example.com',
      phone: '+7 (999) 123-45-67',
      role: 'Администратор',
      twoFactorEnabled: false
    });
    
    const apiKeys = reactive({
      main: 'pk_live_51KjN2pCzV6knDrTM4VRN8GFDf7Xs5n',
      webhook: 'whsec_8fKpq2XLcgDep4wEnRSJNT6W5Nxsb'
    });
    
    const integrations = reactive([
      {
        name: 'Google Analytics',
        description: 'Аналитика и отслеживание',
        enabled: true
      },
      {
        name: 'Telegram',
        description: 'Уведомления в Telegram',
        enabled: false
      },
      {
        name: 'CRM',
        description: 'Интеграция с CRM-системой',
        enabled: true
      },
      {
        name: 'Email-маркетинг',
        description: 'Интеграция с сервисом рассылок',
        enabled: false
      }
    ]);
    
    const teamMembers = [
      {
        name: 'Иван Петров',
        email: 'ivan.petrov@example.com',
        role: 'Администратор',
        department: 'Руководство',
        status: 'Активен',
        lastLogin: 'Сегодня, 10:23'
      },
      {
        name: 'Анна Смирнова',
        email: 'anna.smirnova@example.com',
        role: 'Менеджер',
        department: 'Продажи',
        status: 'Активен',
        lastLogin: 'Сегодня, 09:15'
      },
      {
        name: 'Михаил Иванов',
        email: 'mikhail.ivanov@example.com',
        role: 'Разработчик',
        department: 'Разработка',
        status: 'Активен',
        lastLogin: 'Вчера, 18:42'
      },
      {
        name: 'Елена Козлова',
        email: 'elena.kozlova@example.com',
        role: 'Дизайнер',
        department: 'Дизайн',
        status: 'Отпуск',
        lastLogin: '3 дня назад'
      },
      {
        name: 'Алексей Николаев',
        email: 'alexey.nikolaev@example.com',
        role: 'Аналитик',
        department: 'Аналитика',
        status: 'Активен',
        lastLogin: 'Сегодня, 11:07'
      },
      {
        name: 'Ольга Васильева',
        email: 'olga.vasilyeva@example.com',
        role: 'Маркетолог',
        department: 'Маркетинг',
        status: 'Неактивен',
        lastLogin: '2 недели назад'
      }
    ];
    
    const departments = [
      { name: 'Руководство', memberCount: 2 },
      { name: 'Продажи', memberCount: 8 },
      { name: 'Разработка', memberCount: 12 },
      { name: 'Дизайн', memberCount: 5 },
      { name: 'Маркетинг', memberCount: 6 },
      { name: 'Аналитика', memberCount: 3 }
    ];
    
    const filteredTeamMembers = computed(() => {
      if (!teamSearch.value) return teamMembers;
      const query = teamSearch.value.toLowerCase();
      return teamMembers.filter(member => 
        member.name.toLowerCase().includes(query) ||
        member.email.toLowerCase().includes(query) ||
        member.role.toLowerCase().includes(query) ||
        member.department.toLowerCase().includes(query)
      );
    });
    
    function getInitials(name: string): string {
      return name
        .split(' ')
        .map(part => part.charAt(0))
        .join('')
        .toUpperCase()
        .slice(0, 2);
    }
    
    function saveSettings() {
      // Здесь будет код для сохранения настроек
      console.log('Настройки сохранены');
    }
    
    function resetSettings() {
      // Здесь будет код для сброса изменений
      console.log('Изменения отменены');
    }
    
    return {
      activeTab,
      tabs,
      settings,
      user,
      apiKeys,
      webhookUrl,
      integrations,
      teamMembers,
      departments,
      teamSearch,
      filteredTeamMembers,
      showApiKey,
      showWebhookSecret,
      getInitials,
      saveSettings,
      resetSettings
    };
  }
});
</script>

<style lang="scss" scoped>
.settings-view {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
  margin-bottom: 2rem;
  overflow-x: auto;
  
  .tab {
    padding: 0.75rem 1.25rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-secondary);
    background: none;
    border: none;
    cursor: pointer;
    white-space: nowrap;
    position: relative;
    transition: color 0.2s ease;
    
    &:hover {
      color: var(--text-primary);
    }
    
    &.active {
      color: var(--accent);
      
      &:after {
        content: '';
        position: absolute;
        bottom: -1px;
        left: 0;
        right: 0;
        height: 2px;
        background-color: var(--accent);
      }
    }
  }
}

.panel-header {
  margin-bottom: 2rem;
  
  h2 {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }
  
  p {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
}

.panel-actions {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  
  .search-box {
    input {
      background-color: rgba(255, 255, 255, 0.05);
      border: 1px solid var(--border);
      border-radius: 4px;
      padding: 0.5rem 0.75rem;
      font-size: 0.875rem;
      color: var(--text-primary);
      width: 250px;
      
      &:focus {
        outline: none;
        border-color: var(--accent);
      }
      
      &::placeholder {
        color: var(--text-secondary);
      }
    }
  }
}

.settings-group {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid var(--border);
  
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
  
  h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 1.25rem;
  }
}

.settings-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  
  .settings-item-label {
    display: flex;
    flex-direction: column;
    
    span {
      font-size: 0.875rem;
      font-weight: 500;
    }
    
    .description {
      color: var(--text-secondary);
      font-size: 0.75rem;
      font-weight: 400;
      margin-top: 0.25rem;
    }
  }
  
  .settings-item-control {
    display: flex;
    align-items: center;
  }
}

.api-key-item {
  align-items: flex-start;
  
  .api-key-control {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0.5rem;
  }
  
  .api-key-input {
    display: flex;
    align-items: center;
    
    .api-key {
      width: 260px;
      font-family: monospace;
    }
    
    .icon-button {
      margin-left: 0.5rem;
    }
  }
}

.text-input {
  background-color: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: var(--text-primary);
  
  &:focus {
    outline: none;
    border-color: var(--accent);
  }
}

.select-input {
  background-color: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: var(--text-primary);
  min-width: 150px;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  
  &:focus {
    outline: none;
    border-color: var(--accent);
  }
}

.toggle {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
  
  input {
    opacity: 0;
    width: 0;
    height: 0;
    
    &:checked + .toggle-slider {
      background-color: var(--accent);
    }
    
    &:checked + .toggle-slider:before {
      transform: translateX(18px);
    }
  }
  
  .toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--border);
    transition: .4s;
    border-radius: 34px;
    
    &:before {
      position: absolute;
      content: "";
      height: 16px;
      width: 16px;
      left: 3px;
      bottom: 3px;
      background-color: var(--bg-card);
      transition: .4s;
      border-radius: 50%;
    }
  }
}

.profile-top {
  display: flex;
  margin-bottom: 2rem;
  
  .avatar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-right: 2rem;
    
    .avatar {
      width: 100px;
      height: 100px;
      background-color: var(--accent);
      color: white;
      font-size: 2rem;
      font-weight: 600;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 0.75rem;
    }
    
    .avatar-change-btn {
      background: none;
      border: none;
      color: var(--accent);
      font-size: 0.75rem;
      cursor: pointer;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
  
  .user-info {
    display: flex;
    flex-direction: column;
    justify-content: center;
    
    h3 {
      font-size: 1.5rem;
      font-weight: 600;
      margin-bottom: 0.25rem;
    }
    
    p {
      color: var(--text-secondary);
      margin-bottom: 0.5rem;
    }
    
    .role-badge {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      background-color: rgba(59, 130, 246, 0.1);
      color: var(--accent);
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
      align-self: flex-start;
    }
  }
}

.team-members {
  margin-bottom: 2rem;
  
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
      border-bottom: 1px solid rgba(71, 85, 105, 0.1);
    }
    
    tbody tr:last-child td {
      border-bottom: none;
    }
  }
  
  .member-info {
    display: flex;
    align-items: center;
    
    .member-avatar {
      width: 32px;
      height: 32px;
      background-color: var(--accent);
      color: white;
      font-size: 0.75rem;
      font-weight: 600;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 0.75rem;
    }
    
    .member-name {
      font-weight: 500;
    }
    
    .member-email {
      font-size: 0.75rem;
      color: var(--text-secondary);
    }
  }
  
  .status-badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
    
    &.активен {
      background-color: rgba(16, 185, 129, 0.1);
      color: var(--success);
    }
    
    &.отпуск {
      background-color: rgba(245, 158, 11, 0.1);
      color: var(--warning);
    }
    
    &.неактивен {
      background-color: rgba(107, 114, 128, 0.1);
      color: var(--text-secondary);
    }
  }
}

.departments-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1rem;
  
  .department-item {
    background-color: var(--bg-card);
    border-radius: 0.5rem;
    padding: 1rem;
    border: 1px solid var(--border);
    
    .department-info {
      margin-bottom: 0.75rem;
      
      h4 {
        font-size: 0.875rem;
        font-weight: 600;
        margin-bottom: 0.25rem;
      }
      
      .department-count {
        font-size: 0.75rem;
        color: var(--text-secondary);
      }
    }
  }
}

.primary-button {
  background-color: var(--accent);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  display: inline-flex;
  align-items: center;
  
  &:hover {
    background-color: darken(#3b82f6, 10%);
  }
  
  span {
    margin-right: 0.25rem;
  }
}

.secondary-button {
  background-color: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
    border-color: var(--text-secondary);
  }
  
  &.small {
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
  }
  
  span {
    margin-right: 0.25rem;
  }
}

.icon-button {
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

.settings-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
  margin-top: 2rem;
}

.ml-2 {
  margin-left: 0.5rem;
}

.mt-4 {
  margin-top: 1rem;
}
</style> 