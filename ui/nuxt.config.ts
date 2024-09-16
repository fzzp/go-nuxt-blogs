// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-09-04',
  devtools: { enabled: false },
  modules: [
    '@element-plus/nuxt'
  ],
  components: [{
    path: "~/components",
    pathPrefix: false
  }],
  runtimeConfig: {
    public: {
      Apiprev: process.env.NUXT_PUBLIC_API_BASE,
      ImgPrev: process.env.NUXT_PUBLIC_IMG_BASE
    }
  },
  plugins: [
    {
      src: "~/plugins/checkAuth",
      mode: "client"
    }
  ],
  css: [
    "~/assets/scss/global.scss",
  ],
  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `
            @use "@/assets/scss/element.scss" as element;
          `
        }
      },
    },
  },
  elementPlus: {
    icon: "ElIcon",
    importStyle: "scss",
  },
})
