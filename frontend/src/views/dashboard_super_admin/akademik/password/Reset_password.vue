<script setup lang="ts">
import { ref } from "vue"
import KonfirmasiReset from "./Konfirmasi_reset.vue"

const email = ref("")
const password = ref("")
const showPopup = ref(false)
const passwordError = ref("")
const passwordKonfirmasi = ref("")
const token = localStorage.getItem("token")

const validatePassword = (pwd: string) => {
  // const hasUpper = /[A-Z]/.test(pwd)
  // const hasNumber = /[0-9]/.test(pwd)
  // const hasSymbol = /[^A-Za-z0-9]/.test(pwd)

  if (pwd.length < 8) {
    return "Password minimal 8 karakter"
  } 
  // else if (!hasUpper || !hasNumber || !hasSymbol) {
  //   return "Password harus ada huruf besar, angka, dan simbol"
  // }

  return ""
}

const handleOpenPopup = () => {

  passwordError.value = validatePassword(password.value)

  if (passwordError.value) return

  if (password.value !== passwordKonfirmasi.value) {
    passwordError.value = "Konfirmasi password tidak cocok"
    return
  }

  showPopup.value = true
}

const handleResetPassword = async () => {

  passwordError.value = validatePassword(password.value)

  if (passwordError.value) return

  if (password.value !== passwordKonfirmasi.value) {
    passwordError.value = "Konfirmasi password tidak cocok"
    return
  }

  try {

    const res = await fetch(`/api/auth/reset-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        email: email.value,
        new_password: password.value
      })
    })

    let data: any = {}

    try {
      data = await res.json()
    } catch {
      console.log("Response bukan JSON")
    }

    console.log("RESET RESPONSE:", data)

    if (!res.ok) {
      alert(data.message || "Gagal reset password")
      return
    }

    alert("Password berhasil direset")

    showPopup.value = false

    email.value = ""
    password.value = ""
    passwordKonfirmasi.value = ""

  } catch (err) {
    console.error(err)
    alert("Terjadi error")
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-[#1d1d1d]">
      Reset Password
    </h1>

    <p class="text-gray-500 mt-2 mb-10">
      Ubah password
    </p>

    <!-- FORM -->
    <div class="w-full max-w-xl">

      <!-- EMAIL -->
      <div class="mb-6">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Email
        </label>

        <input
          v-model="email"
          type="email"
          placeholder="Isi Email..."
          class="w-full rounded-xl border border-gray-400 bg-white px-4 py-3 text-sm focus:outline-none focus:border-[#243e90]"
        />

      </div>

      <!-- PASSWORD -->
      <div class="mb-6">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Password Baru
        </label>

        <input
          v-model="password"
          type="text"
          autocomplete="off"
          autocapitalize="off"
          spellcheck="false"
          placeholder="Isi Password baru..."
          class="w-full rounded-xl border border-gray-400 bg-white px-4 py-3 text-sm focus:outline-none focus:border-[#243e90]"
        />

       <p
  v-if="passwordError"
  class="mt-2 text-sm text-red-500"
>
  {{ passwordError }}
</p>

        <p
          v-if="passwordError"
          class="mt-1 text-sm text-red-500"
        >
          {{ passwordError }}
        </p>

      </div>

      <!-- KONFIRMASI -->
      <div class="mb-8">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Konfirmasi Password
        </label>

        <input
          v-model="passwordKonfirmasi"
          type="text"
          autocomplete="off"
          autocapitalize="off"
          spellcheck="false"
          placeholder="Isi Konfirmasi Password..."
          class="w-full rounded-xl border border-gray-400 bg-white px-4 py-3 text-sm focus:outline-none focus:border-[#243e90]"
        />

      </div>

      <!-- BUTTON -->
      <button
        @click="handleOpenPopup"
        class="flex items-center gap-2 rounded-lg bg-[#22c55e] px-5 py-3 text-sm font-semibold text-white hover:bg-[#16a34a] transition-all"
      >

        <!-- ICON -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M5 13l4 4L19 7"
          />
        </svg>

        Simpan

      </button>

    </div>

    <!-- POPUP -->
    <KonfirmasiReset
      v-if="showPopup"
      :email="email"
      @close="showPopup = false"
      @confirm="handleResetPassword"
    />

  </div>
</template>