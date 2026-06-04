<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
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
const perPage = ref(10)
const currentPage = ref(1)

const Konfirmasi_hapus = ref(false)
const selectedEmail = ref("")

const router = useRouter()

const users = ref<User[]>([])
const totalPages = ref(1)

// ================= HEADERS =================
const getHeaders = () => {
  const token = localStorage.getItem("token")

  return {
    Authorization: `Bearer ${token}`,
    accept: "application/json",
  }
}

// ================= GET USERS =================
const getUsers = async () => {
  try {
    const response = await fetch(
      `${API.users}?page=${currentPage.value}&per_page=${perPage.value}`,
      {
        method: "GET",
        headers: getHeaders(),
      }
    )

    const data = await response.json()

    console.log("USERS:", data)

    if (!response.ok) {
      console.log("Gagal mengambil user")
      return
    }

    users.value = data.data.items || []

    totalPages.value =
      data.data.pagination?.total_pages || 1

  } catch (err) {
    console.error("GET USERS ERROR:", err)
  }
}

// ================= FILTER =================
const filteredUsers = computed(() =>
  users.value.filter((item) =>
    item.name?.toLowerCase().includes(search.value.toLowerCase()) ||
    item.email?.toLowerCase().includes(search.value.toLowerCase()) ||
    item.role_name?.toLowerCase().includes(search.value.toLowerCase())
  )
)

// ================= PAGINATION =================
const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getUsers()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getUsers()
  }
}

// ================= DELETE =================
const openDeleteModal = (email: string) => {
  selectedEmail.value = email
  Konfirmasi_hapus.value = true
}

const confirmDelete = () => {
  console.log("hapus:", selectedEmail.value)
  Konfirmasi_hapus.value = false
}

// ================= INIT =================
onMounted(() => {
  getUsers()
})
</script>

<template>
  <div class="min-h-screen bg-[#f5f7fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Akun
    </div>

    <!-- TITLE -->
    <h1 class="mb-1 text-3xl font-bold text-gray-800">
      Akun
    </h1>

    <p class="mb-6 text-gray-500">
      Kelola Data Akun
    </p>

    <!-- HEADER -->
    <div class="mb-4 flex items-center justify-between">

      <!-- SEARCH -->
      <div class="relative">

        <input
          v-model="search"
          type="text"
          placeholder="Cari Akun..."
          class="w-64 rounded-lg border border-gray-200 bg-white py-2 pl-4 pr-10 text-sm outline-none focus:border-blue-500"
        />

        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="absolute right-3 top-2.5 h-4 w-4 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-4.35-4.35m1.85-5.15a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>

      </div>

      <!-- BUTTON -->
      <button
        @click="router.push('/dashboard-superadmin/tambah_akun')"
        class="rounded-lg bg-[#2f4a8a] px-4 py-2 text-sm font-medium text-white hover:bg-[#243b73]"
      >
        + Tambah Akun
      </button>

    </div>

    <!-- TABLE -->
    <div class="rounded-2xl border border-gray-200 bg-white shadow-sm">

      <!-- TABLE HEADER -->
      <div class="border-b border-gray-100 px-5 py-4">

        <h2 class="text-xl font-semibold text-gray-700">
          Data Akun
        </h2>

      </div>

      <!-- TABLE CONTENT -->
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

            <tr
              v-for="(item, index) in filteredUsers"
              :key="item.id"
              class="border-t border-gray-100"
            >

              <td class="px-5 py-4">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <td class="px-5 py-4">
                {{ item.email }}
              </td>

              <td class="px-5 py-4">
                {{ item.name }}
              </td>

              <td class="px-5 py-4">
                {{ item.role_name }}
              </td>

              <td class="px-5 py-4">

                <div class="flex gap-2">

                  <!-- EDIT -->
                  <button
                    @click="router.push(`/dashboard-superadmin/edit_akun/${item.id}`)"
                    class="rounded-md bg-yellow-400 px-3 py-1 text-xs font-medium text-white hover:bg-yellow-500"
                  >
                    ✏ Edit
                  </button>

                  <!-- DELETE -->
                  <button
                    @click="openDeleteModal(item.email)"
                    class="rounded-md bg-red-500 px-3 py-1 text-xs font-medium text-white hover:bg-red-600"
                  >
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
        <div class="flex items-center gap-3 text-sm">

          <button
            @click="prevPage"
            class="text-gray-400 hover:text-black"
          >
            ← Previous
          </button>

          <button
            class="flex h-8 w-8 items-center justify-center rounded bg-[#2f4a8a] text-white"
          >
            {{ currentPage }}
          </button>

          <button
            @click="nextPage"
            class="text-gray-400 hover:text-black"
          >
            Next →
          </button>

        </div>

      </div>

    </div>

    <!-- MODAL -->
    <DeleteConfirmModal
      v-if="Konfirmasi_hapus"
      :email="selectedEmail"
      @close="Konfirmasi_hapus = false"
      @confirm="confirmDelete"
    />

  </div>
</template>