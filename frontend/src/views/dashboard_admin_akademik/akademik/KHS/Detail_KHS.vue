<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()

const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
  "Content-Type": "application/json",
})

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

// ============================
// STATE
// ============================

const allNilai = ref<KhsItem[]>([])
const loading = ref(false)
const error = ref("")

const search = ref("")
const currentPage = ref(1)
const perPage = ref(10)

const mahasiswaName = decodeURIComponent(
  (route.params.id as string) || ""
)

// ============================
// DATA DARI SESSION
// ============================

const khsItem =
  history.state?.khsItem ||
  JSON.parse(
    sessionStorage.getItem("selectedKhs") || "null"
  )

// ============================
// FETCH DATA
// ============================

const fetchByName = async () => {
  loading.value = true
  error.value = ""

  try {
    const requests = Array.from(
      { length: 8 },
      (_, i) =>
        fetch(
          `${BASE_URL}/api/khs?semester=${i + 1}`,
          {
            headers: getHeaders(),
          }
        )
          .then((r) => r.json())
          .then(
            (j) =>
              (j?.data ?? []) as KhsItem[]
          )
          .catch(() => [])
    )

    const results = await Promise.all(
      requests
    )

    const all = results.flat()

    allNilai.value = all.filter(
      (item) =>
        item.mahasiswa_name ===
        mahasiswaName
    )

    console.log(
      "DETAIL KHS:",
      allNilai.value
    )
  } catch (err) {
    console.error(err)

    error.value =
      "Gagal memuat data KHS."
  } finally {
    loading.value = false
  }
}

// ============================
// INITIAL LOAD
// ============================

onMounted(() => {
  if (khsItem) {
    allNilai.value = [khsItem]
  } else {
    fetchByName()
  }
})

// ============================
// HELPER
// ============================

const fmt = (text: string) =>
  (text ?? "-").replace(/-/g, " ")

// ============================
// INFO MAHASISWA
// ============================

const info = computed(
  () => allNilai.value[0] ?? null
)

// ============================
// MATA KULIAH
// ============================

const allMataKuliah = computed(() =>
  allNilai.value.flatMap((semester) =>
    semester.nilai.map((mk) => ({
      ...mk,
      semester:
        semester.semester,
      kelas:
        semester.kelas_name,
    }))
  )
)

const filteredMataKuliah =
  computed(() => {
    if (!search.value)
      return allMataKuliah.value

    const keyword =
      search.value.toLowerCase()

    return allMataKuliah.value.filter(
      (mk) =>
        mk.nama_mk
          .toLowerCase()
          .includes(keyword) ||
        mk.kode_mk
          .toLowerCase()
          .includes(keyword)
    )
  })

// ============================
// PAGINATION
// ============================

watch(
  [search, perPage],
  () => {
    currentPage.value = 1
  }
)

const totalPages = computed(() =>
  Math.max(
    1,
    Math.ceil(
      filteredMataKuliah.value.length /
        perPage.value
    )
  )
)

const paginatedMK = computed(() => {
  const start =
    (currentPage.value - 1) *
    perPage.value

  return filteredMataKuliah.value.slice(
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

// ============================
// RINGKASAN
// ============================

const totalSks = computed(() =>
  allMataKuliah.value.reduce(
    (sum, mk) => sum + mk.sks,
    0
  )
)

// IPS dari data API
const ips = computed(() => {
  if (!allNilai.value.length)
    return "0.00"

  const total =
    allNilai.value.reduce(
      (sum, item) =>
        sum + item.ips,
      0
    )

  return (
    total /
    allNilai.value.length
  ).toFixed(2)
})

// IPK mahasiswa
const ipk = computed(() => {
  if (!info.value)
    return "0.00"

  return info.value.ipk.toFixed(2)
})

// ============================
// GRADE COLOR
// ============================

const gradeColor = (
  grade: string
) => {
  switch (grade) {
    case "A":
      return "text-green-600 font-bold"

    case "B":
      return "text-blue-600 font-bold"

    case "C":
      return "text-yellow-600 font-bold"

    case "D":
      return "text-orange-500 font-bold"

    default:
      return "text-red-600 font-bold"
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="mb-6">
      <p class="text-gray-500 text-sm cursor-pointer hover:underline" @click="router.back()">
        ← Akademik › KHS › {{ mahasiswaName }}
      </p>
      <h1 class="text-4xl font-bold text-gray-700 mt-1">Kartu Hasil Studi</h1>
      <p class="text-gray-500 text-sm mt-1">Detail nilai mahasiswa</p>
    </div>

    <!-- LOADING / ERROR -->
    <div v-if="loading" class="text-center py-14 text-gray-500">Memuat data...</div>
    <div v-else-if="error" class="text-center py-14 text-red-500">
      {{ error }}
      <button @click="fetchByName" class="ml-2 underline text-blue-600">Coba lagi</button>
    </div>

    <template v-else-if="info">

      <!-- INFO MAHASISWA -->
      <div class="bg-white rounded-xl shadow-sm border border-blue-100 p-6 mb-5">
        <h2 class="text-lg font-bold text-gray-700 mb-4">Informasi Mahasiswa</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div class="space-y-3">
            <div><span class="text-gray-500 w-32 inline-block">Nama</span><span class="font-medium">: {{ info.mahasiswa_name }}</span></div>
            <div><span class="text-gray-500 w-32 inline-block">Kelas</span><span class="font-medium">: {{ info.kelas_name }}</span></div>
            <div><span class="text-gray-500 w-32 inline-block">Prodi</span><span class="font-medium capitalize">: {{ fmt(info.prodi_name) }}</span></div>
          </div>
          <div class="space-y-3">
            <div><span class="text-gray-500 w-36 inline-block">Kurikulum</span><span class="font-medium">: {{ info.kurikulum_name }}</span></div>
            <div><span class="text-gray-500 w-36 inline-block">Tahun Akademik</span><span class="font-medium">: {{ info.tahun_akademik }}</span></div>
            <div><span class="text-gray-500 w-36 inline-block">IPK</span><span class="font-bold text-[#243e90]">: {{ info.ipk.toFixed(2) }}</span></div>
          </div>
        </div>
      </div>

      <!-- RINGKASAN per SEMESTER -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-5">
        <div v-for="d in allNilai" :key="d.semester"
          class="bg-white rounded-xl border border-blue-100 shadow-sm p-4 text-center">
          <div class="text-gray-500 text-xs mb-1">Semester {{ d.semester }}</div>
          <div class="text-2xl font-bold text-[#243e90]">{{ d.ips.toFixed(2) }}</div>
          <div class="text-xs text-gray-400 mt-1">IPS</div>
        </div>
      </div>

      <!-- TABEL NILAI -->
      <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">
        <div class="bg-[#243e90] px-5 py-4">
          <h2 class="text-white text-xl font-bold">Detail Nilai Mata Kuliah</h2>
          <p class="text-white text-sm mt-1">Semua mata kuliah yang telah ditempuh</p>
        </div>

        <div class="px-5 pt-4 pb-2">
          <input v-model="search" type="text"
            placeholder="Cari mata kuliah atau kode..."
            class="w-full border border-gray-300 rounded-xl px-4 py-2.5 text-sm outline-none focus:border-blue-500 bg-white" />
        </div>

        <div class="px-5 pb-3 overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="text-gray-600 border-b border-gray-300 text-sm">
                <th class="py-4 px-2 text-center w-10">No</th>
                <th class="py-4 px-2 text-center">Smt</th>
                <th class="py-4 px-2 text-left">Kode MK</th>
                <th class="py-4 px-2 text-left">Nama Mata Kuliah</th>
                <th class="py-4 px-2 text-center">SKS</th>
                <th class="py-4 px-2 text-center">Nilai</th>
                <th class="py-4 px-2 text-center">Grade</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="paginatedMK.length === 0">
                <td colspan="7" class="text-center py-10 text-gray-400">Tidak ada data</td>
              </tr>
              <tr v-for="(mk, index) in paginatedMK" :key="`${mk.kode_mk}-${mk.semester}`"
                class="hover:bg-gray-50 border-b border-gray-100 text-sm">
                <td class="py-3 px-2 text-center">{{ (currentPage - 1) * perPage + index + 1 }}</td>
                <td class="py-3 px-2 text-center">{{ mk.semester }}</td>
                <td class="py-3 px-2 font-mono text-xs text-gray-600">{{ mk.kode_mk }}</td>
                <td class="py-3 px-2 capitalize">{{ fmt(mk.nama_mk) }}</td>
                <td class="py-3 px-2 text-center">{{ mk.sks }}</td>
                <td class="py-3 px-2 text-center">{{ mk.nilai }}</td>
                <td class="py-3 px-2 text-center" :class="gradeColor(mk.grade)">{{ mk.grade }}</td>
              </tr>

              <!-- TOTAL -->
              <tr class="border-t-2 border-gray-300 bg-gray-50 text-sm font-semibold">
                <td colspan="4" class="py-3 px-2 text-right text-gray-600">Total SKS :</td>
                <td class="py-3 px-2 text-center text-[#243e90]">{{ totalSks }}</td>
                <td colspan="2"></td>
              </tr>
              <tr class="bg-gray-50 text-sm font-semibold">
                <td colspan="4" class="py-3 px-2 text-right text-gray-600">IPS (rata-rata) :</td>
                <td colspan="3" class="py-3 px-2 text-[#243e90] font-bold">{{ ips }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- PAGINATION -->
        <div class="flex justify-between items-center px-5 py-4 border-t border-gray-200">
          <select v-model="perPage" class="border border-gray-300 rounded-lg px-3 py-2 text-sm bg-white">
            <option :value="10">10 Baris</option>
            <option :value="25">25 Baris</option>
          </select>
          <div class="flex items-center gap-2 text-sm text-gray-500">
            <button @click="prevPage" :disabled="currentPage === 1"
              class="px-3 py-1.5 rounded-lg border disabled:opacity-40 hover:bg-gray-100">← Prev</button>
            <span class="w-8 h-8 rounded-lg bg-[#243e90] text-white flex items-center justify-center text-sm">{{ currentPage }}</span>
            <button @click="nextPage" :disabled="currentPage === totalPages"
              class="px-3 py-1.5 rounded-lg border disabled:opacity-40 hover:bg-gray-100">Next →</button>
          </div>
        </div>
      </div>

    </template>

    <div v-else class="text-center py-14 text-gray-400">Data tidak ditemukan</div>

  </div>
</template>