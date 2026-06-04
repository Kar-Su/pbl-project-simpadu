<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const role = ref("")
const errorMsg = ref("")
const loading = ref(false)



const handleSimpan = async () => {
  if (!role.value) {
    errorMsg.value = "Role wajib dipilih"
    return
  }

  const token = localStorage.getItem("token")
  if (!token) {
    errorMsg.value = "Token tidak ditemukan, login ulang"
    return
  }

  errorMsg.value = ""
  loading.value = true

console.log("TOKEN:", token)
  try {
const BASE_URL = 'https://be.karlearn.site'
const response = await fetch(`${BASE_URL}/api/roles`, {
  method: "POST",
  credentials: "include",
  headers: {
    "Content-Type": "application/json",
    "Accept": "application/json",
    Authorization: `Bearer ${token}`
  },
  body: JSON.stringify({
    role_name: role.value.trim()
  })
})

    let data: any = {}
    try {
      data = await response.json()
    } catch {
      data = {}
    }

    if (!response.ok) {
      throw new Error(data?.message || "Gagal menambah role")
    }

    alert("Role berhasil ditambahkan")
    router.push("/dashboard-superadmin/role")

  } catch (err: any) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

</script>

<template>
  <div class="min-h-screen bg-[#f5f7fb] p-6">

    <p class="text-sm text-gray-400 mb-2">
      Akademik > Role > Tambah Role
    </p>

    <h1 class="text-3xl font-bold text-gray-800 mb-1">
      Tambah Role
    </h1>

    <p class="text-gray-500 mb-8">
      Tambahkan Role Baru
    </p>

    <!-- ERROR -->
    <p v-if="errorMsg" class="text-red-500 text-sm mb-4">
      {{ errorMsg }}
    </p>

    <!-- ROLE -->
    <div class="mb-8">
      <label>Role</label>

      <input
        v-model="role"
        type="text"
        placeholder="Masukkan role baru..."
        class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
      />
    </div>

    <!-- BUTTON -->
    <button
      @click="handleSimpan"
      :disabled="loading"
      class="bg-green-500 text-white px-5 py-2 rounded"
    >
      Simpan
    </button>

  </div>
</template>