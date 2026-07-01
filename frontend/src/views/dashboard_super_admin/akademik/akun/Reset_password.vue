<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import KonfirmasiReset from "./Konfirmasi_reset.vue"

// ================= STATE =================
const route = useRoute()
const router = useRouter()

const email = ref("")
const password = ref("")
const passwordKonfirmasi = ref("")
const passwordError = ref("")
const showPopup = ref(false)

const token = localStorage.getItem("token")

// ================= INIT =================
onMounted(() => {
  email.value = decodeURIComponent(route.params.email as string) || ""
})

// ================= VALIDASI =================
const validatePassword = (pwd: string): string => {
  if (pwd.length < 8) {
    return "Password minimal 8 karakter"
  }
  return ""
}

// ================= OPEN POPUP =================
const handleOpenPopup = () => {
  passwordError.value = validatePassword(password.value)
  if (passwordError.value) return

  if (password.value !== passwordKonfirmasi.value) {
    passwordError.value = "Konfirmasi password tidak cocok"
    return
  }

  showPopup.value = true
}

// ================= SUBMIT RESET =================
const handleResetPassword = async () => {
  passwordError.value = validatePassword(password.value)
  if (passwordError.value) return

  if (password.value !== passwordKonfirmasi.value) {
    passwordError.value = "Konfirmasi password tidak cocok"
    return
  }

  try {
    const BASE_URL = "https://be.karlearn.site"

    const res = await fetch(`${BASE_URL}/api/auth/reset-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        email: email.value,
        new_password: password.value,
      }),
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
      showPopup.value = false
      return
    }

    alert("Password berhasil direset")
    showPopup.value = false
    router.push("/dashboard-superadmin/akun")

  } catch (err) {
    console.error(err)
    alert("Terjadi error jaringan")
  }
}
</script>

<template>
  <div class="h-full bg-[#f5f7fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-black-800">
      Akademik > Akun > Reset Password
    </div>

    <!-- TITLE -->
    <h1 class="mb-1 text-3xl font-bold text-black-800">
      Reset Password
    </h1>

    <p class="mb-10 text-gray-500">
      Ubah password akun
    </p>

    <!-- FORM -->
    <div class="w-full max-w-xl">

      <!-- EMAIL (read-only) -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Email
        </label>
        <input
          v-model="email"
          type="email"
          readonly
          class="w-full rounded-xl border border-gray-300 bg-gray-100 px-4 py-3 text-sm text-gray-500 cursor-not-allowed focus:outline-none"
        />
      </div>

      <!-- PASSWORD BARU -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Password Baru
        </label>
        <input
          v-model="password"
          type=""
          autocomplete="new-password"
          placeholder="Isi password baru..."
          class="w-full rounded-xl border border-gray-400 bg-white px-4 py-3 text-sm focus:outline-none focus:border-[#243e90]"
        />
        <p v-if="passwordError" class="mt-2 text-sm text-red-500">
          {{ passwordError }}
        </p>
      </div>

      <!-- KONFIRMASI PASSWORD -->
      <div class="mb-8">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Konfirmasi Password
        </label>
        <input
          v-model="passwordKonfirmasi"
          type=""
          autocomplete="new-password"
          placeholder="Isi konfirmasi password..."
          class="w-full rounded-xl border border-gray-400 bg-white px-4 py-3 text-sm focus:outline-none focus:border-[#243e90]"
        />
      </div>

      <!-- BUTTONS -->
      <div class="flex gap-3">

        <!-- KEMBALI -->
        <button
          @click="router.push('/dashboard-superadmin/akun')"
          class="rounded-lg bg-gray-400 px-5 py-3 text-sm font-semibold text-white hover:bg-gray-600 transition-all"
        >
          ← Kembali
        </button>

        <!-- SIMPAN -->
        <button
          @click="handleOpenPopup"
          class="flex items-center gap-2 rounded-lg bg-[#22c55e] px-5 py-3 text-sm font-semibold text-white hover:bg-[#16a34a] transition-all"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          Simpan
        </button>

      </div>

    </div>

    <!-- POPUP KONFIRMASI -->
    <KonfirmasiReset
      v-if="showPopup"
      :email="email"
      @close="showPopup = false"
      @confirm="handleResetPassword"
    />

  </div>
</template>