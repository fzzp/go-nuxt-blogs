<template>
    <div class="uploadFileBox cp" @click="fileRef?.click()">
        <div class="default-body" v-if="!body">
            <img src="@/assets/imgs/file-upload.png" :width="imgWidth" />
            <span class="ft14-h">{{ text }}</span>
        </div>
        <div v-else>
            <slot name="body"></slot>
        </div>
        <input type="file" hidden ref="fileRef" @change="changeFile" :accept="accept">
    </div>
</template>

<script setup lang="ts">
const {$ajax} = useNuxtApp()
const useApi = Api($ajax)

const props = withDefaults(
    defineProps<{ body?: boolean, text?: string, accept?: string, imgWidth?: string, }>(),
    { body: false, text: "", accept: "'image/*'", imgWidth: "60", }
)

const emit = defineEmits(["processHandler"])
const fileRef = ref<HTMLInputElement>()


const changeFile = (event: any) => {
    if (event.target.files[0]) {
        emit("processHandler", { status: "Start", })
        // 限制4m
        useApi.formDataToUploadFile(event.target.files[0], 4 * 1024).then((res: any) => {
            const {data} = res
            debugger
            if(data && data.length && data[0].id) {
                let fileUrl = useRuntimeConfig().public.ImgPrev + data[0].id
                emit("processHandler", { newFileUr: fileUrl, status: "Success" })
                ElMessage.success("上传成功")
            }else {
                emit("processHandler", { status: "Failed" })
                ElMessage.error("上传失败")
            }
        }).catch((err: any) => {
            if (err && err.type === "ErrFileTooLong") {
                ElMessage.error(err.message)
            }
            emit("processHandler", { status: "Failed", })
        })
    } else {
        ElMessage.error("获取本地图片失败！")
    }
}
</script>

<style scoped lang="scss">
@import "~/assets/scss/vars";

.uploadFileBox {
    width: 100%;
    height: 100%;
    display: flex;
    padding: 8px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    border: 1px solid $border;

    .default-body {
        display: flex;
        flex-direction: column;
        text-align: center;
        align-items: center;
        margin: auto;

        span {
            color: $border;
        }
    }
}
</style>