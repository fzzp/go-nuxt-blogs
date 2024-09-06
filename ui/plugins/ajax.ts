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
        await nuxtApp.runWithContext(() => navigateTo('/login'))
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
