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
        <div class="text-h3 text-center m-4 mb-16">书籍名称</div>
        <div class="m-2 text-h6">书籍ID:xxx</div>
        <div class="m-2 text-h6">作者信息:xxx</div>
        <div class="m-2 text-h6">文件信息:xxx</div>
        <div class="m-2 text-h6">上架时间:xxx</div>
        <div class="m-2 text-h6">加入时间:xxx</div>
      </q-card-section>
    </div>
    <q-card-section>
      <div class="flex justify-end">
        <!-- 功能性按钮们 -->
        <q-btn class="m-2 p-3" icon="menu_book">阅读</q-btn>
        <q-btn class="m-2 p-3" icon="download">下载</q-btn>
        <q-btn class="m-2 p-3" icon="favorite">收藏</q-btn>
      </div>
    </q-card-section>

    <q-card-section>
      <!-- 详细介绍内容 -->
      <div class="m-4 p-4">书籍的详细介绍</div>
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

<script>
export default {
  created() {
    console.log('子页面被创建了');
  },
  mounted() {
    console.log('子页面被挂载了');
  },
};
</script>
