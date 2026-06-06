<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { watch } from "vue"

// ================= INTERFACE =================
interface TahunAkademik {
  id: number | string
  tipee_semester?: string
  tahun_awal?: string
  tahun_akhir?: string
  status?: string
}

interface MataKuliah {
  nama: string
  sks: number | null
}

// ================= STATE =================
const namaKurikulum = ref("")
const tahunAkademikId = ref("")
const jurusan = ref("")
const prodi = ref("")
const kodeKurikulum = ref("")

const jurusanList = ref<any[]>([])
const prodiList = ref<any[]>([])
const tahunAkademikList = ref<TahunAkademik[]>([])

const mataKuliahList = ref<MataKuliah[]>([
  { nama: "", sks: null }
])

const filteredProdiList = computed(() => {
  if (!jurusan.value) return []

  return prodiList.value.filter(
    (p) => String(p.jurusan?.id) === String(jurusan.value)
  )
})

watch(jurusan, () => {
  prodi.value = ""
})

// ================= API =================
const getTahunAkademik = async () => {
  try {
    const res = await fetch("/api/tahun-akademik", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`
      }
    })

    const json = await res.json()
    console.log("TAHUN AKADEMIK RAW:", json)

    const data =
      json?.data?.data ??
      json?.data ??
      json?.result ??
      []

    tahunAkademikList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error("GET TAHUN AKADEMIK ERROR:", err)
  }
}

const getJurusan = async () => {
  try {
    const res = await fetch("/api/jurusan/", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })

    const json = await res.json()
    console.log("JURUSAN RAW:", json)

    const data = json?.data ?? json?.result ?? json ?? []
    jurusanList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error("GET JURUSAN ERROR:", err)
  }
}

const getProdi = async () => {
  try {
    const res = await fetch("/api/prodi/", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })

    const json = await res.json()
    console.log("PRODI RAW:", json)

    const data = json?.data ?? json?.result ?? json ?? []
    prodiList.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error("GET PRODI ERROR:", err)
  }
}

// ================= MOUNTED — hanya sekali =================
onMounted(() => {
  getTahunAkademik()
  getJurusan()
  getProdi()
})

// ================= TAMBAH FORM MK =================
const tambahMataKuliah = () => {
  mataKuliahList.value.push({ nama: "", sks: null })
}

// ================= HAPUS FORM MK =================
const hapusMataKuliah = (index: number) => {
  mataKuliahList.value.splice(index, 1)
}

// ================= SIMPAN =================
const simpanKurikulum = async () => {
  try {
    const token = localStorage.getItem("token")

    if (!token) {
      alert("Token tidak ditemukan")
      return
    }

    if (!namaKurikulum.value) {
      alert("Nama kurikulum wajib diisi")
      return
    }

    // Cari objek jurusan & prodi berdasarkan id yang dipilih
    const selectedJurusan = jurusanList.value.find(
      (j) => String(j.id) === String(jurusan.value)
    )

    const selectedProdi = prodiList.value.find(
      (p) => String(p.id) === String(prodi.value)
    )

    // Debug log — semua variabel sudah tersedia di sini
    console.log("KODE:", kodeKurikulum.value)
    console.log("NAMA:", namaKurikulum.value)
    console.log("JURUSAN:", jurusan.value)
    console.log("JURUSAN OBJECT:", selectedJurusan)
    console.log("TA:", tahunAkademikId.value)
    console.log("PRODI VALUE:", prodi.value)
    console.log("PRODI LIST:", prodiList.value)
    console.log("PRODI OBJECT:", selectedProdi)

    const prodiName = selectedProdi?.name || ""

    if (!prodiName) {
      alert("Prodi tidak valid / tidak ditemukan")
      return
    }
    const BASE_URL = 'https://be.karlearn.site'
    const responseKurikulum = await fetch(`${BASE_URL}/api/kurikulum/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        kode: kodeKurikulum.value,
        name: namaKurikulum.value,
        tahun_akademik_id: Number(tahunAkademikId.value),
        jurusan_id: Number(jurusan.value),
        ProdiName: prodiName
      })
    })

    const result = await responseKurikulum.json()
    console.log("RESP KURIKULUM:", result)

    if (!responseKurikulum.ok) {
      alert("Gagal simpan kurikulum")
      return
    }

    alert("Berhasil simpan kurikulum!")

  } catch (error) {
    console.error("SIMPAN ERROR:", error)
    alert("Terjadi error saat simpan")
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Mahasiswa > Kurikulum > Tambah Kurikulum
    </div>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-gray-800">
      Tambah Kurikulum
    </h1>

    <p class="mb-6 mt-1 text-gray-500">
      Pengelolaan Data
    </p>

    <!-- FORM KURIKULUM -->
    <div class="mb-5 rounded-2xl border border-blue-100 bg-white shadow-sm">

      <!-- HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">
          Form Kurikulum
        </h2>
      </div>

      <!-- CONTENT -->
      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-2">

        <!-- NAMA -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Nama Kurikulum
          </label>
          <input v-model="namaKurikulum" type="text" placeholder="Isi Nama Kurikulum ..."
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500" />
        </div>

        <!-- TAHUN AKADEMIK -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Tahun Akademik
          </label>
          <select v-model="tahunAkademikId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500">
            <option value="">Pilih Tahun Akademik</option>
            <option v-for="item in tahunAkademikList" :key="item.id" :value="item.id">
              {{ item.tahun_awal?.split('-')?.[0] || '-' }}/{{ item.tahun_akhir?.split('-')?.[0] || '-' }} {{
                item.tipee_semester }}
            </option>
          </select>
        </div>

        <!-- JURUSAN -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Jurusan
          </label>
          <select v-model="jurusan"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500">
            <option value="">Pilih Jurusan</option>
            <option v-for="j in jurusanList" :key="j.id" :value="j.id">
              {{ j.name }}
            </option>
          </select>
        </div>

        <!-- KODE KURIKULUM -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Kode Kurikulum
          </label>
          <input v-model="kodeKurikulum" type="text"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
            placeholder="Masukkan kode kurikulum" />
        </div>

        <!-- PRODI -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Prodi
          </label>
          <select v-model="prodi"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500">
            <option value="">Pilih Prodi</option>

            <option v-for="p in filteredProdiList" :key="p.id" :value="p.id">
              {{ p.name }}
            </option>
          </select>
        </div>

      </div>

    </div>

    <!-- FORM MATA KULIAH -->
    <div class="rounded-2xl border border-blue-100 bg-white shadow-sm">

      <!-- HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">
          Form Matakuliah
        </h2>
      </div>

      <!-- CONTENT -->
      <div class="p-5">

        <div v-for="(mk, index) in mataKuliahList" :key="index" class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-2">

          <!-- NAMA MK -->
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">
              Nama Matakuliah
            </label>
            <input v-model="mk.nama" type="text" placeholder="Isi Nama Matakuliah ..."
              class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500" />
          </div>

          <!-- SKS -->
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">
              SKS
            </label>
            <input v-model="mk.sks" type="number" placeholder="Isi SKS ..."
              class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500" />
          </div>

          <!-- BUTTON HAPUS -->
          <div class="md:col-span-2 flex justify-end">
            <button v-if="mataKuliahList.length > 1" @click="hapusMataKuliah(index)"
              class="rounded-lg bg-red-500 px-4 py-2 text-sm font-medium text-white hover:bg-red-600">
              Hapus
            </button>
          </div>

        </div>

        <!-- BUTTON TAMBAH -->
        <div class="flex justify-center">
          <button @click="tambahMataKuliah"
            class="flex h-12 w-12 items-center justify-center rounded-full border-2 border-gray-500 text-2xl text-gray-600 transition hover:bg-gray-100">
            +
          </button>
        </div>

      </div>

    </div>

    <!-- BUTTON SIMPAN -->
    <div class="mt-6">
      <button @click="simpanKurikulum"
        class="rounded-xl bg-green-500 px-6 py-3 text-sm font-semibold text-white shadow hover:bg-green-600">
        💾 Simpan
      </button>
    </div>

  </div>
</template>