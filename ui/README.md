# UI

## 笔记

- 安装 Element Plus Nuxt 模块

参考：https://nuxt.com/modules/element-plus

npm i element-plus @element-plus/nuxt -D

- 配置一下，即可使用
```ts
export default defineNuxtConfig({
  modules: [
    '@element-plus/nuxt'
  ],
  elementPlus: { /** Options */ }
})
```
- 使用，如在 app.vue 里
``` vue
<template>
  <div>
    <el-button @click="ElMessage('hello')">button</el-button>
    <ElButton :icon="EditPen" type="success">button</ElButton>
    <LazyElButton type="warning">lazy button</LazyElButton>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus"
import { EditPen } from "@element-plus/icons-vue"
</script>
```
 