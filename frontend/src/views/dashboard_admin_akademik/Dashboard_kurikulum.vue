<script setup lang="ts">
import { ref, onMounted } from 'vue'

// ================= INTERFACE =================
interface MataKuliah {
  id: string
  kode: string
  name: string
  sks: number
}

interface KurikulumMK {
  semester: number
  wajib: boolean
  mata_kuliah: MataKuliah
}

interface Prodi {
  id: number
  name: string
  jenjang: string
  jurusan: { id: number; name: string }
}

interface KurikulumItem {
  id: string
  kode: string
  name: string
  prodi: Prodi
  kurikulum_mk: KurikulumMK[]
}

// ================= STATE =================
const kurikulumList = ref<KurikulumItem[]>([])
const loading = ref(false)

// ================= HEADER =================
const BASE_URL = 'https://be.karlearn.site'
const getHeaders = (): Record<string, string> => ({
  'Content-Type': 'application/json',
  accept: 'application/json',
  Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
})

// ================= GET DATA =================
const getKurikulum = async (): Promise<void> => {
  loading.value = true
  try {
    const res = await fetch(`${BASE_URL}/api/kurikulum?page=1`, {
      headers: getHeaders(),
    })
    const json = await res.json()
    console.log('KURIKULUM:', json)
    kurikulumList.value = json.data?.items ?? []
  } catch (err) {
    console.error('getKurikulum:', err)
  } finally {
    loading.value = false
  }
}

const fmt = (str: string) => (str ?? '-').replace(/-/g, ' ')

onMounted(() => {
  getKurikulum()
})
</script>

<template>
  <div class="p-1">

    <!-- HEADER -->
    <h1 class="text-2xl font-bold text-gray-900 mb-1">Data Kurikulum</h1>
    <p class="text-sm text-gray-500 mb-6">Daftar Kurikulum dan Mata Kuliah</p>

    <!-- CARD -->
    <div class="bg-white rounded-2xl border border-blue-100 shadow-sm p-6">

      <!-- Loading -->
      <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">
        Memuat data...
      </div>

      <!-- DATA -->
      <div v-else class="space-y-8">

        <!-- Kosong -->
        <p v-if="kurikulumList.length === 0" class="text-center py-12 text-gray-400 text-sm">
          Tidak ada data kurikulum
        </p>

        <!-- Per Kurikulum -->
        <div
          v-for="kurikulum in kurikulumList"
          :key="kurikulum.id"
          class="border border-gray-200 rounded-2xl overflow-hidden"
        >

          <!-- Header Kurikulum -->
          <div class="bg-[#1f3c93] text-white px-5 py-4">
            <h2 class="font-bold text-lg capitalize">{{ fmt(kurikulum.name) }}</h2>
            <p class="text-sm opacity-90 mt-0.5">
              Kode: <span class="uppercase">{{ kurikulum.kode }}</span>
              &nbsp;·&nbsp;
              Prodi: <span class="capitalize">{{ fmt(kurikulum.prodi?.name) }}</span>
              &nbsp;·&nbsp;
              {{ kurikulum.prodi?.jenjang }}
            </p>
          </div>

          <!-- Tabel MK -->
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr class="border-b border-gray-200">
                <th class="py-3 px-4 text-left font-semibold text-gray-600">No</th>
                <th class="py-3 px-4 text-left font-semibold text-gray-600">Kode MK</th>
                <th class="py-3 px-4 text-left font-semibold text-gray-600">Nama MK</th>
                <th class="py-3 px-4 text-left font-semibold text-gray-600">SKS</th>
                <th class="py-3 px-4 text-left font-semibold text-gray-600">Semester</th>
                <th class="py-3 px-4 text-left font-semibold text-gray-600">Status</th>
              </tr>
            </thead>

            <tbody>
              <!-- Kosong -->
              <tr v-if="kurikulum.kurikulum_mk.length === 0">
                <td colspan="6" class="text-center py-8 text-gray-400">
                  Tidak ada mata kuliah
                </td>
              </tr>

              <!-- Data -->
              <tr
                v-for="(item, index) in kurikulum.kurikulum_mk"
                :key="item.mata_kuliah.id"
                class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
              >
                <td class="py-3 px-4 text-gray-700">{{ index + 1 }}</td>
                <td class="py-3 px-4 font-medium text-gray-700 uppercase">{{ item.mata_kuliah.kode }}</td>
                <td class="py-3 px-4 text-gray-700 capitalize">{{ fmt(item.mata_kuliah.name) }}</td>
                <td class="py-3 px-4 text-gray-700">{{ item.mata_kuliah.sks }}</td>
                <td class="py-3 px-4 text-gray-700">Semester {{ item.semester }}</td>
                <td class="py-3 px-4">
                  <span
                    :class="item.wajib
                      ? 'bg-green-100 text-green-700'
                      : 'bg-yellow-100 text-yellow-700'"
                    class="px-2 py-1 rounded-full text-xs font-semibold"
                  >
                    {{ item.wajib ? 'Wajib' : 'Pilihan' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>

        </div>

      </div>

    </div>

  </div>
</template>