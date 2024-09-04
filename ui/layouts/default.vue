<template>
    <div class="default-layout">
        <header class="default-layout-header">
            <nav class="navbar">
                <div class="flex max-w">
                    <NuxtLink class="nav-link" to="/">
                        Fzzp Blog
                    </NuxtLink>
                    <ul class="navbar-nav">
                        <li class="nav-item">
                            <NuxtLink class="nav-link" aria-current="page" to="/">首页</NuxtLink>
                        </li>
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link" to="/category">分类</NuxtLink>
                        </li>
                        <li class="nav-item ms-4">
                            <NuxtLink class="nav-link" to="/about">关于</NuxtLink>
                        </li>
                    </ul>
                </div>
            </nav>
        </header>
        <main class="default-layout-content">
            <slot />
        </main>
    </div>
</template>

<script setup lang="ts">

</script>

<style scoped lang="scss">
.default-layout{
    height: 45px;
}
.default-layout-header {
    .navbar {
        position: fixed;
        top: 0;
        left: 0;
        z-index: 99;
        width: 100%;
        padding: 12px;
        background-color: #f5f5f5;
        border-bottom: 1px solid #ccc;
        .flex {
            display: flex;
        }

        .navbar-nav {
            display: flex;
            gap: 30px;
            margin-left: auto;
        }
        .nav-link{
            text-decoration: none;
        }
    }
}
.default-layout-content{
    padding-top: 60px;
}
</style>