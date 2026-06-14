<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

// ─────────────────────────────────────────────
// HELPER HEADER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ─────────────────────────────────────────────
// FILTER STATE
// ─────────────────────────────────────────────
const selectedJurusan = ref<string>("")
const selectedProdi = ref<string>("")
const selectedTahun = ref<string>("")

// ─────────────────────────────────────────────
// MASTER DATA INTERFACES
// ─────────────────────────────────────────────
interface Jurusan {
  id: string
  name: string
}

interface Prodi {
  id: string
  name: string
}

interface TahunAkademik {
  id: string
  tahun_awal: string
  tahun_akhir: string
}

const jurusanList = ref<Jurusan[]>([])
const prodiList = ref<Prodi[]>([])
const filteredProdiList = computed(() => {
  if (!selectedJurusan.value) {
    return prodiList.value
  }

  return prodiList.value.filter((item: any) => {
    return (
      String(item.jurusan?.id ?? "") ===
      selectedJurusan.value
    )
  })
})
const tahunAkademikList = ref<TahunAkademik[]>([])

// ─────────────────────────────────────────────
// DATA KELAS INTERFACES
// ─────────────────────────────────────────────
interface JurusanRef {
  id: number | string
  name: string
}

interface ProdiRef {
  id: number | string
  name: string
  jenjang?: string
  jurusan?: JurusanRef
}

interface TahunAkademikRef {
  id: number | string
  tipe_semester?: string
  tahun_awal: string
  tahun_akhir: string
  status?: string
}

interface KelasItem {
  id?: string | number
  kelas_id?: string | number
  nama_kelas?: string
  nama?: string
  name?: string
  prodi?: ProdiRef
  tahun_akademik?: TahunAkademikRef
  [key: string]: any
}

const kelasList = ref<KelasItem[]>([])

// ─────────────────────────────────────────────
// HELPER TAMPILAN
// ─────────────────────────────────────────────
const getKelasId = (item: KelasItem): string => {
  const id = item.id ?? item.kelas_id ?? ""
  return id ? String(id) : ""
}

const getNamaKelas = (item: KelasItem): string =>
  item.nama_kelas ?? item.nama ?? item.name ?? "-"

const getJurusan = (item: KelasItem): string =>
  item.prodi?.jurusan?.name ?? "-"

const getProdi = (item: KelasItem): string =>
  item.prodi?.name ?? "-"

const getTahunAkademik = (item: KelasItem): string => {
  const ta = item.tahun_akademik
  if (!ta) return "-"
  const awal = ta.tahun_awal?.slice(0, 4) ?? ""
  const akhir = ta.tahun_akhir?.slice(0, 4) ?? ""
  return `${awal}-${akhir}`
}

// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const currentPage = ref<number>(1)
const perPage = ref<number>(5)
const totalItems = ref<number>(0)

const totalPages = computed<number>(() =>
  Math.max(1, Math.ceil(totalItems.value / perPage.value))
)

const pages = computed<(number | string)[]>(() => {
  const total = totalPages.value
  const cur = currentPage.value

  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const result: (number | string)[] = [1, 2]
  if (cur > 4) result.push("...")

  for (let i = Math.max(3, cur - 1); i <= Math.min(total - 2, cur + 1); i++) {
    result.push(i)
  }

  if (cur < total - 3) result.push("...")
  result.push(total - 1, total)

  return [...new Set(result)]
})

// ─────────────────────────────────────────────
// FETCH DATA API
// ─────────────────────────────────────────────
const getJurusanData = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/jurusan`, { headers: getHeaders() })
    if (!res.ok) return
    const data = await res.json()
    const payload = data.data ?? {}
    jurusanList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
  } catch (err) {
    console.error(err)
  }
}

const getProdiData = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, { headers: getHeaders() })
    if (!res.ok) return
    const data = await res.json()
    const payload = data.data ?? {}
    prodiList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
  } catch (err) {
    console.error(err)
  }
}

const getTahunAkademikData = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik`, { headers: getHeaders() })
    if (!res.ok) return
    const data = await res.json()
    const payload = data.data ?? {}
    tahunAkademikList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
  } catch (err) {
    console.error(err)
  }
}

const getKelas = async (): Promise<void> => {
  try {
    let url = `${BASE_URL}/api/kelas?page=${currentPage.value}&per_page=${perPage.value}`
    if (selectedJurusan.value) url += `&jurusan_id=${selectedJurusan.value}`
    if (selectedProdi.value) url += `&prodi_id=${selectedProdi.value}`
    if (selectedTahun.value) url += `&tahun_akademik_id=${selectedTahun.value}`

    const res = await fetch(url, { headers: getHeaders() })
    if (!res.ok) {
      kelasList.value = []
      totalItems.value = 0
      return
    }

    const data = await res.json()
    const payload = data.data ?? {}

    kelasList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
    totalItems.value = payload.pagination?.total ?? data.total ?? kelasList.value.length
  } catch (err) {
    console.error(err)
    kelasList.value = []
    totalItems.value = 0
  }
}

// ─────────────────────────────────────────────
// HANDLERS
// ─────────────────────────────────────────────
const handleFilter = (): void => {
  currentPage.value = 1
  getKelas()
}

const handleTambah = async () => {
  try {
    console.log("klik tambah")

    await router.push("/dashboard-admin/tambah_pesertakelas")

    console.log("berhasil pindah")
  } catch (err) {
    console.error("router error", err)
  }
}

const handleEdit = (item: KelasItem): void => {
  const id = getKelasId(item)

  if (!id) {
    console.warn("handleEdit: id kosong untuk item:", item)
    return
  }

  // Mengarahkan ke route edit disertai dengan ID kelasnya
  router.push(`/dashboard-admin/edit_pesertakelas/${id}`)
}

const handleDelete = (item: KelasItem): void => {
  console.log("Delete:", getKelasId(item))
}

const goToPage = (page: number): void => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  getKelas()
}

watch(selectedJurusan, () => {
  const masihValid =
    filteredProdiList.value.some(
      (item: any) =>
        String(item.id) ===
        selectedProdi.value
    )

  if (!masihValid) {
    selectedProdi.value = ""
  }
})

onMounted((): void => {
  getJurusanData()
  getProdiData()
  getTahunAkademikData()
  getKelas()
})
</script>

<template>
  <div class="min-h-screen bg-[#f8fafc] p-6">
    <div class="overflow-hidden rounded-xl bg-white shadow-md border border-gray-100">
      
      <div class="bg-[#1e3a8a] px-6 py-5">
        <h2 class="text-2xl font-semibold tracking-wide text-white">
          Data Peserta Kelas
        </h2>
        <p class="mt-1 text-sm text-blue-100 opacity-90">
          Kumpulan mahasiswa yang termuat didalam satu kelas
        </p>
      </div>

      <div class="flex flex-wrap items-center justify-between gap-4 px-6 py-5 bg-white">
        
        <div class="flex flex-wrap items-center gap-3 flex-1 min-w-[300px]">
          <div class="relative w-full max-w-[210px]">
            <select
              v-model="selectedJurusan"
              class="h-[42px] w-full appearance-none rounded-lg border border-gray-300 bg-white pl-4 pr-10 text-sm text-gray-500 outline-none transition focus:border-blue-500"
            >
              <option value="">Pilih Jurusan</option>
              <option
  v-for="item in jurusanList"
  :key="item.id"
  :value="String(item.id)"
>
  {{ item.name }}
</option>
            </select>
            <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>

          <div class="relative w-full max-w-[210px]">
            <select
              v-model="selectedProdi"
              class="h-[42px] w-full appearance-none rounded-lg border border-gray-300 bg-white pl-4 pr-10 text-sm text-gray-500 outline-none transition focus:border-blue-500"
            >
              <option value="">Pilih Prodi</option>
              <option
  v-for="item in filteredProdiList"
  :key="item.id"
  :value="String(item.id)"
>
  {{ item.name }}
</option>
            </select>
            <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>

          <div class="relative w-full max-w-[230px]">
            <select
              v-model="selectedTahun"
              class="h-[42px] w-full appearance-none rounded-lg border border-gray-300 bg-white pl-4 pr-10 text-sm text-gray-500 outline-none transition focus:border-blue-500"
            >
              <option value="">Pilih Tahun Akademik</option>
              <option v-for="item in tahunAkademikList" :key="item.id" :value="item.id">
                {{ item.tahun_awal }} - {{ item.tahun_akhir }}
              </option>
            </select>
            <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="handleFilter"
            class="flex h-[42px] items-center gap-2 rounded-lg bg-[#1d357d] px-5 text-sm font-medium text-white transition hover:bg-[#162961]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 7.5L7.5 3m0 0L12 7.5M7.5 3v13.5m13.5 0L16.5 21m0 0L12 16.5m4.5 4.5V7.5" />
            </svg>
            Terapkan
          </button>

          <button
            @click="handleTambah"
            class="flex h-[42px] items-center gap-2 rounded-lg bg-[#1d357d] px-5 text-sm font-medium text-white transition hover:bg-[#162961]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            Tambah
          </button>
        </div>

      </div>

      <div class="overflow-x-auto px-6 pb-6">
        <table class="w-full border-collapse">
          <thead>
            <tr class="text-left text-sm font-semibold text-gray-700 border-b border-gray-200">
              <th class="pb-3 pt-2 w-[60px]">No</th>
              <th class="pb-3 pt-2">Nama Kelas</th>
              <th class="pb-3 pt-2">Jurusan</th>
              <th class="pb-3 pt-2">Prodi</th>
              <th class="pb-3 pt-2">Tahun Akademik</th>
              <th class="pb-3 pt-2 text-left px-4 w-[180px]">Aksi</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-gray-100">
            <tr v-if="kelasList.length === 0">
              <td colspan="6" class="py-10 text-center text-sm text-gray-400">
                Tidak ada data mahasiswa dalam kelas ini
              </td>
            </tr>

            <tr
              v-for="(item, index) in kelasList"
              :key="getKelasId(item) || index"
              class="text-[14px] text-gray-800 hover:bg-gray-50/50"
            >
              <td class="py-4 text-gray-600">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>
              <td class="py-4 font-medium">{{ getNamaKelas(item) }}</td>
              <td class="py-4 text-gray-700">{{ getJurusan(item) }}</td>
              <td class="py-4 text-gray-700">{{ getProdi(item) }}</td>
              <td class="py-4 text-gray-600">{{ getTahunAkademik(item) }}</td>
              <td class="py-4 px-2">
                <div class="flex items-center gap-2">
                  <button
                    @click="handleEdit(item)"
                    class="flex items-center gap-1.5 rounded-md bg-[#f59e0b] px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-amber-600 shadow-sm"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-3.5 h-3.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a2.25 2.25 0 1 1 3.182 3.182L10.582 17.13a4.5 4.5 0 0 1-1.897 1.13L6 19l.74-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Z" />
                    </svg>
                    Edit
                  </button>

                  <button
                    @click="handleDelete(item)"
                    class="flex items-center gap-1.5 rounded-md bg-[#ef4444] px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-red-600 shadow-sm"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-3.5 h-3.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673A2.25 2.25 0 0 1 15.916 21H8.084a2.25 2.25 0 0 1-2.245-1.327L4.772 5.79m14.456 0A48.108 48.108 0 0 0 15.75 5.25m-6.75 0a48.11 48.11 0 0 1 3.478-.459m0 0a48.11 48.11 0 0 1 3.478.459m-3.478 0V4.5a2.25 2.25 0 0 1 2.25-2.25h1.5A2.25 2.25 0 0 1 18.75 4.5v.75" />
                    </svg>
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="kelasList.length > 0" class="flex items-center justify-between border-t border-gray-100 px-6 py-4 bg-gray-50/50">
        <select
          v-model.number="perPage"
          @change="() => { currentPage = 1; getKelas() }"
          class="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-xs text-gray-600 outline-none focus:border-blue-500"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select>

        <div class="flex items-center gap-1">
          <button
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="rounded px-2.5 py-1.5 text-xs font-medium text-gray-500 hover:bg-gray-100 disabled:opacity-40"
          >
            ← Previous
          </button>

          <template v-for="(p, i) in pages" :key="i">
            <span v-if="p === '...'" class="px-2 text-xs text-gray-400">...</span>
            <button
              v-else
              @click="goToPage(p as number)"
              class="flex h-7 w-7 items-center justify-center rounded-md text-xs font-medium transition-colors"
              :class="currentPage === p ? 'bg-[#1d357d] text-white' : 'text-gray-600 hover:bg-gray-100'"
            >
              {{ p }}
            </button>
          </template>

          <button
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="rounded px-2.5 py-1.5 text-xs font-medium text-gray-600 hover:bg-gray-100 disabled:opacity-40"
          >
            Next →
          </button>
        </div>
      </div>

    </div>
  </div>
</template>