<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

// ========================
// TABLE DATA
// ========================
const jurusanData = ref<any[]>([])

const currentPage = ref(1)
const perPage = ref(10)
const totalPages = ref(1)

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ========================
// GET JURUSAN
// ========================
const getJurusan = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan?page=${currentPage.value}&per_page=${perPage.value}`,
      {
        method: "GET",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    const data =
      json?.data?.items ??
      json?.data ??
      []

    jurusanData.value = data.map((item: any) => ({
      id: item.id,
      name: item.name.replace(/-/g, " "),
    }))

    totalPages.value =
      json?.data?.pagination?.total_pages ?? 1
  } catch (err) {
    console.error(err)
  }
}

// ========================
// TAMBAH
// ========================
const tambahData = () => {
  router.push(
    "/dashboard-admin/Tambah_jurusan"
  )
}

// ========================
// EDIT
// ========================
const showEditModal = ref(false)

const editForm = ref({
  oldName: "",
  newName: "",
})

const editData = (item: any) => {
  editForm.value = {
    oldName: item.name,
    newName: item.name,
  }

  showEditModal.value = true
}

const submitEdit = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan/${editForm.value.oldName}`,
      {
        method: "PUT",
        headers: getHeaders(),
        body: JSON.stringify({
          new_name: editForm.value.newName,
        }),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json.message)
      return
    }

    showEditModal.value = false
    await getJurusan()
  } catch (err) {
    console.error(err)
  }
}
const searchJurusan = ref("")

const filteredJurusan = computed(() => {
  if (!searchJurusan.value.trim()) {
    return jurusanData.value
  }

  return jurusanData.value.filter((item) =>
    item.name
      .toLowerCase()
      .includes(searchJurusan.value.toLowerCase())
  )
})
// ========================
// HAPUS
// ========================
const hapusJurusan = async (nama: string) => {
  const konfirmasi = confirm(
    `Hapus jurusan ${nama}?`
  )

  if (!konfirmasi) return

  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan/${nama}`,
      {
        method: "DELETE",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json.message)
      return
    }

    await getJurusan()
  } catch (err) {
    console.error(err)
  }
}

// ========================
// PAGINATION
// ========================
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4)
    return Array.from(
      { length: total },
      (_, i) => i + 1
    )

  const pages: (number | string)[] = []

  pages.push(1)

  const start = Math.max(
    2,
    current - 1
  )

  const end = Math.min(
    total - 1,
    current + 1
  )

  if (start > 2) pages.push("...")

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  if (end < total - 1) {
    pages.push("...")
  }

  pages.push(total)

  return pages
})

const goToPage = async (page: number) => {
  currentPage.value = page
  await getJurusan()
}

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getJurusan()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getJurusan()
  }
}

onMounted(async () => {
  await getJurusan()
})
</script>
<template>
  <div
  class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden"
>

  <!-- HEADER BIRU -->
  <div class="bg-[#243e90] px-5 py-4">
    <h2 class="text-white text-2xl font-bold">
      Data Jurusan
    </h2>

    <p class="text-white text-sm mt-1">
      Data jurusan yang telah dibuat
    </p>
  </div>

  <!-- SEARCH + BUTTON -->
  <div class="p-5 flex flex-col md:flex-row gap-3 justify-between">

    <!-- SEARCH -->
    <div class="relative flex-1">
      <input
        v-model="searchJurusan"
        type="text"
        placeholder="Cari nama jurusan yang diinginkan..."
        class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-4 pr-10 text-sm outline-none"
      />

      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="absolute right-3 top-3.5 h-5 w-5 text-gray-400"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M21 21l-4.35-4.35M11 18a7 7 0 100-14 7 7 0 000 14z"
        />
      </svg>
    </div>

    <!-- TAMBAH -->
    <button
      @click="tambahData"
      class="bg-[#243e90] hover:bg-[#1d377f] text-white px-5 py-3 rounded-lg flex items-center gap-2 font-medium"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="w-4 h-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 4v16m8-8H4"
        />
      </svg>

      Tambah
    </button>

  </div>

  <!-- TABLE -->
  <div class="px-8 py-3 overflow-x-auto">

    <table class="w-full">

      <thead>
<tr class="text-gray-600 border-b border-gray-300">
  <th class="py-4 px-3 text-center w-16">
    No
  </th>

  <th class="px-3 text-center w-40">
    Kode Jurusan
  </th>

  <th class="px-3 text-center">
    Jurusan
  </th>

  <th class="px-3 text-center w-56">
    Aksi
  </th>
</tr>
      </thead>

<tbody>

  <tr v-if="filteredJurusan.length === 0">
    <td
      colspan="4"
      class="text-center py-10 text-gray-400"
    >
      Data jurusan tidak ditemukan
    </td>
  </tr>

  <tr
    v-for="(item,index) in filteredJurusan"
    :key="item.id"
    class="hover:bg-gray-50 border-b border-gray-100"
  >
    <!-- NO -->
    <td class="py-4 px-3 text-center w-16">
      {{
        (currentPage - 1) *
        perPage +
        index +
        1
      }}
    </td>

    <!-- KODE -->
    <td class="py-4 px-3 text-center w-40 font-medium text-gray-600">
      {{ item.kode }}
    </td>

    <!-- JURUSAN -->
    <td class="py-4 px-3 text-center font-medium text-gray-700">
      {{ item.name }}
    </td>

    <!-- AKSI -->
    <td class="py-4 px-3 w-56">
      <div class="flex justify-center items-center gap-2">

        <!-- EDIT -->
        <button
          @click="editData(item)"
          class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-3 py-2 rounded-lg flex items-center gap-1 text-sm font-medium"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L12 15l-4 1 1-4 8.586-8.586z"
            />
          </svg>

          Edit
        </button>

        <!-- HAPUS -->
        <button
          @click="hapusJurusan(item.name)"
          class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg flex items-center gap-1 text-sm font-medium"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5-3h4m-6 3h8"
            />
          </svg>

          Hapus
        </button>

      </div>
    </td>

  </tr>

</tbody>

    </table>

  </div>

  <!-- PAGINATION -->
  <div
    class="flex justify-end items-center gap-2 p-5 border-t border-gray-200"
  >
    <button
      @click="prevPage"
      :disabled="currentPage === 1"
      class="px-3 py-2 rounded-lg border disabled:opacity-50"
    >
      Previous
    </button>

    <template
      v-for="page in visiblePages"
      :key="page"
    >
      <span v-if="page === '...'">
        ...
      </span>

      <button
        v-else
        @click="goToPage(Number(page))"
        :class="
          currentPage === Number(page)
            ? 'bg-[#243e90] text-white'
            : 'bg-white'
        "
        class="w-9 h-9 rounded-lg border"
      >
        {{ page }}
      </button>
    </template>

    <button
      @click="nextPage"
      :disabled="currentPage === totalPages"
      class="px-3 py-2 rounded-lg border disabled:opacity-50"
    >
      Next
    </button>
  </div>

</div>
</template>