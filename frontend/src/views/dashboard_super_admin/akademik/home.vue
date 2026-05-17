<script setup lang="ts">
import { onMounted, ref } from "vue"

const BASE_URL = "http://localhost"

const students = ref<any[]>([])
const kurikulum = ref<any[]>([])
const tahunAkademik = ref<any[]>([])
const totalAkun = ref<number>(0)

  
const getStudents = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/students/`)
    const data = await res.json()
    console.log(data)
    students.value = data?.data || []
  } catch (err) {
    console.error(err)
  }
}

const getKurikulum = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/kurikulum/`)
    const data = await res.json()
    console.log(data)
    kurikulum.value = data?.data || []
  } catch (err) {
    console.error(err)
  }
}

const getTahunAkademik = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik/`)
    const data = await res.json()
    console.log(data)
    tahunAkademik.value = data?.data || []
  } catch (err) {
    console.error(err)
  }
}

const getTotalAkun = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/total-akun/`)
    const data = await res.json()
    console.log(data)
    totalAkun.value = data?.data?.total || 0
  } catch (err) {
    console.error(err)
  }
}

onMounted(() => {
  getStudents()
  getKurikulum()
  getTahunAkademik()
  getTotalAkun()
})
</script>

<template>
  <div>

    <h2 class="text-lg font-bold mb-1">
      Dashboard
    </h2>

    <p class="text-xs text-gray-500 mb-5">
      Selamat Datang Super Admin
    </p>

    <div class="bg-white p-4 rounded shadow w-62.5 mb-6">
      <p class="text-xs text-gray-500">
        TOTAL AKUN
      </p>

      <h2 class="text-xl font-bold text-gray-800">
        {{ totalAkun }}
      </h2>
    </div>

    <div class="grid grid-cols-3 gap-5">

      <!-- TABLE -->
      <div class="col-span-2 bg-white p-4 rounded shadow">

        <h3 class="text-sm font-semibold mb-3">
          Data Akun
        </h3>

        <table class="w-full text-xs">

          <thead>
            <tr class="bg-gray-100 text-left">
              <th class="p-2">No</th>
              <th class="p-2">Nama</th>
              <th class="p-2">Email</th>
            </tr>
          </thead>

          <tbody>

            <tr
              v-for="(s, i) in students"
              :key="i"
              class="border-b"
            >
              <td class="p-2">{{ i + 1 }}</td>
              <td class="p-2">{{ s.name }}</td>
              <td class="p-2">{{ s.email }}</td>
            </tr>

            <tr v-if="students.length === 0">
              <td
                colspan="3"
                class="text-center p-3 text-gray-400"
              >
                Data kosong
              </td>
            </tr>

          </tbody>

        </table>
      </div>

      <!-- SIDE -->
      <div class="space-y-5">

        <!-- TAHUN -->
        <div class="bg-white p-4 rounded shadow">

          <h3 class="text-xs font-semibold mb-2">
            Tahun Akademik
          </h3>

          <div
            v-for="(t, i) in tahunAkademik"
            :key="i"
            class="flex justify-between text-xs border-b py-1"
          >
            <span>{{ t.name }}</span>
          </div>

        </div>

        <!-- KURIKULUM -->
        <div class="bg-white p-4 rounded shadow">

          <h3 class="text-xs font-semibold mb-2">
            Kurikulum
          </h3>

          <div
            v-for="(k, i) in kurikulum"
            :key="i"
            class="flex justify-between text-xs border-b py-1"
          >
            <span>{{ k.name }}</span>
          </div>

        </div>

      </div>

    </div>

  </div>
</template>