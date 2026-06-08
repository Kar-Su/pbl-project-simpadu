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

interface KurikulumItem {
  id: string
  kode: string
  name: string
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
    const res = await fetch('/api/kurikulum?page=1', {
      headers: getHeaders(),
    })

    const data = await res.json()

    kurikulumList.value = data.data.items ?? []
  } catch (err) {
    console.error('getKurikulum:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getKurikulum()
})
</script>

<template>
  <div class="p-1">

    <!-- HEADER -->
    <h1 class="text-2xl font-bold text-gray-900 mb-1">
      Data Kurikulum
    </h1>

    <p class="text-sm text-gray-500 mb-6">
      Daftar Kurikulum dan Mata Kuliah
    </p>

    <!-- CARD -->
    <div class="bg-white rounded-2xl border border-blue-100 shadow-sm p-6">

      <!-- Loading -->
      <div
        v-if="loading"
        class="py-12 text-center text-gray-400 text-sm"
      >
        Memuat data...
      </div>

      <!-- DATA -->
      <div
        v-else
        class="space-y-8"
      >

        <!-- Kurikulum -->
        <div
          v-for="kurikulum in kurikulumList"
          :key="kurikulum.id"
          class="border rounded-2xl overflow-hidden"
        >

          <!-- Header Kurikulum -->
          <div class="bg-[#1f3c93] text-white px-5 py-4">
            <h2 class="font-bold text-lg">
              {{ kurikulum.name }}
            </h2>

            <p class="text-sm opacity-90">
              Kode: {{ kurikulum.kode }}
            </p>
          </div>

          <!-- Tabel MK -->
          <table class="w-full text-sm">

            <thead class="bg-gray-50">
              <tr class="border-b border-gray-200">
                <th class="py-3 px-4 text-left">No</th>
                <th class="py-3 px-4 text-left">Kode MK</th>
                <th class="py-3 px-4 text-left">Nama MK</th>
                <th class="py-3 px-4 text-left">SKS</th>
                <th class="py-3 px-4 text-left">Semester</th>
                <th class="py-3 px-4 text-left">Status</th>
              </tr>
            </thead>

            <tbody>

              <!-- kosong -->
              <tr v-if="kurikulum.kurikulum_mk.length === 0">
                <td
                  colspan="6"
                  class="text-center py-8 text-gray-400"
                >
                  Tidak ada mata kuliah
                </td>
              </tr>

              <!-- data -->
              <tr
                v-for="(item, index) in kurikulum.kurikulum_mk"
                :key="item.mata_kuliah.id"
                class="border-b border-gray-100 hover:bg-gray-50"
              >
                <td class="py-3 px-4">
                  {{ index + 1 }}
                </td>

                <td class="py-3 px-4 font-medium">
                  {{ item.mata_kuliah.kode }}
                </td>

                <td class="py-3 px-4">
                  {{ item.mata_kuliah.name }}
                </td>

                <td class="py-3 px-4">
                  {{ item.mata_kuliah.sks }}
                </td>

                <td class="py-3 px-4">
                  Semester {{ item.semester }}
                </td>

                <td class="py-3 px-4">
                  <span
                    :class="item.wajib
                      ? 'text-green-600 font-semibold'
                      : 'text-yellow-600 font-semibold'"
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