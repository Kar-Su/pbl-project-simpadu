<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface AkademikItem {
  id: number
  tipeSemester: string
  tahunAwal: string
  tahunAkhir: string
  status?: string

  // tambahan
  rawTipeSemester?: string
}

// ================= STATE =================
const filterSemester = ref('')
const filterTahun = ref('')

const currentPage = ref(1)
const perPage = ref(10)

const showModal = ref(false)

const allData = ref<AkademikItem[]>([])
const filteredData = ref<AkademikItem[]>([])

const editingItem = ref<AkademikItem | null>(null)

const editForm = ref({
  tipeSemester: '',
  tahunAwal: '',
  tahunAkhir: '',
  status: '',
})

// ================= PAGINATION NUMBER =================
const displayedPages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  // jika total <= 3 tampil semua
  if (total <= 3) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  // halaman awal
  if (current <= 2) {
    return [1, 2, '...', total]
  }

  // halaman akhir
  if (current >= total - 1) {
    return [1, '...', total - 1, total]
  }

  // halaman tengah
  return [1, '...', current, '...', total]
})

// ================= HEADER =================
const getHeaders = () => ({
  'Content-Type': 'application/json',
  accept: 'application/json',
  Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
})

// ================= FORMAT DATE =================
const formatDate = (date: string): string => {
  if (!date) return '-'

  const d = new Date(date)

  if (isNaN(d.getTime())) return date

  return d.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const getTahunAkademik = async (): Promise<void> => {
  try {
    const res = await fetch(
      'https://be.karlearn.site/api/tahun-akademik/status/aktif',
      {
        method: 'GET',
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    console.log('DATA API:', json)

    const raw = Array.isArray(json.data)
      ? json.data
      : []

    allData.value = raw.map((item: any) => ({
      id: item.id,

      // ambil tipe_semester asli
      tipeSemester:
        item.tipe_semester === 'ganjil'
          ? 'Ganjil'
          : item.tipe_semester === 'genap'
          ? 'Genap'
          : '-',

      rawTipeSemester: item.tipe_semester,

      tahunAwal: item.tahun_awal
        ? item.tahun_awal.split('-')[0]
        : '-',

      tahunAkhir: item.tahun_akhir
        ? item.tahun_akhir.split('-')[0]
        : '-',

      status:
        item.status === 'aktif'
          ? 'Aktif/jalan'
          : 'Tidak Aktif',
    }))

    filteredData.value = [...allData.value]

    console.log(filteredData.value)

  } catch (err) {
    console.error('GET ERROR:', err)
  }
}

onMounted(() => {
  getTahunAkademik()
})

// ================= FILTER =================
const applyFilter = (): void => {
  filteredData.value = allData.value.filter((item) => {
    const semesterMatch = filterSemester.value
      ? item.tipeSemester === filterSemester.value
      : true

    const tahunMatch = filterTahun.value
      ? item.tahunAwal.includes(filterTahun.value)
      : true

    return semesterMatch && tahunMatch
  })

  currentPage.value = 1
}

// ================= SEMESTER OPTIONS =================
const semesterOptions = computed(() => {
  return [...new Set(allData.value.map((item) => item.tipeSemester))]
})

// ================= TOTAL PAGE =================
const totalPages = computed(() => {
  return Math.max(
    1,
    Math.ceil(filteredData.value.length / perPage.value)
  )
})

// ================= PAGINATED DATA =================
const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value

  return filteredData.value.slice(
    start,
    start + perPage.value
  )
})

// ================= VISIBLE PAGE =================
const visiblePages = computed(() => {
  return Array.from(
    { length: totalPages.value },
    (_, i) => i + 1
  )
})

// ================= PREV PAGE =================
const prevPage = (): void => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

// ================= NEXT PAGE =================
const nextPage = (): void => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// ================= EDIT =================
const editItem = (item: AkademikItem): void => {
  editingItem.value = item

  editForm.value = {
    tipeSemester: item.tipeSemester,
    tahunAwal: item.tahunAwal,
    tahunAkhir: item.tahunAkhir,
    status: item.status ?? '',
  }

  showModal.value = true
}

// ================= SAVE EDIT =================
const saveEdit = async (): Promise<void> => {
  if (!editingItem.value) return

  try {

    const payload = {
      id: editingItem.value.id,

      status:
        editForm.value.status === 'Aktif/jalan'
          ? 'aktif'
          : 'tidak_aktif',

      tahun_awal: `${editForm.value.tahunAwal}-01-01`,

      tahun_akhir: `${editForm.value.tahunAkhir}-01-01`,

      tipe_semester:
        editForm.value.tipeSemester.toLowerCase(),
    }

    console.log('PAYLOAD:', payload)

    const res = await fetch(
      `https://be.karlearn.site/api/tahun-akademik/${editingItem.value.id}`,
      {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      }
    )

    const json = await res.json()

    console.log('UPDATE:', json)

    if (!res.ok) {
      alert(json.message || 'Gagal update')
      return
    }

    alert('Berhasil update')

    showModal.value = false

    await getTahunAkademik()

  } catch (err) {
    console.error('UPDATE ERROR:', err)
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 p-6 font-sans">

    <!-- Breadcrumb -->
    <nav class="text-xs text-slate-400 mb-1 flex items-center gap-1">
      <span>Akademik</span>
      <span>›</span>
      <span class="text-slate-700 font-medium">Tahun Akademik</span>
    </nav>

    <!-- Judul -->
    <h1 class="text-2xl font-bold text-slate-800">
      Tahun Akademik
    </h1>

    <p class="text-slate-400 text-sm mb-6">
      Kelola data tahun akademik
    </p>

    <!-- Card -->
    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-6">

      <!-- Tabel -->
      <div class="overflow-x-auto">

        <table class="w-full text-sm">

          <!-- Header -->
          <thead>

            <!-- Judul di dalam tabel -->
            <tr>
              <th
                colspan="6"
                class="text-left text-2xl font-bold text-slate-700 pb-6"
              >
                Tahun Akademik
                {{ filteredData[0]?.tahunAwal }}-{{ filteredData[0]?.tahunAkhir }}
              </th>
            </tr>

            <!-- Header tabel -->
            <tr class="border-b border-slate-200">

              <th class="text-left py-3 px-4 font-semibold text-slate-600 w-12">
                No
              </th>

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Semester
              </th>

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Tahun Awal
              </th>

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Tahun Akhir
              </th>

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Status
              </th>

              <th class="text-center py-3 px-4 font-semibold text-slate-600">
                Aksi
              </th>

            </tr>
          </thead>

          <!-- Body -->
          <tbody>

            <!-- Jika kosong -->
            <tr v-if="paginatedData.length === 0">
              <td
                colspan="6"
                class="py-10 text-center text-slate-400"
              >
                Tidak ada data semester
              </td>
            </tr>

            <!-- Data -->
            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="border-b border-slate-100 hover:bg-slate-50"
            >

              <td class="py-4 px-4">
                {{ index + 1 }}
              </td>

              <td class="py-4 px-4 font-medium text-slate-700">
                {{ item.tipeSemester }}
              </td>

              <td class="py-4 px-4">
                {{ item.tahunAwal }}
              </td>

              <td class="py-4 px-4">
                {{ item.tahunAkhir }}
              </td>

              <td class="py-4 px-4">
                <span
                  class="px-3 py-1 rounded-full text-xs"
                  :class="
                    item.status === 'Aktif/jalan'
                      ? 'bg-green-100 text-green-700'
                      : 'bg-red-100 text-red-700'
                  "
                >
                  {{ item.status }}
                </span>
              </td>

              <td class="py-4 px-4 text-center">
                <button
                  @click="editItem(item)"
                  class="bg-amber-400 hover:bg-amber-500 text-white text-xs px-3 py-1.5 rounded-lg"
                >
                  Edit
                </button>
              </td>

            </tr>

          </tbody>

        </table>
<!-- Pagination -->
<div
  v-if="totalPages > 0"
  class="flex items-center justify-end gap-2 mt-6"
>

  <!-- Previous -->
  <button
    @click="prevPage"
    :disabled="currentPage === 1"
    class="px-4 py-2 text-sm border rounded-lg transition"
    :class="
      currentPage === 1
        ? 'bg-slate-100 text-slate-400 cursor-not-allowed'
        : 'bg-white hover:bg-slate-100 text-slate-700'
    "
  >
    Previous
  </button>

  <!-- Number -->
  <template
    v-for="(page, index) in displayedPages"
    :key="index"
  >

    <!-- Titik -->
    <span
      v-if="page === '...'"
      class="px-2 text-slate-500"
    >
      ...
    </span>

    <!-- Button angka -->
    <button
      v-else
      @click="currentPage = Number(page)"
      class="w-10 h-10 rounded-lg text-sm font-medium transition"
      :class="
        currentPage === page
          ? 'bg-[#1f3c93] text-white'
          : 'bg-white border border-slate-200 text-slate-700 hover:bg-slate-100'
      "
    >
      {{ page }}
    </button>

  </template>

  <!-- Next -->
  <button
    @click="nextPage"
    :disabled="currentPage === totalPages"
    class="px-4 py-2 text-sm border rounded-lg transition"
    :class="
      currentPage === totalPages
        ? 'bg-slate-100 text-slate-400 cursor-not-allowed'
        : 'bg-white hover:bg-slate-100 text-slate-700'
    "
  >
    Next
  </button>

</div>
      </div>
    </div>

    <!-- ═══════════════════════════════════════
         MODAL EDIT
    ════════════════════════════════════════ -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50"
      @click.self="showModal = false"
    >

      <div class="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md mx-4">

        <h3 class="text-base font-semibold text-slate-800 mb-5">
          Edit Tahun Akademik
        </h3>

        <div class="space-y-4">

          <!-- Tipe Semester -->
          <div>
            <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
              Tipe Semester
            </label>

            <select
              v-model="editForm.tipeSemester"
              class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400"
            >
              <option value="">Pilih Semester</option>
              <option>Ganjil</option>
              <option>Genap</option>
            </select>
          </div>

          <!-- Tahun -->
          <div class="flex gap-3">

            <div class="flex-1">
              <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
                Tahun Awal
              </label>

              <input
                v-model="editForm.tahunAwal"
                type="text"
                placeholder="cth: 2023"
                class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400"
              />
            </div>

            <div class="flex-1">
              <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
                Tahun Akhir
              </label>

              <input
                v-model="editForm.tahunAkhir"
                type="text"
                placeholder="cth: 2024"
                class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400"
              />
            </div>

          </div>

          <!-- Status -->
          <div>
            <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
              Status
            </label>

            <select
              v-model="editForm.status"
              class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400"
            >
              <option value="">Pilih Status</option>
              <option>Aktif/jalan</option>
              <option>Tidak Aktif</option>
            </select>
          </div>

        </div>

        <!-- Tombol -->
        <div class="flex justify-end gap-3 mt-6">

          <button
            @click="showModal = false"
            class="px-4 py-2 text-sm text-slate-600 border border-slate-200 hover:bg-slate-50 rounded-lg transition-colors"
          >
            Batal
          </button>

          <button
            @click="saveEdit"
            class="px-4 py-2 text-sm bg-[#1f3c93] hover:bg-blue-800 text-white font-semibold rounded-lg transition-colors"
          >
            Simpan
          </button>

        </div>

      </div>
    </div>

  </div>
</template>