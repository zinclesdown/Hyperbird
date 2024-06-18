<template>
  <q-page class="col items-center justify-evenly p-2">
    <div>Hello world from IndexPage!</div>
    <div>已存储URL对象： {{ store }}</div>
    <div class="p-2">Hello world</div>
  </q-page>
</template>

<script setup lang="ts">
import {} from 'pinia';
import { apiUrlStorage } from './../stores/api-urls';

const store = apiUrlStorage();

console.log(store);
</script>
