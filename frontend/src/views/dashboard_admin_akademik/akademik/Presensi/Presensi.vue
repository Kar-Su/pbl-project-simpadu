<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"

// ================= INTERFACE =================
interface Pegawai {
  detail_id: string
  name: string
  email: string
  status: string
}

interface Sesi {
  sesi_id: string
  pegawai: Pegawai[] | null
  created_at: string
}

// ================= STATE =================
const sesiList = ref<Sesi[]>([])
const pegawaiHariIni = ref<Pegawai[]>([])
const sesiIdHariIni = ref<string>("")

const search = ref("")
const perPage = ref(5)
const currentPage = ref(1)

const today = new Date().toISOString().split("T")[0]
const selectedDate = ref(today)

const isLoading = ref(false)
const isSaving = ref(false)
const saveMessage = ref("")

// ================= HEADERS =================
const getHeaders = () => ({
  "Content-Type": "application/json",
  "Accept": "application/json",
  Authorization: `Bearer ${localStorage.getItem("token")}`
})

// ================= GET PRESENSI =================
const getPresensi = async () => {
  try {
    isLoading.value = true
    const res = await fetch(
      "https://be.karlearn.site/api/presensi/pegawai?page=1",
      { headers: getHeaders() }
    )
    const data = await res.json()
    console.log("PRESENSI:", data)

    sesiList.value = data.data?.items || []
    filterByDate()
  } catch (err) {
    console.error("GET PRESENSI ERROR:", err)
    sesiList.value = []
  } finally {
    isLoading.value = false
  }
}

// ================= FILTER BY DATE =================
// PERBAIKAN: bisa ada lebih dari satu sesi untuk tanggal yang sama
// (misal pegawai baru ditambahkan di sesi berikutnya pada hari yang sama).
// Sebelumnya pakai .find() yang hanya ambil sesi PERTAMA yang cocok,
// jadi pegawai yang baru ditambahkan di sesi berikutnya tidak pernah muncul.
// Sekarang: gabungkan pegawai dari SEMUA sesi tanggal tersebut,
// dedupe berdasarkan detail_id (data dari sesi yang lebih baru menang),
// dan sesi_id aktif diambil dari sesi PALING BARU (terakhir) pada tanggal itu.
const filterByDate = () => {
  currentPage.value = 1

  const sesiSamaTanggal = sesiList.value.filter(
    (s) => s.created_at === selectedDate.value
  )

  if (sesiSamaTanggal.length > 0) {
    const sesiTerakhir = sesiSamaTanggal[sesiSamaTanggal.length - 1]
    sesiIdHariIni.value = sesiTerakhir.sesi_id

    const pegawaiMap = new Map<string, Pegawai>()
    sesiSamaTanggal.forEach((s) => {
      ;(s.pegawai || []).forEach((p) => {
        // deep copy agar perubahan status tidak langsung mutate sumber
        pegawaiMap.set(p.detail_id, { ...p })
      })
    })

    pegawaiHariIni.value = Array.from(pegawaiMap.values())
  } else {
    sesiIdHariIni.value = ""
    pegawaiHariIni.value = []
  }
}

// Watch tanggal berubah → filter ulang
watch(selectedDate, () => {
  filterByDate()
})

// Reset ke halaman 1 setiap kali kata kunci pencarian berubah,
// supaya tidak "nyangkut" di halaman yang sudah tidak ada datanya
watch(search, () => {
  currentPage.value = 1
})

// Reset ke halaman 1 setiap kali jumlah baris per halaman berubah
watch(perPage, () => {
  currentPage.value = 1
})

// ================= UPDATE STATUS (lokal) =================
// PERBAIKAN: sebelumnya pakai perhitungan index global
// (currentPage - 1) * perPage + index, padahal tabel yang ditampilkan
// berasal dari filteredPegawai (hasil pencarian), bukan langsung dari
// pegawaiHariIni. Begitu search diisi, index jadi tidak sinkron dan
// status yang ter-update bisa jadi milik pegawai yang salah.
// Sekarang update langsung berdasarkan detail_id, jadi selalu tepat
// sasaran terlepas dari filter/pencarian/halaman yang aktif.
const updateStatus = (detailId: string, value: string) => {
  const target = pegawaiHariIni.value.find((p) => p.detail_id === detailId)
  if (target) {
    target.status = value
  }
}

// ================= SAVE / PUT =================
const savePresensi = async () => {
  if (!selectedDate.value) {
    saveMessage.value = "Pilih tanggal terlebih dahulu"
    return
  }

  if (pegawaiHariIni.value.length === 0) {
    saveMessage.value = "Tidak ada data pegawai untuk disimpan"
    return
  }

  try {
    isSaving.value = true
    saveMessage.value = ""

    const body = {
      date: selectedDate.value,
      detail: pegawaiHariIni.value.map((p) => ({
        detail_id: p.detail_id,
        status: p.status
      }))
    }

    console.log("PUT BODY:", body)

    const res = await fetch(
      "https://be.karlearn.site/api/presensi/pegawai",
      {
        method: "PUT",
        headers: getHeaders(),
        body: JSON.stringify(body)
      }
    )

    const data = await res.json()
    console.log("PUT RESPONSE:", data)

    if (!res.ok) {
      saveMessage.value = data.message || "Gagal menyimpan presensi"
      return
    }

    saveMessage.value = "✅ Presensi berhasil disimpan"
    // Refresh data
    await getPresensi()
  } catch (err) {
    console.error("PUT PRESENSI ERROR:", err)
    saveMessage.value = "Terjadi error jaringan"
  } finally {
    isSaving.value = false
  }
}

// ================= FILTER SEARCH =================
const filteredPegawai = computed(() => {
  const keyword = search.value.toLowerCase().trim()
  if (!keyword) return pegawaiHariIni.value
  return pegawaiHariIni.value.filter(
    (item) =>
      item.name?.toLowerCase().includes(keyword) ||
      item.email?.toLowerCase().includes(keyword)
  )
})

// ================= PAGINATION =================
const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredPegawai.value.length / perPage.value))
)

// PERBAIKAN: jaga-jaga supaya currentPage tidak pernah melebihi totalPages
// (misalnya setelah data berkurang karena ganti tanggal/sesi)
watch(totalPages, (newTotal) => {
  if (currentPage.value > newTotal) {
    currentPage.value = newTotal
  }
})

const paginatedPegawai = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredPegawai.value.slice(start, start + perPage.value)
})

const nextPage = () => {
  if (currentPage.value < totalPages.value) currentPage.value++
}

const prevPage = () => {
  if (currentPage.value > 1) currentPage.value--
}

// ================= TOTAL =================
const totalHadir = computed(() =>
  pegawaiHariIni.value.filter((p) => p.status === "hadir").length
)

const totalIzin = computed(() =>
  pegawaiHariIni.value.filter((p) => p.status === "izin").length
)

const totalSakit = computed(() =>
  pegawaiHariIni.value.filter((p) => p.status === "sakit").length
)

const totalAlpha = computed(() =>
  pegawaiHariIni.value.filter((p) => p.status === "alpha").length
)

// ================= STATUS COLOR =================
const getStatusClass = (status: string) => {
  switch (status) {
    case "hadir":   return "bg-green-500 text-white"
    case "izin":    return "bg-yellow-400 text-white"
    case "sakit":   return "bg-blue-500 text-white"
    case "alpha":   return "bg-red-500 text-white"
    default:        return "bg-gray-300 text-gray-700"
  }
}

// ================= MOUNTED =================
onMounted(() => {
  getPresensi()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Presensi
    </div>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-gray-800">Presensi</h1>
    <p class="mb-8 mt-1 text-gray-500">Data Presensi Pegawai</p>

<!-- STAT CARDS -->
<div class="mb-6 grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-4">

  <!-- HADIR -->
  <div class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
    <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-[#9db9dc]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-[#4b4b4b]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M9 12l2 2 4-4M12 21a9 9 0 100-18 9 9 0 000 18z" />
      </svg>
    </div>

    <div>
      <p class="text-xs font-bold text-gray-600 uppercase">TOTAL HADIR</p>
      <h2 class="mt-1 text-3xl font-bold text-gray-800">{{ totalHadir }}</h2>
    </div>
  </div>

  <!-- IZIN -->
  <div class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#facc15]">
    <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-[#fde68a]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-[#4b4b4b]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
    </div>

    <div>
      <p class="text-xs font-bold text-gray-600 uppercase">TOTAL IZIN</p>
      <h2 class="mt-1 text-3xl font-bold text-gray-800">{{ totalIzin }}</h2>
    </div>
  </div>

  <!-- SAKIT -->
  <div class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#3b82f6]">
    <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-[#bfdbfe]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-[#4b4b4b]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M19.428 15.428a4 4 0 00-5.656-5.656L12 11.544 10.228 9.77a4 4 0 00-5.656 5.657l1.772 1.772L12 22.856l5.656-5.656 1.772-1.772z" />
      </svg>
    </div>

    <div>
      <p class="text-xs font-bold text-gray-600 uppercase">TOTAL SAKIT</p>
      <h2 class="mt-1 text-3xl font-bold text-gray-800">{{ totalSakit }}</h2>
    </div>
  </div>

  <!-- ALPHA -->
  <div class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#ef4444]">
    <div class="flex h-16 w-16 items-center justify-center rounded-lg bg-[#fecaca]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-[#4b4b4b]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M18.364 5.636l-12.728 12.728M5.636 5.636l12.728 12.728" />
      </svg>
    </div>

    <div>
      <p class="text-xs font-bold text-gray-600 uppercase">TOTAL ALPHA</p>
      <h2 class="mt-1 text-3xl font-bold text-gray-800">{{ totalAlpha }}</h2>
    </div>
  </div>

</div>

    <!-- SEARCH & DATE -->
    <div class="mb-5 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">

      <!-- SEARCH -->
      <div class="relative">
        <input
          v-model="search"
          type="text"
          placeholder="Cari Pegawai..."
          class="w-72 rounded-xl border border-gray-200 bg-white py-3 pl-4 pr-10 text-sm outline-none focus:border-blue-500"
        />
      </div>

      <!-- DATE -->
      <div class="flex items-center gap-3">
        <label class="text-sm font-medium text-gray-600">Pilih Tanggal:</label>
        <input
          v-model="selectedDate"
          type="date"
          class="rounded-xl border-2 border-blue-500 bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700 outline-none"
        />
      </div>

    </div>

    <!-- TABLE -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

      <!-- TABLE HEADER (biru) -->
      <div class="bg-[#243e90] px-5 py-4">
        <h2 class="text-white text-2xl font-bold">Data Presensi</h2>
        <p class="text-white text-sm mt-1">Presensi pegawai tanggal {{ selectedDate }}</p>
      </div>

      <!-- LOADING -->
      <div v-if="isLoading" class="py-10 text-center text-gray-400">
        Memuat data...
      </div>

      <!-- TABLE CONTENT -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">

          <thead>
            <tr class="text-left text-gray-600 border-b border-gray-300">
              <th class="px-6 py-5">No</th>
              <th class="px-6 py-5">Nama</th>
              <th class="px-6 py-5">Email</th>
              <th class="px-6 py-5">Tanggal</th>
              <th class="px-6 py-5">Kehadiran</th>
            </tr>
          </thead>

          <tbody>

            <tr
              v-for="(item, index) in paginatedPegawai"
              :key="item.detail_id"
              class="border-t border-gray-100"
            >

              <!-- NO -->
              <td class="px-6 py-5">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <!-- NAMA -->
              <td class="px-6 py-5 font-medium text-gray-700">
                {{ item.name }}
              </td>

              <!-- EMAIL -->
              <td class="px-6 py-5">{{ item.email }}</td>

              <!-- TANGGAL -->
              <td class="px-6 py-5">
                <span class="rounded-lg bg-blue-100 px-3 py-2 text-xs font-semibold text-blue-700">
                  {{ selectedDate }}
                </span>
              </td>

              <!-- STATUS DROPDOWN -->
              <td class="px-6 py-5">
                <select
                  :value="item.status"
                  @change="updateStatus(item.detail_id, ($event.target as HTMLSelectElement).value)"
                  :class="['rounded-lg px-3 py-2 text-xs font-semibold outline-none cursor-pointer', getStatusClass(item.status)]"
                >
                  <option value="hadir">Hadir</option>
                  <option value="izin">Izin</option>
                  <option value="sakit">Sakit</option>
                  <option value="alpha">Alpha</option>
                </select>
              </td>

            </tr>

            <!-- EMPTY STATE -->
            <tr v-if="paginatedPegawai.length === 0">
              <td colspan="5" class="px-6 py-10 text-center text-gray-400">
                {{
                  sesiList.length === 0
                    ? "Belum ada data presensi"
                    : `Tidak ada data presensi untuk tanggal ${selectedDate}`
                }}
              </td>
            </tr>

          </tbody>

        </table>
      </div>

      <!-- FOOTER -->
      <div class="flex items-center justify-between border-t border-gray-100 px-6 py-4">

        <!-- KIRI: per page -->
        <select
          v-model.number="perPage"
          class="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="20">20 Baris</option>
        </select>

        <!-- KANAN: pagination -->
        <div class="flex items-center gap-4 text-sm">
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="text-gray-400 hover:text-black disabled:opacity-40"
          >
            ← Previous
          </button>

          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#2f4a8a] text-white">
            {{ currentPage }}
          </div>

          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="text-gray-400 hover:text-black disabled:opacity-40"
          >
            Next →
          </button>
        </div>

      </div>

    </div>

    <!-- SAVE BUTTON & MESSAGE -->
    <div class="mt-5 flex items-center gap-4">

      <button
        @click="savePresensi"
        :disabled="isSaving || pegawaiHariIni.length === 0"
        class="flex items-center gap-2 rounded-lg bg-[#22c55e] px-6 py-3 text-sm font-semibold text-white hover:bg-[#16a34a] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <svg v-if="!isSaving" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        {{ isSaving ? "Menyimpan..." : "Simpan Presensi" }}
      </button>

      <!-- PESAN HASIL -->
      <p v-if="saveMessage" :class="saveMessage.startsWith('✅') ? 'text-green-600' : 'text-red-500'" class="text-sm font-medium">
        {{ saveMessage }}
      </p>

    </div>

  </div>
</template>