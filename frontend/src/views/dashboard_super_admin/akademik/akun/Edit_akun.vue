<script setup lang="ts">
import { ref, onMounted } from "vue"

const email = ref("")
const nama = ref("")
const jabatan = ref("")
const password = ref("")

const BASE_URL = "http://localhost"

const getUser = async () => {
  try {
    const token = localStorage.getItem("token")

    const res = await fetch(`${BASE_URL}/api/users/1`, {
      headers: {
        "Authorization": `Bearer ${token}`
      }
    })

    const data = await res.json()

    console.log(data)

    email.value = data.data.email
    nama.value = data.data.name
    jabatan.value = data.data.role

  } catch (err) {
    console.error(err)
  }
}

const handleSimpan = async () => {
  try {
    const token = localStorage.getItem("token")

    const res = await fetch(`${BASE_URL}/api/users/update`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({
        email: email.value,
        name: nama.value,
        role: jabatan.value,
        password: password.value
      })
    })

    const data = await res.json()

    console.log(data)

    if (!res.ok) {
      alert(data.message || "Gagal update akun")
      return
    }

    alert("Akun berhasil diupdate")

  } catch (err) {
    console.error(err)
    alert("Terjadi error")
  }
}

onMounted(() => {
  getUser()
})
</script>

<template>
  <div class="min-h-screen bg-blue-50 p-6">

    <!-- Breadcrumb -->
    <p class="text-sm text-gray-500 mb-2">
      Akademik > Akun > Detail Akun
    </p>

    <!-- Header -->
    <h1 class="text-3xl font-bold text-gray-900 mb-1">
      Edit Akun
    </h1>

    <p class="text-sm text-gray-600 mb-8">
      Tambahkan Data yang Diinginkan
    </p>

    <!-- Email -->
    <div class="mb-5">
      <label class="block text-sm text-gray-700 mb-1">
        Email
      </label>

      <input
        v-model="email"
        type="text"
        class="w-96 border border-gray-400 rounded-lg px-4 py-2.5 text-sm text-gray-700 bg-white focus:outline-none focus:border-blue-400"
      />
    </div>

    <!-- Nama -->
    <div class="mb-5">
      <label class="block text-sm text-gray-700 mb-1">
        Nama
      </label>

      <input
        v-model="nama"
        type="text"
        class="w-96 border border-gray-400 rounded-lg px-4 py-2.5 text-sm text-gray-700 bg-white focus:outline-none focus:border-blue-400"
      />
    </div>

    <!-- Jabatan -->
    <div class="mb-7">
      <label class="block text-sm text-gray-700 mb-1">
        Jabatan
      </label>

      <div class="relative w-96">

        <input
          v-model="jabatan"
          type="text"
          class="w-full border border-gray-400 rounded-lg px-4 py-2.5 text-sm text-gray-700 bg-white focus:outline-none focus:border-blue-400 pr-10"
        />

        <span class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
          ▶
        </span>

      </div>
    </div>
<!-- Password -->
<div class="mb-7">
  <label class="block text-sm text-gray-700 mb-1">
    Password Baru
  </label>

  <input
    v-model="password"
    type="password"
    class="w-96 border border-gray-400 rounded-lg px-4 py-2.5 text-sm text-gray-700 bg-white focus:outline-none focus:border-blue-400"
    placeholder="Kosongkan jika tidak ingin ganti password"
  />
</div>
    <!-- Button -->
    <button
      @click="handleSimpan"
      class="flex items-center gap-2 bg-green-500 hover:bg-green-600 text-white text-sm font-semibold px-5 py-2.5 rounded-lg"
    >
      Simpan
    </button>

  </div>
</template>