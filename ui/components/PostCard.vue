<template>
    <el-card style="position: relative; padding: 16px 32px 16px 16px;" shadow="always">
        <span class="post-attr">随笔</span>
        <h3 class="card-title">{{ post.title }}</h3>
        <p class="card-subtitle">{{ subTitle(post.content) }}</p>
        <div class="card-footer">
            <div class="time">
                <span>{{ post.createdAt }}</span>
            </div>
            <div>
                <NuxtLink v-if="userinfo.id && userinfo.accessToken" :to="`/posts/editor?id=${post.id}`">编辑</NuxtLink>
                <NuxtLink style="margin-left: 15px;" :to="`/posts/${post.id}`">阅读全文</NuxtLink>
            </div>
        </div>
    </el-card>
</template>

<script setup lang="ts">
import type { LoginResponse } from '~/types';

const props = defineProps<{ post: any }>()
const userinfo = ref({} as LoginResponse)

const subTitle = (content: string) => {
    if (content.length < 100) {
        return content
    } else {
        return content.substring(0, 100) + "..."
    }
}

onMounted(()=>{
    const { user } = useAuth()
    userinfo.value = user.value as LoginResponse
})

</script>

<style scoped lang="scss">
@import "~/assets/scss/vars";
@import "~/assets/scss/mixins";

:deep(.el-card__body) {
    padding: 0;
}

.card-title {
    @include ft(18px, $text, 400, 1.2);
}

.card-subtitle {
    @include ft(14px, $subtx, 400, 2.4);
}

.time {
    @include ft(12px, $subtx, 400, 1.2);
}

.card-footer {
    @include flex(space-between, center);
}

.post-attr {
    width: 100px;
    position: absolute;
    top: 0;
    right: 0;
    font-size: 14px;
    font-weight: 600;
    line-height: 28px;
    color: #fff;
    text-align: center;
    background: $primary;
    transform: rotate(45deg);
    transform-origin: 50px 50px;
}
</style>