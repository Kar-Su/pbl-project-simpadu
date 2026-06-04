<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= FILTER =================
const jurusan = ref("")
const prodi = ref("")
const tahunAkademik = ref("")

// ================= DATA =================
const kurikulumData = ref<any[]>([])  

// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(10)

const totalPages = ref(1)
const totalItems = ref(0)

const BASE_URL = "https://be.karlearn.site"

// ================= HEADER =================
const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= GET KURIKULUM =================
const getKurikulum = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/kurikulum?page=${currentPage.value}&per_page=${perPage.value}`,
      {
        method: "GET",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    console.log("KURIKULUM:", json)

    const raw = json.data?.items ?? []

    kurikulumData.value = raw.map((item: any) => ({
      id: item.id,
      nama: item.name ?? "-",
      prodi: item.prodi?.name ?? "-",
      semester: item.kurikulum_mk?.[0]?.semester
        ? `Semester ${item.kurikulum_mk[0].semester}`
        : "-",
      tahun: "-",
    }))

    totalPages.value = json.data?.pagination?.total_pages ?? 1
    totalItems.value = json.data?.pagination?.total_items ?? 0
    perPage.value = json.data?.pagination?.per_page ?? perPage.value
  } catch (err) {
    console.error("GET KURIKULUM ERROR:", err)
  }
}

// ================= PAGINATION LOGIC =================

// NEXT
const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getKurikulum()
  }
}

// PREVIOUS
const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getKurikulum()
  }
}

// GO TO PAGE (NUMBER CLICK)
const goToPage = async (page: number | string) => {
  const p = Number(page)
  if (!p || p < 1 || p > totalPages.value) return

  currentPage.value = p
  await getKurikulum()
}

// ================= WATCH PER PAGE =================
watch(perPage, async () => {
  currentPage.value = 1
  await getKurikulum()
})

// ================= MOUNT =================
onMounted(() => {
  getKurikulum()
})

// ================= PAGINATION UI =================
const pages = computed<(number | string)[]>(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 5) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  if (current <= 3) return [1, 2, 3, "...", total]

  if (current >= total - 2)
    return [1, "...", total - 2, total - 1, total]

  return [1, "...", current, "...", total]
})

// ================= ACTION =================
const pilihData = () => {
  console.log({
    jurusan: jurusan.value,
    prodi: prodi.value,
    tahun: tahunAkademik.value,
  })
}

const tambahData = () => {
  router.push("/dashboard-admin/tambah_kurikulum")
}

const editData = (item: any) => {
  console.log("EDIT:", item)
}
</script>
```


<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Mahasiswa > Kurikulum
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Kurikulum
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-6">
      Pengelolaan Data
    </p>

    <!-- CARD -->
    <div
      class="border-b border-gray-100 px-6 py-5"
    >

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[36px] font-semibold text-[#505050]">
          Data Kurikulum
        </h2>
      </div>

      <!-- FILTER -->
      <div class="px-5 pt-5 flex items-center gap-4 flex-wrap">

        <!-- JURUSAN -->
        <select
          v-model="jurusan"
          class="w-65 h-13.5 border border-gray-300 rounded-xl px-4 text-[18px] text-gray-600 outline-none bg-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Pilih Jurusan</option>
          <option>Teknik Elektro</option>
          <option>Teknik Sipil</option>
        </select>

        <!-- PRODI -->
        <select
          v-model="prodi"
          class="w-65 h-13.5 border border-gray-300 rounded-xl px-4 text-[18px] text-gray-600 outline-none bg-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Pilih Prodi</option>
          <option>Teknik Informatika</option>
          <option>SIKC</option>
        </select>

        <!-- TAHUN AKADEMIK -->
        <select
          v-model="tahunAkademik"
          class="w-65 h-13.5 border border-gray-300 rounded-xl px-4 text-[18px] text-gray-600 outline-none bg-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Pilih Tahun Akademik</option>
          <option>2023-2024</option>
          <option>2024-2025</option>
        </select>

        <!-- BUTTON PILIH -->
        <button
          @click="pilihData"
          class="h-13.5 px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Pilih
        </button>

        <!-- BUTTON TAMBAH -->
        <button
          @click="tambahData"
          class="h-13.5 px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          + Tambah
        </button>
      </div>

      <!-- TABLE -->
      <div class="px-5 pt-8">

        <table class="w-full">

          <!-- HEADER -->
          <thead>
            <tr class="text-left text-gray-600">
              <th class="text-[18px] font-semibold">No</th>
              <th class="text-[18px] font-semibold">Nama Kurikulum</th>
              <th class="text-[18px] font-semibold">Prodi</th>
              <th class="text-[18px] font-semibold">Semester</th>
              <th class="text-[18px] font-semibold">Tahun Akademik</th>
              <th class="text-[18px] font-semibold text-center">Aksi</th>
            </tr>
          </thead>
<tbody>
  <tr
    v-for="(item, index) in kurikulumData"
    :key="item.id"
    class="hover:bg-gray-50"
  >
    <td class="py-4 text-[18px]">
      {{ (currentPage - 1) * perPage + index + 1 }}
    </td>

    <td class="py-4 text-[18px] font-medium">
      {{ item.nama }}
    </td>

    <td class="py-4 text-[18px]">
      {{ item.prodi }}
    </td>

    <td class="py-4 text-[18px]">
      {{ item.semester }}
    </td>

    <td class="py-4 text-[18px]">
      {{ item.tahun }}
    </td>

    <td class="py-4 text-center">
      <button
        type="button"
        @click="editData(item)"
        class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer"
      >
        ✎ Edit
      </button>
    </td>
  </tr>
</tbody>
        </table>
      </div>

<!-- FOOTER -->
<div
  class="flex items-center justify-between px-5 py-5 mt-8"
>

  <!-- BARIS -->
  <select
    v-model="perPage"
    class="w-22.5 h-10.5  border border-[#d8e1f0] bg-white p-5 shadow-sm"
  >
    <option :value="5">5 Baris</option>
    <option :value="10">10 Baris</option>
    <option :value="25">25 Baris</option>
  </select>

  <!-- PAGINATION -->
<div v-if="totalPages > 1" class="flex items-center gap-2">

  <button
    @click="prevPage"
    :disabled="currentPage === 1"
    class="px-3 py-1 border rounded-lg"
  >
    Previous
  </button>

  <template v-for="p in pages" :key="p">

    <span v-if="p === '...'">
      ...
    </span>

    <button
      v-else
      @click="goToPage(Number(p))"
      class="w-8 h-8 rounded-lg"
      :class="currentPage === p
        ? 'bg-blue-500 text-white'
        : 'bg-gray-100'"
    >
      {{ p }}
    </button>

  </template>

  <button
    @click="nextPage"
    :disabled="currentPage === totalPages"
    class="px-3 py-1 border rounded-lg"
  >
    Next
  </button>

</div>
</div>
    </div>
  </div>
</template>
