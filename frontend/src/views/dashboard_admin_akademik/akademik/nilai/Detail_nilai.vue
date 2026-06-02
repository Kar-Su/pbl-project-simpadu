<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

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
const selectedJurusan = ref("")
const selectedProdi = ref("")
const selectedTahun = ref("")

// ─────────────────────────────────────────────
// MASTER DATA
// ─────────────────────────────────────────────
interface Jurusan {
  id: number
  nama: string
}

interface Prodi {
  id: number
  nama: string
}

interface TahunAkademik {
  id: number
  tahun_awal: string
  tahun_akhir: string
}

const jurusanList = ref<Jurusan[]>([])
const prodiList = ref<Prodi[]>([])
const tahunAkademikList = ref<TahunAkademik[]>([])

// ─────────────────────────────────────────────
// DATA NILAI
// ─────────────────────────────────────────────
interface NilaiItem {
  id: number
  nim: string
  nama_mahasiswa: string
  nama_dosen: string
  kode_matakuliah: string
  grade: string
  nilai: number
}

const nilaiList = ref<NilaiItem[]>([])

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

  for (
    let i = Math.max(3, cur - 1);
    i <= Math.min(total - 2, cur + 1);
    i++
  ) {
    result.push(i)
  }

  if (cur < total - 3) result.push("...")

  result.push(total - 1, total)

  return [...new Set(result)]
})

// ─────────────────────────────────────────────
// HIT API JURUSAN
// Endpoint : GET /api/jurusan
// ─────────────────────────────────────────────
const getJurusan = async (): Promise<void> => {
  try {
    const res = await fetch("/api/jurusan", {
      headers: getHeaders(),
    })

    const data = await res.json()

    jurusanList.value = data.data ?? []
  } catch (err) {
    console.error("getJurusan:", err)
  }
}

// ─────────────────────────────────────────────
// HIT API PRODI
// Endpoint : GET /api/prodi
// ─────────────────────────────────────────────
const getProdi = async (): Promise<void> => {
  try {
    const res = await fetch("/api/prodi", {
      headers: getHeaders(),
    })

    const data = await res.json()

    prodiList.value = data.data ?? []
  } catch (err) {
    console.error("getProdi:", err)
  }
}

// ─────────────────────────────────────────────
// HIT API TAHUN AKADEMIK
// Endpoint : GET /api/tahun-akademik
// ─────────────────────────────────────────────
const getTahunAkademik = async (): Promise<void> => {
  try {
    const res = await fetch("/api/tahun-akademik", {
      headers: getHeaders(),
    })

    const data = await res.json()

    tahunAkademikList.value = data.data ?? []
  } catch (err) {
    console.error("getTahunAkademik:", err)
  }
}

// ─────────────────────────────────────────────
// HIT API NILAI
// Endpoint :
// GET /api/nilai?page=1&per_page=5
//
// OPTIONAL FILTER:
// &jurusan_id=
// &prodi_id=
// &tahun_akademik_id=
// ─────────────────────────────────────────────
const getNilai = async (): Promise<void> => {
  try {
    let url = `/api/nilai?page=${currentPage.value}&per_page=${perPage.value}`

    if (selectedJurusan.value) {
      url += `&jurusan_id=${selectedJurusan.value}`
    }

    if (selectedProdi.value) {
      url += `&prodi_id=${selectedProdi.value}`
    }

    if (selectedTahun.value) {
      url += `&tahun_akademik_id=${selectedTahun.value}`
    }

    const res = await fetch(url, {
      headers: getHeaders(),
    })

    const data = await res.json()

    nilaiList.value = data.data ?? []
    totalItems.value = data.meta?.total ?? data.total ?? 0
  } catch (err) {
    console.error("getNilai:", err)
  }
}

// ─────────────────────────────────────────────
// FILTER BUTTON
// ─────────────────────────────────────────────
const handleFilter = (): void => {
  currentPage.value = 1
  getNilai()
}

// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const goToPage = (page: number): void => {
  if (page < 1 || page > totalPages.value) return

  currentPage.value = page
  getNilai()
}

const prevPage = (): void => {
  goToPage(currentPage.value - 1)
}

const nextPage = (): void => {
  goToPage(currentPage.value + 1)
}

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted((): void => {
  getJurusan()
  getProdi()
  getTahunAkademik()
  getNilai()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
      <span>Akademik</span>
      <span>›</span>
      <span>Nilai</span>
      <span>›</span>
      <span class="text-gray-700">Detail Nilai</span>
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold leading-none text-[#333]">
      Detail Nilai
    </h1>

    <p class="mt-3 text-gray-500">
      Detail nilai dari matakuliah yang dipilih
    </p>

    <!-- CARD -->
    <div
      class="mt-8 rounded-2xl border border-[#d8e1f0] bg-white p-5 shadow-sm"
    >

      <!-- HEADER -->
      <h2 class="mb-6 text-[32px] font-bold text-[#444]">
        Data Nilai
      </h2>

      <!-- FILTER -->
      <div class="mb-8 flex flex-wrap items-center gap-5">

        <!-- Jurusan -->
        <div class="relative w-[240px]">
          <select
            v-model="selectedJurusan"
            class="h-[52px] w-full appearance-none rounded-xl border border-[#cbd5e1] bg-white px-4 pr-10 text-[15px] text-gray-600 outline-none focus:border-[#2846a3]"
          >
            <option value="">Pilih Jurusan</option>

            <option
              v-for="item in jurusanList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.nama }}
            </option>
          </select>

          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
            class="pointer-events-none absolute right-4 top-1/2 size-4 -translate-y-1/2 text-gray-500"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m19.5 8.25-7.5 7.5-7.5-7.5"
            />
          </svg>
        </div>

        <!-- Prodi -->
        <div class="relative w-[240px]">
          <select
            v-model="selectedProdi"
            class="h-[52px] w-full appearance-none rounded-xl border border-[#cbd5e1] bg-white px-4 pr-10 text-[15px] text-gray-600 outline-none focus:border-[#2846a3]"
          >
            <option value="">Pilih Prodi</option>

            <option
              v-for="item in prodiList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.nama }}
            </option>
          </select>

          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
            class="pointer-events-none absolute right-4 top-1/2 size-4 -translate-y-1/2 text-gray-500"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m19.5 8.25-7.5 7.5-7.5-7.5"
            />
          </svg>
        </div>

        <!-- Tahun Akademik -->
        <div class="relative w-[240px]">
          <select
            v-model="selectedTahun"
            class="h-[52px] w-full appearance-none rounded-xl border border-[#cbd5e1] bg-white px-4 pr-10 text-[15px] text-gray-600 outline-none focus:border-[#2846a3]"
          >
            <option value="">Pilih Tahun Akademik</option>

            <option
              v-for="item in tahunAkademikList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.tahun_awal }} - {{ item.tahun_akhir }}
            </option>
          </select>

          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
            class="pointer-events-none absolute right-4 top-1/2 size-4 -translate-y-1/2 text-gray-500"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="m19.5 8.25-7.5 7.5-7.5-7.5"
            />
          </svg>
        </div>

        <!-- BUTTON -->
        <button
          @click="handleFilter"
          class="flex h-[52px] items-center gap-2 rounded-xl bg-[#2447a8] px-6 font-semibold text-white transition hover:bg-[#1d3b91]"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="2"
            stroke="currentColor"
            class="size-5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M3 4.5h18m-15 6h12m-9 6h6"
            />
          </svg>

          Pilih
        </button>
      </div>

      <!-- TABLE -->
      <div class="overflow-x-auto">
        <table class="w-full">

          <!-- HEAD -->
          <thead>
            <tr class="text-left text-[15px] font-semibold text-[#555]">
              <th class="pb-4">No</th>
              <th class="pb-4">NIM</th>
              <th class="pb-4">Nama Mahasiswa</th>
              <th class="pb-4">Nama Dosen</th>
              <th class="pb-4">Kode Matakuliah</th>
              <th class="pb-4">Grade</th>
              <th class="pb-4">Nilai</th>
            </tr>
          </thead>

          <!-- BODY -->
          <tbody>

            <!-- EMPTY -->
            <tr v-if="nilaiList.length === 0">
              <td
                colspan="7"
                class="py-12 text-center text-gray-400"
              >
                Tidak ada data
              </td>
            </tr>

            <!-- DATA -->
            <tr
              v-for="(item, index) in nilaiList"
              :key="item.id"
              class="text-[15px] text-[#444]"
            >
              <td class="py-4">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <td class="py-4">
                {{ item.nim }}
              </td>

              <td class="py-4 font-medium">
                {{ item.nama_mahasiswa }}
              </td>

              <td class="py-4">
                {{ item.nama_dosen }}
              </td>

              <td class="py-4">
                {{ item.kode_matakuliah }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.grade }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.nilai }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAGINATION -->
      <div class="mt-28 flex items-center justify-between">

        <!-- SELECT -->
        <select
          v-model.number="perPage"
          @change="() => { currentPage = 1; getNilai() }"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select>

        <!-- PAGINATION BUTTON -->
        <div class="flex items-center gap-2">

          <!-- PREV -->
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="text-sm text-gray-400"
          >
            ← Previous
          </button>

          <!-- PAGE -->
          <template v-for="p in pages" :key="p">

            <span
              v-if="p === '...'"
              class="px-2 text-gray-400"
            >
              ...
            </span>

            <button
              v-else
              @click="goToPage(p as number)"
              class="flex h-9 w-9 items-center justify-center rounded-lg text-sm font-medium"
              :class="
                currentPage === p
                  ? 'bg-[#2447a8] text-white'
                  : 'text-gray-600 hover:bg-gray-100'
              "
            >
              {{ p }}
            </button>

          </template>

          <!-- NEXT -->
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="text-sm text-gray-600"
          >
            Next →
          </button>

        </div>
      </div>
    </div>
  </div>
</template>