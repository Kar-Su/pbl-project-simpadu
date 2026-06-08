<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue"

// ================= INTERFACE =================
interface TahunAkademik {
  id: number | string
  tipee_semester?: string
  tahun_awal?: string
  tahun_akhir?: string
  status?: string
}

interface MataKuliah {
  kode: string
  nama: string
  sks: number | null
}

// ================= STATE =================
const namaKurikulum = ref("")
const tahunAkademikId = ref("")
const jurusan = ref("")
const prodi = ref("")

const sksError = ref("")

const jurusanList = ref<any[]>([])
const prodiList = ref<any[]>([])
const tahunAkademikList = ref<TahunAkademik[]>([])

const mataKuliahList = ref<MataKuliah[]>([
  {
    kode: "",
    nama: "",
    sks: null,
  },
])

// ================= SKS VALIDATION =================
const handleSksInput = (mk: MataKuliah) => {
  const value = String(mk.sks ?? "")

  sksError.value = ""

  if (value === "") return

  if (value.length > 1) {
    mk.sks = Number(value.slice(0, 1))
    sksError.value = "SKS maksimal 1 digit (1-9)"
  }
}

// ================= FILTER PRODI =================
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
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })

    const json = await res.json()

    tahunAkademikList.value =
      Array.isArray(json?.data)
        ? json.data
        : json?.data?.items ?? []
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
    const BASE_URL = "https://be.karlearn.site"

    const res = await fetch(`${BASE_URL}/api/prodi`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })

    const json = await res.json()

    console.log("PRODI RAW:", json)

    prodiList.value =
      json?.data?.items ??
      json?.data ??
      []

    console.log("PRODI LIST:", prodiList.value)
  } catch (err) {
    console.error("GET PRODI ERROR:", err)
  }
}

// ================= MOUNT =================
onMounted(() => {
  getTahunAkademik()
  getJurusan()
  getProdi()
})

// ================= MK =================
const tambahMataKuliah = () => {
  mataKuliahList.value.push({
    kode: "",
    nama: "",
    sks: null,
  })
}

const hapusMataKuliah = (index: number) => {
  mataKuliahList.value.splice(index, 1)
}

// ================= SIMPAN =================
const simpanKurikulum = async () => {
  try {
    const token = localStorage.getItem("token")

    if (!token) return alert("Token tidak ditemukan")
    if (!namaKurikulum.value) return alert("Nama kurikulum wajib diisi")
    if (!tahunAkademikId.value) return alert("Tahun akademik wajib dipilih")
    if (!jurusan.value) return alert("Jurusan wajib dipilih")
    if (!prodi.value) return alert("Prodi wajib dipilih")
    if (!mataKuliahList.value) return alert("Kode kurikulum wajib diisi")

    const adaMKKosong = mataKuliahList.value.some(
      (mk) => !mk.kode || !mk.nama || mk.sks === null
    )

    if (adaMKKosong) {
      return alert("Semua kode, nama, dan SKS matakuliah wajib diisi")
    }

    const adaSKSInvalid = mataKuliahList.value.some(
      (mk) => Number(mk.sks) < 1 || Number(mk.sks) > 9
    )

    if (adaSKSInvalid) {
      return alert("SKS hanya boleh 1 digit (1-9)")
    }

    const selectedProdi = prodiList.value.find(
      (p) => String(p.id) === String(prodi.value)
    )

    console.log("PRODI OBJECT:", selectedProdi)

    // sesuaikan dengan field dari API
    const prodiName =
      selectedProdi?.name ||
      selectedProdi?.nama ||
      selectedProdi?.prodi_name ||
      selectedProdi?.nama_prodi ||
      ""

    if (!prodiName) {
      alert("Nama prodi tidak ditemukan")
      console.log("DATA PRODI:", selectedProdi)
      return
    }

    const payload = {
      kode: mataKuliahList.value,
      name: namaKurikulum.value,
      tahun_akademik_id: Number(tahunAkademikId.value),
      jurusan_id: Number(jurusan.value),
      ProdiName: prodiName,
    }

    console.log("PAYLOAD KURIKULUM:", payload)

    const response = await fetch(
      "https://be.karlearn.site/api/kurikulum",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      }
    )

    const result = await response.json()

    console.log("RESP KURIKULUM:", result)

    if (!response.ok) {
      alert(result?.message || "Gagal simpan kurikulum")
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
            <option value="" disabled>Pilih Tahun Akademik</option>
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
            <option value="" disabled>Pilih Jurusan</option>
            <option v-for="j in jurusanList" :key="j.id" :value="j.id">
              {{ j.name }}
            </option>
          </select>
        </div>

        <!-- KODE KURIKULUM -->
       

<!-- PRODI -->
<div>
  <label class="mb-2 block text-sm font-medium text-gray-700">
    Prodi
  </label>
  <select
    v-model="prodi"
    :disabled="!jurusan"
    class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
    :class="!jurusan ? 'bg-gray-100 cursor-not-allowed text-gray-400' : ''"
  >
    <option value="">{{ jurusan ? 'Pilih Prodi' : 'Pilih Jurusan terlebih dahulu' }}</option>
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

           <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Kode Mata Kuliah
          </label>
          <input v-model="mk.kode" type="text"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
            placeholder="Masukkan kode mata kuliah" />
        </div>
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
<input
  v-model="mk.sks"
  type="number"
  min="1"
  max="9"
  placeholder="Isi SKS ..."
  class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
  @input="handleSksInput(mk)"
/>

<p
  v-if="sksError"
  class="mt-1 text-sm text-red-500"
>
  {{ sksError }}
</p>

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