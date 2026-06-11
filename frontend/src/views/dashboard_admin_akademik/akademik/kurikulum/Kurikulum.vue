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

const prodiMap = ref<Record<number, any>>({})
const tahunMap = ref<Record<number, any>>({})

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
      { method: "GET", headers: getHeaders() }
    )
    const json = await res.json()
    console.log("KURIKULUM:", json)

    const raw = json.data?.items ?? []
    console.log("RAW ITEM:", JSON.stringify(raw[0], null, 2))
    kurikulumData.value = raw.map((item: any) => ({
      id: item.id,
      nama: item.name ?? "-",
      prodi: (item.prodi?.name ?? "-").replace(/-/g, " "),
      semester: item.kurikulum_mk?.[0]?.semester
        ? `Semester ${item.kurikulum_mk[0].semester}`
        : "-",
      tahun: tahunAktif.value,
      rawData: item,
    }))

    totalPages.value = json.data?.pagination?.total_pages ?? 1
    totalItems.value = json.data?.pagination?.total_items ?? 0
    perPage.value = json.data?.pagination?.per_page ?? perPage.value
  } catch (err) {
    console.error("GET KURIKULUM ERROR:", err)
  }
}

const tahunAktif = ref("-")

const getTahunAktif = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/tahun-akademik/status/aktif`,
      { headers: getHeaders() }
    )
    const json = await res.json()

    // data adalah array, ambil index 0
    const data = Array.isArray(json?.data) ? json.data[0] : json?.data
    if (data) {
      const awal = data.tahun_awal?.slice(0, 4) ?? "?"
      const akhir = data.tahun_akhir?.slice(0, 4) ?? "?"
      tahunAktif.value = `${awal}/${akhir}`
    }
  } catch (err) {
    console.error("GET TAHUN AKTIF ERROR:", err)
  }
}

// ================= EDIT MODAL =================
const showEditModal = ref(false)
const editForm = ref({
  id: null as string | null,
  nama: "",
  kode: "",
  kodeAsli: "", 
  prodiId: "",
  semester: "",
})
const editError = ref("")
const editLoading = ref(false)

const prodiListEdit = computed(() =>
  Object.values(prodiMap.value).map((p: any) => ({
    id: String(p.id),
    name: p.name.replace(/-/g, " "),
  }))
)


const editData = (item: any) => {
  editForm.value = {
    id: item.id,
    nama: item.nama,
    kode: (item.rawData?.kode ?? "").replace(/-/g, " "),
    kodeAsli: item.rawData?.kode ?? "",   // ← simpan kode asli
    prodiId: String(item.rawData?.prodi?.id ?? ""),
    semester: String(item.rawData?.kurikulum_mk?.[0]?.semester ?? ""),
  }
  editError.value = ""
  showEditModal.value = true
}

const submitEdit = async () => {
  if (!editForm.value.nama.trim()) {
    editError.value = "Nama kurikulum tidak boleh kosong."
    return
  }
  if (!editForm.value.kode.trim()) {
    editError.value = "Kode kurikulum tidak boleh kosong."
    return
  }
  if (!editForm.value.prodiId) {
    editError.value = "Prodi harus dipilih."
    return
  }

  editLoading.value = true
  editError.value = ""

  const payload = {
    name: editForm.value.nama,
    kode: editForm.value.kode ?? "",
    prodi_id: editForm.value.prodiId,
    semester: Number(editForm.value.semester),
  }

  console.log("PAYLOAD EDIT:", payload)

  try {
const res = await fetch(
  `${BASE_URL}/api/kurikulum/${editForm.value.id}/${editForm.value.kodeAsli}`,  // ← kodeAsli
  {
    method: "PUT",
    headers: getHeaders(),
    body: JSON.stringify({
      name: editForm.value.nama,
      prodi_id: editForm.value.prodiId,
      semester: Number(editForm.value.semester),
    }),
  }
)

    if (!res.ok) {
      const json = await res.json()
      console.error("ERROR UPDATE:", json)
      editError.value = json?.message || JSON.stringify(json) || "Gagal menyimpan perubahan."
      return
    }

    showEditModal.value = false
    await getKurikulum()
  } catch (err) {
    editError.value = "Terjadi kesalahan jaringan."
  } finally {
    editLoading.value = false
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
    await getKurikulum()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getKurikulum()
  }
}

const goToPage = async (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  await getKurikulum()
}

watch(perPage, async () => {
  currentPage.value = 1
  await getKurikulum()
})

// ================= FETCH SUPPORTING DATA =================
const getProdi = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, { headers: getHeaders() })
    const json = await res.json()
    const list = json?.data?.items ?? json?.data ?? []
    list.forEach((p: any) => { prodiMap.value[p.id] = p })
  } catch (err) {
    console.error(err)
  }
}

const getTahunAkademik = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik?per_page=10`, { headers: getHeaders() })
    const json = await res.json()
    console.log("TAHUN AKTIF RAW:", JSON.stringify(json, null, 2))
    const list = Array.isArray(json?.data) ? json.data : json?.data?.items ?? []
    list.forEach((t: any) => { tahunMap.value[t.id] = t })
  } catch (err) {
    console.error(err)
  }
}

// ================= COMPUTED LISTS =================
const jurusanList = computed(() => {
  const map = new Map()
  Object.values(prodiMap.value).forEach((p: any) => {
    const j = p.jurusan
    if (j?.id && !map.has(j.id)) map.set(j.id, { id: String(j.id), name: j.name })
  })
  return Array.from(map.values())
})

const prodiList = computed(() =>
  Object.values(prodiMap.value).map((p: any) => ({ id: String(p.id), name: p.name }))
)

const tahunAkademikList = computed(() =>
  Object.values(tahunMap.value).map((t: any) => ({
    id: String(t.id),
    label: `${t.tahun_awal.slice(0, 4)}/${t.tahun_akhir.slice(0, 4)}`,
  }))
)

// ================= ACTION =================
const pilihData = () => {
  console.log({ jurusan: jurusan.value, prodi: prodi.value, tahun: tahunAkademik.value })
}

const tambahData = () => {
  router.push("/dashboard-admin/tambah_kurikulum")
}

onMounted(async () => {
  await Promise.all([getProdi(), getTahunAkademik(), getTahunAktif()])
  await getKurikulum()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Mahasiswa > Kurikulum
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">Kurikulum</h1>
    <p class="text-gray-500 text-sm mt-2 mb-6">Pengelolaan Data</p>

    <!-- CARD -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[36px] font-semibold text-[#505050]">Data Kurikulum</h2>
      </div>

      <!-- FILTER -->
      <div class="px-5 pt-5 flex items-center gap-4 flex-wrap">

        <select v-model="jurusan" class="w-65 h-13.5 border border-gray-300 rounded-xl px-4">
          <option value="">Pilih Jurusan</option>
          <option v-for="j in jurusanList" :key="j.id" :value="j.id">{{ j.name }}</option>
        </select>

        <select v-model="prodi" class="w-65 h-13.5 border border-gray-300 rounded-xl px-4">
          <option value="">Pilih Prodi</option>
          <option v-for="p in prodiList" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>

        <select v-model="tahunAkademik" class="w-65 h-13.5 border border-gray-300 rounded-xl px-4">
          <option value="">Pilih Tahun Akademik</option>
          <option v-for="t in tahunAkademikList" :key="t.id" :value="t.id">{{ t.label }}</option>
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
      <div class="px-5 pt-8">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-600 border-b border-gray-300">
              <th class="text-[18px] font-semibold">No</th>
              <th class="text-[18px] font-semibold">Nama Kurikulum</th>
              <th class="text-[18px] font-semibold">Prodi</th>
              <th class="text-[18px] font-semibold">Semester</th>
              <th class="text-[18px] font-semibold">Tahun Akademik</th>
              <th class="text-[18px] font-semibold text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in kurikulumData" :key="item.id" class="hover:bg-gray-50">
              <td class="py-4 text-[18px]">{{ (currentPage - 1) * perPage + index + 1 }}</td>
              <td class="py-4 text-[18px]">{{ item.nama }}</td>
              <td class="py-4 text-[18px]">{{ item.prodi }}</td>
              <td class="py-4 text-[18px]">{{ item.semester }}</td>
              <td class="py-4 text-[18px]">{{ item.tahun }}</td>
              <td class="py-4 text-center">
                <button type="button" @click="editData(item)"
                  class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer">
                  ✎ Edit
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- FOOTER / PAGINATION -->
      <div class="flex items-center justify-end w-full px-5 py-5">
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
    <!-- EDIT MODAL — di luar card/table  -->
    <!-- ================================ -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
      @click.self="showEditModal = false">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
        <h2 class="text-xl font-bold text-[#404040] mb-4">Edit Kurikulum</h2>

        <!-- Error -->
        <p v-if="editError" class="text-red-500 text-sm mb-3">{{ editError }}</p>

        <!-- Nama -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Nama Kurikulum</label>
          <input v-model="editForm.nama" type="text"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]" />
        </div>

        <!-- Kode -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Kode Kurikulum</label>
          <input v-model="editForm.kode" type="text" maxlength="3"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]" />
        </div>

        <!-- Prodi -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">Prodi</label>
          <select v-model="editForm.prodiId" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm">
            <option value="">Pilih Prodi</option>
            <option v-for="p in prodiListEdit" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>

        <!-- Semester -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-1">Semester</label>
          <input v-model="editForm.semester" type="number" min="1"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#29479d]" />
        </div>

        <!-- Tombol -->
        <div class="flex justify-end gap-3">
          <button @click="showEditModal = false"
            class="px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-600 hover:bg-gray-100">
            Batal
          </button>
          <button @click="submitEdit" :disabled="editLoading"
            class="px-4 py-2 rounded-lg bg-[#29479d] hover:bg-[#1d377f] text-white text-sm font-semibold disabled:opacity-50">
            {{ editLoading ? "Menyimpan..." : "Simpan" }}
          </button>
        </div>

      </div>
    </div>
    <!-- END EDIT MODAL -->

  </div>
</template>