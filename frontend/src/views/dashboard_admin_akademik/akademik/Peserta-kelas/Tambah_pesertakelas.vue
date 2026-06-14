<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()

const BASE_URL = "https://be.karlearn.site"
const kelasIdFromRoute = route.params.id as string

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

interface Kelas {
  id: string
  name: string
}

interface Mahasiswa {
  id: string
  nama: string
  nim?: string
  email?: string
}

interface Peserta {
  mahasiswa_id: string
  nama: string
  nim?: string
}

const kelasList = ref<Kelas[]>([])
const mahasiswaList = ref<Mahasiswa[]>([])
const pesertaKelas = ref<Peserta[]>([])

const selectedKelas = ref("")
const isDropdownMahasiswaOpen = ref(false)

const isLoading = ref(false)

// Fetching data daftar kelas untuk dropdown pertama
const getKelasData = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/kelas`, { headers: getHeaders() })
    const data = await res.json()
    const payload = data.data ?? []
    kelasList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
  } catch (err) {
    console.error(err)
  }
}

// Fetching data mahasiswa untuk dropdown list checkbox
const getMahasiswaData = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/mahasiswa`, { headers: getHeaders() })
    const data = await res.json()
    const payload = data.data ?? []
    mahasiswaList.value = Array.isArray(payload) ? payload : (payload.items ?? [])
  } catch (err) {
    console.error(err)
  }
}

const getDetailKelas = async () => {
  if (!kelasIdFromRoute) return
  try {
    isLoading.value = true
    const res = await fetch(`${BASE_URL}/api/kelas/${kelasIdFromRoute}`, { headers: getHeaders() })
    const data = await res.json()
    const kelas = data.data
    selectedKelas.value = kelas.id

    pesertaKelas.value = kelas.mahasiswa?.map((m: any) => ({
      mahasiswa_id: m.mahasiswa_id,
      nama: m.name,
      nim: m.email || m.nim, 
    })) ?? []
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

// Mengetahui apakah mahasiswa tertentu sudah dicentang/masuk list
const isChecked = (mhsId: string) => {
  return pesertaKelas.value.some((p) => String(p.mahasiswa_id) === String(mhsId))
}

// Toggle penambahan/penghapusan dari checkbox list
const toggleMahasiswa = (mhs: Mahasiswa) => {
  const index = pesertaKelas.value.findIndex((p) => String(p.mahasiswa_id) === String(mhs.id))
  
  if (index > -1) {
    pesertaKelas.value.splice(index, 1)
  } else {
    pesertaKelas.value.push({
      mahasiswa_id: mhs.id,
      nama: mhs.nama,
      nim: mhs.nim || mhs.email,
    })
  }
}

const removeMahasiswaItem = (index: number) => {
  pesertaKelas.value.splice(index, 1)
}

const handleSimpan = async () => {
  if (!selectedKelas.value) {
    alert("Silakan pilih kelas terlebih dahulu")
    return
  }
  try {
    const payload = {
      kelas_id: selectedKelas.value,
      mahasiswa_ids: pesertaKelas.value.map((item) => item.mahasiswa_id),
    }

    const res = await fetch(`${BASE_URL}/api/peserta-kelas/update-bulk`, {
      method: "PUT",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })

    const data = await res.json()
    if (res.ok) {
      alert("Berhasil menyimpan data peserta kelas")
      router.push("/dashboard-admin/peserta_kelas")
    } else {
      alert(data.message || "Gagal menyimpan")
    }
  } catch (err) {
    console.error(err)
    alert("Terjadi kesalahan sistem")
  }
}

onMounted(async () => {
  await getKelasData()
  await getMahasiswaData()
  await getDetailKelas()

  // Menutup dropdown custom saat klik di luar area dropdown
  window.addEventListener('click', (e: Event) => {
    const target = e.target as HTMLElement
    if (!target.closest('.custom-select-container')) {
      isDropdownMahasiswaOpen.value = false
    }
  })
})
</script>

<template>
  <div class="min-h-screen bg-[#f8fafc] p-6">
    <div class="text-xs text-gray-400 mb-2 font-medium">
      Mahasiswa &gt; Kelas &gt; Tambah Peserta Kelas
    </div>
    
    <h1 class="text-3xl font-bold text-gray-900 mb-1">Tambah Peserta Kelas</h1>
    <p class="text-sm text-gray-500 mb-6">Pengelolaan Data peserta kelas</p>

    <div class="overflow-hidden rounded-xl bg-white shadow-md border border-gray-100 max-w-full">
      
      <div class="bg-[#1e3a8a] px-6 py-4">
        <h2 class="text-xl font-semibold tracking-wide text-white">
          Data Peserta Kelas
        </h2>
        <p class="mt-1 text-xs text-blue-100 opacity-90">
          Penambahan mahasiswa per kelas
        </p>
      </div>

      <div class="p-6 space-y-5 bg-white">
        
        <div>
          <label class="block mb-2 text-sm font-medium text-gray-700">
            Kelas <span class="text-red-500">*</span>
          </label>
          <div class="relative w-full">
            <select
              v-model="selectedKelas"
              class="h-[44px] w-full appearance-none rounded-lg border border-gray-300 bg-white pl-4 pr-10 text-sm text-gray-700 outline-none transition focus:border-blue-500"
            >
              <option value="">Pilih Kelas</option>
              <option v-for="kelas in kelasList" :key="kelas.id" :value="kelas.id">
                {{ kelas.name }}
              </option>
            </select>
            <div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>
        </div>

        <div class="relative custom-select-container">
          <label class="block mb-2 text-sm font-medium text-gray-700">
            Mahasiswa <span class="text-red-500">*</span>
          </label>
          
          <div 
            @click="isDropdownMahasiswaOpen = !isDropdownMahasiswaOpen"
            class="h-[44px] w-full flex items-center justify-between rounded-lg border border-gray-300 bg-white pl-4 pr-3 text-sm text-gray-400 cursor-pointer select-none focus:border-blue-500"
          >
            <span>Pilih Mahasiswa</span>
            <div class="text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>

          <div 
            v-if="isDropdownMahasiswaOpen" 
            class="absolute z-50 left-0 right-0 mt-1 max-h-60 overflow-y-auto rounded-lg border border-gray-200 bg-white shadow-xl py-1 divide-y divide-gray-50"
          >
            <div 
              v-for="mhs in mahasiswaList" 
              :key="mhs.id"
              @click="toggleMahasiswa(mhs)"
              class="flex items-center gap-3 px-4 py-3 hover:bg-slate-50 cursor-pointer select-none text-sm text-gray-700"
            >
              <input 
                type="checkbox" 
                :checked="isChecked(mhs.id)"
                @click.stop="toggleMahasiswa(mhs)"
                class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span>{{ mhs.nama }} / {{ mhs.nim || mhs.email || '-' }}</span>
            </div>
            <div v-if="mahasiswaList.length === 0" class="p-3 text-center text-xs text-gray-400">
              Tidak ada data mahasiswa tersedia.
            </div>
          </div>
        </div>

        <div class="pt-2">
          <div v-if="isLoading" class="text-center py-4 text-sm text-gray-400">
            Memuat data peserta...
          </div>
          
          <div v-else class="space-y-2">
            <div
              v-for="(item, index) in pesertaKelas"
              :key="index"
              class="flex items-center justify-between bg-[#bfdbfe]/60 border border-blue-200 rounded-lg px-4 py-3 transition hover:bg-[#bfdbfe]/80"
            >
              <div class="text-sm font-medium text-slate-700">
                {{ item.nama }} <span v-if="item.nim" class="text-gray-500 font-normal">/ {{ item.nim }}</span>
              </div>

              <button
                @click="removeMahasiswaItem(index)"
                class="flex items-center justify-center rounded-lg bg-[#ef4444] p-2 text-white transition hover:bg-red-600 shadow-sm"
              >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
                  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673A2.25 2.25 0 0 1 15.916 21H8.084a2.25 2.25 0 0 1-2.245-1.327L4.772 5.79m14.456 0A48.108 48.108 0 0 0 15.75 5.25m-6.75 0a48.11 48.11 0 0 1 3.478-.459m0 0a48.11 48.11 0 0 1 3.478.459m-3.478 0V4.5a2.25 2.25 0 0 1 2.25-2.25h1.5A2.25 2.25 0 0 1 18.75 4.5v.75" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-3 pt-4 border-t border-gray-100">
          <button
            @click="handleSimpan"
            class="flex h-[40px] items-center justify-center rounded-lg bg-[#1d357d] px-6 text-sm font-medium text-white transition hover:bg-[#162961]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-4 h-4 mr-1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
            </svg>
            Simpan
          </button>

          <button
            @click="router.push('/dashboard-admin/peserta_kelas')"
            class="flex h-[40px] items-center justify-center rounded-lg bg-gray-200 px-6 text-sm font-medium text-gray-700 transition hover:bg-gray-300"
          >
            Batal
          </button>
        </div>

      </div>
    </div>
  </div>
</template>