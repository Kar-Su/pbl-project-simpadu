<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= FILTER =================
const filterJurusan = ref("")
const filterProdi = ref("")

// ================= DATA =================
const prodiData = ref<any[]>([])
const jurusanMap = ref<Record<number, any>>({})

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

// ================= GET ALL JURUSAN (dari prodi list) =================
const getAllProdi = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, { headers: getHeaders() })
    const json = await res.json()
    console.log("ALL PRODI RAW:", JSON.stringify(json, null, 2))

    // Response bisa array langsung atau object dengan items
    const list = Array.isArray(json?.data)
      ? json.data
      : json?.data?.items ?? []

    // Build jurusanMap dari data prodi
    list.forEach((p: any) => {
      const j = p.jurusan
      if (j?.id && !jurusanMap.value[j.id]) {
        jurusanMap.value[j.id] = { id: j.id, name: j.name }
      }
    })

    return list
  } catch (err) {
    console.error("GET ALL PRODI ERROR:", err)
    return []
  }
}

// ================= GET PRODI (with filter) =================
const getProdi = async () => {
  try {
    let url = `${BASE_URL}/api/prodi`

    // Filter by jurusan name jika dipilih
    if (filterJurusan.value) {
      const jurusan = jurusanMap.value[Number(filterJurusan.value)]
      if (jurusan) {
        url = `${BASE_URL}/api/prodi/jurusan/${encodeURIComponent(jurusan.name)}`
      }
    }

    const res = await fetch(url, { headers: getHeaders() })
    const json = await res.json()
    console.log("PRODI DATA:", JSON.stringify(json, null, 2))

    const raw = Array.isArray(json?.data)
      ? json.data
      : json?.data?.items ?? (json?.data ? [json.data] : [])

    // Filter by prodi name jika ada
    let filtered = raw
    if (filterProdi.value) {
      filtered = raw.filter((p: any) =>
        p.name?.toLowerCase().includes(filterProdi.value.toLowerCase())
      )
    }

    // Pagination manual (karena API mungkin return semua)
    totalItems.value = filtered.length
    totalPages.value = Math.max(1, Math.ceil(filtered.length / perPage.value))

    const start = (currentPage.value - 1) * perPage.value
    const end = start + perPage.value
    const paginated = filtered.slice(start, end)

    prodiData.value = paginated.map((item: any) => ({
      id: item.id,
      nama: (item.name ?? "-").replace(/-/g, " "),
      namaAsli: item.name ?? "",
      jenjang: item.jenjang ?? "-",
      jurusan: (item.jurusan?.name ?? "-").replace(/-/g, " "),
      jurusanId: item.jurusan?.id ?? null,
      rawData: item,
    }))
  } catch (err) {
    console.error("GET PRODI ERROR:", err)
  }
}

// ================= EDIT MODAL =================
const showEditModal = ref(false)
const editForm = ref({
  namaAsli: "",
  nama: "",
  jenjang: "",
  jurusanId: "" as string | number,
})
const editError = ref("")
const editLoading = ref(false)

const jenjangOptions = ["D3", "D4", "S1", "S2", "S3"]

const jurusanList = computed(() =>
  Object.values(jurusanMap.value).map((j: any) => ({
    id: String(j.id),
    name: j.name.replace(/-/g, " "),
    nameAsli: j.name,
  }))
)

const editData = (item: any) => {
  editForm.value = {
    namaAsli: item.namaAsli,
    nama: item.nama,
    jenjang: item.jenjang,
    jurusanId: String(item.jurusanId ?? ""),
  }
  editError.value = ""
  showEditModal.value = true
}

const submitEdit = async () => {
  if (!editForm.value.nama.trim()) {
    editError.value = "Nama prodi tidak boleh kosong."
    return
  }
  if (!editForm.value.jenjang) {
    editError.value = "Jenjang harus dipilih."
    return
  }
  if (!editForm.value.jurusanId) {
    editError.value = "Jurusan harus dipilih."
    return
  }

  editLoading.value = true
  editError.value = ""

  const payload = {
    name: editForm.value.nama.replace(/ /g, "-").toLowerCase(),
    jenjang: editForm.value.jenjang,
    jurusan_id: Number(editForm.value.jurusanId),
  }

  console.log("PAYLOAD EDIT:", payload)
  console.log("URL:", `${BASE_URL}/api/prodi/${editForm.value.namaAsli}`)

  try {
    const res = await fetch(`${BASE_URL}/api/prodi/${editForm.value.namaAsli}`, {
      method: "PUT",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })

    if (!res.ok) {
      const json = await res.json()
      console.error("ERROR UPDATE:", json)
      editError.value = json?.message || JSON.stringify(json) || "Gagal menyimpan perubahan."
      return
    }

    showEditModal.value = false
    await getProdi()
  } catch (err) {
    editError.value = "Terjadi kesalahan jaringan."
  } finally {
    editLoading.value = false
  }
}

// ================= DELETE =================
const showDeleteModal = ref(false)
const deleteTarget = ref<any>(null)
const deleteLoading = ref(false)
const deleteError = ref("")

const confirmDelete = (item: any) => {
  deleteTarget.value = item
  deleteError.value = ""
  showDeleteModal.value = true
}

const submitDelete = async () => {
  if (!deleteTarget.value) return

  deleteLoading.value = true
  deleteError.value = ""

  try {
    const res = await fetch(`${BASE_URL}/api/prodi/${deleteTarget.value.namaAsli}`, {
      method: "DELETE",
      headers: getHeaders(),
    })

    if (!res.ok) {
      const json = await res.json()
      deleteError.value = json?.message || "Gagal menghapus prodi."
      return
    }

    showDeleteModal.value = false
    deleteTarget.value = null

    // Kembali ke halaman 1 jika data di halaman ini habis
    if (prodiData.value.length === 1 && currentPage.value > 1) {
      currentPage.value--
    }
    await getProdi()
  } catch (err) {
    deleteError.value = "Terjadi kesalahan jaringan."
  } finally {
    deleteLoading.value = false
  }
}

// ================= PAGINATION LOGIC =================
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4) return Array.from({ length: total }, (_, i) => i + 1)

  const pages: (number | string)[] = []
  pages.push(1)

  const rangeStart = Math.max(2, current - 1)
  const rangeEnd = Math.min(total - 1, current + 1)

  if (rangeStart > 2) pages.push("...")
  for (let i = rangeStart; i <= rangeEnd; i++) pages.push(i)
  if (rangeEnd < total - 1) pages.push("...")

  pages.push(total)
  return pages
})

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getProdi()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getProdi()
  }
}

const goToPage = async (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  await getProdi()
}

watch(perPage, async () => {
  currentPage.value = 1
  await getProdi()
})

// ================= ACTION =================
const pilihData = async () => {
  currentPage.value = 1
  await getProdi()
}

const tambahData = () => {
  router.push("/dashboard-admin/tambah_prodi")
}

onMounted(async () => {
  await getAllProdi()
  await getProdi()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Prodi
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">Prodi</h1>
    <p class="text-gray-500 text-sm mt-2 mb-6">Pengelolaan Data</p>

    <!-- CARD -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[36px] font-semibold text-[#505050]">Prodi</h2>
      </div>

      <!-- FILTER -->
      <div class="px-5 pt-5 flex items-center gap-4 flex-wrap">

        <select v-model="filterJurusan" class="w-65 h-13.5 border border-gray-300 rounded-xl px-4 bg-white">
          <option value="">Pilih Jurusan</option>
          <option v-for="j in jurusanList" :key="j.id" :value="j.id">{{ j.name }}</option>
        </select>

        <button @click="pilihData"
          class="h-13.5 px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition">
          Pilih
        </button>

        <button @click="tambahData"
          class="h-13.5 px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition">
          + Tambah
        </button>
      </div>

      <!-- TABLE -->
      <div class="px-5 pt-8 overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-600 border-b border-gray-300">
              <th class="text-[18px] font-semibold pb-3">No</th>
              <th class="text-[18px] font-semibold pb-3">Nama Prodi</th>
              <th class="text-[18px] font-semibold pb-3">Jenjang</th>
              <th class="text-[18px] font-semibold pb-3">Jurusan</th>
              <th class="text-[18px] font-semibold pb-3 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="prodiData.length === 0">
              <td colspan="5" class="py-8 text-center text-gray-400 text-sm">Tidak ada data prodi.</td>
            </tr>
            <tr v-for="(item, index) in prodiData" :key="item.id"
              class="border-b border-gray-200 hover:bg-gray-50 transition">
              <td class="py-4 text-[18px]">{{ (currentPage - 1) * perPage + index + 1 }}</td>
              <td class="py-4 text-[18px]">{{ item.nama }}</td>
              <td class="py-4 text-[18px]">{{ item.jenjang }}</td>
              <td class="py-4 text-[18px]">{{ item.jurusan }}</td>
              <td class="py-4 text-center">
                <div class="flex justify-center gap-2">
                  <button type="button" @click="editData(item)"
                    class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer transition">
                    ✎ Edit
                  </button>
                  <button type="button" @click="confirmDelete(item)"
                    class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer transition">
                    🗑 Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- FOOTER / PAGINATION -->
      <div class="flex items-center justify-between w-full px-5 py-5">
        <p class="text-sm text-gray-500">Total: {{ totalItems }} data</p>

        <div class="flex items-center gap-2 text-gray-500 text-sm">

          <button @click="prevPage" :disabled="currentPage === 1"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Previous
          </button>

          <template v-for="item in visiblePages" :key="item">
            <span v-if="item === '...'" class="w-8 h-8 flex items-center justify-center text-gray-400">...</span>
            <button v-else @click="goToPage(Number(item))" :class="currentPage === Number(item)
              ? 'bg-[#1c3277] text-white shadow-md scale-105'
              : 'bg-white text-[#4b4b4b] hover:bg-[#d6ddee]'"
              class="w-8 h-8 rounded-md text-sm font-medium transition-all duration-200">
              {{ item }}
            </button>
          </template>

          <button @click="nextPage" :disabled="currentPage === totalPages"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            Next
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

        </div>
      </div>

    </div><!-- end CARD -->

    <!-- ================================ -->
    <!-- EDIT MODAL                        -->
    <!-- ================================ -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
      @click.self="showEditModal = false">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
        <h2 class="text-xl font-bold text-[#404040] mb-4">Edit Prodi</h2>

        <!-- Error -->
        <p v-if="editError" class="text-red-500 text-sm mb-3 bg-red-50 px-3 py-2 rounded-lg">{{ editError }}</p>

        <!-- Nama Prodi -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Nama Prodi</label>
          <input v-model="editForm.nama" type="text"
            placeholder="contoh: teknik informatika"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]" />
          <p class="text-xs text-gray-400 mt-1">Otomatis dikonversi ke format: teknik-informatika</p>
        </div>

        <!-- Jenjang -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Jenjang</label>
          <select v-model="editForm.jenjang"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]">
            <option value="">Pilih Jenjang</option>
            <option v-for="j in jenjangOptions" :key="j" :value="j">{{ j }}</option>
          </select>
        </div>

        <!-- Jurusan -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-1">Jurusan</label>
          <select v-model="editForm.jurusanId"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]">
            <option value="">Pilih Jurusan</option>
            <option v-for="j in jurusanList" :key="j.id" :value="j.id">{{ j.name }}</option>
          </select>
        </div>

        <!-- Tombol -->
        <div class="flex justify-end gap-3">
          <button @click="showEditModal = false"
            class="px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-600 hover:bg-gray-100 transition">
            Batal
          </button>
          <button @click="submitEdit" :disabled="editLoading"
            class="px-4 py-2 rounded-lg bg-[#29479d] hover:bg-[#1d377f] text-white text-sm font-semibold disabled:opacity-50 transition">
            {{ editLoading ? "Menyimpan..." : "Simpan" }}
          </button>
        </div>

      </div>
    </div>
    <!-- END EDIT MODAL -->

    <!-- ================================ -->
    <!-- DELETE CONFIRM MODAL              -->
    <!-- ================================ -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
      @click.self="showDeleteModal = false">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-sm p-6">
        <h2 class="text-xl font-bold text-[#404040] mb-2">Hapus Prodi</h2>
        <p class="text-gray-600 text-sm mb-4">
          Apakah Anda yakin ingin menghapus prodi
          <span class="font-semibold text-red-600">{{ deleteTarget?.nama }}</span>?
          Tindakan ini tidak dapat dibatalkan.
        </p>

        <p v-if="deleteError" class="text-red-500 text-sm mb-3 bg-red-50 px-3 py-2 rounded-lg">{{ deleteError }}</p>

        <div class="flex justify-end gap-3">
          <button @click="showDeleteModal = false"
            class="px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-600 hover:bg-gray-100 transition">
            Batal
          </button>
          <button @click="submitDelete" :disabled="deleteLoading"
            class="px-4 py-2 rounded-lg bg-red-500 hover:bg-red-600 text-white text-sm font-semibold disabled:opacity-50 transition">
            {{ deleteLoading ? "Menghapus..." : "Hapus" }}
          </button>
        </div>
      </div>
    </div>
    <!-- END DELETE MODAL -->

  </div>
</template>