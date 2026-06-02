<script setup lang="ts">
import { ref, computed, onMounted } from "vue"

// ================= FILTER =================
const kurikulum = ref("")
const matakuliah = ref("")
const tahunAkademik = ref("")

// ================= DATA API =================
interface NilaiItem {
  id: number
  nama_kelas: string
  jurusan: string
  prodi: string
  tahun_akademik: string
}

interface KurikulumItem {
  id: number
  name: string
}

interface TahunAkademikItem {
  id: number
  tahun_awal: string
  tahun_akhir: string
}

const nilaiData = ref<NilaiItem[]>([])
const kurikulumList = ref<KurikulumItem[]>([])
const tahunAkademikList = ref<TahunAkademikItem[]>([])

// ================= HELPER =================
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= GET KURIKULUM =================
const getKurikulum = async (): Promise<void> => {
  try {
    const res = await fetch("/api/kurikulum/", {
      headers: getHeaders(),
    })

    const data = await res.json()

    kurikulumList.value = data.data ?? []
  } catch (err) {
    console.error("GET KURIKULUM:", err)
  }
}

// ================= GET TAHUN AKADEMIK =================
const getTahunAkademik = async (): Promise<void> => {
  try {
    const res = await fetch("/api/tahun-akademik", {
      headers: getHeaders(),
    })

    const data = await res.json()

    tahunAkademikList.value = data.data ?? []
  } catch (err) {
    console.error("GET TAHUN AKADEMIK:", err)
  }
}

// ================= GET NILAI =================
const getNilai = async (): Promise<void> => {
  try {
    const res = await fetch("/api/nilai", {
      headers: getHeaders(),
    })

    const data = await res.json()

    nilaiData.value = data.data ?? []
  } catch (err) {
    console.error("GET NILAI:", err)
  }
}

// ================= FILTER ACTION =================
const pilihData = (): void => {
  console.log({
    kurikulum: kurikulum.value,
    matakuliah: matakuliah.value,
    tahunAkademik: tahunAkademik.value,
  })
}

// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(5)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value

  return nilaiData.value.slice(start, start + perPage.value)
})

// ================= DETAIL =================
const lihatDetail = (item: NilaiItem): void => {
  console.log("DETAIL:", item)
}

// ================= ON MOUNTED =================
onMounted(() => {
  getKurikulum()
  getTahunAkademik()
  getNilai()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Nilai
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Nilai
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-6">
      Kelola nilai tiap mahasiswa
    </p>

    <!-- FILTER CARD -->
    <div
      class="bg-white border border-[#d5e1f0] rounded-2xl p-5 shadow-[0_4px_12px_rgba(0,0,0,0.08)] mb-4"
    >

      <h2 class="text-[28px] font-semibold text-[#202020] mb-5">
        Pilih Matakuliah
      </h2>

      <div class="flex items-center gap-5 flex-wrap">

        <!-- KURIKULUM -->
        <select
          v-model="kurikulum"
          class="w-[310px] h-[58px] border border-gray-400 rounded-2xl px-5 text-[18px] text-gray-600 outline-none bg-white"
        >
          <option value="">
            Pilih Kurikulum
          </option>

          <option
            v-for="item in kurikulumList"
            :key="item.id"
            :value="item.id"
          >
            {{ item.name }}
          </option>
        </select>

        <!-- MATAKULIAH -->
        <select
          v-model="matakuliah"
          class="w-[310px] h-[58px] border border-gray-400 rounded-2xl px-5 text-[18px] text-gray-600 outline-none bg-white"
        >
          <option value="">
            Pilih Matakuliah
          </option>

          <option value="Algoritma">
            Algoritma
          </option>

          <option value="Basis Data">
            Basis Data
          </option>
        </select>

        <!-- TAHUN AKADEMIK -->
        <select
          v-model="tahunAkademik"
          class="w-[310px] h-[58px] border border-gray-400 rounded-2xl px-5 text-[18px] text-gray-600 outline-none bg-white"
        >
          <option value="">
            Pilih Tahun Akademik
          </option>

          <option
            v-for="item in tahunAkademikList"
            :key="item.id"
            :value="item.id"
          >
            {{ item.tahun_awal }} - {{ item.tahun_akhir }}
          </option>
        </select>

        <!-- BUTTON -->
        <button
          @click="pilihData"
          class="h-[58px] px-8 bg-[#29479d] hover:bg-[#1f377b] rounded-2xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Pilih
        </button>

      </div>
    </div>

    <!-- TABLE CARD -->
    <div
      class="bg-white border border-[#d5e1f0] rounded-2xl min-h-[650px] shadow-[0_4px_12px_rgba(0,0,0,0.08)] overflow-hidden"
    >

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[36px] font-semibold text-[#202020]">
          Data Nilai
        </h2>
      </div>

      <!-- TABLE -->
      <div class="px-5 pt-8">

        <table class="w-full border-separate border-spacing-y-6">

          <!-- HEADER -->
          <thead>
            <tr class="text-left text-gray-600">
              <th class="text-[20px] font-semibold">No</th>
              <th class="text-[20px] font-semibold">Nama Kelas</th>
              <th class="text-[20px] font-semibold">Jurusan</th>
              <th class="text-[20px] font-semibold">Prodi</th>
              <th class="text-[20px] font-semibold">Tahun Akademik</th>
              <th class="text-[20px] font-semibold text-center">Aksi</th>
            </tr>
          </thead>

          <!-- BODY -->
          <tbody>

            <tr
              v-if="paginatedData.length === 0"
            >
              <td
                colspan="6"
                class="text-center py-10 text-gray-400 text-lg"
              >
                Tidak ada data
              </td>
            </tr>

            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="text-[#505050]"
            >

              <td class="text-[20px]">
                {{ index + 1 }}
              </td>

              <td class="text-[20px] font-medium">
                {{ item.nama_kelas }}
              </td>

              <td class="text-[20px] font-medium">
                {{ item.jurusan }}
              </td>

              <td class="text-[20px] font-medium">
                {{ item.prodi }}
              </td>

              <td class="text-[20px] font-medium">
                {{ item.tahun_akademik }}
              </td>

              <!-- AKSI -->
              <td class="text-center">

                <button
                  @click="lihatDetail(item)"
                  class="bg-[#29479d] hover:bg-[#1f377b] text-white px-6 py-2 rounded-xl text-[17px] font-semibold shadow-md transition"
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
        class="flex items-center justify-between px-5 pt-28 pb-5"
      >

        <!-- BARIS -->
        <select
          v-model="perPage"
          class="w-[100px] h-[44px] border border-gray-300 rounded-xl px-3 text-sm outline-none"
        >
          <option :value="5">
            5 Baris
          </option>

          <option :value="10">
            10 Baris
          </option>

          <option :value="25">
            25 Baris
          </option>
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