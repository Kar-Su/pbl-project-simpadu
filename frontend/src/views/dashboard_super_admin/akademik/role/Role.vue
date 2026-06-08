<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"
import DeleteConfirmModal from "@/views/dashboard_super_admin/akademik/akun/konfirmasi_hapus.vue"

interface User {
  id: number
  name: string
}

const search = ref("")
const perPage = ref(10)
const currentPage = ref(1)

const Konfirmasi_hapus = ref(false)

const router = useRouter()
const name = ref<User[]>([])
const selectedRole = ref("")
const selectedRoleId = ref<string | null>(null)

const openDeleteModal = (id: number, role: string) => {
  selectedRoleId.value = id.toString()
  selectedRole.value = role
  Konfirmasi_hapus.value = true
}

const getUsers = async () => {
  try {
    const token = localStorage.getItem("token")
    const BASE_URL = 'https://be.karlearn.site'
    const response = await fetch(`${BASE_URL}/api/roles`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    const data = await response.json()

    console.log("ROLE:", data)

    name.value = data?.data?.role || []

  } catch (err) {
    console.log(err)
  }
}

const filteredUsers = computed(() =>
  name.value.filter(item =>
    item.name?.toLowerCase().includes(search.value.toLowerCase())
  )
)

// ================= PAGINATION =================
const totalPages = computed(() =>
  Math.ceil(filteredUsers.value.length / perPage.value)
)

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

  for (let i = rangeStart; i <= rangeEnd; i++) {
    pages.push(i)
  }

  if (rangeEnd < total - 1) pages.push('...')

  pages.push(total)

  return pages
})
const paginatedUsers = computed(() => {

  const start = (currentPage.value - 1) * perPage.value

  return filteredUsers.value.slice(
    start,
    start + perPage.value
  )
})

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}


const confirmDelete = async () => {
  try {
    const token = localStorage.getItem("token")

    const response = await fetch(
      `/api/roles/${selectedRole.value}`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    )

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || "Gagal hapus role")
    }

    alert("Role berhasil dihapus")

    Konfirmasi_hapus.value = false

    getUsers()

  } catch (err: any) {
    alert(err.message)
  }
}

// ================= MOUNTED =================
onMounted(() => {
  getUsers()
})
</script>

<template>
  <div class="min-h-screen bg-[#f5f7fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Role
    </div>

    <!-- TITLE -->
    <h1 class="mb-1 text-3xl font-bold text-gray-800">
      Role
    </h1>

    <p class="mb-6 text-gray-500">
      Kelola Role
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

        

      <!-- BUTTON -->
      <button @click="router.push('/dashboard-superadmin/tambah_role')"
        class="rounded-lg bg-[#2f4a8a] px-4 py-2 text-sm font-medium text-white hover:bg-[#37bd1d]">
        + Tambah Role
      </button>

    </div>

    <!-- TABLE -->
    <div class="bg-[#ececec] rounded-xl shadow-sm
border-l-[4px] border-b-[3px] border-[#9db9dc]">

      <!-- TABLE HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">

        <h2 class="text-xl font-semibold text-gray-700">
          Data Role
        </h2>

      </div>

      <!-- TABLE CONTENT -->
      <div class="overflow-x-auto">

        <table class="w-full text-sm">

          <thead>

            <tr class="text-left text-gray-500">

              <th class="px-5 py-4">No</th>
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
                {{ item.name }}
              </td>

              <td class="px-5 py-4">

                <div class="flex gap-2">

                  <!-- EDIT -->
                  <button @click="router.push(`/dashboard-superadmin/edit_role/${item.id}`)"
                    class="rounded-md bg-yellow-400 px-3 py-1 text-xs font-medium text-white hover:bg-yellow-500">
                    ✏ Edit
                  </button>

                  <!-- DELETE -->
                  <button @click="openDeleteModal(item.id, item.name)"
                    class="rounded-md bg-red-500 px-3 py-1 text-xs font-medium text-white hover:bg-red-600">
                    🗑 Hapus
                  </button>

                </div>

              </td>

            </tr>

          </tbody>

        </table>

      </div>

      <!-- FOOTER -->
      <div class="flex items-center justify-between border-t border-gray-100 px-5 py-4">


<!-- PAGINATION -->
 <div class="flex justify-end mt-10"></div>
        <div class="flex items-center gap-2 text-gray-500 text-sm">

          <button @click="prevPage" :disabled="currentPage === 1"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Previous
          </button>

          <template v-for="item in visiblePages" :key="item">
            <span v-if="item === '...'" class="w-8 h-8 flex items-center justify-center text-gray-400">...</span>
            <button v-else @click="currentPage = Number(item)" :class="currentPage === Number(item)
              ? 'bg-[#1c3277] text-white shadow-md scale-105'
              : 'bg-white text-[#4b4b4b] hover:bg-[#d6ddee]'"
              class="w-8 h-8 rounded-md text-sm font-medium transition-all duration-200">
              {{ item }}
            </button>
          </template>

          <button @click="nextPage" :disabled="currentPage === totalPages"
            class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
            Next
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

        </div>

      </div>

    </div>

    <!-- MODAL -->
<DeleteConfirmModal
  v-if="Konfirmasi_hapus"
  :email="selectedRole"
  @close="Konfirmasi_hapus = false"
  @confirm="confirmDelete"
/>

  </div>
</template>