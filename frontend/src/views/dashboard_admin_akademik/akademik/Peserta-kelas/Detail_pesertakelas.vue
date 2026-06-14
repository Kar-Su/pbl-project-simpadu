<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

// ─────────────────────────────────────────────
// HELPER HEADER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ─────────────────────────────────────────────
// INTERFACE (ID = UUID string, sesuai response API)
// ─────────────────────────────────────────────
interface KelasItem {
  id: string
  nama_kelas: string
}

interface MahasiswaItem {
  id: string
  nama: string
  nim?: string
  email?: string
}

// ─────────────────────────────────────────────
// STATE
// ─────────────────────────────────────────────
const kelasList = ref<KelasItem[]>([])
const mahasiswaList = ref<MahasiswaItem[]>([])

const selectedKelas = ref<string>("")
const selectedMahasiswaList = ref<string[]>([])

const loading = ref(false)
const saving = ref(false)
const error = ref("")
const successMsg = ref("")

const searchMahasiswa = ref("")

// ─────────────────────────────────────────────
// MODE: dari route param (Detail) atau tanpa param (Tambah)
// ─────────────────────────────────────────────
const kelasIdFromRoute = computed(
  () => (route.params.id as string) || ""
)

const isDetailMode = computed(
  () => !!kelasIdFromRoute.value
)

// ─────────────────────────────────────────────
// HELPER: ambil nama tampilan (fallback nama/name)
// ─────────────────────────────────────────────
const displayName = (m: any): string => m?.nama ?? m?.name ?? "-"
const displayNim = (m: any): string => m?.nim ?? m?.email ?? "-"

// ─────────────────────────────────────────────
// HIT API KELAS
// Endpoint : GET /api/kelas
// ─────────────────────────────────────────────
const getKelas = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/kelas`, {
      headers: getHeaders(),
    })

    if (!res.ok) {
      console.error("getKelas HTTP error:", res.status)
      kelasList.value = []
      return
    }

    const data = await res.json()
    kelasList.value = data.data ?? []
  } catch (err) {
    console.error("getKelas:", err)
    kelasList.value = []
  }
}

// ─────────────────────────────────────────────
// HIT API MAHASISWA
// Endpoint : GET /api/mahasiswa
// ─────────────────────────────────────────────
const getMahasiswa = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/mahasiswa`, {
      headers: getHeaders(),
    })

    if (!res.ok) {
      console.error("getMahasiswa HTTP error:", res.status)
      mahasiswaList.value = []
      return
    }

    const data = await res.json()
    mahasiswaList.value = data.data ?? []
  } catch (err) {
    console.error("getMahasiswa:", err)
    mahasiswaList.value = []
  }
}

// ─────────────────────────────────────────────
// HIT API: ambil mahasiswa yang sudah terdaftar
// di kelas ini, supaya checkbox tercentang otomatis
// Endpoint : GET /api/kelas/{id}/mahasiswa
// ─────────────────────────────────────────────
const getPesertaByKelas = async (kelasId: string): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/kelas/${kelasId}/mahasiswa`, {
      headers: getHeaders(),
    })

    if (!res.ok) {
      console.error("getPesertaByKelas HTTP error:", res.status)
      return
    }

    const data = await res.json()
    console.log("getPesertaByKelas response:", data)

    const raw = data?.data ?? []
    const ids: string[] = []

    // backend bisa mengembalikan beberapa bentuk struktur,
    // jadi ditangani secara fleksibel:
    if (Array.isArray(raw)) {
      raw.forEach((entry: any) => {
        // bentuk 1: { mahasiswa_id, name, email }
        if (entry?.mahasiswa_id) {
          ids.push(entry.mahasiswa_id)
        }

        // bentuk 2: { mahasiswa: [{ mahasiswa_id, name, email }], kelas: {...} }
        if (Array.isArray(entry?.mahasiswa)) {
          entry.mahasiswa.forEach((m: any) => {
            if (m?.mahasiswa_id) ids.push(m.mahasiswa_id)
            else if (m?.id) ids.push(m.id)
          })
        }

        // bentuk 3: { id, nama, nim } langsung
        if (entry?.id && !entry?.mahasiswa_id && !entry?.mahasiswa) {
          ids.push(entry.id)
        }
      })
    }

    selectedMahasiswaList.value = [...new Set(ids)]
  } catch (err) {
    console.error("getPesertaByKelas:", err)
  }
}

// ─────────────────────────────────────────────
// SUBMIT
// Endpoint : POST /api/peserta-kelas
// body:
// {
//   kelas_id: string,
//   mahasiswa_ids: string[]
// }
// ─────────────────────────────────────────────
const handleSubmit = async (): Promise<void> => {
  error.value = ""
  successMsg.value = ""

  if (!selectedKelas.value) {
    error.value = "Silakan pilih kelas terlebih dahulu."
    return
  }

  if (selectedMahasiswaList.value.length === 0) {
    error.value = "Pilih minimal satu mahasiswa."
    return
  }

  saving.value = true

  try {
    const payload = {
      kelas_id: selectedKelas.value,
      mahasiswa_ids: selectedMahasiswaList.value,
    }

    const res = await fetch(`${BASE_URL}/api/peserta-kelas`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })

    const data = await res.json()
    console.log("handleSubmit response:", data)

    if (!res.ok || data?.success === false) {
      error.value = data?.message ?? "Gagal menyimpan data peserta kelas."
      return
    }

    successMsg.value = "Data peserta kelas berhasil disimpan."

    // balik ke halaman list setelah berhasil
    setTimeout(() => {
      router.push("/dashboard-admin/peserta_kelas")
    }, 1000)
  } catch (err) {
    console.error("handleSubmit:", err)
    error.value = "Terjadi kesalahan saat menyimpan data."
  } finally {
    saving.value = false
  }
}

// ─────────────────────────────────────────────
// CHECKBOX
// ─────────────────────────────────────────────
const toggleMahasiswa = (id: string): void => {
  if (selectedMahasiswaList.value.includes(id)) {
    selectedMahasiswaList.value =
      selectedMahasiswaList.value.filter((item) => item !== id)
  } else {
    selectedMahasiswaList.value.push(id)
  }
}

// ─────────────────────────────────────────────
// SEARCH / FILTER MAHASISWA (client-side)
// ─────────────────────────────────────────────
const filteredMahasiswa = computed(() => {
  if (!searchMahasiswa.value) return mahasiswaList.value

  const keyword = searchMahasiswa.value.toLowerCase()

  return mahasiswaList.value.filter(
    (m) =>
      displayName(m).toLowerCase().includes(keyword) ||
      displayNim(m).toLowerCase().includes(keyword)
  )
})

// ─────────────────────────────────────────────
// NAMA KELAS TERPILIH (untuk ditampilkan di header)
// ─────────────────────────────────────────────
const selectedKelasName = computed(() => {
  const found = kelasList.value.find(
    (k) => k.id === selectedKelas.value
  )

  return found?.nama_kelas ?? ""
})

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted(async (): Promise<void> => {
  loading.value = true

  await Promise.all([getKelas(), getMahasiswa()])

  // jika datang dari tombol "Detail" -> ada id kelas di route
  if (isDetailMode.value) {
    selectedKelas.value = kelasIdFromRoute.value
    await getPesertaByKelas(kelasIdFromRoute.value)
  }

  loading.value = false
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
      <span class="cursor-pointer hover:underline" @click="router.push('/dashboard-admin/peserta_kelas')">
        Mahasiswa
      </span>
      <span>›</span>
      <span>Kelas</span>
      <span>›</span>
      <span class="text-gray-700">
        {{ isDetailMode ? "Kelola Peserta Kelas" : "Tambah Peserta Kelas" }}
      </span>
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold leading-none text-[#333]">
      Detail Peserta Kelas
    </h1>

    <p class="mt-3 text-gray-500">
      <span v-if="isDetailMode && selectedKelasName">
        Pengelolaan peserta untuk kelas <span class="font-semibold">{{ selectedKelasName }}</span>
      </span>
      <span v-else>
        Pengelolaan Data peserta kelas
      </span>
    </p>

    <!-- LOADING -->
    <div v-if="loading" class="mt-8 text-center text-gray-500">
      Memuat data...
    </div>

    <!-- CARD -->
    <div
      v-else
      class="mt-8 w-full max-w-[980px] rounded-2xl border border-[#d8e1f0] bg-white p-5 shadow-sm"
    >

      <!-- HEADER -->
      <h2 class="mb-6 text-[30px] font-bold text-[#444]">
        Form Kelas
      </h2>

      <!-- ALERT -->
      <div v-if="error" class="mb-4 rounded-xl bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-600">
        {{ error }}
      </div>

      <div v-if="successMsg" class="mb-4 rounded-xl bg-green-50 border border-green-200 px-4 py-3 text-sm text-green-600">
        {{ successMsg }}
      </div>

      <!-- FORM -->
      <div class="space-y-5">

        <!-- KELAS -->
        <div>
          <label class="mb-2 block text-[15px] font-medium text-[#555]">
            Kelas
          </label>

          <div class="relative w-[460px]">
            <select
              v-model="selectedKelas"
              :disabled="isDetailMode"
              class="h-[54px] w-full appearance-none rounded-xl border border-[#bfc8d8] bg-white px-4 pr-10 text-[15px] text-gray-700 outline-none focus:border-[#2447a8] disabled:bg-gray-100 disabled:text-gray-500"
            >
              <option value="">Pilih Kelas</option>

              <option
                v-for="item in kelasList"
                :key="item.id"
                :value="item.id"
              >
                {{ item.nama_kelas }}
              </option>
            </select>

            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="2"
              stroke="currentColor"
              class="pointer-events-none absolute right-4 top-1/2 size-5 -translate-y-1/2 text-gray-500"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="m19.5 8.25-7.5 7.5-7.5-7.5"
              />
            </svg>
          </div>

          <p v-if="isDetailMode" class="mt-1 text-xs text-gray-400">
            Kelas tidak bisa diubah pada mode kelola peserta.
          </p>
        </div>

        <!-- SEARCH MAHASISWA -->
        <div>
          <label class="mb-2 block text-[15px] font-medium text-[#555]">
            Cari Mahasiswa
          </label>

          <input
            v-model="searchMahasiswa"
            type="text"
            placeholder="Cari nama atau NIM..."
            class="h-[54px] w-[460px] rounded-xl border border-[#bfc8d8] bg-white px-4 text-[15px] text-gray-700 outline-none focus:border-[#2447a8]"
          />
        </div>

        <!-- CHECKBOX LIST MAHASISWA -->
        <div>
          <label class="mb-2 block text-[15px] font-medium text-[#555]">
            Mahasiswa ({{ selectedMahasiswaList.length }} dipilih)
          </label>

          <div class="max-h-[420px] overflow-y-auto space-y-3 rounded-xl border border-[#e5e9f2] p-4">

            <div v-if="filteredMahasiswa.length === 0" class="text-center text-gray-400 py-6">
              Tidak ada mahasiswa ditemukan.
            </div>

            <div
              v-for="item in filteredMahasiswa"
              :key="item.id"
              class="flex items-center gap-4"
            >
              <input
                type="checkbox"
                :checked="selectedMahasiswaList.includes(item.id)"
                @change="toggleMahasiswa(item.id)"
                class="h-6 w-6 rounded border-gray-400 text-[#2447a8] focus:ring-[#2447a8]"
              />

              <p class="text-[17px] font-medium text-[#444]">
                {{ displayName(item) }} / {{ displayNim(item) }}
              </p>
            </div>

          </div>
        </div>

        <!-- BUTTON -->
        <div class="pt-6 flex items-center gap-3">
          <button
            @click="handleSubmit"
            :disabled="saving"
            class="rounded-xl bg-[#22c55e] px-8 py-3 text-[16px] font-semibold text-white transition hover:bg-[#16a34a] disabled:opacity-60"
          >
            {{ saving ? "Menyimpan..." : "Simpan" }}
          </button>

          <button
            @click="router.push('/dashboard-admin/peserta_kelas')"
            type="button"
            class="rounded-xl border border-gray-300 px-8 py-3 text-[16px] font-semibold text-gray-600 transition hover:bg-gray-50"
          >
            Batal
          </button>
        </div>

      </div>
    </div>
  </div>
</template>