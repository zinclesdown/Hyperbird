<template>
  <q-page class="m-4 p-4">
    <h2>图书管理面板</h2>
    <p>
      <q-btn label="刷新数据" color="primary" @click="refresh_books(1, 10)" />

      <!-- 显示书籍数据的表格 -->
      <q-table
        title="图书馆书籍数据"
        :rows="books"
        :columns="columns"
        row-key="name"
        virtual-scroll
        :rows-per-page-options="[0]"
        separator="vertical"
      />
    </p>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import * as APIMethods from './../api-methods';
import { BookInfo } from '../api-methods';

const bookNum = ref<number>(0);

import { QTableColumn } from 'quasar';

// 定义列的属性
// https://quasar.dev/vue-components/table
const columns: QTableColumn[] = [
  { name: 'book_id', label: 'ID', field: (book: BookInfo) => book.book_id, required: true, align: 'left' },
  { name: 'book_name', label: '名称', field: (book: BookInfo) => book.book_name, align: 'left' },
  { name: 'author', label: '作者', field: (book: BookInfo) => book.author, align: 'left' },
  { name: 'description', label: '描述', field: (book: BookInfo) => limitString(book.description), align: 'left' },
];
const books = ref<BookInfo[]>([]);

// const rows = books.value;

console.log(bookNum, columns, APIMethods);

// 刷新书籍数据
async function refresh_books(page: number, page_size: number) {
  console.log('refresh_books', page, page_size);
  // alert('刷新数据！');
  // const res = await get_books(page, page_size);
  // console.log(res);
  // bookNum.value = res.data.count;
  // console.log(bookNum);

  console.log(APIMethods.GetBookInfos(0, 10));
  books.value = (await APIMethods.GetBookInfos(0, 10)).books;
}

// 限制字符串长度
function limitString(str: string | undefined, maxLength: number = 40): string {
  if (!str) return '';
  return str.length > maxLength ? str.substring(0, maxLength) + ' ...... ' : str;
}

refresh_books(0, 100);
</script>
