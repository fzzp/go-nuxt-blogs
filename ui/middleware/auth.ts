import { useAuth } from './../composables/useAuth';


export default defineNuxtRouteMiddleware((to, form) => {
  const { user } = useAuth()
  if (!(user.value?.id && user.value.accessToken)){
    return navigateTo("/login", {replace: true})
  }

  return
})