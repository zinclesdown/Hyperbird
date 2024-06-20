<template>
  <q-page class="m-4 p-4">
    <h1>Book Library Manage</h1>
    <p>TODO: Book Library Manage</p>
    API路径: {{ GET_BOOKS_SHORT_INFO }}
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';

console.log('Hello from BookLibraryManage.vue!');

// 遍历库里的所有书籍
import { apiUrlStorage } from './../stores/api-urls';
const urlStore = apiUrlStorage();
const GET_BOOKS_SHORT_INFO = urlStore.bookLibraryGetBooksShortInfo;

const bookNum = ref<number>(0);

console.log(bookNum);

// 刷新书籍数据
// async function refresh_books(page: number, page_size: number) {}
</script>
