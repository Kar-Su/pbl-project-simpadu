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

const rowsPerPage = ref(10)
const currentPage = ref(1)
const totalPages = ref(1)

const user = ref({
  name: "Admin Akademik"
})



const paginatedData = computed(() => akunList.value)

const visiblePages = computed(() => {
  return Array.from(
    { length: totalPages.value },
    (_, i) => i + 1
  )
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
    const res = await fetch(`/api/users/count`, {
      headers: getHeaders()
    })

    const data = await res.json()

    totalAkun.value = data.data ?? 0
  } catch (err) {
    console.error(err)
  }
}

const getAkun = async (page = 1) => {
  try {
   const res = await fetch(
  `/api/users?page=${page}&per_page=${rowsPerPage.value}`,
  {
    headers: {
      ...getHeaders(),
      accept: "application/json"
    }
  }
)

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
          : 'text-[#4b4b4b] hover:bg-[#9fb5d6]'
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
          : 'text-[#4b4b4b] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer mt-2 transition">

          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zM8 11c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5C15 14.17 10.33 13 8 13zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z" />
          </svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Akun
          </span>

        </div>

        <!-- ROLE -->
        <div @click="router.push('/dashboard-superadmin/role')" :class="isActive('/dashboard-superadmin/role')
          ? 'bg-[#243e90] text-white'
          : 'text-[#4b4b4b] hover:bg-[#9fb5d6]'
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
          : 'text-[#4b4b4b] hover:bg-[#9fb5d6]'
          " class="flex items-center gap-3 px-3 py-2 rounded-md cursor-pointer mt-2 transition">

          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path
              d="M12 17a2 2 0 002-2c0-.74-.4-1.38-1-1.72V11a1 1 0 10-2 0v2.28A2 2 0 0012 17zm6-6V9a6 6 0 10-12 0v2a2 2 0 00-2 2v7a2 2 0 002 2h12a2 2 0 002-2v-7a2 2 0 00-2-2zm-8-2a4 4 0 118 0v2h-8V9z" />
          </svg>

          <span v-if="isSidebarOpen" class="text-sm">
            Reset Password
          </span>

        </div>

      </div>

      <!-- LOGOUT -->
      <div @click="showLogoutPopup = true"
        class="p-4 flex items-center gap-3 cursor-pointer text-[#4b4b4b] hover:bg-[#9fb5d6]">

        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path
            d="M10.09 15.59L8.67 17l5 5 5-5-1.41-1.41L15 18.17V10h-2v8.17l-2.91-2.58zM4 4h8V2H4c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h8v-2H4V4z" />
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

        <!-- DASHBOARD -->
        <div v-if="route.path === '/dashboard-superadmin'">

          <!-- TITLE -->
          <h1 class="text-[42px] font-bold text-black">
            Dashboard
          </h1>

          <p class="mt-1 text-[#4b4b4b]">
            Selamat Datang Super Admin
          </p>

          <!-- CARD -->
          <div class="mt-8 bg-[#ececec] rounded-xl w-[320px] p-4 flex items-center gap-4">

            <div class="w-20 h-20 rounded-lg bg-[#9db9dc] flex items-center justify-center">

              <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-[#4c4c4c]" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M17 20h5V4H2v16h5m10 0v-2a4 4 0 00-4-4H9a4 4 0 00-4 4v2m12 0H5" />
              </svg>

            </div>

            <div>

              <div class="text-sm font-bold text-[#4b4b4b]">
                TOTAL AKUN
              </div>

              <div class="text-4xl font-bold mt-3">
                {{ totalAkun }}
              </div>

            </div>

          </div>

          <!-- TABLE -->
          <div class="mt-8 bg-[#ececec] rounded-xl p-5 min-h-125 flex flex-col justify-between">

            <div>

              <h2 class="text-[32px] font-semibold mb-8 text-[#4a4a4a]">
                Data Akun
              </h2>

              <table class="w-full">

                <thead>
                  <tr class="text-[#5b5b5b] text-[15px]">
                    <th class="text-left py-3 font-semibold">No</th>
                    <th class="text-left py-3 font-semibold">Email</th>
                    <th class="text-left py-3 font-semibold">Nama</th>
                    <th class="text-left py-3 font-semibold">Jabatan</th>
                  </tr>
                </thead>

              <tbody>
  <tr
    v-for="(item, index) in paginatedData"
    :key="item.id"
    class="text-[#4b4b4b]"
  >
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
              <div class="flex items-center gap-5 text-gray-500 text-sm">

                <!-- PREVIOUS -->
                <button @click="prevPage" class="text-gray-400 hover:text-black">

                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>

                  Previous

                </button>

                <!-- PAGE -->
                <div class="flex items-center gap-2">

                  <button v-for="page in visiblePages" :key="page" @click="getAkun(page)" :class="currentPage === page
                    ? 'bg-[#1c3277] text-white shadow-md scale-105'
                    : 'bg-white text-[#4b4b4b] hover:bg-[#d6ddee]'"
                    class="w-8 h-8 rounded-md text-sm font-medium transition-all duration-200">
                    {{ page }}
                  </button>

                </div>

                <!-- NEXT -->
                <button @click="nextPage" class="text-gray-400 hover:text-black">

                  Next

                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
                    stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>

                </button>

              </div>

            </div>

          </div>

        </div>

        <router-view />

      </div>

    </div>

    <!-- POPUP -->
    <Konfirmasi_keluar v-if="showLogoutPopup" @close="showLogoutPopup = false" @confirm="handleLogout" />

  </div>
</template>