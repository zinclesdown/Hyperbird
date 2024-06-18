<template>
  <div>DEBUG PAGE</div>

  <div>{{ store }}</div>

  <a href="http://127.0.0.1:8080" target="_blank">API Path</a>
</template>

<script setup lang="ts">
import {} from 'pinia';
import { apiUrlStorage } from './../stores/api-urls';

const store = apiUrlStorage();

console.log(store);
</script>
