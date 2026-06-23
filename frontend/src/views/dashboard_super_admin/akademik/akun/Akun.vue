<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"
import DeleteConfirmModal from "./konfirmasi_hapus.vue"

// ================= TYPE =================
interface User {
  id: string
  email: string
  name: string
  role_name: string
}

// ================= CONFIG API =================
const BASE_URL = "https://be.karlearn.site"

const API = {
  users: `${BASE_URL}/api/users`,
}

// ================= STATE =================
const search = ref("")
const perPage = ref<number>(10)
const currentPage = ref(1)

const Konfirmasi_hapus = ref(false)
const selectedEmail = ref("")
const selectedUserId = ref("")

const router = useRouter()

const users = ref<User[]>([])

const allUsers = ref<User[]>([])

const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const pages: (number | string)[] = []
  pages.push(1)

  const rangeStart = Math.max(2, current - 1)
  const rangeEnd = Math.min(total - 1, current + 1)

  if (rangeStart > 2) pages.push('...')
  for (let i = rangeStart; i <= rangeEnd; i++) pages.push(i)
  if (rangeEnd < total - 1) pages.push('...')

  pages.push(total)
  return pages
})

// ================= HEADERS =================
const getHeaders = () => {
  const token = localStorage.getItem("token")
  return {
    Authorization: `Bearer ${token}`,
    accept: "application/json",
  }
}

// ================= GET ALL USERS (untuk search) =================
const allUsersLoaded = ref(false)
const getAllUsers = async () => {
  try {
    let page = 1
    let totalPagesAll = 1
    const result: User[] = []

    do {
      const response = await fetch(
        `${API.users}?page=${page}&per_page=10`,
        { headers: getHeaders() }
      )
      const data = await response.json()

      if (!response.ok) break

      const items: User[] = data.data.items || []
      result.push(...items)

      totalPagesAll = data.data.pagination?.total_pages || 1
      page++
    } while (page <= totalPagesAll)

    allUsers.value = result
    allUsersLoaded.value = true
    console.log("allUsers loaded:", allUsers.value.length)
  } catch (err) {
    console.error(err)
  }
}

// ================= GET USERS =================
const getUsers = async () => {
  try {
    await getAllUsers()
  } catch (err) {
    console.error(err)
  }
}

// ================= FILTER =================
const filteredUsers = computed(() => {
  const keyword = search.value.toLowerCase().trim()

  if (!keyword) return allUsers.value

  return allUsers.value.filter((item) =>
    item.name?.toLowerCase().includes(keyword) ||
    item.email?.toLowerCase().includes(keyword)
  )
})

// ================= PAGINATION =================
const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredUsers.value.length / perPage.value))
)

const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  const end = start + perPage.value

  return filteredUsers.value.slice(start, end)
})

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}


const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getUsers()
  }
}

// ================= DELETE =================
const openDeleteModal = (id: string, email: string) => {
  selectedUserId.value = id
  selectedEmail.value = email
  Konfirmasi_hapus.value = true
}

const confirmDelete = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/users/${selectedUserId.value}`, {
      method: "DELETE",
      headers: getHeaders(),
    })

    if (!res.ok) {
      const json = await res.json()
      console.error("DELETE ERROR:", json)
      return
    }

    Konfirmasi_hapus.value = false
    await getUsers()
    await getAllUsers()
  } catch (err) {
    console.error("DELETE NETWORK ERROR:", err)
  }
}

// ================= WATCH =================
watch(perPage, async (newVal, oldVal) => {
  if (newVal === oldVal) return
  currentPage.value = 1
  await getUsers()
}, { immediate: false })

watch(search, () => {
  currentPage.value = 1
})

// ================= INIT =================
onMounted(() => {
  getUsers()
  getAllUsers()
})
</script>

<template>
  <div class="h-full bg-[#f5f7fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-black-800">
      Akademik > Akun
    </div>

    <!-- TITLE -->
    <h1 class="mb-1 text-3xl font-bold text-black-800">
      Akun
    </h1>

    <p class="mb-6 text-black-800">
      Kelola Data Akun
    </p>

    <!-- HEADER -->
    <div class="mb-4 flex items-center justify-between">

      <!-- SEARCH -->
      <div class="relative">
        <input v-model="search" type="text" placeholder="Cari Akun..."
          class="w-64 rounded-lg border border-gray-200 bg-white py-2 pl-4 pr-10 text-sm outline-none focus:border-blue-500" />
        <svg xmlns="http://www.w3.org/2000/svg" class="absolute right-3 top-2.5 h-4 w-4 text-gray-400" fill="none"
          viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-4.35-4.35m1.85-5.15a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>

      <!-- BUTTON TAMBAH -->
      <button @click="router.push('/dashboard-superadmin/tambah_akun')"
        class="rounded-lg bg-[#2f4a8a] px-4 py-2 text-sm font-medium text-white hover:bg-[#37bd1d]">
        + Tambah Akun
      </button>

    </div>

    <!-- TABLE -->
    <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

      <!-- TABLE HEADER (biru) -->
      <div class="bg-[#243e90] px-5 py-4">
        <h2 class="text-white text-2xl font-bold">Data Akun</h2>
        <p class="text-white text-sm mt-1">Kelola data akun pengguna</p>
      </div>

      <!-- TABLE CONTENT -->
      <div class="overflow-x-auto">
        <table class="w-full text-sm">

          <thead>
            <tr class="text-left text-black-800 border-b border-gray-200">
              <th class="px-5 py-4">No</th>
              <th class="px-5 py-4">Email</th>
              <th class="px-5 py-4">Nama</th>
              <th class="px-5 py-4">Role</th>
              <th class="px-5 py-4">Aksi</th>
            </tr>
          </thead>

<tbody>
  <tr
  v-for="(item, index) in paginatedUsers"
    :key="item.id"
    class="border-t border-gray-100"
  >
    <td class="px-5 py-4">
      {{
        search
          ? index + 1
          : (currentPage - 1) * perPage + index + 1
      }}
    </td>

    <td class="px-5 py-4">{{ item.email }}</td>
    <td class="px-5 py-4">{{ item.name }}</td>
    <td class="px-5 py-4">{{ item.role_name }}</td>

    <td class="px-5 py-4">
      <div class="flex gap-2">
        <button
          @click="router.push(`/dashboard-superadmin/edit_akun/${item.id}`)"
          class="rounded-md bg-yellow-400 px-3 py-1 text-xs font-medium text-white hover:bg-yellow-900"
        >
          ✏ Edit
        </button>

        <button
          @click="openDeleteModal(item.id, item.email)"
          class="rounded-md bg-red-500 px-3 py-1 text-xs font-medium text-white hover:bg-red-900"
        >
          🗑 Hapus
        </button>

        <button
          @click="router.push(`/dashboard-superadmin/reset_password/${encodeURIComponent(item.email)}`)"
          class="rounded-md bg-blue-500 px-3 py-1 text-xs font-medium text-white hover:bg-blue-900"
        >
          🔑 Reset Password
        </button>
      </div>
    </td>
  </tr>

  <tr v-if="filteredUsers.length === 0">
    <td colspan="5" class="px-5 py-8 text-center text-gray-400">
      Memuat..
    </td>
  </tr>
</tbody>

        </table>
      </div>

      <!-- FOOTER PAGINATION -->
      <div class="flex items-center justify-between border-t border-gray-100 px-5 py-4">

        <!-- KIRI: Tampilkan per halaman -->
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-600">Tampilkan</span>
          <select v-model.number="perPage"
            class="rounded border border-gray-300 bg-white px-2 py-1 text-sm text-gray-700 focus:outline-none focus:border-blue-400">
            <option :value="5">5</option>
            <option :value="10">10</option>
            <option :value="20">20</option>
          </select>
          <span class="text-sm text-gray-600">akun per halaman</span>
        </div>

        <!-- KANAN: Navigasi halaman -->
        <div class="flex items-center gap-2 text-gray-500 text-sm">

          <button @click="prevPage" :disabled="currentPage === 1"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Previous
          </button>

          <template v-for="item in visiblePages" :key="item">
            <span v-if="item === '...'" class="w-8 h-8 flex items-center justify-center text-gray-400">...</span>
            <button v-else @click="currentPage = Number(item); getUsers()" :class="currentPage === Number(item)
              ? 'bg-[#1c3277] text-white shadow-md scale-105'
              : 'bg-white text-[#4b4b4b] hover:bg-[#d6ddee]'"
              class="w-8 h-8 rounded-md text-sm font-medium transition-all duration-200">
              {{ item }}
            </button>
          </template>

          <button @click="nextPage" :disabled="currentPage === totalPages"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            Next
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

        </div>

      </div>

    </div>

    <!-- MODAL HAPUS -->
    <DeleteConfirmModal v-if="Konfirmasi_hapus" :email="selectedEmail" @close="Konfirmasi_hapus = false"
      @confirm="confirmDelete" />

  </div>
</template>>