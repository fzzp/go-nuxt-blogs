<template>
    <el-dialog v-model="dialogVisible" :append-to-body="true" :close-on-click-modal="false" class="PhotoGallery"
        width="800" :z-index="100009" :close-icon="CloseIcon">
        <template #header>
            <div class="head-center">
                <el-radio-group v-model="pos" size="large">
                    <el-radio-button value="Gallery">图库</el-radio-button>
                    <el-radio-button value="Upload">上传</el-radio-button>
                </el-radio-group>
            </div>
        </template>
        <div class="body-container">
            <div v-show="pos === 'Gallery'">
                <div class="d-flex" style="gap: 16px; flex-wrap: wrap">
                    <template v-for="v in listPics" :key="v.id">
                        <PicBox :pic-url="v.picUrl"></PicBox>
                    </template>
                </div>
            </div>
            <div v-show="pos === 'Upload'">
                <div class="d-flex" style="gap: 16px">
                    <div v-if="picUrl">
                        <PicBox :picUrl="picUrl"></PicBox>
                    </div>
                    <div style="width:160px; height: 160px; cursor: pointer;">
                        <FileUploadBox @processHandler="processHandler"></FileUploadBox>
                    </div>
                </div>
            </div>
        </div>
    </el-dialog>
</template>

<script setup lang="ts">
import CloseIcon from './CloseIcon.vue';
const dialogVisible = defineModel<boolean>()
const {$ajax} = useNuxtApp()
const useApi = Api($ajax)

const picUrl = ref("")
const listPics = ref<any[]>([])
const pos = ref("Gallery") // Gallery ｜ Upload

const processHandler = (data: {status: string, newFileUr?: string}) => {
    if(data.status == "Success") {
        picUrl.value = data?.newFileUr || ''
        getListPics()
    }
}

const getListPics = () => {
    useApi.listPicsApi({pageInt: 1, pageSize:100, fileType: 1}).then(res=>{
        try {
            listPics.value = res.data.map((v:any)=>{
            return {
                ...v,
                picUrl: useRuntimeConfig().public.ImgPrev + v.id
            }
        })
        } catch (error) {
            ElMessage.error("处理数据错误")
        }
        
    })
}

onMounted(()=>{
    getListPics()
})

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

    .body-container {
        padding: 16px 16px 50px 16px;
        display: flex;
        justify-content: center;
    }
}
</style>