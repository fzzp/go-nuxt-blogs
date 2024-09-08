<template>
    <el-card style="margin-left: 16px">
        <div class="user">
            <el-avatar :size="80" :src="Avatar" />
            <div style="padding-top: 10px;">{{ userinfo.username || "fzzp" }}</div>
        </div>
        <div class="stat-box">
            <div class="item">
                <div class="item-head">文章数</div>
                <div class="item-footer">{{sumTotal.totalPosts}}</div>
            </div>
            <div class="item">
                <div class="item-head">阅读数</div>
                <div class="item-footer">{{sumTotal.totalViews}}</div>
            </div>
            <div class="item">
                <div class="item-head">评论数</div>
                <div class="item-footer">{{sumTotal.totalComments}}</div>
            </div>
        </div>
        <div style="text-align: center;">
            <div v-if="userinfo.id" style="padding-top: 15px;">
                <NuxtLink to="/posts/editor" target="_blank">去写文章</NuxtLink>
            </div>
        </div>
    </el-card>
</template>

<script setup lang="ts">
import type { LoginResponse } from '~/types';
import Avatar from "~/assets/imgs/avatar.png";

const { $ajax } = useNuxtApp()
const useApi = Api($ajax)
const userinfo = ref({} as LoginResponse)
const sumTotal = ref<any>({})

onMounted(() => {
    const { user } = useAuth()
    userinfo.value = { ...user.value } as LoginResponse

    useApi.getSettings().then(res=>{
        sumTotal.value = res.data
    })
})
</script>

<style scoped lang="scss">
@import "~/assets/scss/vars";
@import "~/assets/scss/mixins";

.user {
    text-align: center;
}

.stat-box {
    display: flex;
    justify-content: space-around;
    align-items: center;
    margin-top: 15px;

    .item {
        div {
            text-align: center;
        }

        .item-head {
            padding-bottom: 3px;
            font-size: 14px;
            color: $subtx;
            font-weight: 500;
        }

        .item-footer {
            font-size: 18px;
            color: $text;
        }
    }
}
</style>