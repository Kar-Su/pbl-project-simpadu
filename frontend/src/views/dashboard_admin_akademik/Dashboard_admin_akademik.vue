<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import KonfirmasiKeluar from "@/views/dashboard_super_admin/akademik/Konfirmasi_keluar.vue"

const router = useRouter()
const route = useRoute()

const goTo = (path: string) => router.push(path)
const isActive = (path: string) => route.path === path

// ─────────────────────────────────────────────
// SIDEBAR STATE
// ─────────────────────────────────────────────
const isSidebarOpen = ref<boolean>(true)

const openAkademik = ref<boolean>(true)
const openMahasiswa = ref<boolean>(false)
const showLogoutPopup = ref<boolean>(false)

// ─────────────────────────────────────────────
// HELPER: ambil token dari localStorage
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  'Content-Type': 'application/json',
  'accept': 'application/json',
  'Authorization': `Bearer ${localStorage.getItem('token') ?? ''}`,
})

// ─────────────────────────────────────────────
// LOGOUT
// ─────────────────────────────────────────────
const handleLogout = (): void => {
  localStorage.removeItem('token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('role')
  router.push('/')
}

// ─────────────────────────────────────────────
// STAT CARDS
// ─────────────────────────────────────────────
const totalPegawai = ref<number>(0)
const totalMahasiswa = ref<number>(0)
const totalDosen = ref<number>(0)

const getAllUsers = async (): Promise<any[]> => {
  try {
    let page = 1
    let lastPage = 1
    let allUsers: any[] = []

    do {
      const res = await fetch(`/api/users?page=${page}`, { headers: getHeaders() })
      const data = await res.json()
      const items = data.data.items ?? []
      allUsers = [...allUsers, ...items]
      lastPage = data.data.pagination?.total_pages ?? 1
      page++
    } while (page <= lastPage)

    return allUsers
  } catch (err) {
    console.error('getAllUsers:', err)
    return []
  }
}

const getTotalPegawai = async (): Promise<void> => {
  try {
    const users = await getAllUsers()
    totalPegawai.value = users.filter((item: any) => {
      const role = item.role_name?.toLowerCase()?.trim()
      return role !== 'mahasiswa'
    }).length
    console.log('TOTAL PEGAWAI:', totalPegawai.value)
  } catch (err) {
    console.error('getTotalPegawai:', err)
  }
}

const getTotalMahasiswa = async (): Promise<void> => {
  try {
    const users = await getAllUsers()
    totalMahasiswa.value = users.filter(
      (item: any) => item.role_name?.toLowerCase() === 'mahasiswa'
    ).length
  } catch (err) {
    console.error('getTotalMahasiswa:', err)
  }
}

const getTotalDosen = async (): Promise<void> => {
  try {
    const users = await getAllUsers()
    totalDosen.value = users.filter((item: any) => {
      const role = item.role_name?.toLowerCase()?.trim()
      return role?.includes('dosen')
    }).length
    console.log('TOTAL DOSEN:', totalDosen.value)
  } catch (err) {
    console.error('getTotalDosen:', err)
  }
}

// ─────────────────────────────────────────────
// TAHUN AKADEMIK
// ─────────────────────────────────────────────
interface TahunAkademik {
  id: number
  tahun_awal: string
  tahun_akhir: string
}
const tahunAkademik = ref<TahunAkademik[]>([])

const getTahunAkademik = async (): Promise<void> => {
  try {
    const res = await fetch('/api/tahun-akademik', {
      headers: getHeaders()
    })
    // const res = await fetch('/api/tahun-akademik/status/aktif', {
    //   headers: getHeaders()
    // })
    const data = await res.json()

    console.log('TAHUN AKADEMIK:', data)

    tahunAkademik.value = Array.isArray(data.data)
      ? data.data
      : []

  } catch (err) {
    console.error('getTahunAkademik:', err)
  }
}

// ─────────────────────────────────────────────
// KURIKULUM
// ─────────────────────────────────────────────
interface Kurikulum {
  id: string
  name: string
}

const kurikulum = ref<Kurikulum[]>([])

const formatYear = (date: string): string => {
  return new Date(date).getFullYear().toString()
}

const getKurikulum = async (): Promise<void> => {
  try {
    const res = await fetch('/api/kurikulum?page=1', { headers: getHeaders() })
    const data = await res.json()
    console.log(data)
    kurikulum.value = data.data.items ?? []
  } catch (err) {
    console.error('getKurikulum:', err)
  }
}

// ─────────────────────────────────────────────
// MAHASISWA
// ─────────────────────────────────────────────
interface MahasiswaItem {
  id: number
  nim: string
  name: string
  kelas: string
  prodi: string
  angkatan: string
}

const akunList = ref<MahasiswaItem[]>([])

const currentPage = ref<number>(1)
const perPage = ref<number>(10)
const totalItems = ref<number>(0)

const totalPages = computed<number>(() =>
  Math.max(1, Math.ceil(totalItems.value / perPage.value))
)


const allMahasiswa = ref<MahasiswaItem[]>([])

const getAkun = async (): Promise<void> => {
  try {

    let page = 1
    let lastPage = 1
    let allUsers: any[] = []

    do {

      const res = await fetch(`/api/users?page=${page}`, {
        headers: getHeaders()
      })

      const data = await res.json()

      const items = data.data.items ?? []

      allUsers = [...allUsers, ...items]

      lastPage = data.data.pagination?.total_pages ?? 1

      page++

    } while (page <= lastPage)

    // FILTER MAHASISWA
    const mahasiswaOnly = allUsers.filter(
      (item: any) =>
        item.role_name?.toLowerCase() === 'mahasiswa'
    )

    // MAP DATA
    allMahasiswa.value = mahasiswaOnly.map((item: any) => ({
      id: item.id,
      nim: item.detail?.nim ?? '-',
      name: item.name ?? '-',
      kelas: item.detail?.kelas ?? '-',
      prodi: item.detail?.prodi ?? '-',
      angkatan: item.detail?.angkatan ?? '-',
    }))

    // TOTAL SETELAH FILTER
    totalItems.value = allMahasiswa.value.length

    // PAGINATION FRONTEND
    const start = (currentPage.value - 1) * perPage.value
    const end = start + perPage.value

    akunList.value = allMahasiswa.value.slice(start, end)

  } catch (err) {
    console.error('getAkun:', err)
  }
}

const pages = computed<(number | string)[]>(() => {
  const total = totalPages.value
  const cur = currentPage.value
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)

  const result: (number | string)[] = [1, 2]
  if (cur > 4) result.push('...')
  for (let i = Math.max(3, cur - 1); i <= Math.min(total - 2, cur + 1); i++) {
    result.push(i)
  }
  if (cur < total - 3) result.push('...')
  result.push(total - 1, total)
  return [...new Set(result)]
})

const goToPage = async (page: number): Promise<void> => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  await getAkun()
}

const prevPage = async (): Promise<void> => {
  await goToPage(currentPage.value - 1)
}

const nextPage = async (): Promise<void> => {
  await goToPage(currentPage.value + 1)
}

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted((): void => {
  getTotalPegawai()
  getTotalMahasiswa()
  getTotalDosen()
  getTahunAkademik()
  getKurikulum()
  getAkun()
})
</script>

<template>
  <div class="bg-[#eef3fb] min-h-screen overflow-y-auto overflow-x-hidden">

    <header class="fixed top-0 left-0 right-0 h-18 bg-[#1f3c93] flex items-center justify-between px-6 z-50 shadow-md">
      <!-- KIRI -->
      <div class="flex items-center gap-5">
        <div class="flex items-center gap-3">
          <img src="@/assets/images/logo.png" alt="logo" class="w-10 h-10 object-contain" />
          <h1 class="text-2xl font-bold text-white tracking-wide">SABAR</h1>
        </div>

        <button @click="isSidebarOpen = !isSidebarOpen" class="text-white hover:bg-white/10 p-2 rounded-lg transition">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"
            class="size-7">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
          </svg>
        </button>
      </div>

      <!-- KANAN -->
      <div class="flex items-center gap-3 cursor-pointer hover:opacity-80 transition"
        @click="goTo('/dashboard-admin/profile')">
        <div class="w-10 h-10 rounded-full bg-white overflow-hidden">
          <img src="https://i.pravatar.cc/100" class="w-full h-full object-cover" />
        </div>

        <div class="text-white font-medium">
          Admin Akademik
        </div>
      </div>
    </header>

    <!-- ═══════════════════════════════════════
         SIDEBAR
    ════════════════════════════════════════ -->
    <aside :class="[
      'fixed left-0 top-18 bottom-0 bg-[#c8d8ee] flex flex-col justify-between transition-all duration-300 overflow-hidden',
      isSidebarOpen ? 'w-62.5' : 'w-20'
    ]">

      <!-- MENU -->
      <div class="p-3 space-y-1">

        <!-- DASHBOARD -->
        <div @click="goTo('/dashboard-admin')" :class="isActive('/dashboard-admin')
          ? 'bg-[#1f3c93] text-white'
          : 'text-gray-700 hover:bg-[#b8c9e2]'"
          class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 transition">

          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
            stroke="currentColor" class="size-5 shrink-0">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3.75 3v11.25M3.75 18.75h16.5M9.75 3v11.25M15.75 3v11.25" />
          </svg>

          <span v-if="isSidebarOpen">
            Dashboard
          </span>

        </div>

        <!-- AKADEMIK -->
        <div>

          <!-- HEADER -->
          <div @click="openAkademik = !openAkademik"
            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 text-gray-700 hover:bg-[#b8c9e2] transition">

            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
              stroke="currentColor" class="size-5 shrink-0">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 14 9-3m0 0-3 3m3-3 3 3M3 20.25h18" />
            </svg>

            <span v-if="isSidebarOpen" class="flex-1">
              Akademik
            </span>

            <svg v-if="isSidebarOpen" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
              stroke-width="1.5" stroke="currentColor" :class="[
                'size-4 transition-transform',
                openAkademik ? 'rotate-180' : ''
              ]">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25 12 15.75 4.5 8.25" />
            </svg>

          </div>

          <!-- SUB MENU -->
          <div v-if="openAkademik && isSidebarOpen" class="ml-3 mt-1 space-y-1">

            <div v-for="item in [
              {
                path: '/dashboard-admin/tahun_akademik',
                label: 'Tahun Akademik',
                icon: '📅'
              },
              {
                path: '/dashboard-admin/presensi',
                label: 'Presensi',
                icon: '📝'
              },
              {
                path: '/dashboard-admin/kurikulum',
                label: 'Kurikulum',
                icon: '📚'
              },
              {
                path: '/dashboard-admin/kelas',
                label: 'Kelas',
                icon: '🏫'
              },
              {
                path: '/dashboard-admin/nilai',
                label: 'Nilai',
                icon: '📊'
              },
              {
                path: '/dashboard-admin/khs',
                label: 'KHS',
                icon: '📄'
              },
              {
                path: '/dashboard-admin/dosen',
                label: 'Dosen',
                icon: '👨‍🏫'
              }
            ]" :key="item.path" @click="goTo(item.path)" :class="isActive(item.path)
              ? 'bg-[#1f3c93] text-white'
              : 'text-gray-700 hover:bg-[#b8c9e2]'"
              class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">

              <span>{{ item.icon }}</span>
              <span>{{ item.label }}</span>

            </div>

          </div>

        </div>

        <!-- MAHASISWA -->
        <div>

          <div @click="openMahasiswa = !openMahasiswa"
            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 text-gray-700 hover:bg-[#b8c9e2] transition">

            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
              stroke="currentColor" class="size-5 shrink-0">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M15 19.128a9.38 9.38 0 0 0 2.625.372A9.337 9.337 0 0 0 21 18.75" />
            </svg>

            <span v-if="isSidebarOpen" class="flex-1">
              Mahasiswa
            </span>

            <svg v-if="isSidebarOpen" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
              stroke-width="1.5" stroke="currentColor" :class="[
                'size-4 transition-transform',
                openMahasiswa ? 'rotate-180' : ''
              ]">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25 12 15.75 4.5 8.25" />
            </svg>

          </div>

          <div v-if="openMahasiswa && isSidebarOpen" class="ml-3 mt-1 space-y-1">

            <div @click="goTo('/dashboard-admin/peserta_kelas')" :class="isActive('/dashboard-admin/peserta_kelas')
              ? 'bg-[#1f3c93] text-white'
              : 'text-gray-700 hover:bg-[#b8c9e2]'"
              class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">

              <span>👨‍🎓</span>
              <span>Peserta Kelas</span>

            </div>

          </div>

        </div>

      </div>

      <!-- LOGOUT -->
      <div @click="showLogoutPopup = true"
        class="border-t border-[#b0c4de] flex cursor-pointer items-center gap-3 p-5 text-gray-700 hover:bg-[#9fb5d6] hover:text-red-500 transition">

        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
          class="size-5">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6A2.25 2.25 0 0 0 5.25 5.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
        </svg>

        <span v-if="isSidebarOpen">
          Keluar
        </span>

      </div>

    </aside>

    <!-- CONTENT -->
    <main :class="[
      'pt-22.5 p-6 transition-all duration-300',
      isSidebarOpen ? 'ml-62.5' : 'ml-20'
    ]">

      <!-- DASHBOARD HOME -->
      <div v-if="route.path === '/dashboard-admin'">

        <h1 class="text-3xl font-bold mb-1">Dashboard</h1>
        <p class="text-gray-600 mb-6">Selamat Datang Admin Akademik</p>

        <!-- STAT CARDS -->
        <div class="grid grid-cols-3 gap-4 mb-6">

          <!-- Total Pegawai -->
          <div class="bg-white rounded-xl p-4 flex items-center gap-4 shadow-[0_8px_25px_rgba(59,130,246,0.25)]">
            <div class="bg-[#a9c3e3] p-4 rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                stroke="currentColor" class="size-7 text-[#1f3c93]">
                <path stroke-linecap="round" stroke-linejoin="round"
                  d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
              </svg>
            </div>
            <div>
              <p class="text-sm font-semibold text-gray-500">TOTAL PEGAWAI</p>
              <h2 class="text-3xl font-bold">{{ totalPegawai }}</h2>
            </div>
          </div>

          <!-- Total Mahasiswa -->
          <div class="bg-white rounded-xl p-4 flex items-center gap-4 shadow-[0_8px_25px_rgba(59,130,246,0.25)]">
            <div class="bg-[#a9c3e3] p-4 rounded-lg">
            <i class="fi fi-rr-student text-3xl text-[#1f3c93]"></i>
            </div>
            <div>
              <p class="text-sm font-semibold text-gray-500">TOTAL MAHASISWA</p>
              <h2 class="text-3xl font-bold">{{ totalMahasiswa }}</h2>
            </div>
          </div>

          <!-- Total Dosen -->
          <div class="bg-white rounded-xl p-4 flex items-center gap-4 shadow-[0_8px_25px_rgba(59,130,246,0.25)]">
            <div class="bg-[#a9c3e3] p-4 rounded-lg">
              <svg data-slot="icon" fill="none" stroke-width="1.5" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" class="size-8 text-[#1f3c93]"></path>
</svg>
            </div>
            <div>
              <p class="text-sm font-semibold text-gray-500">TOTAL DOSEN</p>
              <h2 class="text-3xl font-bold">{{ totalDosen }}</h2>
            </div>
          </div>

        </div>

        <!-- CONTENT GRID -->
        <div class="grid grid-cols-4 gap-4">

          <!-- TABLE -->
          <div class="col-span-3 bg-white rounded-xl p-5 shadow-[0_8px_25px_rgba(59,130,246,0.25)]">

            <h2 class="text-xl font-bold mb-4">Data Mahasiswa</h2>

            <table class="w-full text-sm border-collapse">
              <thead>
                <tr class="bg-gray-50">
                  <th class="py-3 px-2 text-left">No</th>
                  <th class="py-3 px-2 text-left">NIM</th>
                  <th class="py-3 px-2 text-left">Nama</th>
                  <th class="py-3 px-2 text-left">Kelas</th>
                  <th class="py-3 px-2 text-left">Prodi</th>
                  <th class="py-3 px-2 text-left">Angkatan</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="akunList.length === 0">
                  <td colspan="6" class="py-6 text-center text-gray-400">Tidak ada data</td>
                </tr>
                <tr v-for="(item, index) in akunList" :key="item.id" class="hover:bg-gray-50">
                  <td class="py-3 px-2">{{ (currentPage - 1) * perPage + index + 1 }}</td>
                  <td class="py-3 px-2">{{ item.nim }}</td>
                  <td class="py-3 px-2">{{ item.name }}</td>
                  <td class="py-3 px-2">{{ item.kelas }}</td>
                  <td class="py-3 px-2">{{ item.prodi }}</td>
                  <td class="py-3 px-2">{{ item.angkatan }}</td>
                </tr>
              </tbody>
            </table>

            <!-- PAGINATION -->
            <div class="flex justify-end mt-5 pt-4">

              <div class="flex items-center gap-2">

                <button @click="prevPage" :disabled="currentPage === 1"
                  class="px-3 py-1 border rounded-lg bg-white hover:bg-gray-100 disabled:opacity-50">
                  Previous
                </button>

                <template v-for="p in pages" :key="p">

                  <span v-if="p === '...'">...</span>

                  <button v-else @click="goToPage(p as number)" class="w-8 h-8 rounded-lg" :class="currentPage === p
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 hover:bg-gray-200'
                    ">
                    {{ p }}
                  </button>

                </template>

                <button @click="nextPage" :disabled="currentPage === totalPages"
                  class="px-3 py-1 border rounded-lg bg-white hover:bg-gray-100 disabled:opacity-50">
                  Next
                </button>

              </div>

            </div>

          </div>

          <!-- SIDE PANEL -->
          <div class="space-y-4">

            <!-- Tahun Akademik -->

            <div class="bg-white rounded-xl p-5 shadow">

              <!-- HEADER -->
              <div class="flex items-center justify-between mb-3">

                <h2 class="font-bold">
                  Tahun Akademik
                </h2>

                <button @click="goTo('/dashboard-admin/dashboard_tahunakademik')"
                  class="text-blue-600 hover:text-blue-800">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="size-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M2.25 12s3.75-7.5 9.75-7.5 9.75 7.5 9.75 7.5-3.75 7.5-9.75 7.5S2.25 12 2.25 12Z" />
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 15.75A3.75 3.75 0 1 0 12 8.25a3.75 3.75 0 0 0 0 7.5Z" />
                  </svg>
                </button>

              </div>

              <!-- DATA -->
              <div v-for="item in tahunAkademik.slice(0, 4)" :key="item.id" class="bg-gray-50 rounded-lg p-3 mb-2">
                {{ formatYear(item.tahun_awal) }} -
                {{ formatYear(item.tahun_akhir) }}
              </div>

              <div v-if="tahunAkademik.length === 0" class="text-gray-400 text-sm">
                Tidak ada data
              </div>

            </div>

            <!-- Kurikulum -->
            <div class="bg-white rounded-xl p-5 shadow">

              <!-- HEADER -->
              <div class="flex items-center justify-between mb-3">

                <h2 class="font-bold">
                  Kurikulum
                </h2>

                <button @click="goTo('/dashboard-admin/dashboard_kurikulum')" class="text-blue-600 hover:text-blue-800">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="size-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M2.25 12s3.75-7.5 9.75-7.5 9.75 7.5 9.75 7.5-3.75 7.5-9.75 7.5S2.25 12 2.25 12Z" />
                    <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 15.75A3.75 3.75 0 1 0 12 8.25a3.75 3.75 0 0 0 0 7.5Z" />
                  </svg>
                </button>

              </div>

              <!-- DATA -->
              <div v-for="item in kurikulum.slice(0, 4)" :key="item.id" class="bg-gray-50 rounded-lg p-3 mb-2">
                {{ item.name }}
              </div>

              <div v-if="kurikulum.length === 0" class="text-gray-400 text-sm">
                Tidak ada data
              </div>

            </div>

          </div>

        </div>

      </div>

      <!-- ROUTER VIEW -->
      <div v-else>
        <router-view />
      </div>

    </main>

    <!-- POPUP -->
    <KonfirmasiKeluar v-if="showLogoutPopup" @close="showLogoutPopup = false" @confirm="handleLogout" />

  </div>
</template>