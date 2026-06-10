<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

const jurusanName = ref("")
const loading = ref(false)
const error = ref("")

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

const submitData = async () => {
  error.value = ""

  if (!jurusanName.value.trim()) {
    error.value = "Nama jurusan wajib diisi"
    return
  }

  loading.value = true

  try {
    const res = await fetch(`${BASE_URL}/api/jurusan`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({
        jurusan_name: jurusanName.value,
      }),
    })

    const json = await res.json()

    if (!res.ok) {
      error.value =
        json?.message || "Gagal menambahkan jurusan"
      return
    }

    router.push("/dashboard-admin/jurusan")
  } catch (err) {
    error.value = "Terjadi kesalahan jaringan"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- Breadcrumb -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Jurusan > Tambah Jurusan
    </div>

    <!-- Judul -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Jurusan
    </h1>

    <p class="text-gray-500 text-sm mt-3 mb-6">
      Tempat menambahkan jurusan
    </p>

    <!-- Card -->
    <div
      class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] max-w-4xl"
    >
      <div class="p-4 md:p-5">

        <h2 class="text-[32px] font-semibold text-[#505050] mb-5">
          Form Tambah Jurusan
        </h2>

        <!-- Error -->
        <div
          v-if="error"
          class="mb-4 rounded-lg bg-red-100 text-red-600 px-4 py-2 text-sm"
        >
          {{ error }}
        </div>

        <!-- Input -->
        <div class="mb-5">
          <label
            class="block text-[18px] text-[#505050] font-medium mb-2"
          >
            Nama Jurusan
            <span class="text-red-500">*</span>
          </label>

          <input
            v-model="jurusanName"
            type="text"
            placeholder="Isi Nama Jurusan ..."
            class="w-full h-12 px-4 rounded-lg border border-gray-500 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <!-- Button -->
        <button
          @click="submitData"
          :disabled="loading"
          class="flex items-center gap-2 bg-[#22c55e] hover:bg-[#16a34a] disabled:opacity-50 text-white px-5 py-2.5 rounded-lg font-medium shadow-sm transition"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-5 h-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M5 13l4 4L19 7"
            />
          </svg>

          {{ loading ? "Menyimpan..." : "Simpan" }}
        </button>

      </div>
    </div>

  </div>
</template>