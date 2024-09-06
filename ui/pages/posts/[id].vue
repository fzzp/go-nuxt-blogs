<template>
    <Page>
        <main class="post-detail-page">
            <div class="post-container">
                <MdPreview :editorId="id" :modelValue="post.content" :codeFoldable="false" code-theme="a11y" />
            </div>
            <el-button type="primary" v-if="userinfo.id && userinfo.accessToken" :icon="Edit"
                @click="goEdit">去编辑</el-button>
            <!-- <div class="postCatalog">
            <MdCatalog :editorId="id" :scrollElement="scrollElement" />
        </div> -->
        </main>
    </Page>
</template>

<script setup lang="ts">
import { Edit } from "@element-plus/icons-vue";
import { ref } from 'vue';
import { MdPreview, MdCatalog } from 'md-editor-v3';
// preview.css相比style.css少了编辑器那部分样式
import 'md-editor-v3/lib/preview.css';
import type { LoginResponse } from "~/types";
const { $ajax } = useNuxtApp()
const useApi = Api($ajax)
const route = useRoute()
const post = ref<any>({})
const userinfo = ref({} as LoginResponse)

const id = 'preview-only';
const scrollElement = ".postCatalog";

const goEdit = () => {
    navigateTo("/posts/editor?id=" + route.params.id)
}

onMounted(async () => {
    useApi.getPostDetailApi(Number(route.params.id)).then(res => {
        post.value = res.data
    })

    const { user } = useAuth()
    userinfo.value = user.value as LoginResponse

})

</script>

<style scoped lang="scss">
.post-detail-page {
    .post-container {
        margin: 0 auto;
    }

    :deep(.md-editor-preview-wrapper) {
        padding: 10px 0;
    }
}
</style>