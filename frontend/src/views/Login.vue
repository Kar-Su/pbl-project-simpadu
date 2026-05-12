<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import image from '@/assets/images/image.png'
import logo from '@/assets/images/logo.png'
  
const router = useRouter()

const username = ref('')
const password = ref('')
const showPass = ref(false)

// const ENDPOINT = "/api/auth/login"
const ENDPOINT = "/api/auth/login"

const handleLogin = async () => {
  if (!username.value || !password.value) {
    alert("Isi username & password")
    return
  }

  try {
    const res = await fetch(`${ENDPOINT}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "accept": "application/json",
      },
      body: JSON.stringify({
        email: username.value,
        password: password.value,
      }),
    })

    const data = await res.json()

    console.log("LOGIN RESPONSE:", data)

    if (!res.ok) {
      alert(data.message || "Login gagal")
      return
    }

    const token = data.data.access_token
    const refreshToken = data.data.refresh_token

    const role = (
      data.data.role_name ||
      data.data.role ||
      data.data.user?.role_name ||
      data.data.user?.role ||
      ""
    ).toLowerCase()

    localStorage.setItem("token", token)
    localStorage.setItem("refresh_token", refreshToken)
    localStorage.setItem("role", role)

    console.log("ROLE:", role)

    if (role === "super-admin") {
      router.push("/dashboard-superadmin")
    } else if (
      ["admin-akademik","admin-mahasiswa","admin-keuangan",
       "admin-pegawai","dummy-dosen","dummy-mahasiswa","tumbal"]
      .includes(role)
    ) {
      router.push("/dashboard-admin")
    } else {
      alert("Role tidak dikenali: " + role)
    }
  } catch (error) {
    console.error("ERROR LOGIN:", error)
    alert("Terjadi error saat login")
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-200 flex items-center justify-center">
    
    <div class="bg-white rounded-xl shadow-lg overflow-hidden flex w-120">

      <!-- KIRI -->
      <div class="relative w-48 shrink-0">
        <img
          :src="image"
          class="w-full h-full object-cover"
          alt="Kampus"
        />
        <div class="absolute inset-0 bg-transparent"></div>

        <p class="absolute bottom-3 left-3 text-left text-white text-xs font-bold tracking-widest uppercase">
          POLIBAN
        </p>
      </div>

      <!-- KANAN -->
      <div class="flex-1 px-7 py-8 flex flex-col justify-center">

        <!-- LOGO -->
        <div class="text-center mb-5">
          <img
            :src="logo"
            class="mx-auto w-10 h-10 object-contain mb-2"
          />

          <h2 class="text-base font-bold text-gray-800">
            Selamat Datang
          </h2>

          <p class="text-[11px] text-gray-400 tracking-widest uppercase mt-0.5">
            Silakan Login
          </p>
        </div>

        <!-- Username -->
        <div class="mb-3">
          <label class="block text-xs font-semibold text-gray-600 mb-1">
            Username
          </label>

          <input
            v-model="username"
            type="text"
            class="w-full border border-gray-300 rounded px-3 py-2 text-xs outline-none"
            placeholder="Email"
          />
        </div>

       <!-- Password -->
<div class="mb-5">
  <label class="block text-xs font-semibold text-gray-600 mb-1">
    Password
  </label>

  <div class="relative">
    <input
      v-model="password"
      :type="showPass ? 'text' : 'password'"
      class="w-full border border-gray-300 rounded px-3 py-2 text-xs outline-none pr-9"
      placeholder="Password"
      @keydown.enter="handleLogin"
    />

    <button
      type="button"
      @click="showPass = !showPass"
      class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
      tabindex="-1"
    >
      <!-- Mata terbuka -->
      <svg v-if="showPass" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none"
        viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        <path stroke-linecap="round" stroke-linejoin="round"
          d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
      </svg>

      <!-- Mata tertutup -->
      <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none"
        viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round"
          d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 4.411m0 0L21 21" />
      </svg>
    </button>
  </div>
</div>

        <button
          @click="handleLogin"
          class="w-full bg-[#1a5fd4] hover:bg-[#1550b8] text-white text-xs font-semibold py-2.5 rounded"
        >
          Masuk
        </button>

      </div>
    </div>
  </div>
</template>

<style scoped>
</style>