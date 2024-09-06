<template>
    <div class="b_picview">
      <div class="b_picview-close" @click="emit('close')">
        <el-icon><CloseBold /></el-icon>
      </div>
      <img :src="picUrl" >
    </div>
    </template>
    
    <script setup lang="ts">
    import { CloseBold } from "@element-plus/icons-vue";
    const props = defineProps({picUrl: String})
    const emit = defineEmits(["close"])
    </script>
    
    <style scoped lang="scss">
    .b_picview {
      position: fixed;
      width: 100vw;
      height: 100vh;
      top: 0;
      left: 0;
      z-index: 9998;
      display: flex;
      justify-content: center;
      align-items: center;
      background: rgba($color: #000000, $alpha: .4);
      transition: all .2s;
      img{
        max-width: 90vw;
        object-fit: contain;
        transition: all .2s;
      }
      .b_picview-close {
        position: fixed;
        z-index: 9999;
        top: 30px;
        right: 30px;
        width: 34px;
        height: 34px;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 18px;
        color: #fff;
        background: rgba($color: #fff, $alpha: .5);
        border-radius: 40px;
      }
    }
    </style>