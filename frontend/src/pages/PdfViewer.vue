<template>
  <h1>PDF Viewer</h1>

  <p>Book ID: {{ curBookId }}</p>
  <p>API Path: {{ store.bookLibraryGetServedBookfileById }}</p>
  <p>PDF URL for plugin: {{ pdfUrl }}</p>
</template>

<script setup lang="ts">
import {} from 'pinia';
import { apiUrlStorage } from './../stores/api-urls';
import { ref } from 'vue';

const store = apiUrlStorage();
let curBookId = ref<string>();

// 从URL参数的book_id中读取书籍ID
const urlParams = new URLSearchParams(window.location.search);
const book_id = urlParams.get('book_id');
if (book_id) {
  curBookId.value = book_id;
  console.log('从URL参数里找到了URL book_id:', book_id);
}

const pdfUrl = ref<string>();

// 格式化该URL, 供q-pdf-viewer插件使用
// 使用URLSearchParams对象

if (book_id != null) {
  let url = new URL(store.bookLibraryGetServedBookfileById);
  url.searchParams.append('book_id', book_id);
  pdfUrl.value = url.toString();
} else {
  console.error('book_id is null!');
}

console.log(pdfUrl.value);
</script>
