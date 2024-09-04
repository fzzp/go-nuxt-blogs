<template>
    <el-dialog v-model="dialogVisible" :append-to-body="true" :close-on-click-modal="false" class="PhotoGallery"
        width="800" :z-index="100009"
        :close-icon="CloseIcon"    
    >
        <template #header>
            <div class="head-center">
                <el-radio-group v-model="pos" size="large">
                    <el-radio-button value="Gallery">图库</el-radio-button>
                    <el-radio-button value="Upload">上传</el-radio-button>
                </el-radio-group>
            </div>
        </template>
        <div class="body-container">
            <div v-show="pos === 'Gallery'">Gallery</div>
            <div v-show="pos === 'Upload'">Upload</div>
        </div>
    </el-dialog>
</template>

<script setup lang="ts">
import CloseIcon from './CloseIcon.vue';
const dialogVisible = defineModel<boolean>()

const pos = ref("Gallery") // Gallery ｜ Upload

</script>

<style lang="scss" coped>
.PhotoGallery {
    &.el-dialog {
        --el-dialog-padding-primary: 0;
    }

    .el-radio-button__inner {
        --el-border-radius-base: 6px;
        --el-font-size-base: 15px;
    }

    .head-center {
        display: flex;
        justify-content: center;
        padding: 30px;
    }

    .body-container{

    }
}
</style>