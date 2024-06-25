<template>
  <div class="book-image m-2 rounded-2xl">
    <canvas id="pdf-canvas" class="w-full h-full"></canvas>
  </div>
</template>

<script setup lang="ts">
// 测试首页PDF渲染，使用PDFJS
import * as pdfjsLib from 'pdfjs-dist';
import { PDFDocumentProxy } from 'pdfjs-dist/types/src/display/api';

// 使用Vite的import.meta.url来正确设置workerSrc的路径
// 假设你已经将pdf.worker.js复制到了public目录
// const workerSrc = new URL('public/pdf.worker.mjs', import.meta.url).href;
pdfjsLib.GlobalWorkerOptions.workerSrc = 'pdf.worker.mjs';

async function renderPDFPage(pdfFile: string) {
  // 在canvas id="pdf-canvas"的元素上渲染PDF页面
  try {
    const pdf: PDFDocumentProxy = await pdfjsLib.getDocument(pdfFile).promise;
    const page = await pdf.getPage(1);
    const canvasElement = document.getElementById('pdf-canvas') as HTMLCanvasElement;
    if (!canvasElement) {
      console.error('Canvas element not found');
      return;
    }
    const context = canvasElement.getContext('2d');
    if (!context) {
      console.error('Unable to get canvas context');
      return;
    }
    const viewport = page.getViewport({ scale: 1.5 });
    canvasElement.height = viewport.height;
    canvasElement.width = viewport.width;

    const renderContext = {
      canvasContext: context,
      viewport: viewport,
    };
    page.render(renderContext);
  } catch (error) {
    console.error('Error rendering PDF page:', error);
  }
}

renderPDFPage('http://127.0.0.1:8080/api/book_library/get_book_first_page_pdf?book_id=TestBooksID3withfile2');
</script>
