import type { LoginResponse } from '~/types';
import { useAuth } from './../composables/useAuth';
export default defineNuxtPlugin((nuxtApp) => {
  const ajax = $fetch.create({
    baseURL: useRuntimeConfig().public.Apiprev,
    onRequest({ request, options, error }) {
      const { user } = useAuth()
      const headers: any = {
        "Accept": "application/json",
        "Authorization": user.value?.accessToken ? `Bearer ${user.value?.accessToken}` : undefined
      }
      options.headers = headers
    },
    onResponse({ request, options, error }) {
      // 做些什么...
    },
    async onResponseError({ response }) {
      if (response.status === 401) {
        const { user } = useAuth()
        user.value = {} as LoginResponse // 登陆失效，清空个人信息
        if(!response.url.includes("auth/getUserInfo")){
          await nuxtApp.runWithContext(() => navigateTo('/login'))
        }
      }
    }
  })

  // Expose to useNuxtApp().$request
  return {
    provide: {
      ajax
    }
  }
})
