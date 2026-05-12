# Kurikulum.vue

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const jurusan = ref('')
const prodi = ref('')
const tahunAkademik = ref('')

const dataKurikulum = ref([
  {
    id: 1,
    nama: 'Merdeka',
    prodi: 'Teknik Informatika',
    semester: 'Semester 1',
    tahun: '2023-2024'
  },
  {
    id: 2,
    nama: 'Merdeka',
    prodi: 'SIKC',
    semester: 'Semester 2',
    tahun: '2023-2024'
  }
])

const currentPage = ref(1)
const perPage = ref(5)

const totalPages = computed(() => {
  return Math.ceil(dataKurikulum.value.length / perPage.value)
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  const end = start + perPage.value
  return dataKurikulum.value.slice(start, end)
})

const goToTambahAkademik = () => {
  router.push('/Tambah_akademik')
}

const editData = (id: number) => {
  console.log('Edit data:', id)
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">
    
    <!-- Breadcrumb -->
    <p class="text-sm text-gray-500 mb-2">
      Mahasiswa > Kurikulum
    </p>

    <!-- Title -->
    <h1 class="text-3xl font-bold text-gray-700">
      Kurikulum
    </h1>

    <p class="text-gray-500 mt-1 mb-6">
      Pengelolaan Data
    </p>

    <!-- Card -->
    <div class="bg-white rounded-2xl border border-blue-200 shadow-sm p-5">
      
      <!-- Header -->
      <h2 class="text-2xl font-semibold text-gray-700 mb-6">
        Data Kurikulum
      </h2>

      <!-- Filter -->
      <div class="flex flex-wrap gap-4 items-center mb-6">
        
        <!-- Jurusan -->
        <select
          v-model="jurusan"
          class="border-b-2 border-blue-300 outline-none px-3 py-2 min-w-[150px] text-gray-600"
        >
          <option value="">Jurusan</option>
          <option value="TI">Teknik Informatika</option>
          <option value="SI">Sistem Informasi</option>
        </select>

        <!-- Prodi -->
        <select
          v-model="prodi"
          class="border-b-2 border-blue-300 outline-none px-3 py-2 min-w-[150px] text-gray-600"
        >
          <option value="">Prodi</option>
          <option value="TI">Teknik Informatika</option>
          <option value="SIKC">SIKC</option>
        </select>

        <!-- Tahun Akademik -->
        <select
          v-model="tahunAkademik"
          class="border-b-2 border-blue-300 outline-none px-3 py-2 min-w-[170px] text-gray-600"
        >
          <option value="">Tahun Akademik</option>
          <option value="2023-2024">2023-2024</option>
          <option value="2024-2025">2024-2025</option>
        </select>

        <!-- Tombol -->
        <div class="flex gap-3 ml-auto">
          
          <button
            class="bg-blue-800 hover:bg-blue-900 text-white px-5 py-2 rounded-lg flex items-center gap-2 transition"
          >
            <span>⇅</span>
            Pilih
          </button>

          <button
            @click="goToTambahAkademik"
            class="bg-blue-800 hover:bg-blue-900 text-white px-5 py-2 rounded-lg flex items-center gap-2 transition"
          >
            <span>+</span>
            Tambah Kelas
          </button>

        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left border-separate border-spacing-y-4">
          <thead>
            <tr class="text-gray-500 text-sm">
              <th>No</th>
              <th>Nama Kurikulum</th>
              <th>Prodi</th>
              <th>Semester</th>
              <th>Tahun Akademik</th>
              <th class="text-center">Aksi</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="text-gray-700 font-medium"
            >
              <td>{{ index + 1 }}</td>
              <td>{{ item.nama }}</td>
              <td>{{ item.prodi }}</td>
              <td>{{ item.semester }}</td>
              <td>{{ item.tahun }}</td>

              <td class="text-center">
                <button
                  @click="editData(item.id)"
                  class="bg-orange-400 hover:bg-orange-500 text-white px-4 py-1 rounded-lg transition"
                >
                  ✏ Edit
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Footer -->
      <div class="flex flex-wrap items-center justify-between mt-16">
        
        <!-- Per Page -->
        <select
          v-model="perPage"
          class="border rounded-lg px-3 py-2 text-sm text-gray-600"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="20">20 Baris</option>
        </select>

        <!-- Pagination -->
        <div class="flex items-center gap-3 text-sm text-gray-500">
          
          <button class="hover:text-blue-700">
            ← Previous
          </button>

          <button
            class="w-8 h-8 rounded-lg bg-blue-800 text-white"
          >
            1
          </button>

          <button class="hover:text-blue-700">2</button>
          <button class="hover:text-blue-700">3</button>

          <span>...</span>

          <button class="hover:text-blue-700">67</button>
          <button class="hover:text-blue-700">68</button>

          <button class="hover:text-blue-700">
            Next →
          </button>

        </div>
      </div>
    </div>
  </div>
</template>
```

---

# Router (`src/router/index.ts`)

```ts
import { createRouter, createWebHistory } from 'vue-router'
import Kurikulum from '@/views/Kurikulum.vue'
import Tambah_akademik from '@/views/Tambah_akademik.vue'

const routes = [
  {
    path: '/',
    name: 'Kurikulum',
    component: Kurikulum
  },
  {
    path: '/Tambah_akademik',
    name: 'Tambah_akademik',
    component: Tambah_akademik
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

---

# Tambah_akademik.vue

```vue
<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">
    <div class="bg-white rounded-2xl shadow-sm p-6">
      <h1 class="text-3xl font-bold text-gray-700 mb-4">
        Tambah Akademik
      </h1>

      <p class="text-gray-500">
        Halaman tambah data akademik.
      </p>
    </div>
  </div>
</template>
```

---

# Install Tailwind CSS

```bash
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

---

# tailwind.config.js

```js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

---

# src/style.css

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```
