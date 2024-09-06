import type { NitroFetchRequest, $Fetch } from 'nitropack'

import type {
    ApiEcho,
    LoginResponse,
    NewPostRequest,
    PageInfo,
} from '~/types'

export const Api = <T>(fetch: $Fetch<T, NitroFetchRequest>) => ({
    async postLoginApi(body: { email: string, password: string }): Promise<ApiEcho<LoginResponse>> {
        return fetch('/v1/login', { "method": "POST", body })
    },

    async createPostApi(body: NewPostRequest): Promise<ApiEcho<number>> {
        return fetch('/v1/auth/createPost', { "method": "POST", body })
    },

    async updatePostApi(body: any): Promise<ApiEcho<any>> {
        return fetch('v1/auth/updatePost', { "method": "POST", body })
    },

    async listPostsApi(query: PageInfo): Promise<ApiEcho<any[]>> {
        return fetch('/v1/posts', { "method": "GET", query })
    },

    async listPicsApi(query: PageInfo & {fileType: number}): Promise<ApiEcho<any>> {
        return fetch('/v1/auth/files', { "method": "GET", query })
    },

    async listAttributesApi(): Promise<ApiEcho<any>> {
        return fetch('/v1/attributes', { "method": "GET" })
    },

    async listTagsApi(query: PageInfo): Promise<ApiEcho<any>> {
        return fetch('/v1/tags', { "method": "GET", query })
    },

    async getPostDetailApi(id: number): Promise<ApiEcho<any>> {
        return fetch(`/v1/posts/${id}`, { "method": "GET" })
    },

    async formDataToUploadFile(file: any, maxSize: number) {
        // 校验文件大小, file.size 单位是字节，因此需要转换为kb
        if ((file.size / 1024) > maxSize) {
            const unitMB = Math.floor(maxSize / 1024)
            Promise.reject({ message: `上传图片失败,该图片大小超过${unitMB}M,请压缩后上传!`, type: "ErrFileTooLong" })
        }

        const formData = new FormData()
        // 接口是支持多文件上传的，这里只做单文件上传
        formData.append("files", file) // NOTE: 这里最终类似{files: [file]}
        formData.append("fileType", "IMAGE")

        const { user } = useAuth()

        // return fetch('/v1/auth/savefile', { "method": "POST", formData })
        return $fetch(useRuntimeConfig().public.Apiprev + "/v1/auth/savefile", {
            method: 'POST',
            body: formData,
            headers: {
                'Authorization': 'Bearer ' + user.value?.accessToken
            }
        })
    }
})

