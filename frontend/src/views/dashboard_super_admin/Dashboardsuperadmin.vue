<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue"
import { useRouter, useRoute } from "vue-router"
import Konfirmasi_keluar from "./akademik/konfirmasi_keluar.vue"

const router = useRouter()
const route = useRoute()

// ================= STATE =================
// ================= STATE =================
const showLogoutPopup = ref(false)
const isSidebarOpen = ref(true)

const akunList = ref<any[]>([])
const totalAkun = ref(0)
const totalRole = ref(0)

const rowsPerPage = ref(10)
const currentPage = ref(1)
const totalPages = ref(1)

const user = ref({
  name: "Admin Akademik"
})



const paginatedData = computed(() => akunList.value)

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

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    await getAkun(currentPage.value + 1)
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    await getAkun(currentPage.value - 1)
  }
}

// ================= ACTIVE MENU =================
const isActive = (path: string) => route.path === path

// ================= API =================
const getHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem("token")}`
})

const getTotalAkun = async () => {
  try {
    const BASE_URL = 'https://be.karlearn.site'
    const res = await fetch(`${BASE_URL}/api/users/count`, {
      headers: getHeaders()
    })

    const data = await res.json()

    totalAkun.value = data.data ?? 0
  } catch (err) {
    console.error(err)
  }
}

const getTotalRole = async () => {
  try {
    const BASE_URL = "https://be.karlearn.site"

    const res = await fetch(`${BASE_URL}/api/roles`, {
      headers: getHeaders()
    })

    const data = await res.json()

    totalRole.value = data?.data?.total_items || 0

  } catch (err) {
    console.error(err)
  }
}

const getAkun = async (page = 1) => {
  try {
    const BASE_URL = 'https://be.karlearn.site'
    const res = await fetch(
      `${BASE_URL}/api/users?page=${page}&per_page=${rowsPerPage.value}`,
      {
        headers: {
          ...getHeaders(),
          accept: "application/json"
        }
      }
    )

    console.log("REQUEST URL:", res.url)

    const data = await res.json()

    akunList.value = data.data.items || []

    totalPages.value = data.data.pagination.total_pages || 1

    currentPage.value = data.data.pagination.page || 1

  } catch (err) {
    console.error(err)
  }
}

watch(rowsPerPage, () => {
  getAkun(1)
})


// ================= LOGOUT =================
const handleLogout = () => {
  localStorage.clear()
  router.push("/")
}

onMounted(() => {
  getTotalAkun()
  getTotalRole()
  getAkun()
})
</script>

<template>
  <div class="flex h-screen bg-[#dfe6ef] overflow-hidden">

    <!-- SIDEBAR -->
    <aside :class="isSidebarOpen ? 'w-55' : 'w-17.5'" class="bg-[#b8c9e3] transition-all duration-300 flex flex-col">

      <!-- LOGO -->
      <div class="h-16 bg-[#243e90] flex items-center px-4 gap-3">

        <img src="@/assets/images/logo.png" class="w-10 h-10" />

        <h1 v-if="isSidebarOpen" class="text-white font-bold text-[34px]">
          Sabar
        </h1>

      </div>

      <!-- MENU -->
      <div class="flex-1 px-2 py-3">

        <!-- DASHBOARD -->
        <div @click="router.push('/dashboard-superadmin')" :class="isActive('/dashboard-superadmin')
          ? 'bg-[#243e90] text-white'
          : 'text-[#000000] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer transition">

          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M3 13h8V3H3v10zm10 8h8V11h-8v10zM3 21h8v-6H3v6zm10-18v6h8V3h-8z" />
          </svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Dashboard
          </span>

        </div>

        <!-- AKUN -->
         
        <div @click="router.push('/dashboard-superadmin/akun')" :class="isActive('/dashboard-superadmin/akun')
          ? 'bg-[#243e90] text-white'
          : 'text-[#000000] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer mt-2 transition">

          <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
  <path d="M0 0h24v24H0z" fill="none" />
  <path fill="currentColor" d="M15 14c-2.67 0-8 1.33-8 4v2h16v-2c0-2.67-5.33-4-8-4m-9-4V7H4v3H1v2h3v3h2v-3h3v-2m6 2a4 4 0 0 0 4-4a4 4 0 0 0-4-4a4 4 0 0 0-4 4a4 4 0 0 0 4 4" />
</svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Akun
          </span>

        </div>

        <!-- ROLE -->
        <div @click="router.push('/dashboard-superadmin/role')" :class="isActive('/dashboard-superadmin/role')
          ? 'bg-[#243e90] text-white'
          : 'text-[#000000] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer mt-2 transition">

          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M19 3H5c-1.1 0-2 .9-2 2v14a2 2 0 002 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14H7v-2h5v2zm5-4H7v-2h10v2zm0-4H7V7h10v2z" />
          </svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Role
          </span>

        </div>

        <!-- RESET PASSWORD -->
        <div @click="router.push('/dashboard-superadmin/reset_password')" :class="isActive('/dashboard-superadmin/reset_password')
          ? 'bg-[#243e90] text-white'
          : 'text-[#000000] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer mt-2 transition">

      <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
  <path d="M0 0h24v24H0z" fill="none" />
  <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
    <path d="M2.586 17.414A2 2 0 0 0 2 18.828V21a1 1 0 0 0 1 1h3a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h1a1 1 0 0 0 1-1v-1a1 1 0 0 1 1-1h.172a2 2 0 0 0 1.414-.586l.814-.814a6.5 6.5 0 1 0-4-4z" />
    <circle cx="16.5" cy="7.5" r=".5" fill="currentColor" />
  </g>
</svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Reset Password
          </span>

        </div>

      </div>

      <!-- LOGOUT -->
      <div @click="showLogoutPopup = true"
        class="p-4 flex items-center gap-3 cursor-pointer text-[#000000] hover:bg-[#9fb5d6]">

        <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 14 14">
  <path d="M0 0h14v14H0z" fill="none" />
  <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" d="M9.5 10.5v2a1 1 0 0 1-1 1h-7a1 1 0 0 1-1-1v-11a1 1 0 0 1 1-1h7a1 1 0 0 1 1 1v2M6.5 7h7m-2-2l2 2l-2 2" />
</svg>

        <span v-if="isSidebarOpen" class="text-sm">
          Keluar
        </span>

      </div>

    </aside>

    <!-- MAIN -->
    <div class="flex-1 flex flex-col">

      <!-- TOPBAR -->
      <div class="h-16 bg-[#243e90] flex items-center justify-between px-5">

        <!-- LEFT -->
        <div class="flex items-center gap-4">

          <!-- HAMBURGER -->
          <button @click="isSidebarOpen = !isSidebarOpen" class="text-white">

            <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>

          </button>

        </div>

        <!-- USER -->
        <div class="flex items-center gap-3">

          <img src="https://i.pravatar.cc/40" class="w-10 h-10 rounded-full border-2 border-white" />

          <div class="text-white font-medium text-sm">
            {{ user.name }}
          </div>

        </div>

      </div>

      <!-- CONTENT -->
      <div class="flex-1 overflow-auto p-4">
        <router-view />

        <!-- DASHBOARD -->
        <div v-if="route.path === '/dashboard-superadmin'">

          <!-- TITLE -->
          <h1 class="text-[42px] font-bold text-black">
            Dashboard
          </h1>

          <p class="mt-1 text-[#000000]">
            Selamat Datang Super Admin
          </p>

          <!-- CARD -->
          <div class="mt-8 flex gap-6">

            <!-- TOTAL AKUN -->
            <div class="relative bg-[#ececec] rounded-xl w-[320px] p-4 flex items-center gap-4 shadow-sm
         border-l-[4px] border-b-[3px] border-[#9db9dc]">
              <!-- Icon -->
              <div class="w-20 h-20 rounded-lg bg-[#9db9dc] flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-[#4c4c4c]" viewBox="0 0 24 24"
                  fill="currentColor">
                  <path d="M0 0h24v24H0z" fill="none" />
                  <path
                    d="M11 10v2H9v2H7v-2H5.8c-.4 1.2-1.5 2-2.8 2c-1.7 0-3-1.3-3-3s1.3-3 3-3c1.3 0 2.4.8 2.8 2zm-8 0c-.6 0-1 .4-1 1s.4 1 1 1s1-.4 1-1s-.4-1-1-1m13 4c2.7 0 8 1.3 8 4v2H8v-2c0-2.7 5.3-4 8-4m0-2c-2.2 0-4-1.8-4-4s1.8-4 4-4s4 1.8 4 4s-1.8 4-4 4" />
                </svg>
              </div>

              <div>
                <div class="text-sm font-bold text-[#000000]">
                  TOTAL AKUN
                </div>

                <div class="text-4xl font-bold mt-3">
                  {{ totalAkun }}
                </div>
              </div>
            </div>



            <!-- TOTAL ROLE -->
            <div
              class="bg-[#ececec] rounded-xl w-[320px] p-4 flex items-center gap-4 shadow-sm border-l-[3px] border-b-[2px] border-[#9db9dc]">

              <!-- Icon -->
              <div class="w-20 h-20 rounded-lg bg-[#9db9dc] flex items-center justify-center">

                <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-[#000000]" viewBox="0 0 256 256"
                  fill="currentColor">
                  <path d="M0 0h256v256H0z" fill="none" />
                  <path
                    d="M213.74 42.26a76 76 0 0 0-125.23 79.58l-57 57a11.93 11.93 0 0 0-3.51 8.47V216a12 12 0 0 0 12 12h32a4 4 0 0 0 4-4v-20h20a4 4 0 0 0 4-4v-20h20a4 4 0 0 0 2.83-1.17l11.33-11.34A75.7 75.7 0 0 0 160 172h.1a76 76 0 0 0 53.64-129.74m14.22 56c-1.15 36.22-31.6 65.72-67.87 65.77H160a67.5 67.5 0 0 1-25.21-4.83a4 4 0 0 0-4.45.83l-12 12H96a4 4 0 0 0-4 4v20H72a4 4 0 0 0-4 4v20H40a4 4 0 0 1-4-4v-28.72a4.06 4.06 0 0 1 1.17-2.83L96 125.66a4 4 0 0 0 .83-4.45A67.5 67.5 0 0 1 92 95.91c0-36.27 29.55-66.72 65.77-67.91A68 68 0 0 1 228 98.23ZM188 76a8 8 0 1 1-8-8a8 8 0 0 1 8 8" />
                </svg>

              </div>

              <!-- Text -->
              <div class="flex flex-col justify-center">

                <div class="text-sm font-bold text-[#000000]">
                  TOTAL ROLE
                </div>

                <div class="text-4xl font-bold mt-3 text-[#000000]">
                  {{ totalRole }}
                </div>

              </div>

            </div>

          </div>

        </div>

        <!-- TABLE -->
<div v-if="route.path === '/dashboard-superadmin'" class="mt-8 bg-[#ececec] rounded-xl p-5 min-h-125 flex flex-col justify-between border-l-[3px] border-b-[2px] border-[#9db9dc] shadow-sm">
          <div>

            <h2 class="text-[32px] font-semibold mb-8 text-[#000000]">
              Data Akun
            </h2>

            <table class="w-full">

              <thead>
                <tr class="text-[#000000] text-[15px]">
                  <th class="text-left py-3 font-semibold">No</th>
                  <th class="text-left py-3 font-semibold">Email</th>
                  <th class="text-left py-3 font-semibold">Nama</th>
                  <th class="text-left py-3 font-semibold">Jabatan</th>
                </tr>
              </thead>

              <tbody>
                <tr v-for="(item, index) in paginatedData" :key="item.id" class="text-[#000000]">
                  <td class="py-4">
                    {{
                      (currentPage - 1) * rowsPerPage + index + 1
                    }}
                  </td>

                  <td>{{ item.email }}</td>

                  <td>{{ item.name }}</td>

                  <td>{{ item.role_name }}</td>
                </tr>
              </tbody>

            </table>

          </div>

          <!-- FOOTER -->
          <div class="flex items-center justify-between mt-10">

            <!-- PAGINATION -->
             <div class="flex justify-end mt-10"></div>
<div class="flex items-center gap-2 text-gray-500 text-sm">

  <!-- PREVIOUS -->
  <button @click="prevPage" :disabled="currentPage === 1"
    class="flex items-center gap-1 px-2 py-1 rounded hover:text-black disabled:opacity-40">
    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
    </svg>
    Previous
  </button>

  <!-- PAGE NUMBERS -->
  <template v-for="item in visiblePages" :key="item">
    <span v-if="item === '...'" class="w-8 h-8 flex items-center justify-center text-gray-400">...</span>
    <button v-else @click="getAkun(Number(item))" :class="currentPage === Number(item)
      ? 'bg-[#1c3277] text-white shadow-md scale-105'
      : 'bg-white text-[#4b4b4b] hover:bg-[#d6ddee]'"
      class="w-8 h-8 rounded-md text-sm font-medium transition-all duration-200">
      {{ item }}
    </button>
  </template>

  <!-- NEXT -->
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

      </div>

    </div>

  </div>

  <!-- POPUP -->
  <Konfirmasi_keluar v-if="showLogoutPopup" @close="showLogoutPopup = false" @confirm="handleLogout" />

</template>