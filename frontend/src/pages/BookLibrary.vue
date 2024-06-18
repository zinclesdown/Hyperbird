<template>
  <div>
    <h3>Book Library</h3>
    <p>图书馆系统</p>

    <!-- 显示图书的组件。调用 get_books_shortinfo(page:int, page_size:int)来获取书籍ID信息 -->
    <!-- 该方法返回 数组[BookShortInfo{book_id:string, book_name:string}] -->
    <div id="books_displayer">
      {{ curpageBooksInfo }}

      <div v-for="bookInfo in curpageBooksInfo" :key="bookInfo.book_id">
        <div>书籍信息</div>
        <div>{{ bookInfo.book_id }}</div>
        <div>{{ bookInfo.book_name }}</div>
      </div>
    </div>

    <!-- 以卡片格式显示书籍 -->
    <div class="q-pa-md">
      <div class="q-gutter-md row items-start">
        <q-card v-for="book in curpageBooksInfo" :key="book.book_id" flat bordered class="my-card">
          <q-card-section>
            <div class="text-h6">{{ book.book_id }}</div>
            <div class="text-h6">{{ book.book_name }}</div>
          </q-card-section>

          <q-card-section>
            <q-img :src="book.book_image_path" alt="book image" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- 按钮，手动触发刷新 -->
    <button @click="refreshPage(0, 10)">刷新</button>
  </div>
</template>

<style>
.my-card {
  width: 200px;
  height: 300px;
}
</style>

<script lang="ts" setup>
// 定义接口
interface BookShortInfo {
  book_id: string;
  book_name: string;
}

import {} from 'pinia';
import { apiUrlStorage } from './../stores/api-urls';
import axios from 'axios';
import { Ref, onMounted, ref } from 'vue';

// Pinia, 读取urlStore的API地址
const urlStore = apiUrlStorage();

const GET_BOOK_URL = urlStore.bookLibraryGetBooksShortInfo;

// 书籍结构体信息：
// {book_id:int, book_name:string, book_image_path:string}
// 书籍ID，书籍名称，书籍封面图片路径

// 接口
interface BookShortInfo {
  book_id: string;
  book_name: string;
  book_image_path: string;
}

interface BookShortInfoResponse {
  books: BookShortInfo[];
}

// 获取书籍信息。第一页为page=0，每页显示10本书
async function get_books_shortinfo(page: number = 0, page_size: number = 10): Promise<BookShortInfoResponse> {
  try {
    const response = await axios.get<BookShortInfoResponse>(GET_BOOK_URL, {
      params: {
        page: page,
        page_size: page_size,
      },
    });
    console.log('获取书籍信息：', response.data);
    return response.data;
  } catch (error) {
    console.error('在向服务器请求书籍信息时发生错误：', error);
    return { books: [] };
  }
}

// 记录书籍信息
const curpageBooksInfo: Ref<BookShortInfo[]> = ref([]);

async function refreshPage(page: number, page_size: number) {
  const data = await get_books_shortinfo(page, page_size);
  curpageBooksInfo.value = data.books; // 使用新的接口

  // 更新vue内的显示
  console.log(curpageBooksInfo.value);

  // 更新显示
}

onMounted(() => {
  refreshPage(0, 10);
  console.log('mounted');
});
</script>
