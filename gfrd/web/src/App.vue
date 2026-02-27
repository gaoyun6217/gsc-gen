<template>
  <n-config-provider :theme="lightTheme">
    <n-message-provider>
      <router-view />
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { lightTheme } from 'naive-ui'
</script>
