<template>
    <Page>
        <div class="Home">
            <div class="article">
                <div class="article-item" v-for="x in listPosts" :key="x.id">
                    <PostCard :post="x" />
                </div>
                <div class="pagedata">
                    <ClientOnly>
                        <el-pagination @current-change="changePageInt" background layout="prev, pager, next" :page-count="metadata.lastPage" />
                    </ClientOnly>
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
const metadata = ref<any>({})
const pagedata = {
    pageInt: 1,
    pageSize: 5
}

type PageData = typeof pagedata

const changePageInt = (pageInt: number) => {
    pagedata.pageInt = pageInt
    getPosts(pagedata)
}

const getPosts = (pd: PageData) => {
    useApi.listPostsApi({pageInt: pd.pageInt, pageSize: pd.pageSize}).then(res=>{
        listPosts.value = res.data.list
        metadata.value = res.data.metadata
    })
}

onMounted(()=>{
    getPosts(pagedata)
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

    .pagedata{
       display: flex;
       justify-content: center;
       padding: 20px 0 50px 0;
    }

    .additional {
        width: 280px;
    }

}
</style>
