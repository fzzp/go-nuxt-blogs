<template>
	<ClientOnly>
		<main>
			<MdEditor v-model="text" :page-fullscreen="true" :toolbars="toolbars" @onSave="onSave">
				<template #defToolbars>
					<NormalToolbar title="图片" @onClick="handler">
						<template #trigger>
							<svg class="md-editor-icon" aria-hidden="true">
								<use xlink:href="#md-editor-icon-image"></use>
							</svg>
						</template>
					</NormalToolbar>
				</template>
			</MdEditor>
			<PhotoGallery v-model="dialogVisible" />
		</main>
	</ClientOnly>
</template>

<script setup lang="ts">
import { MdEditor, NormalToolbar } from 'md-editor-v3';
import { toolbars } from "@/data/toolbar";
import 'md-editor-v3/lib/style.css';

const text = ref("Hello Editor")
const dialogVisible = ref(false)

definePageMeta({
	layout: false
})

const onSave = (v: string, h: Promise<string>) => {
	console.log("-> ", v);

	h.then((html) => {
		console.log(html);
	});
};

const handler = () => {
	console.log('自定义图片!');
	dialogVisible.value = true
};

</script>

<style scoped lang="scss"></style>