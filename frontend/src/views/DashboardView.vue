<template>
  <DashboardComponent />
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import DashboardComponent from '@/components/DashboardComponent.vue';

export default defineComponent({
  name: 'DashboardView',
  components: {
    DashboardComponent
  }
});
</script> 