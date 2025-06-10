<template>
  <div class="wh-full flex-col bg-[url(@/assets/images/login_bg.webp)] bg-cover">
    <div
      class="m-auto max-w-700 min-w-345 f-c-c rounded-8 bg-opacity-20 bg-cover p-12 card-shadow auto-bg"
    >
      <div class="hidden w-380 px-20 py-35 md:block">
        <img src="@/assets/images/login_banner.webp" class="w-full" alt="login_banner">
      </div>

      <div class="w-320 flex-col px-20 py-32">
        <h2 class="f-c-c text-24 text-#6a6a6a font-normal">
          <img src="@/assets/images/logo.png" class="mr-12 h-50">
          {{ title }}
        </h2>
        <n-input
          v-model:value="loginInfo.username"
          autofocus
          class="mt-32 h-40 items-center"
          placeholder="请输入用户名"
          :maxlength="20"
        >
          <template #prefix>
            <i class="i-fe:user mr-12 opacity-20" />
          </template>
        </n-input>
        <n-input
          v-model:value="loginInfo.password"
          class="mt-20 h-40 items-center"
          type="password"
          show-password-on="mousedown"
          placeholder="请输入密码"
          :maxlength="20"
          @keydown.enter="handleLogin()"
        >
          <template #prefix>
            <i class="i-fe:lock mr-12 opacity-20" />
          </template>
        </n-input>
        <n-input
          v-model:value="loginInfo.repassword"
          class="mt-20 h-40 items-center"
          type="password"
          show-password-on="mousedown"
          placeholder="请再次输入密码"
          :maxlength="20"
          @keydown.enter="handleLogin()"
        >
          <template #prefix>
            <i class="i-fe:lock mr-12 opacity-20" />
          </template>
        </n-input>

        <div class="mt-20 flex items-center">
          <n-button
            class="ml-32 h-40 flex-1 rounded-5 text-16"
            type="primary"
            :loading="loading"
            @click="handleLogin()"
          >
            注册
          </n-button>
        </div>
      </div>
    </div>

    <TheFooter class="py-12" />
  </div>
</template>

<script setup>
import { useAuthStore } from '@/store'
import { lStorage, request, throttle } from '@/utils'
import { useStorage } from '@vueuse/core'
import api from './api'

const router = useRouter()
const title = import.meta.env.VITE_TITLE

const loginInfo = ref({
  username: '',
  password: '',
  repassword: '',
})

const captchaUrl = ref('')
const captchaId = ref('')
const initCaptcha = throttle(async () => {
  // captchaUrl.value = `${import.meta.env.VITE_AXIOS_BASE_URL}/admin/captcha?${Date.now()}`
  try {
    const response = await request.get(`/admin/captcha?${Date.now()}`, { responseType: 'blob' })

    // // 从响应头中获取验证码ID
    captchaId.value = response.headers['x-captcha-id']

    // 显示验证码图片
    const imageUrl = URL.createObjectURL(response.data)
    captchaUrl.value = imageUrl
  }
  catch (error) {
    console.error('获取验证码失败:', error)
  }
}, 500)
// initCaptcha()

const loading = ref(false)
async function handleLogin() {
  const { username, password, repassword } = loginInfo.value
  if (!username || !password)
    return $message.warning('请输入用户名和密码')
  if (password !== repassword)
    return $message.warning('两次密码不一致')
  try {
    loading.value = true
    $message.loading('正在验证，请稍后...', { key: 'login' })
    await api.register({
      type: 'adduser',
      user: { username, password: password.toString() },
    }, captchaId.value)
    router.push('/')
  }
  catch (error) {
    // 10003为验证码错误专属业务码
    if (error?.code === 10003) {
      // 为防止爆破，验证码错误则刷新验证码
      initCaptcha()
    }
    $message.destroy('login')
    console.error(error)
  }
  loading.value = false
}
</script>
