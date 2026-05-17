<script setup lang="ts">
import { ref } from 'vue'

const tahunAkademik = ref('')
const jurusan = ref('')
const prodi = ref('')

const namaKelas = ref('')
const dosen = ref('')
const kurikulum = ref('')
const semester = ref('')

const simpanData = () => {
  const data = {
    tahunAkademik: tahunAkademik.value,
    jurusan: jurusan.value,
    prodi: prodi.value,
    namaKelas: namaKelas.value,
    dosen: dosen.value,
    kurikulum: kurikulum.value,
    semester: semester.value
  }

  console.log('Data Kelas:', data)
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- Breadcrumb -->
    <p class="text-sm text-gray-500 mb-2">
      Mahasiswa > Kelas > Tambah Kelas
    </p>

    <!-- Title -->
    <h1 class="text-4xl font-bold text-gray-800">
      Detail Kelas
    </h1>

    <p class="text-gray-500 mt-1 mb-6">
      Pengelolaan Data
    </p>

    <!-- FORM AKADEMIK -->
    <div class="bg-white border border-blue-200 rounded-xl p-5 shadow-sm mb-6">
      
      <h2 class="text-2xl font-semibold text-gray-700 mb-6">
        Form Akademik
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

        <!-- Tahun Akademik -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Tahun Akademik
          </label>

          <select
            v-model="tahunAkademik"
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">Pilih Tahun Akademik</option>
            <option value="2023-2024">2023-2024</option>
            <option value="2024-2025">2024-2025</option>
          </select>
        </div>

        <!-- Jurusan -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Jurusan
          </label>

          <select
            v-model="jurusan"
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">Pilih Jurusan</option>
            <option value="Teknik Informatika">
              Teknik Informatika
            </option>
            <option value="Sistem Informasi">
              Sistem Informasi
            </option>
          </select>
        </div>

        <!-- Prodi -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Prodi
          </label>

          <select
            v-model="prodi"
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">Pilih Prodi</option>
            <option value="TI">Teknik Informatika</option>
            <option value="SIKC">SIKC</option>
          </select>
        </div>

      </div>
    </div>

    <!-- FORM KELAS -->
    <div class="bg-white border border-blue-200 rounded-xl p-5 shadow-sm">

      <h2 class="text-2xl font-semibold text-gray-700 mb-6">
        Form Kelas
      </h2>

      <!-- Nama Kelas -->
      <div class="mb-6">
        <label class="block text-gray-700 mb-2 font-medium">
          Nama Kelas
        </label>

        <input
          v-model="namaKelas"
          type="text"
          placeholder="Isi Nama Kelas..."
          class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
        />
      </div>

      <!-- Dosen -->
      <div class="mb-6 max-w-xl">
        <label class="block text-gray-700 mb-2 font-medium">
          Dosen
        </label>

        <select
          v-model="dosen"
          class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
        >
          <option value="">Pilih Dosen Pengampu</option>
          <option value="Dosen 1">Dosen 1</option>
          <option value="Dosen 2">Dosen 2</option>
        </select>
      </div>

      <!-- Kurikulum & Semester -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">

        <!-- Kurikulum -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Kurikulum
          </label>

          <select
            v-model="kurikulum"
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">Pilih Kurikulum</option>
            <option value="Merdeka">Merdeka</option>
            <option value="Kampus Merdeka">Kampus Merdeka</option>
          </select>
        </div>

        <!-- Semester -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Semester
          </label>

          <select
            v-model="semester"
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          >
            <option value="">Pilih Semester</option>
            <option value="Semester 1">Semester 1</option>
            <option value="Semester 2">Semester 2</option>
          </select>
        </div>

      </div>
    </div>

    <!-- Tombol Simpan -->
    <button
      @click="simpanData"
      class="mt-6 bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg flex items-center gap-2 transition"
    >
      💾 Simpan
    </button>

  </div>
</template>