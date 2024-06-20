<template>
  <q-page class="m-4 p-4">
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
    <h4>已有书籍</h4>
    <div class="q-pa-md">
      <!-- <div class="q-gutter-md row items-start"> -->
      <div class="q-pa-md row q-col-gutter-md">
        <router-link
          v-for="book in curpageBooksInfo"
          :key="book.book_id"
          :to="{
            path: '/book_library/info',
            query: { book_id: book.book_id },
          }"
        >
          <!-- 自定义的卡片样式 -->
          <q-card flat bordered class="my-card">
            <q-card-section>
              <div class="rounded-2xl">
                <q-img src="https://cdn.quasar.dev/img/parallax2.jpg" alt="book image" class="absolute-full z-10">
                  <div class="absolute-bottom text-subtitle2 text-center">
                    <div class="text-on-image">ID:{{ book.book_id }}</div>
                    <div class="text-on-image">NAME:{{ book.book_name }}</div>
                  </div>
                </q-img>
              </div>
            </q-card-section>
          </q-card>

          <!--  -->
        </router-link>
      </div>
    </div>

    <!-- 按钮，手动触发刷新 -->
    <button @click="refreshPage(0, 10)">刷新</button>
  </q-page>
</template>

<style scoped>
/* 调整卡片样式 */
.my-card {
  width: 205px;
  height: 288px;
  overflow: hidden;
}

/* 调整图书卡片中的图片 */
.q-card__section {
  height: 288px; /* Adjust this to your needs */
}

/* 调整图片文件的尺寸。quasar和tailwind都没有现成的方法可用，只能写css了。 */
.book-image .q-img {
  object-fit: cover;
  width: 100%;
  height: 100%;
}

.text-on-image {
  z-index: -10;
}
</style>

<script lang="ts" setup>
import { GetBooksShortInfo, BookShortInfo } from '../api-methods';
import { Ref, onMounted, ref } from 'vue';
// 记录书籍信息
const curpageBooksInfo: Ref<BookShortInfo[]> = ref([]);

async function refreshPage(page: number, page_size: number) {
  const data = await GetBooksShortInfo(page, page_size);
  curpageBooksInfo.value = data.books; // 使用新的接口

  // 更新vue内的显示
  console.log(curpageBooksInfo.value);

  // 更新显示
}

onMounted(() => {
  1720;

  refreshPage(0, 10);
  console.log('mounted');
});
</script>
