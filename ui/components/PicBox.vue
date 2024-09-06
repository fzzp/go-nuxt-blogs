<template>
    <div class="b-picbox">
        <div class="pic">
            <img :src="picUrl" :width="width" :height="height" @dblclick="showBigPic = true">
            <span class="copy cp" @click="handleCopy">复制链接</span>
        </div>
        <PicView :picUrl="picUrl" v-show="showBigPic" @close="showBigPic = false" />
    </div>
</template>


<script setup lang="ts">
const props = withDefaults(defineProps<{
    picUrl: string;
    width?: string;
    height?: string;
}>(), {
    width: "150px",
    height: "auto"
})

const handleCopy = () => {
    createBtnCopy(props.picUrl)
}

const showBigPic = ref(false)

onMounted(() => {

})

</script>

<style scoped lang="scss">
.b-picbox {
    border: 1px solid #18171c1a;
    overflow: hidden;

    .pic {
        padding: 5px;
        position: relative;

        img {
            object-fit: cover;
            vertical-align: middle;
        }

        .copy {
            position: absolute;
            top: 0;
            right: 0;
            display: flex;
            justify-content: center;
            padding: 4px;
            font-size: 12px;
            color: #fff;
            background-color: #C167EFe0;
            border-bottom-left-radius: 3px;
            border-bottom-right-radius: 3px;
        }
    }
}
</style>