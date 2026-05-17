<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import DeleteConfirmModal from "./konfirmasi_hapus.vue"

interface User {
  id: string
  email: string
  name: string
  role_name: string  // ← ganti dari jabatan
}

const search = ref("")
const perPage = ref(5)
const currentPage = ref(1)
const Konfirmasi_hapus = ref(false)
const selectedEmail = ref("")
const router = useRouter()
const users = ref<User[]>([])

const getUsers = async () => {
  try {
    const token = localStorage.getItem("token")
    // const roles = ['admin-akademik', 'admin-keuangan', 'admin-mahasiswa', 'admin-pegawai', 'dosen', 'dummy', 'mahasiswa']
    const roles = ['dosen', 'dummy', 'mahasiswa']
    const results = await Promise.all(
      roles.map(role =>
        fetch(`/api/user/role/${role}`, {
          headers: { "Authorization": `Bearer ${token}` }
        }).then(r => r.json()).then(d => Array.isArray(d.data) ? d.data : [])
      )
    )
    users.value = results.flat()
  } catch (err) {
    console.error("GET USERS ERROR:", err)
  }
}

const filteredUsers = computed(() =>
  users.value.filter(item =>
    item.role_name?.toLowerCase().includes(search.value.toLowerCase())
  )
)
const totalPages = computed(() => Math.ceil(filteredUsers.value.length / perPage.value))
const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredUsers.value.slice(start, start + perPage.value)
})
const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }
const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
const openDeleteModal = (email: string) => { selectedEmail.value = email; Konfirmasi_hapus.value = true }
const confirmDelete = () => { console.log("hapus:", selectedEmail.value); Konfirmasi_hapus.value = false }

onMounted(() => { getUsers() })
</script>

<template>
  <div class="min-h-screen bg-[#f5f7fb] p-6">

    <!-- Breadcrumb -->
    <div class="text-sm text-gray-400 mb-2">
      Akademik > Akun
    </div>

    <!-- Title -->
    <h1 class="text-3xl font-bold text-gray-800 mb-1">
      Akun
    </h1>

    <p class="text-gray-500 mb-6">
      Kelola Data Akun
    </p>

    <!-- Header -->
    <div class="flex items-center justify-between mb-4">

      <!-- Search -->
      <div class="relative">
        <input v-model="search" type="text" placeholder="Cari Akun..."
          class="w-64 rounded-lg border border-gray-200 bg-white py-2 pl-4 pr-10 text-sm outline-none focus:border-blue-500" />

        <svg xmlns="http://www.w3.org/2000/svg" class="absolute right-3 top-2.5 h-4 w-4 text-gray-400" fill="none"
          viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-4.35-4.35m1.85-5.15a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>

      <!-- Button -->
<button @click="router.push('/dashboard-superadmin/tambah_akun')"
  class="flex items-center gap-2 rounded-lg bg-[#2f4a8a] px-4 py-2 text-sm font-medium text-white hover:bg-[#243b73]">
  + Tambah Akun
</button>
    </div>

    <!-- Table Card -->
    <div class="rounded-2xl border border-gray-200 bg-white shadow-sm">

      <!-- Table Header -->
      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-xl font-semibold text-gray-700">
          Data Akun
        </h2>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-sm">

          <thead>
            <tr class="text-left text-gray-500">
              <th class="px-5 py-4">No</th>
              <th class="px-5 py-4">Email</th>
              <th class="px-5 py-4">Nama</th>
              <th class="px-5 py-4">Role</th>
              <th class="px-5 py-4">Aksi</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="(item, index) in paginatedUsers" :key="item.id" class="border-t border-gray-100">
              <td class="px-5 py-4">
                {{ index + 1 }}
              </td>

              <td class="px-5 py-4">
                {{ item.email }}
              </td>

              <td class="px-5 py-4">
                {{ item.name }}
              </td>

              
                <td class="px-5 py-4">{{ item.role_name }}</td>

              <td class="px-5 py-4">
                <div class="flex gap-2">

                  <button @click="router.push('/dashboard-superadmin/edit_akun')"
                    class="rounded-md bg-yellow-400 px-3 py-1 text-xs font-medium text-white hover:bg-yellow-500">
                    ✏ Edit
                  </button>

                  <!-- Delete -->
                  <button @click="openDeleteModal(item.email)"
                    class="rounded-md bg-red-500 px-3 py-1 text-xs font-medium text-white hover:bg-red-600">
                    🗑 Hapus
                  </button>

                </div>
              </td>
            </tr>
          </tbody>

        </table>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between border-t border-gray-100 px-5 py-4">

        <!-- Per Page -->
        <select v-model="perPage" class="rounded-md border border-gray-300 px-2 py-1 text-sm outline-none">
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="20">20 Baris</option>
        </select>

        <!-- Pagination -->
        <div class="flex items-center gap-3 text-sm">

          <button @click="prevPage" class="text-gray-400 hover:text-black">
            ← Previous
          </button>

          <button class="flex h-8 w-8 items-center justify-center rounded bg-[#2f4a8a] text-white">
            {{ currentPage }}
          </button>

          <button @click="nextPage" class="text-gray-400 hover:text-black">
            Next →
          </button>

        </div>
      </div>

    </div>

  </div>
  <DeleteConfirmModal v-if="Konfirmasi_hapus" :email="selectedEmail" @kembali="Konfirmasi_hapus = false"
    @yakin="confirmDelete" />
</template>