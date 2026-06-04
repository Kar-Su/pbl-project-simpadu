<script setup lang="ts">
import { ref, computed } from "vue"

// ================== FILTER ==================
const semester = ref("")
const tahunAkademik = ref("")

// ================== DATA ==================
const kelasData = ref([
  {
    id: 1,
    kelas: "4A",
    jurusan: "Teknik Elektro",
    prodi: "Teknik Informatika",
    tahun: "2023-2024",
  },
  {
    id: 2,
    kelas: "4B",
    jurusan: "Teknik Elektro",
    prodi: "SIKC",
    tahun: "2023-2024",
  },
])

// ================== PAGINATION ==================
const currentPage = ref(1)
const perPage = ref(5)

const BASE_URL = "https://be.karlearn.site"

const totalPages = computed(() => 68)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return kelasData.value.slice(start, start + perPage.value)
})

// ================== ACTION ==================
const pilihFilter = () => {
  console.log("Semester:", semester.value)
  console.log("Tahun:", tahunAkademik.value)
}

const lihatDetail = (item: any) => {
  console.log("DETAIL:", item)
}
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-1 text-sm text-gray-500 font-medium">
      Akademik > KHS
    </div>

    <!-- TITLE -->
    <h1 class="text-[40px] font-bold text-[#404040] leading-none">
      Kartu Hasil Studi
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-5">
      Data hasil studi mahasiswa
    </p>

    <!-- FILTER -->
    <div
      class="bg-white rounded-xl border border-[#d7e3f3] shadow-[0_4px_10px_rgba(0,0,0,0.08)] p-5 mb-4"
    >
      <h2 class="text-[38px] font-semibold text-[#404040] mb-5">
        Pilih Semester
      </h2>

      <div class="flex items-center gap-3">

        <!-- SEMESTER -->
        <select
          v-model="semester"
          class="w-[260px] h-[56px] border border-gray-300 rounded-xl px-4 text-[20px] text-gray-600 bg-white outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Pilih Semester</option>
          <option>Semester 1</option>
          <option>Semester 2</option>
          <option>Semester 3</option>
        </select>

        <!-- TAHUN AKADEMIK -->
        <select
          v-model="tahunAkademik"
          class="w-[260px] h-[56px] border border-gray-300 rounded-xl px-4 text-[20px] text-gray-600 bg-white outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Pilih Tahun Akademik</option>
          <option>2023-2024</option>
          <option>2024-2025</option>
        </select>

        <!-- BUTTON -->
        <button
          @click="pilihFilter"
          class="h-[56px] px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Pilih
        </button>

      </div>
    </div>

    <!-- TABLE CARD -->
    <div
      class="bg-white rounded-2xl border border-[#d7e3f3] shadow-[0_4px_12px_rgba(0,0,0,0.08)] overflow-hidden"
    >

      <!-- HEADER -->
      <div class="px-6 pt-5 pb-3">
        <h2 class="text-[38px] font-semibold text-[#404040]">
          Data Kelas
        </h2>
      </div>

      <!-- TABLE -->
      <div class="px-5">

        <table class="w-full border-separate border-spacing-y-4">

          <thead>
            <tr class="text-left text-gray-600">
              <th class="font-semibold text-[18px]">No</th>
              <th class="font-semibold text-[18px]">Nama Kelas</th>
              <th class="font-semibold text-[18px]">Jurusan</th>
              <th class="font-semibold text-[18px]">Prodi</th>
              <th class="font-semibold text-[18px]">Tahun Akademik</th>
              <th class="font-semibold text-[18px] text-center">Aksi</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="text-[#4b4b4b]"
            >
              <td class="text-[18px] py-2">
                {{ index + 1 }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.kelas }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.jurusan }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.prodi }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.tahun }}
              </td>

              <td class="text-center">
                <button
                  @click="lihatDetail(item)"
                  class="bg-[#29479d] hover:bg-[#1d377f] text-white px-5 py-2 rounded-xl text-[16px] font-semibold shadow-md transition"
                >
                  Lihat
                </button>
              </td>
            </tr>
          </tbody>

        </table>
      </div>

      <!-- FOOTER -->
      <div
        class="flex items-center justify-between px-5 py-6 mt-10"
      >

        <!-- SELECT BARIS -->
        <select
          v-model="perPage"
          class="w-[90px] h-[42px] border border-gray-300 rounded-lg px-3 text-sm outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select>

        <!-- PAGINATION -->
        <div class="flex items-center gap-5 text-gray-500">

          <button class="text-gray-400">
            ← Previous
          </button>

          <button
            class="w-9 h-9 rounded-lg bg-[#29479d] text-white font-semibold"
          >
            1
          </button>

          <button>2</button>
          <button>3</button>

          <span>...</span>

          <button>67</button>
          <button>68</button>

          <button class="text-gray-700 font-medium">
            Next →
          </button>

        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
table {
  border-collapse: separate;
}
</style>