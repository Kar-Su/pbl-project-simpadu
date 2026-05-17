<template>
  <div class="min-h-screen bg-slate-100 p-6 font-sans">
    <!-- Breadcrumb -->
    <nav class="text-sm text-slate-500 mb-1">
      <span>Akademik</span>
      <span class="mx-1">›</span>
      <span class="text-slate-700">Tahun Akademik</span>
    </nav>

    <!-- Page Title -->
    <h1 class="text-2xl font-bold text-slate-800">Tahun Akademik</h1>
    <p class="text-slate-400 text-sm mb-6">lorem ipsum</p>

    <!-- Card -->
    <div class="bg-white rounded-2xl shadow-sm p-6">
      <h2 class="text-base font-semibold text-slate-700 mb-4">Data Tahun Akademik</h2>

      <!-- Filters -->
      <div class="flex flex-wrap gap-3 mb-6">
        <select
          v-model="filterTahunAkademik"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm text-slate-600 bg-white focus:outline-none focus:ring-2 focus:ring-blue-400 min-w-[150px]"
        >
          <option value="">Tahun Akademik</option>
          <option value="2024/2025">2024/2025</option>
          <option value="2025/2026">2025/2026</option>
          <option value="2026/2027">2026/2027</option>
        </select>

        <select
          v-model="filterTahunAwal"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm text-slate-600 bg-white focus:outline-none focus:ring-2 focus:ring-blue-400 min-w-[150px]"
        >
          <option value="">Tahun Awal</option>
          <option value="2024">2024</option>
          <option value="2025">2025</option>
          <option value="2026">2026</option>
        </select>

        <select
          v-model="filterTahunAkhir"
          class="border border-slate-200 rounded-lg px-3 py-2 text-sm text-slate-600 bg-white focus:outline-none focus:ring-2 focus:ring-blue-400 min-w-[150px]"
        >
          <option value="">Tahun Akhir</option>
          <option value="2025">2025</option>
          <option value="2026">2026</option>
          <option value="2027">2027</option>
        </select>

        <button
          @click="applyFilter"
          class="flex items-center gap-2 bg-blue-900 hover:bg-blue-800 text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2a1 1 0 01-.293.707L13 13.414V19a1 1 0 01-.553.894l-4 2A1 1 0 017 21v-7.586L3.293 6.707A1 1 0 013 6V4z" />
          </svg>
          Pilih
        </button>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-slate-100">
              <th class="text-left py-3 px-4 font-semibold text-slate-600 w-12">No</th>
              <th class="text-center py-3 px-4 font-semibold text-slate-600">Tipe Semester</th>
              <th class="text-center py-3 px-4 font-semibold text-slate-600">Tahun Awal</th>
              <th class="text-center py-3 px-4 font-semibold text-slate-600">Tahun Akhir</th>
              <th class="text-center py-3 px-4 font-semibold text-slate-600">Status</th>
              <th class="text-center py-3 px-4 font-semibold text-slate-600">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="border-b border-slate-50 hover:bg-slate-50 transition-colors"
            >
              <td class="py-3 px-4 text-slate-600">{{ (currentPage - 1) * perPage + index + 1 }}</td>
              <td class="py-3 px-4 text-center font-medium text-slate-700">{{ item.tipeSemester }}</td>
              <td class="py-3 px-4 text-center text-slate-700">{{ item.tahunAwal }}</td>
              <td class="py-3 px-4 text-center text-slate-700">{{ item.tahunAkhir }}</td>
              <td class="py-3 px-4 text-center">
                <span
                  :class="item.status === 'Aktif/jalan'
                    ? 'text-green-600 font-semibold'
                    : 'text-slate-500'"
                >
                  {{ item.status }}
                </span>
              </td>
              <td class="py-3 px-4 text-center">
                <button
                  @click="editItem(item)"
                  class="flex items-center gap-1.5 bg-amber-400 hover:bg-amber-500 text-white text-xs font-medium px-3 py-1.5 rounded-lg transition-colors mx-auto"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Edit
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Footer: rows per page + pagination -->
      <div class="flex items-center justify-between mt-5 flex-wrap gap-3">
        <!-- Rows per page -->
        <select
          v-model="perPage"
          @change="currentPage = 1"
          class="border border-slate-200 rounded-lg px-3 py-1.5 text-sm text-slate-600 bg-white focus:outline-none focus:ring-2 focus:ring-blue-400"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="20">20 Baris</option>
        </select>

        <!-- Pagination -->
        <div class="flex items-center gap-1 text-sm">
          <button
            @click="currentPage > 1 && currentPage--"
            :disabled="currentPage === 1"
            class="px-3 py-1.5 rounded-lg text-slate-500 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            ← Previous
          </button>

          <template v-for="page in visiblePages" :key="page">
            <span v-if="page === '...'" class="px-2 py-1.5 text-slate-400">...</span>
            <button
              v-else
              @click="currentPage = Number(page)"
              :class="[
                'w-8 h-8 rounded-lg font-medium transition-colors',
                currentPage === page
                  ? 'bg-blue-900 text-white'
                  : 'text-slate-600 hover:bg-slate-100'
              ]"
            >
              {{ page }}
            </button>
          </template>

          <button
            @click="currentPage < totalPages && currentPage++"
            :disabled="currentPage === totalPages"
            class="px-3 py-1.5 rounded-lg text-slate-500 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            Next →
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
      @click.self="showModal = false"
    >
      <div class="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md mx-4">
        <h3 class="text-base font-semibold text-slate-800 mb-4">Edit Tahun Akademik</h3>
        <div class="space-y-3">
          <div>
            <label class="text-xs font-medium text-slate-500 mb-1 block">Tipe Semester</label>
            <select v-model="editForm.tipeSemester" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400">
              <option>Ganjil</option>
              <option>Genap</option>
            </select>
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="text-xs font-medium text-slate-500 mb-1 block">Tahun Awal</label>
              <input v-model="editForm.tahunAwal" type="number" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
            </div>
            <div class="flex-1">
              <label class="text-xs font-medium text-slate-500 mb-1 block">Tahun Akhir</label>
              <input v-model="editForm.tahunAkhir" type="number" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400" />
            </div>
          </div>
          <div>
            <label class="text-xs font-medium text-slate-500 mb-1 block">Status</label>
            <select v-model="editForm.status" class="w-full border border-slate-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400">
              <option>Aktif/jalan</option>
              <option>Tidak Aktif</option>
            </select>
          </div>
        </div>
        <div class="flex justify-end gap-3 mt-5">
          <button @click="showModal = false" class="px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">Batal</button>
          <button @click="saveEdit" class="px-4 py-2 text-sm bg-blue-900 hover:bg-blue-800 text-white font-medium rounded-lg transition-colors">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface AkademikItem {
  id: number
  tipeSemester: string
  tahunAwal: number
  tahunAkhir: number
  status: string
}

// --- State ---
const filterTahunAkademik = ref('')
const filterTahunAwal = ref('')
const filterTahunAkhir = ref('')
const currentPage = ref(1)
const perPage = ref(5)
const showModal = ref(false)
const editingItem = ref<AkademikItem | null>(null)
const editForm = ref<Omit<AkademikItem, 'id'>>({
  tipeSemester: '',
  tahunAwal: 0,
  tahunAkhir: 0,
  status: '',
})

// --- Data ---
const allData = ref<AkademikItem[]>([
  { id: 1, tipeSemester: 'Ganjil', tahunAwal: 2024, tahunAkhir: 2025, status: 'Aktif/jalan' },
  { id: 2, tipeSemester: 'Genap', tahunAwal: 2024, tahunAkhir: 2025, status: 'Tidak Aktif' },
  { id: 3, tipeSemester: 'Ganjil', tahunAwal: 2025, tahunAkhir: 2026, status: 'Tidak Aktif' },
  { id: 4, tipeSemester: 'Genap', tahunAwal: 2025, tahunAkhir: 2026, status: 'Tidak Aktif' },
  { id: 5, tipeSemester: 'Ganjil', tahunAwal: 2026, tahunAkhir: 2027, status: 'Tidak Aktif' },
  { id: 6, tipeSemester: 'Genap', tahunAwal: 2026, tahunAkhir: 2027, status: 'Tidak Aktif' },
  { id: 7, tipeSemester: 'Ganjil', tahunAwal: 2027, tahunAkhir: 2028, status: 'Tidak Aktif' },
  // Generate more for pagination demo
  ...Array.from({ length: 330 }, (_, i) => ({
    id: i + 8,
    tipeSemester: i % 2 === 0 ? 'Ganjil' : 'Genap',
    tahunAwal: 2028 + Math.floor(i / 2),
    tahunAkhir: 2029 + Math.floor(i / 2),
    status: 'Tidak Aktif',
  })),
])

const filteredData = ref<AkademikItem[]>([...allData.value])

// --- Computed ---
const totalPages = computed(() => Math.ceil(filteredData.value.length / perPage.value))

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredData.value.slice(start, start + perPage.value)
})

const visiblePages = computed((): (number | string)[] => {
  const total = totalPages.value
  const cur = currentPage.value
  const pages: (number | string)[] = []

  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  pages.push(1)
  if (cur > 3) pages.push('...')
  for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++) {
    pages.push(i)
  }
  if (cur < total - 2) pages.push('...')
  pages.push(total)

  return pages
})

// --- Methods ---
function applyFilter() {
  currentPage.value = 1
  filteredData.value = allData.value.filter((item) => {
    const matchAwal = filterTahunAwal.value ? item.tahunAwal === Number(filterTahunAwal.value) : true
    const matchAkhir = filterTahunAkhir.value ? item.tahunAkhir === Number(filterTahunAkhir.value) : true
    return matchAwal && matchAkhir
  })
}

function editItem(item: AkademikItem) {
  editingItem.value = item
  editForm.value = { ...item }
  showModal.value = true
}

function saveEdit() {
  if (!editingItem.value) return
  const idx = allData.value.findIndex((d) => d.id === editingItem.value!.id)
  if (idx !== -1) {
    allData.value[idx] = { id: editingItem.value.id, ...editForm.value }
  }
  applyFilter()
  showModal.value = false
}
</script>