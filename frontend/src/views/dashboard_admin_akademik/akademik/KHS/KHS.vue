<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

interface KhsClass {
  id: number
  kelas: string
  tahunAkademik: string
  semester: string
  prodi: string
}

const search = ref("")
const selectedSemester = ref("")
const selectedProdi = ref("")
const semesterFilter = ref("")
const prodiFilter = ref("")

const kelasList = ref<KhsClass[]>([])
const loading = ref(false)
const error = ref("")

const currentPage = ref(1)
const perPage = ref(10)

const fetchKelas = async () => {
  loading.value = true
  error.value = ""
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(`${BASE_URL}/api/khs/kelas`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const json = await res.json()
    // Sesuaikan field mapping dengan response API
    kelasList.value = (json.data ?? json) as KhsClass[]
  } catch (e: any) {
    error.value = "Gagal memuat data kelas."
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchKelas)

const applyFilter = () => {
  semesterFilter.value = selectedSemester.value
  prodiFilter.value = selectedProdi.value
  currentPage.value = 1
}

const filteredData = computed(() =>
  kelasList.value.filter((item) => {
    const matchSearch = item.kelas
      .toLowerCase()
      .includes(search.value.toLowerCase())
    const matchSemester =
      !semesterFilter.value || item.semester === semesterFilter.value
    const matchProdi =
      !prodiFilter.value || item.prodi === prodiFilter.value
    return matchSearch && matchSemester && matchProdi
  })
)

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredData.value.length / perPage.value))
)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredData.value.slice(start, start + perPage.value)
})

const goDetail = (id: number) => {
  router.push(`/dashboard-admin/detail_khs/${id}`)
}

const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }
</script>

<template>
  <div class="p-4">

    <!-- HEADER -->
    <div class="mb-6">
      <p class="text-gray-500 text-sm">Akademik › KHS</p>
      <h1 class="text-4xl font-bold text-gray-700">Kartu Hasil Studi</h1>
      <p class="text-gray-500">Data hasil studi mahasiswa</p>
    </div>

    <!-- FILTER -->
    <div class="bg-white rounded-xl shadow overflow-hidden mb-5">
      <div class="bg-[#1f3c93] text-white px-5 py-4">
        <h2 class="text-2xl font-bold">Saring Data</h2>
        <p>Memilah data berdasarkan kurikulum, matakuliah, dan tahun akademik</p>
      </div>

      <div class="p-4 flex flex-col lg:flex-row gap-3">
        <select v-model="selectedSemester" class="flex-1 border rounded-lg px-4 py-3">
          <option value="">Pilih Semester</option>
          <option value="Semester 1">Semester 1</option>
          <option value="Semester 2">Semester 2</option>
        </select>

        <select v-model="selectedProdi" class="flex-1 border rounded-lg px-4 py-3">
          <option value="">Pilih Prodi</option>
          <option value="Teknik Informatika">Teknik Informatika</option>
          <option value="Teknik Elektro">Teknik Elektro</option>
        </select>

        <button
          @click="applyFilter"
          class="bg-[#1f3c93] text-white px-6 rounded-lg"
        >
          Terapkan
        </button>
      </div>
    </div>

    <!-- TABLE -->
    <div class="bg-white rounded-xl shadow overflow-hidden">
      <div class="bg-[#1f3c93] text-white px-5 py-4">
        <h2 class="text-2xl font-bold">Data KHS</h2>
        <p>Kumpulan data hasil studi mahasiswa tiap semester</p>
      </div>

      <!-- SEARCH -->
      <div class="p-4">
        <input
          v-model="search"
          type="text"
          placeholder="Cari Kelas yang diinginkan. contoh(TI-4A)"
          class="w-full border rounded-lg px-4 py-3"
        />
      </div>

      <!-- LOADING / ERROR -->
      <div v-if="loading" class="text-center py-10 text-gray-500">
        Memuat data...
      </div>

      <div v-else-if="error" class="text-center py-10 text-red-500">
        {{ error }}
        <button @click="fetchKelas" class="ml-2 underline text-blue-600">Coba lagi</button>
      </div>

      <table v-else class="w-full">
        <thead>
          <tr class="text-left text-gray-700">
            <th class="px-6 py-4">No</th>
            <th>Nama Kelas</th>
            <th>Tahun Akademik</th>
            <th>Semester</th>
            <th>Prodi</th>
            <th>Aksi</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(item, index) in paginatedData"
            :key="item.id"
          >
            <td class="px-6 py-4">
              {{ (currentPage - 1) * perPage + index + 1 }}
            </td>
            <td>{{ item.kelas }}</td>
            <td>{{ item.tahunAkademik }}</td>
            <td>{{ item.semester }}</td>
            <td>{{ item.prodi }}</td>
            <td>
              <button
                @click="goDetail(item.id)"
                class="bg-[#1f3c93] text-white px-4 py-2 rounded-lg"
              >
                Lihat
              </button>
            </td>
          </tr>

          <tr v-if="paginatedData.length === 0">
            <td colspan="6" class="text-center py-10 text-gray-500">
              Data tidak ditemukan
            </td>
          </tr>
        </tbody>
      </table>

      <!-- PAGINATION -->
      <div class="flex justify-between items-center p-5 mt-20">
        <select v-model="perPage" class="border rounded-lg px-3 py-2">
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
          <option :value="50">50 Baris</option>
        </select>

        <div class="flex items-center gap-4 text-gray-500">
          <button @click="prevPage" :disabled="currentPage === 1">← Previous</button>
          <button class="w-8 h-8 rounded bg-[#1f3c93] text-white">{{ currentPage }}</button>
          <button @click="nextPage" :disabled="currentPage === totalPages">Next →</button>
        </div>
      </div>
    </div>

  </div>
</template>