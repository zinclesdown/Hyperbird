<template>
  <div class="book-image m-2 rounded-2xl">Hello world/</div>

  <q-btn @click="UserMethods.printUserStatus()">显示用户</q-btn>

  <q-btn @click="UserMethods.DebugFakeLogin()">虚拟登陆</q-btn>
</template>

<script setup lang="ts">
console.log('Hello world');

import * as UserMethods from '../user-methods';
</script>
