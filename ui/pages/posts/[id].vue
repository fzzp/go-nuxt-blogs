<template>
    <Page>
        <main class="post-detail-page">
            <div class="post-container">
                <MdPreview :editorId="id" :modelValue="text" />
            </div>
            <!-- <div class="postCatalog">
            <MdCatalog :editorId="id" :scrollElement="scrollElement" />
        </div> -->
        </main>
    </Page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { MdPreview, MdCatalog } from 'md-editor-v3';
// preview.css相比style.css少了编辑器那部分样式
import 'md-editor-v3/lib/preview.css';

const id = 'preview-only';
const text = ref('# Hello Editor');
const scrollElement = ".postCatalog";

onMounted(async () => {
    console.log(1111111)
    const { data } = await useFetch<string>("/example.md")
    text.value = data.value as string
})

</script>

<style scoped lang="scss">
.post-detail-page {
    .post-container {
        width: 60vw;
        margin: 0 auto;
    }
}
</style>