<script setup lang="ts">
import { ref } from 'vue'

const namaKurikulum = ref('')
const tahunAkademik = ref('')
const jurusan = ref('')
const prodi = ref('')

const mataKuliah = ref([
  {
    nama: '',
    sks: ''
  }
])

const tambahMataKuliah = () => {
  mataKuliah.value.push({
    nama: '',
    sks: ''
  })
}

const simpanData = () => {
  const data = {
    namaKurikulum: namaKurikulum.value,
    tahunAkademik: tahunAkademik.value,
    jurusan: jurusan.value,
    prodi: prodi.value,
    mataKuliah: mataKuliah.value
  }

  console.log(data)
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">
    
    <!-- Breadcrumb -->
    <p class="text-sm text-gray-500 mb-2">
      Mahasiswa > Kurikulum > Tambah Kurikulum
    </p>

    <!-- Title -->
    <h1 class="text-4xl font-bold text-gray-700">
      Kurikulum
    </h1>

    <p class="text-gray-500 mt-1 mb-6">
      Pengelolaan Data
    </p>

    <!-- FORM KURIKULUM -->
    <div class="bg-white border border-blue-200 rounded-xl p-5 max-w-3xl shadow-sm">
      
      <h2 class="text-2xl font-semibold text-gray-700 mb-6">
        Form Kurikulum
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        <!-- Nama Kurikulum -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Nama Kurikulum
          </label>

          <input
            v-model="namaKurikulum"
            type="text"
            placeholder="Isi Nama Kurikulum ..."
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

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

    <!-- FORM MATAKULIAH -->
    <div class="bg-white border border-blue-200 rounded-xl p-5 max-w-3xl shadow-sm mt-6">
      
      <h2 class="text-2xl font-semibold text-gray-700 mb-6">
        Form Matakuliah
      </h2>

      <div
        v-for="(item, index) in mataKuliah"
        :key="index"
        class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-4"
      >
        
        <!-- Nama Matakuliah -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            Nama Matakuliah
          </label>

          <input
            v-model="item.nama"
            type="text"
            placeholder="Isi Nama Matakuliah ..."
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

        <!-- SKS -->
        <div>
          <label class="block text-gray-700 mb-2 font-medium">
            SKS
          </label>

          <input
            v-model="item.sks"
            type="number"
            placeholder="Isi SKS ..."
            class="w-full border rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>

      </div>

      <!-- Tombol Tambah -->
      <div class="flex justify-center mt-4">
        <button
          @click="tambahMataKuliah"
          class="w-12 h-12 rounded-full border-2 border-gray-400 text-2xl text-gray-600 hover:bg-gray-100 transition"
        >
          +
        </button>
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