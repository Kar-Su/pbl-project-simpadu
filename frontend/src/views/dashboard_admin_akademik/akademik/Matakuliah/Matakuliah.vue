<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import KonfirmasiHapus from "@/views/dashboard_super_admin/akademik/akun/konfirmasi_hapus.vue"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

// ================= DATA =================
const mataKuliahData = ref<any[]>([])
const search = ref("")
const showEditModal = ref(false)
const editLoading = ref(false)
const editError = ref("")
const showDeleteModal = ref(false)
const deleteTarget = ref<any>(null)

const editForm = ref({
  kodeLama: "",
  kode: "",
  name: "",
  sks: 0,
})

// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(10)
const totalPages = ref(1)
const totalItems = ref(0)

// ================= HEADER =================
const getHeaders = () => ({
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= GET MATA KULIAH =================
const getMataKuliah = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/mata-kuliah?page=${currentPage.value}`,
      {
        method: "GET",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

mataKuliahData.value = (json.data?.items ?? []).map((item: any) => ({
  id: item.id,
  kode: item.kode.replace(/-/g, " "),
  nama: item.name.replace(/-/g, " "),
  sks: item.sks,
  rawData: item,
}))

    totalPages.value = json.data?.pagination?.total_pages ?? 1
    totalItems.value = json.data?.pagination?.total_items ?? 0
    perPage.value = json.data?.pagination?.per_page ?? 10
  } catch (err) {
    console.error("GET MK ERROR:", err)
  }
}

// ================= SEARCH =================
const filteredData = computed(() => {
  if (!search.value) return mataKuliahData.value

  return mataKuliahData.value.filter((item) =>
    item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
    item.kode.toLowerCase().includes(search.value.toLowerCase())
  )
})

// ================= PAGINATION =================
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const pages: (number | string)[] = [1]

  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)

  if (start > 2) pages.push("...")

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  if (end < total - 1) pages.push("...")

  pages.push(total)

  return pages
})

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getMataKuliah()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getMataKuliah()
  }
}

const goToPage = async (page: number) => {
  currentPage.value = page
  await getMataKuliah()
}

// ================= ACTION =================
const tambahData = () => {
  router.push("/dashboard-admin/tambah_matakuliah")
}

const submitEdit = async () => {
  editLoading.value = true
  editError.value = ""

  try {
    const res = await fetch(
      `${BASE_URL}/api/mata-kuliah/${editForm.value.kodeLama}`,
      {
        method: "PUT",
        headers: {
          ...getHeaders(),
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          kode: editForm.value.kode,
          name: editForm.value.name,
          sks: Number(editForm.value.sks),
        }),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      editError.value = json.message || "Gagal mengupdate data"
      return
    }

    showEditModal.value = false
    await getMataKuliah()
  } catch (err) {
    editError.value = "Terjadi kesalahan jaringan"
  } finally {
    editLoading.value = false
  }
}
const hapusData = (item: any) => {
  deleteTarget.value = item
  showDeleteModal.value = true
}

const submitDelete = async () => {
  if (!deleteTarget.value) return

  try {
    const res = await fetch(
      `${BASE_URL}/api/mata-kuliah/${deleteTarget.value.kode}`,
      {
        method: "DELETE",
        headers: getHeaders(),
      }
    )

    let result: any = {}

    try {
      result = await res.json()
    } catch {
      result = {}
    }

    if (!res.ok) {
      alert(
        result?.message ||
        result?.error ||
        "Gagal menghapus mata kuliah"
      )
      return
    }

    showDeleteModal.value = false
    deleteTarget.value = null

    await getMataKuliah()
  } catch (err) {
    console.error("DELETE ERROR:", err)
    alert("Terjadi kesalahan saat menghapus data")
  }
}
const editData = (item: any) => {
  editForm.value = {
    kodeLama: item.kode,
    kode: item.kode,
    name: item.nama,
    sks: item.sks,
  }

  editError.value = ""
  showEditModal.value = true
}

onMounted(() => {
  getMataKuliah()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- Breadcrumb -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Mata Kuliah
    </div>

    <!-- Title -->
    <h1 class="text-[42px] font-bold text-[#404040]">
      Mata Kuliah
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-6">
      Pengelolaan Data Mata Kuliah
    </p>

<div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

  <!-- HEADER BIRU -->
  <div class="bg-[#243e90] px-5 py-4">
    <h2 class="text-white text-2xl font-bold">
      Data Mata Kuliah
    </h2>

    <p class="text-white text-sm mt-1">
      Kumpulan data matakuliah beserta SKS
    </p>
  </div>

      <!-- Search -->
      <div class="px-6 py-5 flex justify-between items-center">

        <input
          v-model="search"
          type="text"
          placeholder="Cari kode atau nama mata kuliah..."
          class="w-80 border border-gray-300 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-[#29479d]"
        />

        <button
          @click="tambahData"
          class="px-6 py-3 bg-[#29479d] hover:bg-[#1d377f] text-white rounded-xl font-semibold transition"
        >
          + Tambah
        </button>

      </div>

      <!-- Table -->
      <div class="px-6 pb-6 overflow-x-auto">

        <table class="w-full">

          <thead class="">

            <tr class="text-left text-black-600 border-b border-gray-300">
              <th class="text-left py-4 px-4">No</th>
              <th class="text-left py-4 px-4">Kode MK</th>
              <th class="text-left py-4 px-4">Nama Mata Kuliah</th>
              <th class="text-left py-4 px-4">SKS</th>
              <th class="text-center py-4 px-4">Aksi</th>
            </tr>

          </thead>

          <tbody>

            <tr
              v-for="(item, index) in filteredData"
              :key="item.id"
              class=" hover:bg-blue-50 transition"
            >
              <td class="px-4 py-4">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <td class="px-4 py-4 font-medium">
                {{ item.kode }}
              </td>

              <td class="px-4 py-4">
                {{ item.nama }}
              </td>

              <td class="px-4 py-4">
                {{ item.sks }}
              </td>

<td class="px-4 py-4 text-center">
  <div class="flex justify-center gap-2">

    <button
      @click="editData(item)"
      class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg text-sm"
    >
      Edit
    </button>

    <button
      @click="hapusData(item)"
      class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg text-sm"
    >
      Hapus
    </button>

  </div>
</td>
            </tr>

            <tr v-if="filteredData.length === 0">
              <td colspan="5" class="text-center py-10 text-gray-500">
                Data mata kuliah tidak ditemukan
              </td>
            </tr>

          </tbody>

        </table>

      </div>

      <!-- Pagination -->
      <div class="flex justify-end items-center px-6 pb-6 gap-2">

        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-3 py-2 border rounded-lg disabled:opacity-40"
        >
          Previous
        </button>

        <template
          v-for="page in visiblePages"
          :key="page"
        >
          <span
            v-if="page === '...'"
            class="px-2"
          >
            ...
          </span>

          <button
            v-else
            @click="goToPage(Number(page))"
            class="w-9 h-9 rounded-lg"
            :class="
              currentPage === Number(page)
                ? 'bg-[#29479d] text-white'
                : 'border'
            "
          >
            {{ page }}
          </button>
        </template>

        <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="px-3 py-2 border rounded-lg disabled:opacity-40"
        >
          Next
        </button>

<div
  v-if="showEditModal"
  class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
>
  <div class="bg-white rounded-xl w-[500px] p-6">

    <h2 class="text-xl font-bold mb-4">
      Edit Mata Kuliah
    </h2>

    <p
      v-if="editError"
      class="text-red-500 mb-3"
    >
      {{ editError }}
    </p>

    <div class="space-y-4">

      <div>
        <label class="block mb-1 font-medium">
          Kode MK
        </label>

        <input
          v-model="editForm.kode"
          class="w-full border rounded-lg px-3 py-2"
        />
      </div>

      <div>
        <label class="block mb-1 font-medium">
          Nama Mata Kuliah
        </label>

        <input
          v-model="editForm.name"
          class="w-full border rounded-lg px-3 py-2"
        />
      </div>

      <div>
        <label class="block mb-1 font-medium">
          SKS
        </label>

        <input
          v-model="editForm.sks"
          type="number"
          class="w-full border rounded-lg px-3 py-2"
        />
      </div>

    </div>

    <div class="flex justify-end gap-2 mt-6">

      <button
        @click="showEditModal = false"
        class="px-4 py-2 border rounded-lg"
      >
        Batal
      </button>

      <button
        @click="submitEdit"
        class="px-4 py-2 bg-[#29479d] text-white rounded-lg"
      >
        {{ editLoading ? "Menyimpan..." : "Simpan" }}
      </button>

    </div>

  </div>
</div>

      </div>

    </div>

  </div>
  <KonfirmasiHapus
  v-if="showDeleteModal"
  :message="`Apakah anda yakin ingin menghapus mata kuliah '${deleteTarget?.nama}'?`"
  @close="showDeleteModal = false"
  @confirm="submitDelete"
/>
</template>