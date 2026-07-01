<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"


const router = useRouter()


const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") || ""}`,
})

// ================= TYPES =================

interface NilaiItem {
  kode_mk: string
  nama_mk: string
  sks: number
  nilai: number
  grade: string
}

interface KhsItem {
  mahasiswa_name: string
  semester: number
  prodi_name: string
  kelas_name: string
  tahun_akademik: number
  kurikulum_name: string
  ips: number
  ipk: number
  nilai: NilaiItem[]
}

// ================= STATE =================

const allData = ref<KhsItem[]>([])
const loading = ref(false)
const error = ref("")

// ================= FILTER =================

const search = ref("")
const filterSemester = ref("")
const filterProdi = ref("")
const filterTahun = ref("")

// ================= PAGINATION =================

const currentPage = ref(1)
const perPage = ref(10)

// ================= FETCH =================

const fetchKHS = async () => {
  loading.value = true
  error.value = ""

  try {
    const token = localStorage.getItem("token")

    console.log("TOKEN =", token)

    const response = await fetch(
      `${BASE_URL}/api/khs?semester=1`,
      {
        method: "GET",
        headers: {
          // Accept: "application/json",
          Authorization: `Bearer ${token}`,
        },
      }
    )

    console.log("STATUS =", response.status)

    const result = await response.json()

    console.log("RESULT =", result)

    if (!result.success) {
      error.value = result.message || "API gagal"
      return
    }

    allData.value = result.data || []

    console.log(
      "JUMLAH DATA =",
      allData.value.length
    )
  } catch (err) {
    console.error(err)
    error.value = "Gagal memuat data"
  } finally {
    loading.value = false
  }
}

onMounted(fetchKHS)

// ================= OPTIONS =================

const semesterOptions = computed(() => {
  const set = new Set(
    allData.value.map(
      (item) => item.semester
    )
  )

  return Array.from(set).sort(
    (a, b) => a - b
  )
})

const prodiOptions = computed(() => {
  const set = new Set(
    allData.value.map(
      (item) => item.prodi_name
    )
  )

  return Array.from(set).sort()
})

const tahunOptions = computed(() => {
  const set = new Set(
    allData.value.map((item) =>
      String(item.tahun_akademik)
    )
  )

  return Array.from(set).sort()
})

// ================= HELPER =================

const fmt = (text: string) =>
  (text ?? "-").replace(/-/g, " ")

// ================= FILTER DATA =================

const filteredData = computed(() => {
  return allData.value.filter(
    (item) => {
      const keyword =
        search.value.toLowerCase()

      const matchSearch =
        !search.value ||
        item.mahasiswa_name
          .toLowerCase()
          .includes(keyword) ||
        item.kelas_name
          .toLowerCase()
          .includes(keyword)

      const matchSemester =
        !filterSemester.value ||
        String(item.semester) ===
          filterSemester.value

      const matchProdi =
        !filterProdi.value ||
        item.prodi_name ===
          filterProdi.value

      const matchTahun =
        !filterTahun.value ||
        String(
          item.tahun_akademik
        ) === filterTahun.value

      return (
        matchSearch &&
        matchSemester &&
        matchProdi &&
        matchTahun
      )
    }
  )
})

// ================= PAGINATION =================

watch(
  [
    filterSemester,
    filterProdi,
    filterTahun,
    search,
    perPage,
  ],
  () => {
    currentPage.value = 1
  }
)

const totalPages = computed(() =>
  Math.max(
    1,
    Math.ceil(
      filteredData.value.length /
        perPage.value
    )
  )
)

const paginatedData = computed(() => {
  const start =
    (currentPage.value - 1) *
    perPage.value

  return filteredData.value.slice(
    start,
    start + perPage.value
  )
})

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (
    currentPage.value <
    totalPages.value
  ) {
    currentPage.value++
  }
}

const visiblePages = computed(() => {
  const total =
    totalPages.value

  const current =
    currentPage.value

  if (total <= 5) {
    return Array.from(
      { length: total },
      (_, i) => i + 1
    )
  }

  const pages: (
    | number
    | string
  )[] = [1]

  const start = Math.max(
    2,
    current - 1
  )

  const end = Math.min(
    total - 1,
    current + 1
  )

  if (start > 2) {
    pages.push("...")
  }

  for (
    let i = start;
    i <= end;
    i++
  ) {
    pages.push(i)
  }

  if (end < total - 1) {
    pages.push("...")
  }

  pages.push(total)

  return pages
})

// ================= RESET =================

const resetFilter = () => {
  search.value = ""
  filterSemester.value = ""
  filterProdi.value = ""
  filterTahun.value = ""
  currentPage.value = 1
}

// ================= DETAIL =================

const goDetail = (
  item: KhsItem
) => {
  sessionStorage.setItem(
    "selectedKhs",
    JSON.stringify(item)
  )

  router.push(
    `/dashboard-admin/detail_khs/${encodeURIComponent(
      item.mahasiswa_name
    )}`
  )
}
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB + TITLE -->
    <div class="mb-6">
      <p class="text-gray-500 text-sm">Akademik › KHS</p>
      <h1 class="text-4xl font-bold text-gray-700">Kartu Hasil Studi</h1>
      <p class="text-gray-500 text-sm mt-1">Data hasil studi mahasiswa</p>
    </div>

    <!-- FILTER CARD -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden mb-5">
      <div class="bg-[#243e90] px-5 py-4">
        <h2 class="text-white text-xl font-bold">Saring Data</h2>
        <p class="text-white text-sm mt-1">Memilah data berdasarkan semester, prodi, dan tahun akademik</p>
      </div>

      <div class="p-5 flex flex-wrap gap-3 items-end">
        <!-- Semester -->
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-600">Semester</label>
          <select v-model="filterSemester"
            class="border border-gray-300 rounded-xl px-4 py-2.5 text-sm outline-none focus:border-blue-500 bg-white min-w-[160px]">
            <option value="">Semua Semester</option>
            <option v-for="s in semesterOptions" :key="s" :value="String(s)">Semester {{ s }}</option>
          </select>
        </div>

        <!-- Prodi -->
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-600">Prodi</label>
          <select v-model="filterProdi"
            class="border border-gray-300 rounded-xl px-4 py-2.5 text-sm outline-none focus:border-blue-500 bg-white min-w-[180px]">
            <option value="">Semua Prodi</option>
            <option v-for="p in prodiOptions" :key="p" :value="p" class="capitalize">{{ fmt(p) }}</option>
          </select>
        </div>

        <!-- Tahun Akademik -->
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-600">Tahun Akademik</label>
          <select v-model="filterTahun"
            class="border border-gray-300 rounded-xl px-4 py-2.5 text-sm outline-none focus:border-blue-500 bg-white min-w-[160px]">
            <option value="">Semua Tahun</option>
            <option v-for="t in tahunOptions" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <button @click="resetFilter"
          class="px-5 py-2.5 bg-gray-400 hover:bg-gray-500 text-white rounded-xl text-sm font-semibold">
          Reset
        </button>
      </div>
    </div>

    <!-- TABLE CARD -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">
      <div class="bg-[#243e90] px-5 py-4">
        <h2 class="text-white text-xl font-bold">Data KHS</h2>
        <p class="text-white text-sm mt-1">Kumpulan data hasil studi mahasiswa tiap semester</p>
      </div>

      <!-- SEARCH -->
      <div class="px-5 pt-4 pb-2">
        <div class="relative">
          <input v-model="search" type="text" placeholder="Cari nama mahasiswa atau nama kelas..."
            class="w-full border border-gray-300 rounded-xl px-4 py-3 pr-10 text-sm outline-none focus:border-blue-500 bg-white" />
          <svg xmlns="http://www.w3.org/2000/svg" class="absolute right-3 top-3.5 h-5 w-5 text-gray-400" fill="none"
            viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-4.35-4.35M11 18a7 7 0 100-14 7 7 0 000 14z" />
          </svg>
        </div>
      </div>

      <!-- LOADING / ERROR -->
      <div v-if="loading" class="text-center py-14 text-gray-500">Memuat data...</div>
      <div v-else-if="error" class="text-center py-14 text-red-500">
        {{ error }}
        <button @click="fetchKHS" class="ml-2 underline text-blue-600">Coba lagi</button>
      </div>

      <!-- TABLE -->
      <div v-else class="px-5 pb-3 overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-gray-600 border-b border-gray-300 text-sm">
              <th class="py-4 px-2 text-center w-12">No</th>
              <th class="py-4 px-2 text-left">Nama Mahasiswa</th>
              <th class="py-4 px-2 text-left">Kelas</th>
              <th class="py-4 px-2 text-left">Prodi</th>
              <th class="py-4 px-2 text-center">Semester</th>
              <th class="py-4 px-2 text-center">Tahun Akademik</th>
              <th class="py-4 px-2 text-center">IPS</th>
              <th class="py-4 px-2 text-center">IPK</th>
              <th class="py-4 px-2 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="paginatedData.length === 0">
              <td colspan="9" class="text-center py-12 text-gray-400">Data tidak ditemukan</td>
            </tr>
            <tr v-for="(item, index) in paginatedData" :key="`${item.mahasiswa_name}-${item.semester}`"
              class="hover:bg-gray-50 border-b border-gray-100 text-sm">
              <td class="py-3 px-2 text-center">{{ (currentPage - 1) * perPage + index + 1 }}</td>
              <td class="py-3 px-2">{{ item.mahasiswa_name }}</td>
              <td class="py-3 px-2">{{ item.kelas_name }}</td>
              <td class="py-3 px-2 capitalize">{{ fmt(item.prodi_name) }}</td>
              <td class="py-3 px-2 text-center">{{ item.semester }}</td>
              <td class="py-3 px-2 text-center">{{ item.tahun_akademik }}</td>
              <td class="py-3 px-2 text-center font-medium">{{ item.ips.toFixed(2) }}</td>
              <td class="py-3 px-2 text-center font-medium">{{ item.ipk.toFixed(2) }}</td>
              <td class="py-3 px-2 text-center">
                <button @click="goDetail(item)"
                  class="bg-[#243e90] hover:bg-[#1d377f] text-white px-4 py-1.5 rounded-lg text-sm font-medium">
                  Detail
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAGINATION -->
      <div class="flex justify-between items-center px-5 py-4 border-t border-gray-200">
        <select v-model="perPage" class="border border-gray-300 rounded-lg px-3 py-2 text-sm bg-white">
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
          <option :value="50">50 Baris</option>
        </select>

        <div class="flex items-center gap-2 text-sm text-gray-500">
          <button @click="prevPage" :disabled="currentPage === 1"
            class="px-3 py-1.5 rounded-lg border disabled:opacity-40 hover:bg-gray-100">← Previous</button>

          <template v-for="page in visiblePages" :key="page">
            <span v-if="page === '...'" class="px-1">...</span>
            <button v-else @click="currentPage = Number(page)"
              :class="currentPage === Number(page) ? 'bg-[#243e90] text-white' : 'bg-white hover:bg-gray-100'"
              class="w-8 h-8 rounded-lg border text-sm font-medium">
              {{ page }}
            </button>
          </template>

          <button @click="nextPage" :disabled="currentPage === totalPages"
            class="px-3 py-1.5 rounded-lg border disabled:opacity-40 hover:bg-gray-100">Next →</button>
        </div>
      </div>

    </div>
  </div>
</template>