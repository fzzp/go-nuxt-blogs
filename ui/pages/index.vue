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
const { $ajax } = useNuxtApp()
const useApi = Api($ajax)

const listPosts = ref<any[]>([])

onMounted(()=>{
    useApi.listPostsApi({pageInt: 1, pageSize: 10}).then(res=>{
        console.log("res: ",res)
        listPosts.value = res.data
    })
})
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
        width: 280px;
    }
}
</style>
