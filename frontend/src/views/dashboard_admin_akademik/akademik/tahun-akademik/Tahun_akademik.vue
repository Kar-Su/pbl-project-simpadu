<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// ================= TYPE =================
interface AkademikItem {
  id: number
  tipeeSemester: string
  tahunAwal: string
  tahunAkhir: string
  status?: string
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
  tipeeSemester: '',
  tahunAwal: '',
  tahunAkhir: '',
  status: '',
})

// ================= WATCH TAHUN =================
watch(() => editForm.value.tahunAwal, (val) => {
  const angka = Number(val)
  if (val && val.length === 4 && !isNaN(angka)) {
    editForm.value.tahunAkhir = String(angka + 1)
  } else {
    editForm.value.tahunAkhir = ''
  }
})

// ================= HEADER =================
const getHeaders = () => ({
  'Content-Type': 'application/json',
  accept: 'application/json',
  Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
})

// ================= FETCH API =================
const getTahunAkademik = async (): Promise<void> => {
  try {
    const res = await fetch(
      'https://be.karlearn.site/api/tahun-akademik',
      {
        method: 'GET',
        headers: getHeaders(),
      }
    )

    const json = await res.json()
    console.log('RESPONSE:', json)

    const raw = Array.isArray(json.data) ? json.data : []

    allData.value = raw.map((item: any) => ({
      id: item.id,
      tipeeSemester:
        item.tipee_semester === 'ganjil'
          ? 'Ganjil'
          : item.tipee_semester === 'genap'
            ? 'Genap'
            : '-',

      rawTipeSemester: item.tipe_semester,

      tahunAwal: item.tahun_awal
        ? item.tahun_awal.split('-')[0]
        : '-',

      tahunAkhir: item.tahun_akhir
        ? item.tahun_akhir.split('-')[0]
        : '-',

      status: item.status === 'aktif'
        ? 'Aktif/jalan'
        : 'Non Aktif',
    }))
console.log("ALL DATA:", raw)
    filteredData.value = [...allData.value]
    currentPage.value = 1
  } catch (err) {
    console.error('GET ERROR:', err)
  }
}

onMounted(() => {
  getTahunAkademik()
})

// ================= FILTER =================
const applyFilter = () => {
  filteredData.value = allData.value.filter((item) => {
    const semesterMatch = filterSemester.value
      ? item.tipeeSemester === filterSemester.value
      : true

    const tahunMatch = filterTahun.value
      ? item.tahunAwal.includes(filterTahun.value)
      : true

    return semesterMatch && tahunMatch
  })

  currentPage.value = 1
}

// ================= PAGINATION =================
const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredData.value.length / perPage.value))
)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredData.value.slice(start, start + perPage.value)
})

// ================= PAGE NAV =================
const nextPage = () => {
  if (currentPage.value < totalPages.value) currentPage.value++
}

const prevPage = () => {
  if (currentPage.value > 1) currentPage.value--
}

// ================= DISPLAY PAGE NUMBER =================
const displayedPages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 3) return Array.from({ length: total }, (_, i) => i + 1)

  if (current <= 2) return [1, 2, '...', total]

  if (current >= total - 1) return [1, '...', total - 1, total]

  return [1, '...', current, '...', total]
})

const isDuplicateData = computed(() => {
  if (!editingItem.value) return false

  return allData.value.some((item) => {
    return (
      item.id !== editingItem.value?.id &&
      item.tipeeSemester.toLowerCase() ===
      editForm.value.tipeeSemester.toLowerCase() &&
      item.tahunAwal === editForm.value.tahunAwal
    )
  })
})

const isInvalidYear = computed(() => {
  return (
    editForm.value.tahunAwal.length > 0 &&
    editForm.value.tahunAwal.length < 4
  )
})

// ================= EDIT =================
const editItem = (item: AkademikItem) => {
  editingItem.value = item

  editForm.value = {
    tipeeSemester: item.tipeeSemester,
    tahunAwal: item.tahunAwal,
    tahunAkhir: item.tahunAkhir,
    status: item.status ?? '',
  }

  showModal.value = true
}

// ================= SAVE EDIT =================
const saveEdit = async () => {
  if (!editingItem.value) return
  if (isDuplicateData.value) {
    alert('Semester dan tahun awal sudah terdaftar')
    return
  }
  if (isInvalidYear.value) {
    alert('Tahun harus 4 digit')
    return
  }

  if (isDuplicateData.value) {
    alert('Semester dan tahun awal sudah terdaftar')
    return
  }
  try {
    const payload = {
      id: editingItem.value.id,
      status:
        editForm.value.status === 'Aktif/jalan'
          ? 'aktif'
          : 'nonaktif',
      tahun_awal: `${editForm.value.tahunAwal}-01-01`,
      tahun_akhir: `${editForm.value.tahunAkhir}-01-01`,
      tipe_semester: editForm.value.tipeeSemester.toLowerCase(),
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
    console.log('PUT RESPONSE:', json)

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
// ================= TOGGLE STATUS =================
const toggleStatus = async (item: AkademikItem) => {
  try {
    const payload = {
      id: item.id,
      status: item.status === 'Aktif/jalan'
        ? 'nonaktif'
        : 'aktif',
      tahun_awal: `${item.tahunAwal}-01-01`,
      tahun_akhir: `${item.tahunAkhir}-01-01`,
      tipe_semester:
        item.rawTipeSemester ??
        item.tipeeSemester.toLowerCase(),
    }

    const res = await fetch(
      `https://be.karlearn.site/api/tahun-akademik/${item.id}`,
      {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      }
    )

    const json = await res.json()
    console.log("STATUS:", res.status)
    console.log("RESPONSE:", json)
    console.log("PAYLOAD:", payload)

    if (!res.ok) {
      alert(json.message || 'Gagal update status')
      return
    }

    await getTahunAkademik()
  } catch (err) {
    console.error('TOGGLE ERROR:', err)
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
    <div class="col-span-3 bg-[#ececec] rounded-xl p-5 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">

      <!-- Tabel -->
      <div class="overflow-x-auto">

        <table class="w-full text-sm">

          <!-- Header -->
          <thead>

            <!-- Judul di dalam tabel -->
            <tr>
              <th colspan="6" class="text-left text-2xl font-bold text-black-800 pb-6">
                Tahun Akademik
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
              <td colspan="6" class="py-10 text-center text-slate-400">
                Tidak ada data semester
              </td>
            </tr>

            <!-- Data -->
            <tr v-for="(item, index) in paginatedData" :key="item.id"
              class="border-b border-slate-100 hover:bg-slate-50">

              <td class="py-4 px-4">
                {{ index + 1 }}
              </td>

              <td class="py-4 px-4 font-medium text-slate-700">
                {{ item.tipeeSemester }}
              </td>

              <td class="py-4 px-4">
                {{ item.tahunAwal }}
              </td>

              <td class="py-4 px-4">
                {{ item.tahunAkhir }}
              </td>

              <td class="py-4 px-4">
                <span class="px-3 py-1 rounded-full text-xs" :class="item.status === 'Aktif/jalan'
                  ? 'bg-green-100 text-green-700'
                  : 'bg-red-100 text-red-700'
                  ">
                  {{ item.status }}
                </span>
              </td>
              <td class="py-4 px-4 text-center">
                <div class="flex items-center justify-center gap-2">

                  <!-- Tombol Toggle Status -->
                  <button @click="toggleStatus(item)" class="text-xs px-3 py-1.5 rounded-lg text-white transition"
                    :class="item.status === 'Aktif/jalan'
                      ? 'bg-red-400 hover:bg-red-500'
                      : 'bg-green-500 hover:bg-green-600'
                      ">
                    {{ item.status === 'Aktif/jalan' ? 'Nonaktifkan' : 'Aktifkan' }}
                  </button>

                  <!-- Tombol Edit -->
                  <button @click="editItem(item)"
                    class="bg-amber-400 hover:bg-amber-500 text-white text-xs px-3 py-1.5 rounded-lg">
                    Edit
                  </button>

                </div>
              </td>

            </tr>

          </tbody>

        </table>
        <!-- Pagination -->
        <div v-if="totalPages > 0" class="flex items-center justify-end gap-2 mt-6">

          <!-- Previous -->
          <button @click="prevPage" :disabled="currentPage === 1" class="px-4 py-2 text-sm border rounded-lg transition"
            :class="currentPage === 1
              ? 'bg-slate-100 text-slate-400 cursor-not-allowed'
              : 'bg-white hover:bg-slate-100 text-slate-700'
              ">
            Previous
          </button>

          <!-- Number -->
          <template v-for="(page, index) in displayedPages" :key="index">

            <!-- Titik -->
            <span v-if="page === '...'" class="px-2 text-slate-500">
              ...
            </span>

            <!-- Button angka -->
            <button v-else @click="currentPage = Number(page)"
              class="w-10 h-10 rounded-lg text-sm font-medium transition" :class="currentPage === page
                ? 'bg-[#1f3c93] text-white'
                : 'bg-white border border-slate-200 text-slate-700 hover:bg-slate-100'
                ">
              {{ page }}
            </button>

          </template>

          <!-- Next -->
          <button @click="nextPage" :disabled="currentPage === totalPages"
            class="px-4 py-2 text-sm border rounded-lg transition" :class="currentPage === totalPages
              ? 'bg-slate-100 text-slate-400 cursor-not-allowed'
              : 'bg-white hover:bg-slate-100 text-slate-700'
              ">
            Next
          </button>

        </div>
      </div>
    </div>

    <!-- ═══════════════════════════════════════
         MODAL EDIT
    ════════════════════════════════════════ -->
    <div v-if="showModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50"
      @click.self="showModal = false">

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
            <select v-model="editForm.tipeeSemester"
              class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400">
              <option value="" disabled>Pilih Semester</option>
              <option>Ganjil</option>
              <option>Genap</option>
            </select>
            <p v-if="isDuplicateData" class="text-red-500 text-xs mt-1">
              Kombinasi semester dan tahun awal sudah digunakan.
            </p>
          </div>

          <!-- Tahun -->
          <div class="flex gap-3">

            <!-- Tahun Awal -->
            <div class="flex-1">
              <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
                Tahun Awal
              </label>

              <input v-model="editForm.tahunAwal" type="text" inputmode="numeric" maxlength="4" placeholder="cth: 2024"
                @input="editForm.tahunAwal = editForm.tahunAwal.replace(/\D/g, '').slice(0, 4)"
                class="w-full border rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400"
                :class="(isDuplicateData || editForm.tahunAwal.length > 0 && editForm.tahunAwal.length < 4)
                  ? 'border-red-500'
                  : 'border-slate-200'" />

              <!-- Error 4 digit -->
              <p v-if="editForm.tahunAwal.length > 0 && editForm.tahunAwal.length < 4"
                class="text-red-500 text-xs mt-1">
                Tahun harus terdiri dari 4 digit.
              </p>

              <!-- Error duplikat -->
              <p v-else-if="isDuplicateData" class="text-red-500 text-xs mt-1">
                Tahun awal dan semester tersebut sudah terdaftar.
              </p>
            </div>

            <!-- Tahun Akhir -->
            <div class="flex-1">
              <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
                Tahun Akhir
              </label>

              <input :value="editForm.tahunAkhir || '-'" type="text" disabled
                class="w-full border border-slate-200 bg-slate-100 rounded-lg px-3 py-2.5 text-sm text-slate-400 cursor-not-allowed" />
            </div>

          </div>


        </div>

        <!-- Tombol -->
        <div class="flex justify-end gap-3 mt-6">

          <button @click="showModal = false"
            class="px-4 py-2 text-sm text-slate-600 border border-slate-200 hover:bg-slate-50 rounded-lg transition-colors">
            Batal
          </button>

          <button @click="saveEdit" :disabled="isDuplicateData || isInvalidYear"
            class="px-4 py-2 text-sm text-white font-semibold rounded-lg transition-colors" :class="isDuplicateData || isInvalidYear
              ? 'bg-slate-400 cursor-not-allowed'
              : 'bg-[#1f3c93] hover:bg-blue-800'
              ">
            Simpan
          </button>

        </div>

      </div>
    </div>

  </div>
</template>