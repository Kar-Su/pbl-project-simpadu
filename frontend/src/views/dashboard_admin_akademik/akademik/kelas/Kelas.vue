<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"
import KonfirmasiHapus from "@/views/dashboard_super_admin/akademik/akun/konfirmasi_hapus.vue"

const router = useRouter()

// ================= FILTER =================
const jurusan = ref<string>("")
const prodi = ref<string>("")
const tahunAkademik = ref<string>("")

const prodiMap = ref<Record<number, any>>({})
const tahunMap = ref<Record<number, any>>({})

// ================= DATA =================
const allKelasData = ref<any[]>([])
const kelasData = ref<any[]>([])


const showDeleteModal = ref(false)
const deleteTarget = ref<any>(null)
// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(5)
const totalPages = ref(1)
const totalItems = ref(0)

// ================= FETCH =================
const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= GET KELAS (semua halaman) =================
const getKelas = async () => {
  try {
    let allItems: any[] = []
    let page = 1
    let lastPage = 1

    do {
      const res = await fetch(
        `${BASE_URL}/api/kelas?page=${page}&per_page=100`,
        { method: "GET", headers: getHeaders() }
      )
      const json = await res.json()
      const items = json?.data?.items ?? []
      allItems = [...allItems, ...items]
      lastPage = json?.data?.pagination?.total_pages ?? 1
      page++
    } while (page <= lastPage)

    allKelasData.value = allItems.map((item: any) => ({
      id: item.id,
      name: item.name ?? "-",

      jurusanId: String(item.prodi?.jurusan?.id ?? ""),
      jurusanName: (item.prodi?.jurusan?.name ?? "-").replace(/-/g, " "),

      prodiId: String(item.prodi?.id ?? ""),
      prodiName: (item.prodi?.name ?? "-").replace(/-/g, " "),

      tahunId: String(item.tahun_akademik?.id ?? ""),
      tahunLabel: item.tahun_akademik
        ? `${item.tahun_akademik.tahun_awal?.slice(0, 4) ?? "?"}/${item.tahun_akademik.tahun_akhir?.slice(0, 4) ?? "?"}`
        : "-",

      rawData: item,
    }))

    applyFilter()
  } catch (err) {
    console.error("GET KELAS ERROR:", err)
    allKelasData.value = []
    applyFilter()
  }
}

// ================= APPLY FILTER (lokal) =================
const applyFilter = () => {
  const filtered = allKelasData.value.filter((item) => {
    const matchJurusan = !jurusan.value || item.jurusanId === jurusan.value
    const matchProdi = !prodi.value || item.prodiId === prodi.value
    const matchTahun = !tahunAkademik.value || item.tahunId === tahunAkademik.value
    return matchJurusan && matchProdi && matchTahun
  })

  totalItems.value = filtered.length
  totalPages.value = Math.max(1, Math.ceil(filtered.length / perPage.value))

  if (currentPage.value > totalPages.value) currentPage.value = 1

  const start = (currentPage.value - 1) * perPage.value
  kelasData.value = filtered.slice(start, start + perPage.value)
}

// ================= FETCH SUPPORTING DATA =================
const getProdi = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, { headers: getHeaders() })
    const json = await res.json()
    const list = json?.data?.items ?? json?.data ?? []
    list.forEach((p: any) => {
      prodiMap.value[p.id] = p
    })
  } catch (err) {
    console.error(err)
  }
}

const getTahunAkademik = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik?per_page=100`, { headers: getHeaders() })
    const json = await res.json()
    const list = Array.isArray(json?.data) ? json.data : json?.data?.items ?? []
    list.forEach((t: any) => {
      tahunMap.value[t.id] = t
    })
  } catch (err) {
    console.error(err)
  }
}

// ================= DROPDOWN LIST =================
const jurusanList = computed(() => {
  const map = new Map()
  Object.values(prodiMap.value).forEach((p: any) => {
    const j = p?.jurusan
    if (j?.id && !map.has(j.id)) {
      map.set(j.id, { id: String(j.id), name: (j.name ?? "").replace(/-/g, " ") })
    }
  })
  return Array.from(map.values())
})

const prodiList = computed(() => {
  return Object.values(prodiMap.value).map((p: any) => ({
    id: String(p.id),
    name: (p.name ?? "").replace(/-/g, " "),
    jurusanId: String(p.jurusan?.id ?? ""),
  }))
})

const filteredProdiList = computed(() => {
  if (!jurusan.value) return prodiList.value
  return prodiList.value.filter((p) => p.jurusanId === jurusan.value)
})

const tahunAkademikList = computed(() => {
  return Object.values(tahunMap.value).map((t: any) => ({
    id: String(t.id),
    label: `${t.tahun_awal?.slice(0, 4) ?? "?"}/${t.tahun_akhir?.slice(0, 4) ?? "?"}`,
  }))
})

// ================= WATCH =================
watch(jurusan, () => {
  const stillValid = filteredProdiList.value.some((p) => p.id === prodi.value)
  if (!stillValid) prodi.value = ""
})

watch(perPage, () => {
  currentPage.value = 1
  applyFilter()
})

// ================= HELPER =================
const rowNumber = (index: number) =>
  (currentPage.value - 1) * perPage.value + index + 1

// ================= PAGINATION LOGIC =================
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4) return Array.from({ length: total }, (_, i) => i + 1)

  const pages: (number | string)[] = [1]
  const rangeStart = Math.max(2, current - 1)
  const rangeEnd = Math.min(total - 1, current + 1)

  if (rangeStart > 2) pages.push("...")
  for (let i = rangeStart; i <= rangeEnd; i++) pages.push(i)
  if (rangeEnd < total - 1) pages.push("...")
  pages.push(total)

  return pages
})

const goToPage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  applyFilter()
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    applyFilter()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    applyFilter()
  }
}

// ================= ACTION =================
const pilihData = () => {
  currentPage.value = 1
  applyFilter()
}

const resetFilter = () => {
  jurusan.value = ""
  prodi.value = ""
  tahunAkademik.value = ""
  currentPage.value = 1
  applyFilter()
}

const tambahData = () => {
  router.push("/dashboard-admin/tambah_kelas")
}

const editData = (item: any) => {
  sessionStorage.setItem(
    "selectedKelas",
    JSON.stringify(item.rawData)
  )

  router.push(`/dashboard-admin/edit_kelas/${item.id}`)
}

const hapusData = (item: any) => {
  deleteTarget.value = item
  showDeleteModal.value = true
}

const submitDelete = async () => {
  if (!deleteTarget.value) return

  try {
    const res = await fetch(
      `${BASE_URL}/api/kelas/${deleteTarget.value.id}`,
      {
        method: "DELETE",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json?.message || "Gagal menghapus kelas")
      return
    }

    showDeleteModal.value = false
    deleteTarget.value = null

    await getKelas()
  } catch (err) {
    console.error(err)
    alert("Terjadi kesalahan jaringan")
  }
}

// ================= INIT =================
onMounted(async () => {
  await Promise.all([getProdi(), getTahunAkademik()])
  await getKelas()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Kelas
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Kelas
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-6">
      Pengelolaan Data
    </p>

    <!-- CARD -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

      <!-- HEADER BIRU -->
      <div class="bg-[#243e90] px-5 py-4">
        <h2 class="text-white text-2xl font-bold">
          Data Kelas
        </h2>

        <p class="text-white text-sm mt-1">
          Kumpulan data kelas yang tersimpan
        </p>
      </div>

      <!-- FILTER -->
      <div class="px-5 pt-5 flex items-center gap-4 flex-wrap">

        <!-- JURUSAN -->
        <select
          v-model="jurusan"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="" disabled>Pilih Jurusan</option>
          <option v-for="j in jurusanList" :key="j.id" :value="j.id">
            {{ j.name }}
          </option>
        </select>

        <!-- PRODI -->
        <select
          v-model="prodi"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="" disabled>Pilih Prodi</option>
          <option v-for="p in filteredProdiList" :key="p.id" :value="p.id">
            {{ p.name }}
          </option>
        </select>

        <!-- TAHUN AKADEMIK -->
        <select
          v-model="tahunAkademik"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="" disabled>Pilih Tahun Akademik</option>
          <option v-for="t in tahunAkademikList" :key="t.id" :value="t.id">
            {{ t.label }}
          </option>
        </select>

        <!-- BUTTON PILIH -->
        <button
          @click="pilihData"
          class="h-[54px] px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Pilih
        </button>

        <!-- BUTTON RESET -->
        <button
          @click="resetFilter"
          class="h-[54px] px-6 bg-gray-400 hover:bg-gray-500 rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Reset
        </button>

        <!-- BUTTON TAMBAH -->
        <button
          @click="tambahData"
          class="h-[54px] px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          + Tambah
        </button>

      </div>

      <!-- TABLE -->
      <div class="px-5 pt-8">
        <table class="w-full border-separate border-spacing-y-5">

          <thead>
            <tr class="text-left text-gray-600 border-b border-gray-300">
              <th class="text-[18px] font-semibold">No</th>
              <th class="text-[18px] font-semibold">Nama Kelas</th>
              <th class="text-[18px] font-semibold">Jurusan</th>
              <th class="text-[18px] font-semibold">Prodi</th>
              <th class="text-[18px] font-semibold">Tahun Akademik</th>
              <th class="text-[18px] font-semibold text-center">Aksi</th>
            </tr>
          </thead>

          <tbody>

            <tr
              v-for="(item, index) in kelasData"
              :key="item.id"
              class="text-[#505050]"
            >

              <td class="text-[18px]">
                {{ rowNumber(index) }}
              </td>

              <td class="text-[18px] font-medium capitalize">
                {{ item.name }}
              </td>

              <td class="text-[18px] font-medium capitalize">
                {{ item.jurusanName }}
              </td>

              <td class="text-[18px] font-medium capitalize">
                {{ item.prodiName }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.tahunLabel }}
              </td>

              <td class="flex items-center justify-center gap-3">

                <button
                  @click="editData(item)"
                  class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-5 py-2 rounded-xl text-[16px] font-semibold shadow-md transition"
                >
                  ✎ Edit
                </button>

                <button
                  @click="hapusData(item)"
                  class="bg-[#ef4d43] hover:bg-[#d93d34] text-white px-5 py-2 rounded-xl text-[16px] font-semibold shadow-md transition"
                >
                  🗑 Hapus
                </button>

              </td>

            </tr>

            <tr v-if="kelasData.length === 0">
              <td
                colspan="6"
                class="text-center text-gray-400 py-10 text-[16px]"
              >
                Tidak ada data kelas
              </td>
            </tr>

          </tbody>

        </table>
      </div>

      <!-- FOOTER / PAGINATION -->
      <div class="flex items-center justify-end px-5 pt-10 pb-5">

        <div class="flex items-center gap-2 text-gray-500 text-sm">

          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="px-3 py-1 border rounded-lg disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Previous
          </button>

          <template v-for="item in visiblePages" :key="item">
            <span v-if="item === '...'" class="px-1 text-gray-400">...</span>
            <button
              v-else
              @click="goToPage(Number(item))"
              class="w-8 h-8 rounded-lg"
              :class="currentPage === Number(item)
                ? 'bg-blue-500 text-white'
                : 'bg-gray-100'"
            >
              {{ item }}
            </button>
          </template>

          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 border rounded-lg disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Next
          </button>

        </div>

      </div>

    </div>

  </div>
  <KonfirmasiHapus
  v-if="showDeleteModal"
  :message="`Apakah anda yakin ingin menghapus kelas '${deleteTarget?.name}'?`"
  @close="showDeleteModal = false"
  @confirm="submitDelete"
/>
</template>