<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const email = ref<string>('')
const password = ref<string>('')
const showPass = ref<boolean>(false)
const loading = ref<boolean>(false)
const errorMsg = ref<string>('')

const BASE_URL = 'https://be.karlearn.site'
const ENDPOINT = `${BASE_URL}/api/auth/login`

const handleLogin = async (): Promise<void> => {
  errorMsg.value = ''

  if (!email.value || !password.value) {
    errorMsg.value = 'Isi username & password'
    return
  }

  loading.value = true

  try {
    const res = await fetch(ENDPOINT, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'accept': 'application/json',
      },
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    const data = await res.json()

    console.log('LOGIN RESPONSE:', data)
    console.log('PATH ENDPOINT:', res.url)

    if (!res.ok) {
      errorMsg.value = data.message || 'Login gagal'
      errorMsg.value = "Username atau Kata Sandi salah"
      return
    }

    if (!data?.data) {
      errorMsg.value = 'Response login tidak valid'
      return
    }

    const role: string = (data.data.role_name ?? '').toLowerCase()

    const token: string = data.data.access_token
    const refreshToken: string = data.data.refresh_token

    localStorage.setItem('token', token)
    localStorage.setItem('refresh_token', refreshToken)
    localStorage.setItem('role', role)

    console.log('ROLE:', role)

    if (role === 'super-admin') {
      router.push('/dashboard-superadmin')
    } else if (
      [
        'admin-akademik',
        'admin-mahasiswa',
        'admin-keuangan',
        'admin-pegawai',
        'dummy-dosen',
        'dummy-mahasiswa',
        'tumbal'
      ].includes(role)
    ) {
      router.push('/dashboard-admin')
    } else {
      errorMsg.value = 'Role tidak dikenali: ' + role
    }

  } catch (error) {
    console.error('ERROR LOGIN:', error)
    errorMsg.value = 'Terjadi error saat login'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#dde3ec] flex items-center justify-center p-8">

    <div class="bg-white rounded-2xl shadow-2xl overflow-hidden flex w-160">

      <!-- KIRI: Gambar Kampus -->
      <div class="relative w-80 shrink-0 bg-[#2a5a8a]">
        <img
          src="@/assets/images/bg-poliban.png"
          alt="Kampus SIMPADU"
          class="w-full h-full object-cover block"
        />

        <div class="absolute inset-0 bg-linear-to-t from-black/60 via-black/10 to-transparent" />
        <!-- <p class="absolute bottom-14 left-3 text-white text-[11px] font-bold tracking-[0.15em] uppercase">
          SIMPADU
        </p> -->
       
        <p class="absolute bottom-3 left-3 text-white text-[11px] font-bold tracking-[0.15em] uppercase">
          SIMPADU
        </p>
      </div>

      <!-- KANAN: Form Login -->
      <div class="flex-1 px-8 py-9 flex flex-col justify-center">

        <!-- Logo & Heading -->
        <div class="text-center mb-6">
          <div class="w-13 h-13 bg-blue-50 rounded-full inline-flex items-center justify-center mb-3">
            <img
              src="@/assets/images/logo.png"
              alt="Logo"
              class="w-10 h-10 object-contain"
            />
          </div>

          <h1 class="text-[17px] font-bold text-gray-900 mb-1">
            Selamat Datang
          </h1>

          <p class="text-[11px] text-gray-400 leading-relaxed max-w-52.5 mx-auto">
            <strong class="text-[#1a3a7a] font-bold">Silakan login!</strong>
            Nikmati Kemudahan Sistem Autentikasi Tunggal Untuk Mengakses Semua Layanan Dengan Satu Akun
          </p>
        </div>

        <!-- Field Email -->
        <div class="mb-3">
          <label class="block text-[11px] font-semibold text-gray-500 mb-1.5">
            Email
          </label>
          <div class="relative flex items-center">
            <span class="absolute left-3 text-gray-400 pointer-events-none">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="4" width="20" height="16" rx="2"/>
                <path d="m2 7 10 7 10-7"/>
              </svg>
            </span>
            <input
              v-model="email"
              type="text"
              placeholder="SuperAdmin@sabar.id"
              class="w-full border border-gray-200 bg-gray-50 rounded-lg pl-8 pr-3 py-2.5 text-xs text-gray-800 placeholder-gray-300 outline-none focus:border-blue-500 focus:bg-white transition-colors"
            />
          </div>
        </div>

        <!-- Field Password -->
        <div class="mb-5">
          <label class="block text-[11px] font-semibold text-gray-500 mb-1.5">
            Password
          </label>
          <div class="relative flex items-center">
            <span class="absolute left-3 text-gray-400 pointer-events-none">
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
            </span>
            <input
              v-model="password"
              :type="showPass ? 'text' : 'password'"
              placeholder="Admin12345"
              class="w-full border border-gray-200 bg-gray-50 rounded-lg pl-8 pr-9 py-2.5 text-xs text-gray-800 placeholder-gray-300 outline-none focus:border-blue-500 focus:bg-white transition-colors"
              @keydown.enter="handleLogin"
            />
            <button
              type="button"
              tabindex="-1"
              class="absolute right-3 text-gray-400 hover:text-gray-600 transition-colors"
              @click="showPass = !showPass"
            >
              <!-- Mata terbuka -->
              <svg
                v-if="showPass"
                class="w-3.5 h-3.5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>

              <!-- Mata tertutup -->
              <svg
                v-else
                class="w-3.5 h-3.5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Pesan Error -->
        <p
          v-if="errorMsg"
          class="text-[11px] text-red-500 font-semibold text-center mb-3 -mt-2"
        >
          {{ errorMsg }}
        </p>

        <!-- Tombol Login -->
        <button
          type="button"
          :disabled="loading"
          class="w-full bg-[#1a5fd4] hover:bg-[#1550b8] disabled:opacity-60 disabled:cursor-not-allowed text-white text-xs font-bold py-3 rounded-lg transition-colors active:scale-[0.98]"
          @click="handleLogin"
        >
          {{ loading ? 'Memproses...' : 'Masuk' }}
        </button>

      </div>
    </div>
  </div>
</template>