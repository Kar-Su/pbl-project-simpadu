<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const goToTambahTahunAkademik = () => {
  router.push('/dashboard-admin/tambah_tahunakademik')
}

// ================= TYPE =================
interface AkademikItem {
  id: number
  tipeSemester: string
  rawTipeSemester?: string
  awalTanggal: number
  awalBulan: number
  awalTahun: number
  akhirTanggal: number
  akhirBulan: number
  akhirTahun: number
  status?: string
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
  awalTanggal: '',
  awalBulan: '',
  awalTahun: '',
  akhirTanggal: '',
  akhirBulan: '',
  akhirTahun: '',
  status: '',
})

// PERINGATAN: tahun akademik yang seharusnya non aktif tapi gagal
// dinonaktifkan otomatis (misal endpoint menolak / tidak ada hak akses).
// Selama daftar ini tidak kosong, banner peringatan akan tampil
// supaya admin tahu harus minta tim backend membenarkan data secara manual.
const duplicateActiveWarning = ref<string[]>([])
const isFixingDuplicates = ref(false)

// ================= NAMA BULAN (untuk tampilan) =================
const NAMA_BULAN = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

// ================= HELPER TANGGAL =================
const pad2 = (n: number) => String(n).padStart(2, '0')

const isTahunKabisat = (tahun: number): boolean => {
  if (!tahun) return false
  return (tahun % 4 === 0 && tahun % 100 !== 0) || tahun % 400 === 0
}

// Mengembalikan jumlah hari maksimal untuk bulan & tahun tertentu
// (Februari otomatis menyesuaikan tahun kabisat)
const getMaxTanggal = (bulan: number, tahun: number): number => {
  if (!bulan || bulan < 1 || bulan > 12) return 31
  const hariPerBulan = [
    31, isTahunKabisat(tahun) ? 29 : 28, 31, 30, 31, 30,
    31, 31, 30, 31, 30, 31,
  ]
  return hariPerBulan[bulan - 1]
}

// Format tampilan tanpa tanda "-", contoh: "12 April 2023"
const formatTanggal = (tanggal?: number, bulan?: number, tahun?: number): string => {
  if (!tanggal || !bulan || !tahun) return '-'
  return `${tanggal} ${NAMA_BULAN[bulan - 1] ?? '-'} ${tahun}`
}

// Pecah string tanggal dari API ("YYYY-MM-DD") jadi tanggal/bulan/tahun
const parseTanggalAPI = (str?: string) => {
  if (!str) return { tanggal: 0, bulan: 0, tahun: 0 }
  const [tahun, bulan, tanggal] = str.split('-').map(Number)
  return { tahun: tahun || 0, bulan: bulan || 0, tanggal: tanggal || 0 }
}

// Batasi input bulan: hanya angka, max 2 digit, tidak boleh > 12
const clampBulan = (val: string): string => {
  let v = val.replace(/\D/g, '').slice(0, 2)
  if (v && Number(v) > 12) v = '12'
  return v
}

// Batasi input tanggal: hanya angka, max 2 digit, tidak boleh melebihi
// jumlah hari pada bulan (dan tahun, khusus Februari) yang sedang dipilih
const clampTanggal = (val: string, bulan: string, tahun: string): string => {
  let v = val.replace(/\D/g, '').slice(0, 2)
  const maxTanggal = getMaxTanggal(Number(bulan), Number(tahun))
  if (v && Number(v) > maxTanggal) v = String(maxTanggal)
  return v
}

// ================= WATCH: HITUNG TAHUN AKHIR OTOMATIS =================
// Tahun Akhir selalu sama tanggal & bulannya dengan Tahun Awal,
// hanya tahunnya yang ditambah 1. Field Tahun Akhir di form selalu disabled.
watch(
  () => [editForm.value.awalTanggal, editForm.value.awalBulan, editForm.value.awalTahun],
  () => {
    const { awalTanggal, awalBulan, awalTahun } = editForm.value
    if (awalTanggal && awalBulan && awalTahun.length === 4) {
      editForm.value.akhirTanggal = awalTanggal
      editForm.value.akhirBulan = awalBulan
      editForm.value.akhirTahun = String(Number(awalTahun) + 1)
    } else {
      editForm.value.akhirTanggal = ''
      editForm.value.akhirBulan = ''
      editForm.value.akhirTahun = ''
    }
  }
)

// ================= INPUT HANDLER (TANGGAL AWAL) =================
const onAwalTanggalInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  editForm.value.awalTanggal = clampTanggal(target.value, editForm.value.awalBulan, editForm.value.awalTahun)
}

const onAwalBulanInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  editForm.value.awalBulan = clampBulan(target.value)
  // bulan berubah -> jumlah hari maksimal bisa berubah, cek ulang tanggal yang sudah diisi
  editForm.value.awalTanggal = clampTanggal(editForm.value.awalTanggal, editForm.value.awalBulan, editForm.value.awalTahun)
}

const onAwalTahunInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  editForm.value.awalTahun = target.value.replace(/\D/g, '').slice(0, 4)
  // tahun berubah -> kemungkinan kabisat/non-kabisat berubah, cek ulang tanggal (khusus Februari)
  editForm.value.awalTanggal = clampTanggal(editForm.value.awalTanggal, editForm.value.awalBulan, editForm.value.awalTahun)
}

// Jumlah hari maksimal untuk bulan & tahun yang sedang dipilih di form (buat hint)
const maxTanggalAwal = computed(() =>
  getMaxTanggal(Number(editForm.value.awalBulan), Number(editForm.value.awalTahun))
)

// ================= HEADER =================
const getHeaders = () => ({
  'Content-Type': 'application/json',
  accept: 'application/json',
  Authorization: `Bearer ${localStorage.getItem('token') ?? ''}`,
})

const formatLabel = (item: AkademikItem) =>
  `${item.tipeSemester} ${formatTanggal(item.awalTanggal, item.awalBulan, item.awalTahun)} - ${formatTanggal(item.akhirTanggal, item.akhirBulan, item.akhirTahun)}`

// ================= FETCH (tanpa cek duplikat) =================
const fetchAndSetData = async (): Promise<void> => {
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

  allData.value = raw.map((item: any) => {
    const awal = parseTanggalAPI(item.tahun_awal)
    const akhir = parseTanggalAPI(item.tahun_akhir)

    return {
      id: item.id,
      tipeSemester:
        item.tipe_semester === 'ganjil'
          ? 'Ganjil'
          : item.tipe_semester === 'genap'
            ? 'Genap'
            : '-',

      rawTipeSemester: item.tipe_semester,

      awalTanggal: awal.tanggal,
      awalBulan: awal.bulan,
      awalTahun: awal.tahun,

      akhirTanggal: akhir.tanggal,
      akhirBulan: akhir.bulan,
      akhirTahun: akhir.tahun,

      status: item.status === 'aktif'
        ? 'Aktif/jalan'
        : 'Non Aktif',
    }
  })

  console.log('ALL DATA:', raw)
  filteredData.value = [...allData.value]
  currentPage.value = 1
}

// ================= NONAKTIFKAN 1 ITEM (dipakai untuk auto-fix & toggle) =================
const deactivateItem = async (item: AkademikItem): Promise<boolean> => {
  try {
    const payload = {
      id: item.id,
      status: 'nonaktif',
      tahun_awal: `${item.awalTahun}-${pad2(item.awalBulan)}-${pad2(item.awalTanggal)}`,
      tahun_akhir: `${item.akhirTahun}-${pad2(item.akhirBulan)}-${pad2(item.akhirTanggal)}`,
      tipe_semester: item.rawTipeSemester ?? item.tipeSemester.toLowerCase(),
    }

    const res = await fetch(
      `https://be.karlearn.site/api/tahun-akademik/${item.id}`,
      {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(payload),
      }
    )

    return res.ok
  } catch (err) {
    console.error('AUTO DEACTIVATE ERROR:', err)
    return false
  }
}

// ================= CEK & PERBAIKI TAHUN AKADEMIK AKTIF GANDA =================
// Idealnya hanya 1 tahun akademik yang aktif. Kalau backend mengirim lebih
// dari 1 (misal karena belum ada validasi unique-active di backend),
// yang pertama di list dipertahankan aktif, sisanya otomatis dinonaktifkan
// lewat endpoint PUT yang sama dengan tombol "Nonaktifkan".
// Kalau ada yang gagal dinonaktifkan (endpoint menolak / bukan hak akses kita),
// tampilkan peringatan supaya diminta diperbaiki manual di backend.
const checkAndFixMultipleActive = async () => {
  duplicateActiveWarning.value = []

  const activeItems = allData.value.filter((i) => i.status === 'Aktif/jalan')
  if (activeItems.length <= 1) return

  const dipertahankan = activeItems[0]
  const harusDinonaktifkan = activeItems.slice(1)

  console.warn(
    `Ditemukan ${activeItems.length} tahun akademik aktif bersamaan. ` +
    `Mempertahankan "${formatLabel(dipertahankan)}" aktif, mencoba menonaktifkan sisanya otomatis.`
  )

  isFixingDuplicates.value = true
  const gagal: AkademikItem[] = []

  for (const item of harusDinonaktifkan) {
    const ok = await deactivateItem(item)
    if (!ok) gagal.push(item)
  }

  isFixingDuplicates.value = false

  if (gagal.length > 0) {
    // Tidak semua berhasil dinonaktifkan otomatis -> butuh perbaikan manual di backend
    duplicateActiveWarning.value = gagal.map(formatLabel)
  } else {
    // Semua berhasil dinonaktifkan otomatis, ambil data terbaru tanpa alert mengganggu
    await fetchAndSetData()
  }
}

// ================= FETCH API (entrypoint utama) =================
const getTahunAkademik = async (): Promise<void> => {
  try {
    await fetchAndSetData()
    await checkAndFixMultipleActive()
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
      ? item.tipeSemester === filterSemester.value
      : true

    const tahunMatch = filterTahun.value
      ? String(item.awalTahun).includes(filterTahun.value)
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

// Duplikat dicek berdasarkan kombinasi semester + tanggal awal (tanggal, bulan, tahun)
const isDuplicateData = computed(() => {
  if (!editingItem.value) return false

  return allData.value.some((item) => {
    return (
      item.id !== editingItem.value?.id &&
      item.tipeSemester.toLowerCase() === editForm.value.tipeSemester.toLowerCase() &&
      String(item.awalTahun) === editForm.value.awalTahun &&
      String(item.awalBulan) === editForm.value.awalBulan &&
      String(item.awalTanggal) === editForm.value.awalTanggal
    )
  })
})

const isInvalidYear = computed(() => {
  return (
    editForm.value.awalTahun.length > 0 &&
    editForm.value.awalTahun.length < 4
  )
})

// Tanggal & bulan harus diisi berpasangan (tidak boleh hanya salah satu)
const isInvalidDate = computed(() => {
  const { awalTanggal, awalBulan } = editForm.value
  return (!!awalTanggal && !awalBulan) || (!awalTanggal && !!awalBulan)
})

// Form belum lengkap (belum bisa disimpan)
const isFormIncomplete = computed(() => {
  const { tipeSemester, awalTanggal, awalBulan, awalTahun } = editForm.value
  return !tipeSemester || !awalTanggal || !awalBulan || awalTahun.length !== 4
})

// ================= EDIT =================
const editItem = (item: AkademikItem) => {
  editingItem.value = item

  editForm.value = {
    tipeSemester: item.tipeSemester,
    awalTanggal: item.awalTanggal ? String(item.awalTanggal) : '',
    awalBulan: item.awalBulan ? String(item.awalBulan) : '',
    awalTahun: item.awalTahun ? String(item.awalTahun) : '',
    akhirTanggal: item.akhirTanggal ? String(item.akhirTanggal) : '',
    akhirBulan: item.akhirBulan ? String(item.akhirBulan) : '',
    akhirTahun: item.akhirTahun ? String(item.akhirTahun) : '',
    status: item.status ?? '',
  }

  showModal.value = true
}

// ================= SAVE EDIT =================
const saveEdit = async () => {
  if (!editingItem.value) return
  if (isFormIncomplete.value) {
    alert('Tanggal, bulan, dan tahun awal wajib diisi lengkap')
    return
  }
  if (isDuplicateData.value) {
    alert('Semester dan tanggal awal sudah terdaftar')
    return
  }
  if (isInvalidYear.value) {
    alert('Tahun harus 4 digit')
    return
  }
  if (isInvalidDate.value) {
    alert('Tanggal dan bulan harus diisi bersamaan')
    return
  }

  try {
    const payload = {
      id: editingItem.value.id,
      status:
        editForm.value.status === 'Aktif/jalan'
          ? 'aktif'
          : 'nonaktif',
      tahun_awal: `${editForm.value.awalTahun}-${pad2(Number(editForm.value.awalBulan))}-${pad2(Number(editForm.value.awalTanggal))}`,
      tahun_akhir: `${editForm.value.akhirTahun}-${pad2(Number(editForm.value.akhirBulan))}-${pad2(Number(editForm.value.akhirTanggal))}`,
      tipe_semester: editForm.value.tipeSemester.toLowerCase(),
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
// PERBAIKAN: tahun akademik tidak boleh aktif lebih dari 1 secara bersamaan.
// Kalau user klik "Aktifkan" sementara ada tahun akademik lain yang masih
// aktif, akan dikonfirmasi dulu lalu yang lama otomatis dinonaktifkan
// sebelum yang baru diaktifkan. Kalau proses nonaktifkan otomatis gagal
// (misal hak akses endpoint terbatas), proses dibatalkan dan user diberi
// tahu untuk minta tim backend menonaktifkannya secara manual.
const toggleStatus = async (item: AkademikItem) => {
  try {
    const akanDiaktifkan = item.status !== 'Aktif/jalan'

    if (akanDiaktifkan) {
      const aktifLain = allData.value.filter(
        (i) => i.id !== item.id && i.status === 'Aktif/jalan'
      )

      if (aktifLain.length > 0) {
        const namaLain = aktifLain.map(formatLabel).join(', ')

        const konfirmasi = confirm(
          `Tahun akademik "${namaLain}" sedang aktif.\n` +
          `Hanya boleh ada 1 tahun akademik aktif dalam satu waktu.\n\n` +
          `Mengaktifkan "${formatLabel(item)}" akan menonaktifkan "${namaLain}" secara otomatis. Lanjutkan?`
        )

        if (!konfirmasi) return

        for (const lain of aktifLain) {
          const ok = await deactivateItem(lain)
          if (!ok) {
            alert(
              `Gagal menonaktifkan otomatis "${formatLabel(lain)}". ` +
              `Mohon hubungi tim backend untuk menonaktifkannya secara manual ` +
              `sebelum mengaktifkan tahun akademik baru.`
            )
            return
          }
        }
      }
    }

    const payload = {
      id: item.id,
      status: akanDiaktifkan ? 'aktif' : 'nonaktif',
      tahun_awal: `${item.awalTahun}-${pad2(item.awalBulan)}-${pad2(item.awalTanggal)}`,
      tahun_akhir: `${item.akhirTahun}-${pad2(item.akhirBulan)}-${pad2(item.akhirTanggal)}`,
      tipe_semester:
        item.rawTipeSemester ??
        item.tipeSemester.toLowerCase(),
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

    <!-- BANNER: lagi memperbaiki duplikat aktif -->
    <div v-if="isFixingDuplicates" class="mb-4 rounded-lg border border-blue-200 bg-blue-50 px-4 py-3 text-sm text-blue-700">
      Mendeteksi lebih dari satu tahun akademik aktif, sedang mencoba menonaktifkan otomatis...
    </div>

    <!-- BANNER: ada duplikat aktif yang GAGAL diperbaiki otomatis -->
    <div v-if="duplicateActiveWarning.length > 0" class="mb-4 rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-700">
      <p class="font-semibold mb-1">
        ⚠️ Ditemukan tahun akademik aktif lebih dari 1, dan sistem gagal menonaktifkannya secara otomatis:
      </p>
      <ul class="list-disc list-inside mb-1">
        <li v-for="(nama, idx) in duplicateActiveWarning" :key="idx">{{ nama }}</li>
      </ul>
      <p>
        Mohon hubungi tim backend untuk menonaktifkan data di atas secara manual,
        karena seharusnya hanya 1 tahun akademik yang boleh aktif dalam satu waktu.
      </p>
    </div>

    <!-- Card -->
<div class="col-span-3 bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

  <!-- HEADER BIRU -->
  <div class="bg-[#243e90] px-5 py-4">
    <h2 class="text-white text-2xl font-bold">
      Tahun Akademik
    </h2>
    <p class="text-white text-sm mt-1">
      Kelola data tahun akademik
    </p>
  </div>

  <!-- ISI -->
  <!-- ISI -->
<div class="p-5">

  <!-- Tombol Tambah -->
  <div class="flex justify-end mb-4">
    <button
      @click="goToTambahTahunAkademik"
      class="bg-[#1f3c93] hover:bg-blue-800 text-white px-4 py-2 rounded-lg text-sm font-medium transition"
    >
      + Tambah
    </button>
  </div>

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
            <tr class="text-left text-gray-600 border-b border-gray-300">

              <th class="text-left py-3 px-4 font-semibold text-slate-600 w-12">
                No
              </th>

              <!-- <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Semester
              </th> -->

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Tanggal Awal
              </th>

              <th class="text-left py-3 px-4 font-semibold text-slate-600">
                Tanggal Akhir
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

              <!-- <td class="py-4 px-4 font-medium text-slate-700">
                {{ item.tipeSemester }}
              </td> -->

              <td class="py-4 px-4">
                {{ formatTanggal(item.awalTanggal, item.awalBulan, item.awalTahun) }}
              </td>

              <td class="py-4 px-4">
                {{ formatTanggal(item.akhirTanggal, item.akhirBulan, item.akhirTahun) }}
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
        </div>
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
            <select v-model="editForm.tipeSemester"
              class="w-full border border-slate-200 rounded-lg px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-400">
              <option value="" disabled>Pilih Semester</option>
              <option>Ganjil</option>
              <option>Genap</option>
            </select>
            <p v-if="isDuplicateData" class="text-red-500 text-xs mt-1">
              Kombinasi semester dan tanggal awal sudah digunakan.
            </p>
          </div>

          <!-- Tanggal Awal -->
          <div>
            <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
              Tanggal Awal
            </label>

            <div class="flex gap-2">
              <!-- Tanggal -->
              <input :value="editForm.awalTanggal" @input="onAwalTanggalInput" type="text" inputmode="numeric"
                maxlength="2" placeholder="Tgl"
                class="w-1/3 border rounded-lg px-3 py-2.5 text-sm text-center outline-none focus:ring-2 focus:ring-blue-400"
                :class="isDuplicateData ? 'border-red-500' : 'border-slate-200'" />

              <!-- Bulan -->
              <input :value="editForm.awalBulan" @input="onAwalBulanInput" type="text" inputmode="numeric"
                maxlength="2" placeholder="Bln"
                class="w-1/3 border rounded-lg px-3 py-2.5 text-sm text-center outline-none focus:ring-2 focus:ring-blue-400"
                :class="isDuplicateData ? 'border-red-500' : 'border-slate-200'" />

              <!-- Tahun -->
              <input :value="editForm.awalTahun" @input="onAwalTahunInput" type="text" inputmode="numeric"
                maxlength="4" placeholder="Thn"
                class="w-1/3 border rounded-lg px-3 py-2.5 text-sm text-center outline-none focus:ring-2 focus:ring-blue-400"
                :class="(isDuplicateData || isInvalidYear) ? 'border-red-500' : 'border-slate-200'" />
            </div>

            <!-- Hint max tanggal untuk bulan terpilih -->
            <p v-if="editForm.awalBulan" class="text-slate-400 text-xs mt-1">
              Bulan {{ editForm.awalBulan }} maksimal tanggal {{ maxTanggalAwal }}.
            </p>

            <!-- Error tahun -->
            <p v-if="isInvalidYear" class="text-red-500 text-xs mt-1">
              Tahun harus terdiri dari 4 digit.
            </p>

            <!-- Error tanggal & bulan harus berpasangan -->
            <p v-else-if="isInvalidDate" class="text-red-500 text-xs mt-1">
              Tanggal dan bulan harus diisi bersamaan.
            </p>

            <!-- Error duplikat -->
            <p v-else-if="isDuplicateData" class="text-red-500 text-xs mt-1">
              Semester dan tanggal awal tersebut sudah terdaftar.
            </p>
          </div>

          <!-- Tanggal Akhir (otomatis, disabled) -->
          <div>
            <label class="text-xs font-semibold text-slate-500 mb-1.5 block">
              Tanggal Akhir (otomatis)
            </label>

            <div class="flex gap-2">
              <input :value="editForm.akhirTanggal || '-'" type="text" disabled
                class="w-1/3 border border-slate-200 bg-slate-100 rounded-lg px-3 py-2.5 text-sm text-center text-slate-400 cursor-not-allowed" />

              <input :value="editForm.akhirBulan || '-'" type="text" disabled
                class="w-1/3 border border-slate-200 bg-slate-100 rounded-lg px-3 py-2.5 text-sm text-center text-slate-400 cursor-not-allowed" />

              <input :value="editForm.akhirTahun || '-'" type="text" disabled
                class="w-1/3 border border-slate-200 bg-slate-100 rounded-lg px-3 py-2.5 text-sm text-center text-slate-400 cursor-not-allowed" />
            </div>

            <p class="text-slate-400 text-xs mt-1">
              Tahun akhir otomatis mengikuti tanggal & bulan awal, tahun ditambah 1.
            </p>
          </div>

        </div>

        <!-- Tombol -->
        <div class="flex justify-end gap-3 mt-6">

          <button @click="showModal = false"
            class="px-4 py-2 text-sm text-slate-600 border border-slate-200 hover:bg-slate-50 rounded-lg transition-colors">
            Batal
          </button>

          <button @click="saveEdit" :disabled="isDuplicateData || isInvalidYear || isInvalidDate || isFormIncomplete"
            class="px-4 py-2 text-sm text-white font-semibold rounded-lg transition-colors" :class="(isDuplicateData || isInvalidYear || isInvalidDate || isFormIncomplete)
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