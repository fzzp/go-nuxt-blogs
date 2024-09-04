<template>
    <Page>
        <div class="Home">
            <div class="article">
                <div class="article-item" v-for="x in listPosts" :key="x.id">
                    <PostCard :post="x" />
                </div>
            </div>

            <div class="additional">
                <Additional />
            </div>
        </div>
    </Page>
</template>

<script setup lang="ts">
const listPosts = [
    {
        id: 1,
        title: "标题标题",
        content: "内容内容内容内容内容内容内容内容内容内容内容内容内容内容",
        createdAt: "2024-09-04"
    },
    {
        id: 2,
        title: "标题标题",
        content: "内容内容内容内容内容内容内容内容内容内容内容内容内容内容",
        createdAt: "2024-09-04"
    },
    {
        id: 3,
        title: "标题标题",
        content: "内容内容内容内容内容内容内容内容内容内容内容内容内容内容",
        createdAt: "2024-09-04"
    },
]
</script>

<style scoped lang="scss">
.Home {
    display: flex;
    position: relative;

    .article {
        flex-grow: 1;
        .article-item{
            margin-bottom: 16px;
        }
    }

    .additional {
        width: 240px;
    }
}
</style>
