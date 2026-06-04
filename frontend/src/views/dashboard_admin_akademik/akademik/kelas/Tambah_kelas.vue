<script setup lang="ts">
import { ref, onMounted } from "vue"

// ================= INTERFACE =================
interface TahunAkademik {
  id: string
  nama?: string
  tahun?: string
}

interface Jurusan {
  id: string
  nama?: string
}

interface Dosen {
  id: string
  name: string
}

interface Kurikulum {
  id: string
  nama?: string
}

// ================= STATE =================
const tahunAkademikList = ref<TahunAkademik[]>([])
const jurusanList = ref<Jurusan[]>([])
const dosenList = ref<Dosen[]>([])
const kurikulumList = ref<Kurikulum[]>([])

// ================= FORM =================
const tahunAkademikId = ref("")
const jurusanId = ref("")
const prodi = ref("")

const namaKelas = ref("")
const dosenId = ref("")
const kurikulumId = ref("")
const semester = ref("")

// ================= GET TAHUN AKADEMIK =================
const getTahunAkademik = async () => {

  try {

    const token = localStorage.getItem("token")
    const BASE_URL = 'https://be.karlearn.site'
    const response = await fetch(`${BASE_URL}/api/tahun-akademik`, {

      method: "GET",

      headers: {
        Authorization: `Bearer ${token}`
      }

    })

    const result = await response.json()

    tahunAkademikList.value =
      Array.isArray(result.data)
        ? result.data
        : []

  } catch (error) {

    console.error("GET TAHUN AKADEMIK ERROR:", error)

  }

}

// ================= GET JURUSAN =================
const getJurusan = async () => {

  try {

    const token = localStorage.getItem("token")

    const response = await fetch("/api/jurusan", {

      method: "GET",

      headers: {
        Authorization: `Bearer ${token}`
      }

    })

    const result = await response.json()

    jurusanList.value =
      Array.isArray(result.data)
        ? result.data
        : []

  } catch (error) {

    console.error("GET JURUSAN ERROR:", error)

  }

}

// ================= GET DOSEN =================
const getDosen = async () => {

  try {

    const token = localStorage.getItem("token")

    const response = await fetch("/api/user/role/dosen", {

      method: "GET",

      headers: {
        Authorization: `Bearer ${token}`
      }

    })

    const result = await response.json()

    dosenList.value =
      Array.isArray(result.data)
        ? result.data
        : []

  } catch (error) {

    console.error("GET DOSEN ERROR:", error)

  }

}

// ================= GET KURIKULUM =================
const getKurikulum = async () => {

  try {

    const token = localStorage.getItem("token")

    const response = await fetch("/api/kurikulum", {

      method: "GET",

      headers: {
        Authorization: `Bearer ${token}`
      }

    })

    const result = await response.json()

    kurikulumList.value =
      Array.isArray(result.data)
        ? result.data
        : []

  } catch (error) {

    console.error("GET KURIKULUM ERROR:", error)

  }

}

// ================= SIMPAN =================
const simpanKelas = async () => {

  try {

    const token = localStorage.getItem("token")

    const payload = {

      tahun_akademik_id: tahunAkademikId.value,
      jurusan_id: jurusanId.value,
      prodi: prodi.value,

      nama_kelas: namaKelas.value,
      dosen_id: dosenId.value,
      kurikulum_id: kurikulumId.value,
      semester: semester.value

    }

    console.log("PAYLOAD:", payload)

    // ganti endpoint jika backend kelas sudah ada
    const response = await fetch("/api/kelas", {

      method: "POST",

      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },

      body: JSON.stringify(payload)

    })

    const result = await response.json()

    console.log("CREATE KELAS:", result)

    alert("Kelas berhasil disimpan!")

  } catch (error) {

    console.error("SIMPAN KELAS ERROR:", error)

  }

}

// ================= MOUNTED =================
onMounted(() => {

  getTahunAkademik()
  getJurusan()
  getDosen()
  getKurikulum()

})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Kelas > Tambah Kelas
    </div>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-gray-800">
      Detail Kelas
    </h1>

    <p class="mb-6 mt-1 text-gray-500">
      Pengelolaan Data
    </p>

    <!-- FORM AKADEMIK -->
    <div class="mb-5 rounded-2xl border border-blue-100 bg-white shadow-sm">

      <!-- HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">

        <h2 class="text-2xl font-semibold text-gray-700">
          Form Akademik
        </h2>

      </div>

      <!-- CONTENT -->
      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-3">

        <!-- TAHUN AKADEMIK -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Tahun Akademik
          </label>

          <select
            v-model="tahunAkademikId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >

            <option value="">
              Pilih Tahun Akademik
            </option>

            <option
              v-for="item in tahunAkademikList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.nama || item.tahun }}
            </option>

          </select>

        </div>

        <!-- JURUSAN -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Jurusan
          </label>

          <select
            v-model="jurusanId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >

            <option value="">
              Pilih Jurusan
            </option>

            <option
              v-for="item in jurusanList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.nama }}
            </option>

          </select>

        </div>

        <!-- PRODI -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Prodi
          </label>

          <input
            v-model="prodi"
            type="text"
            placeholder="Isi Prodi"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          />

        </div>

      </div>

    </div>

    <!-- FORM KELAS -->
    <div class="rounded-2xl border border-blue-100 bg-white shadow-sm">

      <!-- HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">

        <h2 class="text-2xl font-semibold text-gray-700">
          Form Kelas
        </h2>

      </div>

      <!-- CONTENT -->
      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-2">

        <!-- NAMA KELAS -->
        <div class="md:col-span-2">

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Nama Kelas
          </label>

          <input
            v-model="namaKelas"
            type="text"
            placeholder="Isi Nama Kelas..."
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          />

        </div>

        <!-- DOSEN -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Dosen
          </label>

          <select
            v-model="dosenId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >

            <option value="">
              Pilih Dosen Pengampu
            </option>

            <option
              v-for="item in dosenList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.name }}
            </option>

          </select>

        </div>

        <!-- KURIKULUM -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Kurikulum
          </label>

          <select
            v-model="kurikulumId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >

            <option value="">
              Pilih Kurikulum
            </option>

            <option
              v-for="item in kurikulumList"
              :key="item.id"
              :value="item.id"
            >
              {{ item.nama }}
            </option>

          </select>

        </div>

        <!-- SEMESTER -->
        <div>

          <label class="mb-2 block text-sm font-medium text-gray-700">
            Semester
          </label>

          <select
            v-model="semester"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >

            <option value="">
              Pilih Semester
            </option>

            <option value="1">Semester 1</option>
            <option value="2">Semester 2</option>
            <option value="3">Semester 3</option>
            <option value="4">Semester 4</option>
            <option value="5">Semester 5</option>
            <option value="6">Semester 6</option>
            <option value="7">Semester 7</option>
            <option value="8">Semester 8</option>

          </select>

        </div>

      </div>

    </div>

    <!-- BUTTON -->
    <div class="mt-5">

      <button
        @click="simpanKelas"
        class="rounded-xl bg-green-500 px-6 py-3 text-sm font-semibold text-white shadow hover:bg-green-600"
      >
        💾 Simpan
      </button>

    </div>

  </div>
</template>