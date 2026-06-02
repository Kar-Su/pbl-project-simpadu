<script setup lang="ts">
import { ref, onMounted } from 'vue'

// ─────────────────────────────────────────────
// INTERFACES
// ─────────────────────────────────────────────
interface TahunAkademikItem {
  id: number
  tahun_awal: string
  tahun_akhir: string
}

interface SemesterItem {
  id: number
  nama_semester: string   // "Genap" | "Ganjil"
  tahun_awal: string
  tahun_akhir: string
  kurikulum: string
  tahun_akademik_id: number
}

// ─────────────────────────────────────────────
// STATE
// ─────────────────────────────────────────────
const tahunAkademikList = ref<TahunAkademikItem[]>([])
const semesterList = ref<SemesterItem[]>([])
const loading = ref<boolean>(false)

// ─────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  'Content-Type': 'application/json',
  'accept': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token') ?? ''}`,
})

// ─────────────────────────────────────────────
// HIT API: Ambil semua Tahun Akademik
// Endpoint : GET /api/tahun-akademik
// Response : { data: TahunAkademikItem[] }
// ─────────────────────────────────────────────
// const getTahunAkademik = async (): Promise<void> => {
//   try {

//     loading.value = true

//     const res = await fetch(
//       '/api/tahun-akademik',
//       {
//         headers: getHeaders()
//       }
//     )

//     const data = await res.json()

//     tahunAkademikList.value = data.data ?? []

//     // kosongkan semester
//     semesterList.value = []

//     // loop semua tahun akademik
//     for (const tahun of tahunAkademikList.value) {

//       const semesterRes = await fetch(
//         `/api/tahun-akademik/${tahun.id}/semester`,
//         {
//           headers: getHeaders()
//         }
//       )

//       const semesterData = await semesterRes.json()

//       const items = Array.isArray(semesterData)
//         ? semesterData
//         : semesterData.data ?? []

//       const mapped = items.map((item: any) => ({
//         id: item.id,
//         nama_semester: item.tipe_semester ?? '-',
//         tahun_awal: item.tahun_awal ?? '-',
//         tahun_akhir: item.tahun_akhir ?? '-',
//         kurikulum: item.kurikulum?.name ?? '-',
//         tahun_akademik_id: item.tahun_akademik_id ?? 0,
//       }))

//       semesterList.value.push(...mapped)
//     }

//   } catch (err) {

//     console.error('getTahunAkademik:', err)

//   } finally {

//     loading.value = false

//   }
// }

const getTahunAkademik = async (): Promise<void> => {
  try {

    loading.value = true

    const res = await fetch(
      '/api/tahun-akademik/status/aktif',
      {
        headers: getHeaders()
      }
    )

    const json = await res.json()

    console.log('DATA API:', json)

    // ✅ ambil dari json.data
    const raw: any[] = Array.isArray(json.data) ? json.data : []

    // ✅ mapping langsung ke semesterList (tidak perlu loop tahun lagi)
    semesterList.value = raw.map((item: any) => ({
      id: item.id,
      nama_semester: item.type
        ? item.type.charAt(0).toUpperCase() + item.type.slice(1)
        : '-',
      tahun_awal: item.tahun_awal ?? '-',
      tahun_akhir: item.tahun_akhir ?? '-',
      kurikulum: item.kurikulum?.name ?? '-',
      tahun_akademik_id: item.id,
    }))

    console.log('SEMESTER LIST:', semesterList.value)

  } catch (err) {

    console.error('getTahunAkademik:', err)

  } finally {

    loading.value = false

  }
}

const formatYear = (date: string): string => {
  return new Date(date).getFullYear().toString()
}

// const selectTahun = async (tahun: TahunAkademikItem): Promise<void> => {
//   selectedTahun.value = tahun
//   loading.value = true
//   semesterList.value = []

//   try {

//     const res = await fetch(
//       `/api/tahun-akademik/${tahun.id}/semester`,
//       {
//         headers: getHeaders()
//       }
//     )

//     const data = await res.json()

//     console.log('SEMESTER:', data)

//     const items = Array.isArray(data)
//       ? data
//       : data.data ?? []

//     semesterList.value = items.map((item: any) => ({
//       id: item.id,
//       nama_semester: item.tipe_semester ?? '-',
//       tahun_awal: item.tahun_awal ?? '-',
//       tahun_akhir: item.tahun_akhir ?? '-',
//       kurikulum: item.kurikulum?.name ?? '-',
//       tahun_akademik_id: item.tahun_akademik_id ?? 0,
//     }))

//   } catch (err) {
//     console.error('getSemester:', err)
//   } finally {
//     loading.value = false
//   }
// }

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted((): void => {
  getTahunAkademik()
})
</script>

<template>
  <div class="p-1">

    <!-- HEADER -->
    <h1 class="text-2xl font-bold text-gray-900 mb-1">Tahun Akademik</h1>
    <p class="text-sm text-gray-500 mb-6">Detail dari Tahun Akademik</p>

    <!-- TABS TAHUN AKADEMIK -->
    <!-- <div class="flex gap-2 mb-4 flex-wrap">
      <button
        v-for="tahun in tahunAkademikList"
        :key="tahun.id"
        @click="selectTahun(tahun)"
        class="px-4 py-1.5 rounded-full text-sm font-semibold border transition-colors"
        :class="selectedTahun?.id === tahun.id
          ? 'bg-[#1f3c93] text-white border-[#1f3c93]'
          : 'bg-white text-gray-600 border-gray-300 hover:border-[#1f3c93] hover:text-[#1f3c93]'"
      >
        {{ formatYear(tahun.tahun_awal) }} –
{{ formatYear(tahun.tahun_akhir) }}
      </button>
    </div> -->

    <!-- CARD TABEL -->
    <div class="bg-white rounded-2xl border border-blue-100 shadow-sm p-6">

      <!-- Judul tabel -->
<h2 class="text-base font-bold text-gray-800 mb-5">
  Data Tahun Akademik
</h2>

      <!-- Loading -->
      <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">
        Memuat data...
      </div>

      <!-- Tabel -->
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-200">
            <th class="py-3 px-2 text-left font-semibold text-gray-600 w-12">No</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Semester</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Tahun Awal</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Tahun Akhir</th>
            <th class="py-3 px-2 text-left font-semibold text-gray-600">Kurikulum</th>
          </tr>
        </thead>

        <tbody>
          <!-- Kosong -->
          <tr v-if="semesterList.length === 0 && !loading">
            <td colspan="5" class="py-12 text-center text-gray-400">
              Tidak ada data semester
            </td>
          </tr>

          <!-- Baris data -->
          <tr
            v-for="(item, index) in semesterList"
            :key="item.id"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="py-4 px-2 text-gray-700">{{ index + 1 }}</td>
            <td class="py-4 px-2 text-gray-700 font-medium">{{ item.nama_semester }}</td>
            <td class="py-4 px-2 text-gray-700">
  {{ formatYear(item.tahun_awal) }}
</td>

<td class="py-4 px-2 text-gray-700">
  {{ formatYear(item.tahun_akhir) }}
</td>
            <td class="py-4 px-2 text-gray-700">{{ item.kurikulum }}</td>
          </tr>
        </tbody>
      </table>

    </div>

  </div>
</template>