<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= STATE =================
const semester = ref("")
const tahunAwal = ref("")
const tahunAkhir = ref("")

const loading = ref(false)

// ================= HEADER =================
const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= SIMPAN DATA =================
const simpanData = async () => {
  try {
    // VALIDASI
    if (!semester.value || !tahunAwal.value || !tahunAkhir.value) {
      alert("Semua field wajib diisi")
      return
    }

    loading.value = true

    const res = await fetch("/api/tahun-akademik", {
      method: "POST",

      headers: getHeaders(),

body: JSON.stringify({
  id: Number(
    `${tahunAwal.value}${semester.value === 'ganjil' ? '1' : '2'}`
  ),

  tipe_semester: semester.value,

  tahun_awal:
    semester.value === 'ganjil'
      ? `${tahunAwal.value}-01-01`
      : `${tahunAwal.value}-07-01`,

  tahun_akhir:
    semester.value === 'ganjil'
      ? `${tahunAwal.value}-06-30`
      : `${tahunAwal.value}-12-30`,

  status: "aktif",
}),
    })

    const json = await res.json()

    console.log("RESPONSE:", json)

    if (!res.ok) {
      alert("Gagal tambah data")
      return
    }

    alert("Berhasil tambah tahun akademik")

    // RESET FORM
    semester.value = ""
    tahunAwal.value = ""
    tahunAkhir.value = ""

    // PINDAH HALAMAN
    router.push("/dashboard-admin/tahun-akademik")

  } catch (err) {
    console.error("ERROR:", err)
    alert("Terjadi kesalahan")
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-6">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Tahun Akademik > Tambah Tahun Akademik
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Tahun Akademik
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-8">
      lorem ipsum
    </p>

    <!-- CARD -->
    <div
      class="bg-white rounded-xl border border-[#d6e2f1] shadow-[0_4px_10px_rgba(0,0,0,0.08)] overflow-visible max-w-[1040px]"
    >

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[34px] font-semibold text-[#505050]">
          Tambah Data Tahun Akademik
        </h2>
      </div>

      <!-- FORM -->
      <div class="px-5 pt-4 pb-6 flex items-start gap-5 relative">

        <!-- SEMESTER -->
        <div class="flex flex-col relative">

          <label class="text-[20px] font-medium text-[#555] mb-2">
            Tahun Akademik
          </label>

          <select
            v-model="semester"
            class="w-[320px] h-[58px] border border-gray-400 rounded-lg px-4 text-[18px] text-gray-600 bg-white outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Pilih Tahun Akademik</option>
            <option value="ganjil">Ganjil</option>
            <option value="genap">Genap</option>
          </select>

        </div>

        <!-- TAHUN AWAL -->
        <div class="flex flex-col">

          <label class="text-[20px] font-medium text-[#555] mb-2">
            Tahun Awal
          </label>

          <input
            v-model="tahunAwal"
            type="text"
            placeholder="isi Tahun Awal ..."
            class="w-[320px] h-[58px] border border-gray-400 rounded-lg px-4 text-[18px] outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <!-- TAHUN AKHIR -->
        <div class="flex flex-col">

          <label class="text-[20px] font-medium text-[#555] mb-2">
            Tahun Akhir
          </label>

          <input
            v-model="tahunAkhir"
            type="text"
            placeholder="Isi Tahun Akhir ..."
            class="w-[320px] h-[58px] border border-gray-400 rounded-lg px-4 text-[18px] outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
    </div>

    <!-- BUTTON -->
<button
  @click="simpanData"
  :disabled="loading"
  class="mt-5 flex items-center gap-2 bg-[#1fb85a] hover:bg-[#159548] disabled:bg-gray-400 text-white px-5 py-3 rounded-lg shadow-md transition font-semibold text-[18px]"
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    class="w-5 h-5"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    stroke-width="2"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      d="M5 13l4 4L19 7"
    />
  </svg>

  {{ loading ? 'Menyimpan...' : 'Simpan' }}
</button>
  </div>
</template>

<style scoped>
select {
  appearance: auto;
}
</style>