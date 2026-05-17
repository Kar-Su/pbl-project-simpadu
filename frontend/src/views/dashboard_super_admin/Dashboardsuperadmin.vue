<script setup lang="ts">
import { onMounted, ref } from "vue"
import { useRouter, useRoute } from "vue-router"
import Konfirmasi_keluar from "./akademik/konfirmasi_keluar.vue"

const showLogoutPopup = ref(false)
const router = useRouter()
const route = useRoute()
const openAkademik = ref(false)
const isSidebarOpen = ref(true)
const user = ref<any>(null)
const kurikulum = ref<any[]>([])
const tahunAkademik = ref<any[]>([])
const akunList = ref<any[]>([])
const totalAkun = ref(0)

const isActive = (path: string) => route.path === path

const getHeaders = () => ({
  "Authorization": `Bearer ${localStorage.getItem("token")}`
})

const getKurikulum = async () => {
  try {
    const res = await fetch(`/api/kurikulum/`, { headers: getHeaders() })
    const data = await res.json()
    kurikulum.value = data.data ?? []
  } catch (err) { console.error(err) }
}

const getTahunAkademik = async () => {
  try {
    const res = await fetch(`/api/tahun-akademik`, { headers: getHeaders() })
    const data = await res.json()
    tahunAkademik.value = data.data ?? []
  } catch (err) { console.error(err) }
}

const getTotalAkun = async () => {
  try {
    const res = await fetch(`/api/user/count`, { headers: getHeaders() })
    const data = await res.json()
    totalAkun.value = data.data ?? 0  // ← data.data langsung angka 8
  } catch (err) { console.error(err) }
}
const getAkun = async () => {
  try {
    //  const roles = ['admin-akademik', 'admin-keuangan', 'admin-mahasiswa', 'admin-pegawai', 'dosen', 'dummy', 'mahasiswa']
    const roles = ['dosen', 'dummy', 'mahasiswa']
    const results = await Promise.all(
      roles.map(role =>
        fetch(`/api/user/role/${role}`, { headers: getHeaders() })
          .then(r => r.json())
          .then(d => Array.isArray(d.data) ? d.data : [])
      )
    )
    akunList.value = results.flat()
    console.log("TOTAL USER:", akunList.value.length)
  } catch (err) { console.error(err) }
}

const handleLogout = () => {
  localStorage.clear()
  router.push("/")
}

onMounted(() => {
  getKurikulum()
  getTahunAkademik()
  getAkun()
  getTotalAkun()
})
</script>

<template>
  <div class="flex h-screen bg-gray-100">

    <!-- SIDEBAR -->
    <aside
      :class="isSidebarOpen ? 'w-60' : 'w-16'"
      class="bg-[#2f4a8a] text-white flex flex-col p-4 transition-all duration-300"
    >
      <h1 class="text-lg font-bold mb-6">Logo</h1>

      <nav class="space-y-2 text-sm">
        <div
          @click="router.push('/dashboard-superadmin')"
          :class="isActive('/dashboard-superadmin') ? 'bg-blue-500' : 'hover:bg-blue-700'"
          class="px-3 py-2 rounded cursor-pointer"
        >
          Dashboard
        </div>

        <div
          @click="openAkademik = !openAkademik"
          class="px-3 py-2 rounded cursor-pointer hover:bg-blue-700 flex justify-between items-center"
        >
          <span>Akademik</span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-5 h-5 transition-transform duration-300"
            :class="{ 'rotate-180': openAkademik }"
            viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round"
          >
            <polyline points="6 9 12 15 18 9" />
          </svg>
        </div>

        <div
          class="ml-4 flex flex-col space-y-1 border-l border-blue-300 pl-3 overflow-hidden transition-all duration-300"
          :class="openAkademik ? 'max-h-40 opacity-100' : 'max-h-0 opacity-0'"
        >
          <div
            @click="router.push('/dashboard-superadmin/akun')"
            :class="isActive('/dashboard-superadmin/akun') ? 'bg-blue-500' : 'hover:bg-blue-700'"
            class="px-3 py-2 rounded cursor-pointer"
          >
            Akun
          </div>
          <div
            @click="router.push('/dashboard-superadmin/reset_password')"
            :class="isActive('/dashboard-superadmin/reset_password') ? 'bg-blue-500' : 'hover:bg-blue-700'"
            class="px-3 py-2 rounded cursor-pointer"
          >
            Reset Password
          </div>
        </div>
      </nav>

      <button
        @click="showLogoutPopup = true"
        class="mt-auto text-left text-xs opacity-80 hover:opacity-100"
      >
        Keluar
      </button>
    </aside>

    <!-- MAIN -->
    <div class="flex-1 flex flex-col overflow-hidden">

      <!-- TOPBAR -->
      <div class="bg-[#2f4a8a] text-white px-6 py-3 flex justify-between items-center">
        <div class="text-lg font-semibold">Dashboard Super Admin</div>
        <div class="text-sm">{{ user?.name || 'Admin' }}</div>
      </div>

      <!-- CONTENT -->
      <div class="flex-1 overflow-auto bg-[#eef3fb]">

        <!-- Dashboard utama -->
        <div v-if="route.path === '/dashboard-superadmin'" class="p-6">

          <div class="bg-white rounded-xl shadow p-5 w-72 mb-6">
            <h2 class="text-gray-500 text-sm font-semibold">TOTAL AKUN</h2>
            <div class="text-4xl font-bold text-[#2f4a8a] mt-2">{{ totalAkun }}</div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

            <div class="lg:col-span-2 bg-white rounded-xl shadow p-5">
              <h2 class="text-xl font-bold mb-4">Data Akun</h2>
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b text-left">
                    <th class="py-2">No</th>
                    <th class="py-2">Nama</th>
                    <th class="py-2">Email</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, index) in akunList" :key="item.id" class="border-b">
                    <td class="py-3">{{ index + 1 }}</td>
                    <td class="py-3">{{ item.nama || item.name }}</td>
                    <td class="py-3">{{ item.email }}</td>
                  </tr>
                  <tr v-if="akunList.length === 0">
                    <td colspan="3" class="py-4 text-center text-gray-400">Tidak ada data</td>
                  </tr>
                </tbody>
              </table>
            </div>

<div class="space-y-6">
              <div class="bg-white rounded-xl shadow p-5">
                <h2 class="text-lg font-bold mb-4">Tahun Akademik</h2>
                <div v-for="item in tahunAkademik" :key="item.id" class="border-b py-2">
                  {{ item.tahun_awal?.slice(0, 4) }} - {{ item.tahun_akhir?.slice(0, 4) }} ({{ item.tipee_semester }})
                </div>
                <div v-if="tahunAkademik.length === 0" class="text-gray-400 text-sm">Tidak ada data</div>
              </div>

<div class="bg-white rounded-xl shadow p-5">
  <h2 class="text-lg font-bold mb-4">Kurikulum</h2>
  <div v-for="item in kurikulum" :key="item.id" class="border-b py-2">
    {{ item.nama || item.name || item.kode }}
  </div>
  <div v-if="kurikulum.length === 0" class="text-gray-400 text-sm">Tidak ada data</div>
</div>
            </div>

          </div>
        </div>

        <!-- Halaman child (Akun / Reset Password) -->
        <router-view v-else />

      </div>
    </div>

  </div>

  <Konfirmasi_keluar
    v-if="showLogoutPopup"
    @close="showLogoutPopup = false"
    @confirm="handleLogout"
  />
</template>