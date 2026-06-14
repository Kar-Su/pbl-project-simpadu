<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"

const route = useRoute()

const kurikulum = ref<any>(null)
const matkulList = ref<any[]>([])
const loading = ref(false)
const errorMsg = ref("")

const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

const fmt = (str: string) => (str ?? "-").replace(/-/g, " ")

// ================= MAP DATA KE matkulList =================
const setKurikulumData = (data: any) => {
  kurikulum.value = data

  matkulList.value =
    data.kurikulum_mk?.map((item: any) => ({
      kode: item.mata_kuliah?.kode ?? "-",
      nama: item.mata_kuliah?.name ?? "-",
      sks: item.mata_kuliah?.sks ?? "-",
      semester: item.semester,
      wajib: item.wajib,
    })) ?? []
}

// ================= FALLBACK: AMBIL DARI LIST API =================
const fetchFromList = async () => {
  loading.value = true

  try {
    const id = route.params.id
    let allItems: any[] = []
    let page = 1
    let lastPage = 1

    do {
      const res = await fetch(
        `${BASE_URL}/api/kurikulum?page=${page}&per_page=100`,
        { method: "GET", headers: getHeaders() }
      )

      if (!res.ok) {
        throw new Error(`Gagal memuat data (status ${res.status})`)
      }

      const json = await res.json()
      const items = json.data?.items ?? []
      allItems = [...allItems, ...items]
      lastPage = json.data?.pagination?.total_pages ?? 1
      page++
    } while (page <= lastPage)

    const found = allItems.find((item: any) => item.id === id)

    if (!found) {
      errorMsg.value = "Data kurikulum tidak ditemukan"
      return
    }

    setKurikulumData(found)
  } catch (err) {
    console.error("FETCH KURIKULUM LIST ERROR:", err)
    errorMsg.value = "Terjadi kesalahan saat memuat data kurikulum"
  } finally {
    loading.value = false
  }
}

// ================= GET DETAIL =================
const getKurikulumDetail = async () => {
  errorMsg.value = ""

  // 1. Coba ambil dari state navigasi (dikirim dari halaman list)
  const stateData = history.state?.kurikulum
  if (stateData) {
    setKurikulumData(stateData)
    return
  }

  // 2. Fallback: kalau halaman direfresh / dibuka langsung,
  //    ambil dari list API lalu cari sesuai id
  await fetchFromList()
}

onMounted(() => {
  getKurikulumDetail()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <div class="text-sm text-gray-500 mb-2">
      Mahasiswa > Kurikulum
    </div>

    <h1 class="text-[42px] font-bold text-[#404040]">
      Kurikulum
    </h1>

    <p class="text-gray-500 text-sm mb-6">
      Data matakuliah yang tercantum dalam kurikulum
    </p>

    <!-- Loading -->
    <div v-if="loading" class="py-12 text-center text-gray-400 text-sm">
      Memuat data...
    </div>

    <!-- Error -->
    <div v-else-if="errorMsg" class="py-12 text-center text-red-400 text-sm">
      {{ errorMsg }}
    </div>

    <div
      v-else
      class="bg-white rounded-xl shadow-md overflow-hidden border border-gray-200"
    >

      <div class="bg-[#243e90] px-6 py-5">
        <h2 class="text-white text-3xl font-bold capitalize">
          {{ fmt(kurikulum?.name) }}
        </h2>

        <p class="text-white mt-2 uppercase">
          kode :
          {{ kurikulum?.kode }}
        </p>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">

          <thead>
            <tr class="border-b border-gray-300 bg-white">
              <th class="p-4 text-left">No</th>
              <th class="p-4 text-left">Kode Matakuliah</th>
              <th class="p-4 text-left">Nama</th>
              <th class="p-4 text-left">SKS</th>
              <th class="p-4 text-left">Semester</th>
              <th class="p-4 text-left">Status</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(item, index) in matkulList"
              :key="index"
              class=" hover:bg-gray-50"
            >
              <td class="p-4">
                {{ index + 1 }}
              </td>

              <td class="p-4 uppercase">
                {{ item.kode }}
              </td>

              <td class="p-4 capitalize">
                {{ item.nama }}
              </td>

              <td class="p-4">
                {{ item.sks }}
              </td>

              <td class="p-4">
                Semester {{ item.semester }}
              </td>

              <td class="p-4">
                <span
                  :class="
                    item.wajib
                      ? 'text-green-500 font-semibold'
                      : 'text-orange-400 font-semibold'
                  "
                >
                  {{ item.wajib ? 'Wajib' : 'Pilihan' }}
                </span>
              </td>
            </tr>

            <tr v-if="matkulList.length === 0">
              <td
                colspan="6"
                class="text-center py-8 text-gray-400"
              >
                Tidak ada data mata kuliah
              </td>
            </tr>
          </tbody>

        </table>
      </div>

    </div>

  </div>
</template>