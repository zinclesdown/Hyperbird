<template>
  <canvas ref="pdfCanvas"></canvas>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as pdfjsLib from 'pdfjs-dist';
import { PDFDocumentProxy } from 'pdfjs-dist/types/src/display/api';
import * as APIMethods from '../api-methods';

const props = defineProps({
  bookId: {
    type: String,
    required: true,
  },
});

const pdfCanvas = ref<HTMLCanvasElement | null>(null);

pdfjsLib.GlobalWorkerOptions.workerSrc = '/public/pdf.worker.mjs';

async function renderPDFPage(pdfFile: string) {
  if (!pdfCanvas.value) {
    console.error('Canvas element not found');
    return;
  }
  try {
    const pdf: PDFDocumentProxy = await pdfjsLib.getDocument(pdfFile).promise;
    const page = await pdf.getPage(1);
    const context = pdfCanvas.value.getContext('2d');
    if (!context) {
      console.error('Unable to get canvas context');
      return;
    }
    const viewport = page.getViewport({ scale: 1.5 });
    pdfCanvas.value.height = viewport.height;
    pdfCanvas.value.width = viewport.width;

    const renderContext = {
      canvasContext: context,
      viewport: viewport,
    };
    page.render(renderContext);
  } catch (error) {
    console.error('Error rendering PDF page:', error);
  }
}

onMounted(() => {
  APIMethods.GetFirstPagePDFUrl(props.bookId).then((url) => {
    if (url == null) {
      console.error('获取首页PDF的URL失败！');
      return;
    } else {
      renderPDFPage(url);
    }
  });
});
</script>
