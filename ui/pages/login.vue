<template>
    <Page>
        <el-form :model="form" :rules="rules" label-width="auto" style="max-width: 400px; margin: 50px auto;">
            <el-form-item label="邮箱" prop="email">
                <el-input v-model="form.email" />
            </el-form-item>

            <el-form-item label="密码" prop="password">
                <el-input v-model="form.password" type="password" show-password />
            </el-form-item>

            <el-form-item>
                <el-button type="primary" @click="onSubmit">登陆</el-button>
            </el-form-item>
        </el-form>
    </Page>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
const { $ajax } = useNuxtApp()
const useApi = Api($ajax)

const form = reactive({
    email: '',
    password: '',
})

const rules = reactive({
    email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
    ],
})

const onSubmit = () => {
    const { user } = useAuth()
    useApi.postLoginApi(form).then(res=>{
        user.value = res.data        
        ElMessage.success(res.message)
        navigateTo("/")
    }).catch(err=>{
        console.log(err)
        ElMessage.error(err?.data?.message||"请求失败")
    })
}
</script>

<style scoped lang="scss"></style>