<template>
    <el-card style="position: relative; padding: 16px 32px 16px 16px;" shadow="always">
        <span class="post-attr">随笔</span>
        <h3 class="card-title">{{ post.title }}</h3>
        <p class="card-subtitle">{{ subTitle(post.content) }}</p>
        <div class="card-footer">
            <div class="">
                <span>{{ post.createdAt }}</span>
            </div>
            <NuxtLink :to="`/posts/${post.id}`">
                <!-- <el-button type="primary" link> -->
                    阅读全文
                <!-- </el-button> -->
            </NuxtLink>
        </div>
    </el-card>
</template>

<script setup lang="ts">
const props = defineProps<{ post: any }>()
const subTitle = (content: string) => {
    if (content.length < 100) {
        return content
    } else {
        return content.substring(0, 100) + "..."
    }
}

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
    line-height: 24px;
    color: #fff;
    text-align: center;
    background: $primary;
    transform: rotate(45deg);
    transform-origin: 50px 50px;
}
</style>