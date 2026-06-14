<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()

const BASE_URL = "https://be.karlearn.site"

// Ambil kelas_id dari params route: /dashboard-admin/detail_khs/:id
const kelasId = route.params.id as string

interface Mahasiswa {
  nim: string
  nama: string
  prodi: string
}

interface KelasInfo {
  nama: string
}

const search = ref("")
const mahasiswa = ref<Mahasiswa[]>([])
const kelasInfo = ref<KelasInfo>({ nama: "" })
const loading = ref(false)
const error = ref("")

const currentPage = ref(1)
const perPage = ref(10)

const fetchMahasiswa = async () => {
  loading.value = true
  error.value = ""
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(
      `${BASE_URL}/api/khs/mahasiswa?kelas_id=${kelasId}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      }
    )
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const json = await res.json()
    // Sesuaikan field mapping dengan response API
    mahasiswa.value = (json.data ?? json) as Mahasiswa[]
    kelasInfo.value = { nama: json.kelas ?? kelasId }
  } catch (e: any) {
    error.value = "Gagal memuat data mahasiswa."
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchMahasiswa)

const filteredData = computed(() =>
  mahasiswa.value.filter(
    (item) =>
      item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.nim.toLowerCase().includes(search.value.toLowerCase())
  )
)

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredData.value.length / perPage.value))
)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredData.value.slice(start, start + perPage.value)
})

const lihatDetail = (nim: string) => {
  router.push(`/dashboard-admin/detail_khs_mahasiswa/${nim}`)
}

const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }
</script>

<template>
  <div class="p-4">

    <!-- BREADCRUMB -->
    <div class="mb-5">
      <p class="text-sm text-gray-500">
        Akademik › KHS › {{ kelasInfo.nama || kelasId }}
      </p>
      <h1 class="text-4xl font-bold text-gray-700">Kartu Hasil Studi</h1>
      <p class="text-gray-500">Data hasil studi mahasiswa</p>
    </div>

    <!-- CARD -->
    <div class="bg-white rounded-xl shadow overflow-hidden">

      <!-- HEADER -->
      <div class="bg-[#1f3c93] text-white px-5 py-4">
        <h2 class="text-2xl font-bold">Data Kartu Hasil Studi</h2>
        <p>Kelas {{ kelasInfo.nama || kelasId }}</p>
      </div>

      <!-- SEARCH -->
      <div class="p-4">
        <input
          v-model="search"
          type="text"
          placeholder="Cari mahasiswa berdasarkan NIM atau Nama"
          class="w-full border rounded-lg px-4 py-3"
        />
      </div>

      <!-- LOADING / ERROR -->
      <div v-if="loading" class="text-center py-10 text-gray-500">
        Memuat data...
      </div>

      <div v-else-if="error" class="text-center py-10 text-red-500">
        {{ error }}
        <button @click="fetchMahasiswa" class="ml-2 underline text-blue-600">Coba lagi</button>
      </div>

      <!-- TABLE -->
      <table v-else class="w-full">
        <thead>
          <tr class="text-left text-gray-700">
            <th class="px-6 py-4">No</th>
            <th>NIM</th>
            <th>Nama Mahasiswa</th>
            <th>Prodi</th>
            <th>Aksi</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(item, index) in paginatedData"
            :key="item.nim"
          >
            <td class="px-6 py-4">
              {{ (currentPage - 1) * perPage + index + 1 }}
            </td>
            <td>{{ item.nim }}</td>
            <td>{{ item.nama }}</td>
            <td>{{ item.prodi }}</td>
            <td>
              <button
                @click="lihatDetail(item.nim)"
                class="bg-[#1f3c93] text-white px-4 py-2 rounded-lg hover:opacity-90"
              >
                Lihat
              </button>
            </td>
          </tr>

          <tr v-if="paginatedData.length === 0">
            <td colspan="5" class="text-center py-10 text-gray-500">
              Data mahasiswa tidak ditemukan
            </td>
          </tr>
        </tbody>
      </table>

      <!-- FOOTER -->
      <div class="flex justify-between items-center p-5 mt-24">
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