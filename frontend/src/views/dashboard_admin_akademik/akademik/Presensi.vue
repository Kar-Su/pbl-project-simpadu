<script setup lang="ts">
import { ref, computed, onMounted } from "vue"

// ================= INTERFACE =================
interface Dosen {
  id: string
  nip: string
  name: string
  email: string
  role_name: string
  status: string
}

// ================= STATE =================
const dosenList = ref<Dosen[]>([])

const search = ref("")
const perPage = ref(5)
const currentPage = ref(1)

// tanggal hari ini
const today = new Date().toISOString().split("T")[0]

const selectedDate = ref(today)

// ================= GET DOSEN =================
const getDosen = async () => {
  try {
    const token = localStorage.getItem("token")

    const response = await fetch("/api/roles", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        Authorization: `Bearer ${token}`
      }
    })

    const text = await response.text()
const result = JSON.parse(text.replace(/^[^{]*/, ""))
    console.log("DOSEN API:", result)

    const data = Array.isArray(result.data?.items)
  ? result.data.items
  : []

    dosenList.value = data
      .filter((item: any) => {
        const role =
          item.role_name?.toLowerCase() ||
          item.role?.name?.toLowerCase()

        return role?.includes("dosen")
      })
      .map((item: any) => ({
        id: item.id ?? "",
        nip: item.nip ?? item.nomor_induk ?? "-",
        name: item.name ?? item.nama ?? "-",
        email: item.email ?? "-",
        role_name: item.role_name ?? item.role?.name ?? "Dosen",
        status: item.status ?? "Hadir"
      }))

  } catch (error) {
    console.error("GET DOSEN ERROR:", error)
    dosenList.value = []
  }
}
// ================= UPDATE STATUS =================
const updateStatus = (
  index: number,
  value: string
) => {

  dosenList.value[index].status = value

}

// ================= FILTER =================
const filteredDosen = computed(() => {

  return dosenList.value.filter((item) =>

    item.name?.toLowerCase().includes(search.value.toLowerCase()) ||
    item.email?.toLowerCase().includes(search.value.toLowerCase()) ||
    item.nip?.toLowerCase().includes(search.value.toLowerCase())

  )

})

// ================= PAGINATION =================
const totalPages = computed(() => {

  return Math.ceil(
    filteredDosen.value.length / perPage.value
  )

})

const paginatedDosen = computed(() => {

  const start =
    (currentPage.value - 1) * perPage.value

  return filteredDosen.value.slice(
    start,
    start + perPage.value
  )

})

const nextPage = () => {

  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }

}

const prevPage = () => {

  if (currentPage.value > 1) {
    currentPage.value--
  }

}

// ================= TOTAL =================
const totalHadir = computed(() =>
  dosenList.value.filter(
    (item) => item.status === "Hadir"
  ).length
)

const totalIzin = computed(() =>
  dosenList.value.filter(
    (item) =>
      item.status === "Izin" ||
      item.status === "Sakit"
  ).length
)

const totalTidakHadir = computed(() =>
  dosenList.value.filter(
    (item) => item.status === "Tidak Hadir"
  ).length
)

// ================= STATUS COLOR =================
const getStatusClass = (status: string) => {

  switch (status) {

    case "Hadir":
      return "bg-green-500 text-white"

    case "Izin":
      return "bg-yellow-400 text-white"

    case "Sakit":
      return "bg-blue-500 text-white"

    case "Tidak Hadir":
      return "bg-red-500 text-white"

    default:
      return "bg-gray-300"

  }

}

// ================= MOUNTED =================
onMounted(() => {
  getDosen()
})
</script>

<template>

  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Presensi
    </div>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-gray-800">
      Presensi
    </h1>

    <p class="mb-8 mt-1 text-gray-500">
      Data Presensi Pegawai
    </p>

    <!-- CARD -->
    <div class="mb-6 grid grid-cols-1 gap-5 md:grid-cols-3">

      <!-- HADIR -->
      <div class="rounded-2xl border border-blue-100 bg-white p-5 shadow-sm">

        <p class="text-sm font-semibold text-gray-500">
          TOTAL HADIR
        </p>

        <h2 class="mt-2 text-3xl font-bold text-gray-800">
          {{ totalHadir }}
        </h2>

      </div>

      <!-- IZIN -->
      <div class="rounded-2xl border border-blue-100 bg-white p-5 shadow-sm">

        <p class="text-sm font-semibold text-gray-500">
          TOTAL IZIN / SAKIT
        </p>

        <h2 class="mt-2 text-3xl font-bold text-gray-800">
          {{ totalIzin }}
        </h2>

      </div>

      <!-- TIDAK HADIR -->
      <div class="rounded-2xl border border-blue-100 bg-white p-5 shadow-sm">

        <p class="text-sm font-semibold text-gray-500">
          TOTAL TIDAK HADIR
        </p>

        <h2 class="mt-2 text-3xl font-bold text-gray-800">
          {{ totalTidakHadir }}
        </h2>

      </div>

    </div>

    <!-- SEARCH & DATE -->
    <div class="mb-5 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">

      <!-- SEARCH -->
      <div class="relative">

        <input v-model="search" type="text" placeholder="Cari Dosen..."
          class="w-72 rounded-xl border border-gray-200 bg-white py-3 pl-4 pr-10 text-sm outline-none focus:border-blue-500" />

      </div>

      <!-- DATE -->
      <div class="flex items-center gap-3">

        <label class="text-sm font-medium text-gray-600">
          Pilih Tanggal:
        </label>

        <input v-model="selectedDate" type="date"
          class="rounded-xl border-2 border-blue-500 bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700 outline-none" />

      </div>

    </div>

    <!-- TABLE -->
    <div class="overflow-hidden rounded-2xl border border-blue-100 bg-white shadow-sm">

      <!-- HEADER -->
      <div class="border-b border-gray-100 px-6 py-5">

        <h2 class="text-2xl font-semibold text-gray-700">
          Data Presensi
        </h2>

      </div>

<!-- TABLE -->
<div class="overflow-x-auto">

  <table class="w-full text-sm">

    <thead>
      <tr class="text-left text-gray-500">

        <th class="px-6 py-5">No</th>
        <th class="px-6 py-5">NIP</th>
        <th class="px-6 py-5">Nama</th>
        <th class="px-6 py-5">Email</th>
        <th class="px-6 py-5">Role</th>
        <th class="px-6 py-5">Tanggal</th>
        <th class="px-6 py-5">Kehadiran</th>

      </tr>
    </thead>

    <tbody>

      <tr
        v-for="(item, index) in paginatedDosen"
        :key="item.id"
        class="border-t border-gray-100"
      >

        <!-- NO -->
        <td class="px-6 py-5">
          {{ (currentPage - 1) * perPage + index + 1 }}
        </td>

        <!-- NIP -->
        <td class="px-6 py-5">
          {{ item.nip || "-" }}
        </td>

        <!-- NAMA -->
        <td class="px-6 py-5 font-medium text-gray-700">
          {{ item.name }}
        </td>

        <!-- EMAIL -->
        <td class="px-6 py-5">
          {{ item.email }}
        </td>

        <!-- ROLE -->
        <td class="px-6 py-5">
          {{ item.role_name }}
        </td>

        <!-- TANGGAL -->
        <td class="px-6 py-5">

          <span
            class="rounded-lg bg-blue-100 px-3 py-2 text-xs font-semibold text-blue-700"
          >
            {{ selectedDate }}
          </span>

        </td>

        <!-- STATUS -->
        <td class="px-6 py-5">

          <select
            :value="item.status"
            @change="updateStatus(index, ($event.target as HTMLSelectElement).value)"
            :class="[
              'rounded-lg px-3 py-2 text-xs font-semibold outline-none',
              getStatusClass(item.status)
            ]"
          >

            <option value="Hadir">
              Hadir
            </option>

            <option value="Izin">
              Izin
            </option>

            <option value="Sakit">
              Sakit
            </option>

            <option value="Tidak Hadir">
              Tidak Hadir
            </option>

          </select>

        </td>

      </tr>

      <!-- KOSONG -->
      <tr v-if="paginatedDosen.length === 0">

        <td colspan="7" class="px-6 py-10 text-center text-gray-400">
          Tidak ada data dosen
        </td>

      </tr>

    </tbody>

  </table>

</div>

      <!-- FOOTER -->
      <div class="flex items-center justify-between border-t border-gray-100 px-6 py-4">

        <!-- PER PAGE -->
        <select v-model="perPage" class="rounded-lg border border-gray-300 px-3 py-2 text-sm outline-none">
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="20">20 Baris</option>
        </select>

        <!-- PAGINATION -->
        <div class="flex items-center gap-4 text-sm">

          <button @click="prevPage" class="text-gray-400 hover:text-black">
            ← Previous
          </button>

          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#2f4a8a] text-white">
            {{ currentPage }}
          </div>

          <button @click="nextPage" class="text-gray-400 hover:text-black">
            Next →
          </button>

        </div>

      </div>

    </div>

  </div>
</template>