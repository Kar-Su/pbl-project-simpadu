<script setup lang="ts">
import { ref } from "vue"
import KonfirmasiReset from "./Konfirmasi_reset.vue"

const email = ref("")
const password = ref("")
const showPopup = ref(false)
const passwordError = ref("")

const validatePassword = (pwd: string) => {
  const hasUpper = /[A-Z]/.test(pwd)
  const hasNumber = /[0-9]/.test(pwd)
  const hasSymbol = /[^A-Za-z0-9]/.test(pwd)
  if (!hasUpper || !hasNumber || !hasSymbol) {
    return "Password harus ada huruf besar, angka, dan simbol"
  }
  return ""
}

const handleOpenPopup = () => {
  passwordError.value = validatePassword(password.value)
  if (passwordError.value) return
  showPopup.value = true
}

const handleResetPassword = async () => {
  passwordError.value = validatePassword(password.value)
  if (passwordError.value) return

  const token = localStorage.getItem("token")
  if (!token) {
    alert("Token tidak ditemukan, silakan login ulang")
    return
  }

  try {
    const res = await fetch(`/api/auth/reset-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "accept": "application/json",
        "Authorization": `Bearer ${token}`
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

  } catch (err) {
    console.error(err)
    alert("Terjadi error")
  }
}
</script>

<template>

  <!-- Email -->
  <div class="mb-4">
    <label class="block text-sm text-gray-700 mb-1">
      Email
    </label>

    <input
      v-model="email"
      type="email"
      placeholder="Isi Email..."
      class="w-80 border border-gray-300 rounded px-3 py-2 text-sm bg-white focus:outline-none focus:border-blue-400"
    />
  </div>

  <!-- Password -->
  <div class="mb-6">
    <label class="block text-sm text-gray-700 mb-1">
      Password Baru
    </label>

    <input
      v-model="password"
      type="text"
      autocomplete="off"
      autocapitalize="off"
      spellcheck="false"
      placeholder="Isi Password..."
      class="w-80 border border-gray-300 rounded px-3 py-2 text-sm bg-white focus:outline-none focus:border-blue-400"
    />

    <p v-if="passwordError" class="text-red-500 text-xs mt-1">
      {{ passwordError }}
    </p>
  </div>

  <!-- BUTTON -->
<button
  @click="handleOpenPopup"
  class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded text-sm"
>
  Simpan
</button>

  <!-- POPUP -->
  <KonfirmasiReset
    v-if="showPopup"
    :email="email"
    @close="showPopup = false"
    @confirm="handleResetPassword"
  />

</template>