<template>
    <Page>
        <h1>Tag</h1>
        <div v-for="t in tags">
            <div :key="t.id">
                <NuxtLink>{{ t.tagName }}</NuxtLink>
            </div>
        </div>
    </Page>
</template>

<script setup lang="ts">
const { $ajax } = useNuxtApp()
const useApi = Api($ajax)

const tags = ref<any[]>([])

onMounted(() => {
    useApi.listTagsApi({ pageInt: 1, pageSize: 100 }).then(res => {
        tags.value = res.data
    })
})
</script>

<style scoped lang="scss"></style>