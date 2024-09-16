

// plugins/my-plugin.ts
export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook('app:created', () => {
    // Vue 应用实例已创建， 可以做些初始化事情...
    const { $ajax } = useNuxtApp()

    // 检查是否登陆了
    const { user } = useAuth()
    if (user.value?.id && user.value.accessToken) {
      const useApi = Api($ajax)
      useApi.getUserInfo() // 请求检查一下token是否过期
    }
  });
});
