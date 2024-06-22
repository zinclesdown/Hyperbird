<template>
  <q-card class="m-8">
    <div class="row">
      <q-card-section class="flex items-center m-4">
        <!-- 图像的显示区域，位于左上侧。 -->
        <div class="book-image m-2 rounded-2xl">
          <q-img src="https://cdn.quasar.dev/img/parallax2.jpg" alt="book image" class="w-full h-full object-fill" />
        </div>
      </q-card-section>

      <q-card-section class="m-8 p-4 flex-auto">
        <!-- 重要信息的描述内容，位于图片的右侧 -->
        <div class="text-sm text-center m-1">你正在阅读：</div>
        <div class="text-h3 text-center m-4 mb-16">书籍名称:{{ curBookInfo?.bookname }}</div>
        <div class="m-2 text-h6">书籍ID:{{ curBookInfo?.bookid }}</div>
        <div class="m-2 text-h6">作者信息:xxx</div>
        <div class="m-2 text-h6">文件信息:xxx</div>
        <div class="m-2 text-h6">上架时间:{{ curBookInfo?.CreatedAt }}</div>
        <div class="m-2 text-h6">加入时间:{{ curBookInfo?.CreatedAt }}</div>
        <div class="m-2 text-h6">图像链接:{{ curBookInfo?.bookimagepath }}</div>
      </q-card-section>
    </div>
    <q-card-section>
      <div class="flex justify-end">
        <!-- 功能性按钮们 -->
        <!-- <q-btn class="m-2 p-3" icon="menu_book" @click="_on_online_read_clicked">在线阅读</q-btn> -->
        <q-btn class="m-2 p-3" icon="menu_book" @click="_on_browser_read_clicked">浏览器阅读</q-btn>
        <q-btn class="m-2 p-3" icon="menu_book" @click="_on_check_firstpage_clicked">[DEBUG]查看封面PDF</q-btn>

        <!-- <q-btn class="m-2 p-3" icon="download">下载</q-btn> -->
        <q-btn class="m-2 p-3" icon="favorite">收藏</q-btn>
      </div>
    </q-card-section>

    <q-card-section>
      <!-- 详细介绍内容 -->
      <div class="m-4 p-4">书籍的详细介绍</div>
      <div>{{ curBookInfo }}</div>
    </q-card-section>
  </q-card>
</template>

<!-- 调整图片文件的尺寸。quasar和tailwind都没有现成的方法可用，只能写css了。 -->
<style scoped>
.book-image {
  width: 256px; /* 25% of the viewport width */
  height: 361px; /* 141% of the width, which is the aspect ratio of A4 paper */
  overflow: hidden;
}
</style>

<script lang="ts" setup>
// 书籍信息的数据结构。从后端API获取到的数据结构如下：
// {
//   "book": {
//     "ID": 1,
//     "CreatedAt": "2024-06-19T02:17:58.133751431+08:00",
//     "UpdatedAt": "2024-06-19T02:17:58.133751431+08:00",
//     "DeletedAt": null,
//     "bookid": "TestBooksID",
//     "bookname": "第一本测试书籍",
//     "bookimagepath": "",
//     "author": "",
//     "description": "",
//     "bookfiletype": "",
//     "bookfilehash": "4df6245b1adb0031f00415d06f723f80e61dbb60a51309fc8cc06dc39a8888cc",
//     "availablegroups": ""
//   }
// }

import { apiUrlStorage } from './../stores/api-urls';
import { BookInfo, GetBookInfoById } from './../api-methods';
import { onMounted, ref } from 'vue';

var curBookId = ref<string>('1');
var curBookInfo = ref<BookInfo>();

// Pinia, 读取urlStore的API地址
const urlStore = apiUrlStorage();

// 从URL参数的book_id中读取书籍ID
const urlParams = new URLSearchParams(window.location.search);
const book_id = urlParams.get('book_id');
if (book_id) {
  curBookId.value = book_id;
  console.log('找到book_id了！ :', book_id);
}

// 页面加载时，获取书籍信息
onMounted(async () => {
  curBookInfo.value = await GetBookInfoById(curBookId.value);
});

//
// ======== PDF 浏览器相关
// 格式化URL, 供pdf浏览器使用。使用URLSearchParams对象

// function _on_online_read_clicked() {
//   alert('在线阅读功能尚未实现！');
//   console.log('在线阅读功能尚未实现！');
// }

function _on_browser_read_clicked() {
  let pdfFileUrl = ref<string>();
  if (book_id != null) {
    let _pdfFileURL = new URL(urlStore.bookLibraryGetServedBookfileById);
    _pdfFileURL.searchParams.append('book_id', book_id);
    pdfFileUrl.value = _pdfFileURL.toString();
  } else {
    console.error('book_id is null!');
  }

  console.log('欲访问PDF文件的URL为:', pdfFileUrl.value);

  console.log('pdfFileUrl:', pdfFileUrl.value);
  window.open(pdfFileUrl.value, '_blank');
}

function _on_check_firstpage_clicked() {
  let pdfFileUrl = ref<string>();
  if (book_id != null) {
    let _pdfFileURL = new URL(urlStore.BookLibraryGetBookFirstPagePdf);
    _pdfFileURL.searchParams.append('book_id', book_id);
    pdfFileUrl.value = _pdfFileURL.toString();
  } else {
    console.error('book_id is null!');
  }

  console.log('欲访问PDF文件的URL为:', pdfFileUrl.value);

  console.log('pdfFileUrl:', pdfFileUrl.value);
  window.open(pdfFileUrl.value, '_blank');
}
</script>
