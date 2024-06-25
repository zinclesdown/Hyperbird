<template>
  <q-card class="m-8">
    <div class="row">
      <q-card-section class="flex items-center m-4">
        <!-- 图像的显示区域，位于左上侧。 -->
        <!-- <div class="book-image m-2 rounded-2xl">
          <q-img src="https://cdn.quasar.dev/img/parallax2.jpg" alt="book image" class="w-full h-full object-fill" />
        </div> -->
        <div class="book-image m-2 rounded-2xl shadow-lg">
          <!-- <canvas id="pdf-canvas" class="w-full h-full"></canvas> -->
          <PdfPreview :bookId="curBookId" class="w-full h-full" />
        </div>
      </q-card-section>

      <q-card-section class="m-8 p-4 flex-auto bg-gray-50 rounded-3xl shadow-lg">
        <!-- 重要信息的描述内容，位于图片的右侧 -->
        <div class="text-sm text-center m-1">你正在阅读：</div>
        <div class="text-h3 text-center m-4 mb-16">{{ curBookInfo?.book_name }}</div>
        <div class="m-2 text-h6">ID:{{ curBookInfo?.book_id }}</div>
        <div class="m-2 text-h6">作者信息:{{ curBookInfo?.author }}</div>
        <div class="m-2 text-h6">文件信息:{{ curBookInfo?.book_file_hash }}</div>
        <!-- <div class="m-2 text-h6">上架时间:{{ curBookInfo?.CreatedAt }}</div>
        <div class="m-2 text-h6">加入时间:{{ curBookInfo?.CreatedAt }}</div>
        <div class="m-2 text-h6">图像链接:{{ curBookInfo?.bookimagepath }}</div> -->
      </q-card-section>
    </div>
    <q-card-section>
      <!-- <div class="q-pa-md q-gutter-sm">
        <q-btn color="primary" icon="mail" label="左侧" />
        <q-btn color="secondary" icon-right="mail" label="On Right" />
        <q-btn color="red" icon="mail" icon-right="send" label="On Left and Right" />
        <br />
        <q-btn icon="phone" label="Stacked" stack glossy color="purple" />
      </div> -->

      <div class="flex justify-end q-pa-md q-gutter-sm">
        <!-- 功能性按钮们 -->
        <q-btn class="m-2 p-3" push color="primary" icon="menu_book" @click="_on_browser_read_clicked" label="浏览器阅读" />
        <q-btn class="m-2 p-3" push color="secondary" icon="menu_book" @click="_on_check_firstpage_clicked" label="查看封面" />
        <q-btn class="m-2 p-3" push color="red" icon="favorite">收藏</q-btn>
      </div>
    </q-card-section>

    <q-card-section>
      <!-- 详细介绍内容 -->
      <div class="m-4 p-4 text-sm">书籍的详细介绍: {{ curBookInfo?.description }}</div>
      <div class="text-gray-500 m-4 p-8 bg-gray-50 rounded-3xl">{{ curBookInfo }}</div>
    </q-card-section>
  </q-card>
</template>

<!-- 调整图片文件的尺寸。quasar和tailwind都没有现成的方法可用，只能写css了。 -->
<style scoped>
.book-image {
  width: 320px; /* 25% of the viewport width */
  height: 450px; /* 141% of the width, which is the aspect ratio of A4 paper */
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
import { BookInfo, GetBookInfoById } from './../api-methods';
import { onMounted, ref } from 'vue';

// import * as pdfjsLib from 'pdfjs-dist';
// import { PDFDocumentProxy } from 'pdfjs-dist/types/src/display/api';
import * as APIMethods from '../api-methods';

var curBookId = ref<string>('1');
var curBookInfo = ref<BookInfo>();

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

async function _on_browser_read_clicked() {
  let url = await APIMethods.GetBookPDFUrl(curBookId.value);
  window.open(url, '_blank');
}

async function _on_check_firstpage_clicked() {
  let url = await APIMethods.GetFirstPagePDFUrl(curBookId.value);
  window.open(url, '_blank');
}

// // 预览首页PDF
// pdfjsLib.GlobalWorkerOptions.workerSrc = '/public/pdf.worker.mjs';
// async function renderPDFPage(pdfFile: string) {
//   // 在canvas id="pdf-canvas"的元素上渲染PDF页面
//   try {
//     const pdf: PDFDocumentProxy = await pdfjsLib.getDocument(pdfFile).promise;
//     const page = await pdf.getPage(1);
//     const canvasElement = document.getElementById('pdf-canvas') as HTMLCanvasElement;
//     if (!canvasElement) {
//       console.error('Canvas element not found');
//       return;
//     }
//     const context = canvasElement.getContext('2d');
//     if (!context) {
//       console.error('Unable to get canvas context');
//       return;
//     }
//     const viewport = page.getViewport({ scale: 1.5 });
//     canvasElement.height = viewport.height;
//     canvasElement.width = viewport.width;

//     const renderContext = {
//       canvasContext: context,
//       viewport: viewport,
//     };
//     page.render(renderContext);
//   } catch (error) {
//     console.error('Error rendering PDF page:', error);
//   }
// }

// // 渲染首页PDF
// APIMethods.GetFirstPagePDFUrl(curBookId.value).then((url) => {
//   if (url == null) {
//     console.error('获取首页PDF的URL失败！');
//     return;
//   } else {
//     renderPDFPage(url);
//   }
// });
import PdfPreview from 'src/components/PdfPreview.vue';
</script>
