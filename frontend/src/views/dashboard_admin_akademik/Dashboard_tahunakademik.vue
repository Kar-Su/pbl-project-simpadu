<script setup lang="ts">
import { ref, onMounted } from 'vue'

// ─────────────────────────────────────────────
// INTERFACES
// ─────────────────────────────────────────────
interface TahunAkademikItem {
  id: number
  tipee_semester: string
  tahun_awal: string
  tahun_akhir: string
  status: string
}

// ─────────────────────────────────────────────
// STATE
// ─────────────────────────────────────────────
const tahunAkademikList = ref<TahunAkademikItem[]>([])
const loading = ref(false)

const BASE_URL = 'https://be.karlearn.site'

// ─────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  'Content-Type': 'application/json',
  accept: 'application/json',
  Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
})

const formatYear = (date: string): string => {
  return new Date(date).getFullYear().toString()
}

const capitalize = (str: string): string => {
  if (!str) return '-'
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// ─────────────────────────────────────────────
// FETCH
// ─────────────────────────────────────────────
const getTahunAkademik = async (): Promise<void> => {
  try {
    loading.value = true

    const res = await fetch(`${BASE_URL}/api/tahun-akademik`, {
      headers: getHeaders(),
    })

    const json = await res.json()
    console.log('TAHUN AKADEMIK:', json)

    tahunAkademikList.value = Array.isArray(json.data) ? json.data : []
  } catch (err) {
    console.error('getTahunAkademik:', err)
  } finally {
    loading.value = false
  }
}

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted(() => {
  getTahunAkademik()
})
</script>

<template>
  <div class="p-1">

    <!-- HEADER -->
    <h1 class="text-2xl font-bold text-gray-900 mb-1">Tahun Akademik</h1>
    <p class="text-sm text-gray-500 mb-6">Detail dari Tahun Akademik</p>

    <!-- CARD TABEL -->
    <div class="bg-white rounded-2xl border border-blue-100 shadow-sm p-6">

      <!-- Judul tabel -->
      <h2 class="text-base font-bold text-gray-800 mb-5">Data Tahun Akademik</h2>

      <!-- Loading -->
      <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">
        Memuat data...
      </div>

      <!-- Tabel -->
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-200">
            <th class="py-3 px-2 text-left font-semibold text-gray-600 w-12">No</th>
            <!-- <th class="py-3 px-2 text-left font-semibold text-gray-600">Semester</th> -->
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Tahun Awal</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Tahun Akhir</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Status</th>
          </tr>
        </thead>

        <tbody>
          <!-- Kosong -->
          <tr v-if="tahunAkademikList.length === 0">
            <td colspan="5" class="py-12 text-center text-gray-400">
              Tidak ada data tahun akademik
            </td>
          </tr>

          <!-- Baris data -->
          <tr
            v-for="(item, index) in tahunAkademikList"
            :key="item.id"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="py-4 px-2 text-gray-700">{{ index + 1 }}</td>
            <!-- <td class="py-4 px-2 text-gray-700 font-medium">{{ capitalize(item.tipee_semester) }}</td> -->
            <td class="py-4 px-2 text-gray-700">{{ formatYear(item.tahun_awal) }}</td>
            <td class="py-4 px-2 text-gray-700">{{ formatYear(item.tahun_akhir) }}</td>
            <td class="py-4 px-2">
              <span
                :class="item.status === 'aktif'
                  ? 'bg-green-100 text-green-700'
                  : 'bg-red-100 text-red-500'"
                class="px-3 py-1 rounded-full text-xs font-semibold capitalize"
              >
                {{ item.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

    </div>

  </div>
</template>