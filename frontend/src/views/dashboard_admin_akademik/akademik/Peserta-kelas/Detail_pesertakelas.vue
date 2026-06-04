<script setup lang="ts">
import { ref, onMounted } from "vue"

// ─────────────────────────────────────────────
// HELPER HEADER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ─────────────────────────────────────────────
// INTERFACE
// ─────────────────────────────────────────────
interface KelasItem {
  id: number
  nama_kelas: string
}

interface MahasiswaItem {
  id: number
  nama: string
  nim: string
}

// ─────────────────────────────────────────────
// STATE
// ─────────────────────────────────────────────
const kelasList = ref<KelasItem[]>([])
const mahasiswaList = ref<MahasiswaItem[]>([])

const selectedKelas = ref("")
const selectedMahasiswa = ref("")

const selectedMahasiswaList = ref<number[]>([])

// ─────────────────────────────────────────────
// HIT API KELAS
// Endpoint : GET /api/kelas
// ─────────────────────────────────────────────
const getKelas = async (): Promise<void> => {
  try {
    const BASE_URL = 'https://be.karlearn.site'
    const res = await fetch(`${BASE_URL}/api/kelas`, {
      headers: getHeaders(),
    })

    const data = await res.json()

    kelasList.value = data.data ?? []
  } catch (err) {
    console.error("getKelas:", err)
  }
}

// ─────────────────────────────────────────────
// HIT API MAHASISWA
// Endpoint : GET /api/mahasiswa
// ─────────────────────────────────────────────
const getMahasiswa = async (): Promise<void> => {
  try {
    const res = await fetch("/api/mahasiswa", {
      headers: getHeaders(),
    })

    const data = await res.json()

    mahasiswaList.value = data.data ?? []
  } catch (err) {
    console.error("getMahasiswa:", err)
  }
}

// ─────────────────────────────────────────────
// SUBMIT
// Endpoint : POST /api/peserta-kelas
// body:
// {
//   kelas_id: number,
//   mahasiswa_ids: number[]
// }
// ─────────────────────────────────────────────
const handleSubmit = async (): Promise<void> => {
  try {
    const payload = {
      kelas_id: Number(selectedKelas.value),
      mahasiswa_ids: selectedMahasiswaList.value,
    }

    const res = await fetch("/api/peserta-kelas", {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })

    const data = await res.json()

    console.log(data)
  } catch (err) {
    console.error("handleSubmit:", err)
  }
}

// ─────────────────────────────────────────────
// CHECKBOX
// ─────────────────────────────────────────────
const toggleMahasiswa = (id: number): void => {
  if (selectedMahasiswaList.value.includes(id)) {
    selectedMahasiswaList.value =
      selectedMahasiswaList.value.filter((item) => item !== id)
  } else {
    selectedMahasiswaList.value.push(id)
  }
}

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted((): void => {
  getKelas()
  getMahasiswa()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
      <span>Mahasiswa</span>
      <span>›</span>
      <span>Kelas</span>
      <span>›</span>
      <span class="text-gray-700">Tambah Peserta Kelas</span>
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold leading-none text-[#333]">
      Detail Peserta Kelas
    </h1>

    <p class="mt-3 text-gray-500">
      Pengelolaan Data peserta kelas
    </p>

    <!-- CARD -->
    <div
      class="mt-8 w-full max-w-[980px] rounded-2xl border border-[#d8e1f0] bg-white p-5 shadow-sm"
    >

      <!-- HEADER -->
      <h2 class="mb-6 text-[30px] font-bold text-[#444]">
        Form Kelas
      </h2>

      <!-- FORM -->
      <div class="space-y-5">

        <!-- KELAS -->
        <div>
          <label class="mb-2 block text-[15px] font-medium text-[#555]">
            Kelas
          </label>

          <div class="relative w-[460px]">
            <select
              v-model="selectedKelas"
              class="h-[54px] w-full appearance-none rounded-xl border border-[#bfc8d8] bg-white px-4 pr-10 text-[15px] text-gray-700 outline-none focus:border-[#2447a8]"
            >
              <option value="">Pilih Kelas</option>

              <option
                v-for="item in kelasList"
                :key="item.id"
                :value="item.id"
              >
                {{ item.nama_kelas }}
              </option>
            </select>

            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="2"
              stroke="currentColor"
              class="pointer-events-none absolute right-4 top-1/2 size-5 -translate-y-1/2 text-gray-500"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m19.5 8.25-7.5 7.5-7.5-7.5"
              />
            </svg>
          </div>
        </div>

        <!-- MAHASISWA -->
        <div>
          <label class="mb-2 block text-[15px] font-medium text-[#555]">
            Mahasiswa
          </label>

          <div class="relative w-[460px]">
            <select
              v-model="selectedMahasiswa"
              class="h-[54px] w-full appearance-none rounded-xl border border-[#bfc8d8] bg-white px-4 pr-10 text-[15px] text-gray-700 outline-none focus:border-[#2447a8]"
            >
              <option value="">Pilih Mahasiswa</option>

              <option
                v-for="item in mahasiswaList"
                :key="item.id"
                :value="item.id"
              >
                {{ item.nama }}
              </option>
            </select>

            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="2"
              stroke="currentColor"
              class="pointer-events-none absolute right-4 top-1/2 size-5 -translate-y-1/2 text-gray-500"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m19.5 8.25-7.5 7.5-7.5-7.5"
              />
            </svg>
          </div>
        </div>

        <!-- CHECKBOX -->
        <div class="space-y-4 pt-2">

          <div
            v-for="item in mahasiswaList"
            :key="item.id"
            class="flex items-center gap-4"
          >
            <input
              type="checkbox"
              :checked="selectedMahasiswaList.includes(item.id)"
              @change="toggleMahasiswa(item.id)"
              class="h-6 w-6 rounded border-gray-400 text-[#2447a8] focus:ring-[#2447a8]"
            />

            <p class="text-[17px] font-medium text-[#444]">
              {{ item.nama }} / {{ item.nim }}
            </p>
          </div>

        </div>

        <!-- BUTTON -->
        <div class="pt-6">
          <button
            @click="handleSubmit"
            class="rounded-xl bg-[#22c55e] px-8 py-3 text-[16px] font-semibold text-white transition hover:bg-[#16a34a]"
          >
            Simpan
          </button>
        </div>
        

      </div>
    </div>
  </div>
</template>