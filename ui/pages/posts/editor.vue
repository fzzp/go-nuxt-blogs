<template>
	<ClientOnly>
		<main>
			<MdEditor v-model="form.content" :page-fullscreen="true" 
				:toolbars="toolbars" placeholder="写点什么..."
				@onChange="onChange" @onGetCatalog="onGetCatalog"
				:codeFoldable="false"
				code-theme="a11y"
				@onSave="drawer = true">
				<template #defToolbars>
					<NormalToolbar title="图片" @onClick="handler">
						<template #trigger>
							<svg class="md-editor-icon" aria-hidden="true">
								<use xlink:href="#md-editor-icon-image"></use>
							</svg>
						</template>
					</NormalToolbar>
					<NormalToolbar title="状态" @onClick="handler">
						<template #trigger>
							<span v-show="sLoading">
								<el-button type="success" size="small" text :loading="sLoading">
									正在保存...
								</el-button>
							</span>
						</template>
					</NormalToolbar>
				</template>
			</MdEditor>
		</main>
		<PhotoGallery v-model="dialogVisible" />
		<el-drawer v-model="drawer" :close-on-click-modal="false">
			<el-form :model="form" label-width="auto" style="max-width: 600px">
				<el-form-item label="标题">
					<el-input v-model="form.title" type="textarea" placeholder="标题" />
				</el-form-item>
			</el-form>

			<el-form-item label="属性">
				<el-select v-model="form.attrId">
					<el-option v-for="v in attributes" :label="v.attrName" :value="v.id" :key="v.id" />
				</el-select>
			</el-form-item>

			<el-form-item label="标签">
				<el-select v-model="form.tags" multiple>
					<el-option v-for="v in tags" :label="v.tagName" :value="v.id" :key="v.id" />
				</el-select>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" @click="onSubmit()">保存</el-button>
			</el-form-item>
		</el-drawer>
	</ClientOnly>
</template>

<script setup lang="ts">
import { MdEditor, NormalToolbar } from 'md-editor-v3';
import { toolbars } from "@/data/toolbar";
import 'md-editor-v3/lib/style.css';

definePageMeta({
	middleware: ['auth'] // FIXME: 感觉好像不生效
})

const { $ajax } = useNuxtApp()
const useApi = Api($ajax)
const dialogVisible = ref(false)
const drawer = ref(false)
const attributes = ref<any[]>([])
const tags = ref<any[]>([])
const route = useRoute()
const sLoading = ref(false)
let timer: any = null
const formKey = "new_post"

const form = reactive({
	id: 0,
	title: "",
	attrId: 1,
	content: "",
	tags: [],
})

definePageMeta({
	layout: false
})

const onChange = (text: string) => {
	form.content = text

	// 定时存储
	if (timer) {
		clearTimeout(timer)
	}
	sLoading.value = true
	timer = setTimeout(() => {
		onSubmit(true)
	}, 3000)

}

const onGetCatalog = (list: any[]) => {
	list.forEach((v: any) => {
		if (v.level == 1) {
			form.title = v.text
		}
	})
};

const onSubmit = (auto?: boolean) => {
	if (form.id > 0) {
		// 编辑
		useApi.updatePostApi(form).then(res => {
			!auto && ElMessage.success("更新成功")
			drawer.value = false
		}).finally(() => { sLoading.value = false })
	} else {
		// 新增
		const store = useStore(formKey, form)
		if (auto) {
			// 自动的，则不新增,保存到localstorage里
			store.value = form
			sLoading.value = false
			return
		}
		useApi.createPostApi(form).then(res => {
			!auto && ElMessage.success(res.message)
			// 提交成功了，则重置
			store.value = {} as any
			drawer.value = true
			navigateTo("/")
		}).finally(() => { sLoading.value = false })
	}
}

const handler = () => {
	dialogVisible.value = true
};

onMounted(() => {
	// 多做一次校验
	const { user } = useAuth()
	if(!(user?.value?.id && user.value?.accessToken)) {
		return navigateTo("/login", {replace: true})
	}
	let id = Number(route.query.id)
	if (!isNaN(id) && id > 0) {
		// 编辑
		useApi.getPostDetailApi(id).then(res => {
			const { data } = res
			if (data) {
				let tags = data.tags || []
				form.id = data.id
				form.title = data.title
				form.attrId = data.attrId
				form.content = data.content
				form.tags = tags.map((v: any) => v.id)
			}
		})
	}else {
		// 新文章，从localstorage取，如果上次有没提交数据的话
		const store = useStore(formKey, {} as any)
		let formJson = JSON.parse(JSON.stringify(form))
		let mergeJson = {...formJson, ...store.value}
		form.title = mergeJson.title
		form.attrId = mergeJson.attrId
		form.content = mergeJson.content
		form.tags = mergeJson.tags
		if(store.value!.content) {
			ElMessage.success("加载上次未提交保存文章")
		}
	}

	useApi.listAttributesApi().then(res => {
		attributes.value = res.data
	})

	useApi.listTagsApi({ pageInt: 1, pageSize: 100 }).then(res => {
		tags.value = res.data
	})
})

onUnmounted(()=>{
	// 消除编辑器的副作用，返回其他页面无法滚动
	document.body.style.overflow = "auto"
})

</script>

<style scoped lang="scss"></style>