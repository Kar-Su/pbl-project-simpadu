<script setup lang="ts">
import { ref, computed, onMounted } from "vue"

// ─────────────────────────────────────────────
// HELPER HEADER
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})
// ─────────────────────────────────────────────
// DATA DOSEN
// ─────────────────────────────────────────────
interface DosenItem {
  id: number
  nip: string
  nama_dosen: string
  jurusan: string
}

const dosenList = ref<DosenItem[]>([])
const allDosen = ref<DosenItem[]>([])

// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const currentPage = ref<number>(1)
const perPage = ref<number>(10)
const totalItems = ref<number>(0)

const totalPages = computed<number>(() =>
  Math.max(1, Math.ceil(totalItems.value / perPage.value))
)

const pages = computed<(number | string)[]>(() => {
  const total = totalPages.value
  const cur = currentPage.value

  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const result: (number | string)[] = [1, 2]

  if (cur > 4) result.push("...")

  for (
    let i = Math.max(3, cur - 1);
    i <= Math.min(total - 2, cur + 1);
    i++
  ) {
    result.push(i)
  }

  if (cur < total - 3) result.push("...")

  result.push(total - 1, total)

  return [...new Set(result)]
})


// ─────────────────────────────────────────────
// GET DOSEN DARI USERS
// FILTER ROLE DOSEN
// ─────────────────────────────────────────────
const getDosen = async (): Promise<void> => {
  try {

    let page = 1
    let lastPage = 1
    let allUsers: any[] = []

    do {
      const BASE_URL = 'https://be.karlearn.site'
      const res = await fetch(`${BASE_URL}/api/users?page=${page}`, {
        headers: getHeaders(),
      })

      const data = await res.json()

      const items = data.data.items ?? []

      allUsers = [...allUsers, ...items]

      lastPage = data.data.pagination?.total_pages ?? 1

      page++

    } while (page <= lastPage)

    // FILTER DOSEN
    const dosenOnly = allUsers.filter((item: any) => {
      const role = item.role_name?.toLowerCase()?.trim()
      return role?.includes("dosen")
    })

    // MAP DATA
    allDosen.value = dosenOnly.map((item: any) => ({
      id: item.id,
      nip: item.detail?.nip ?? "-",
      nama_dosen: item.name ?? "-",
      jurusan: item.detail?.jurusan ?? "-",
    }))

    // TOTAL DATA
    totalItems.value = allDosen.value.length

    // PAGINATION FRONTEND
    const start = (currentPage.value - 1) * perPage.value
    const end = start + perPage.value

    dosenList.value = allDosen.value.slice(start, end)

  } catch (err) {
    console.error("getDosen:", err)
  }
}


// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const goToPage = (page: number): void => {
  if (page < 1 || page > totalPages.value) return

  currentPage.value = page
  getDosen()
}

const prevPage = (): void => {
  goToPage(currentPage.value - 1)
}

const nextPage = (): void => {
  goToPage(currentPage.value + 1)
}



// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted((): void => {
  getDosen()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
      <span>Akademik</span>
      <span>›</span>
      <span class="text-gray-700">Dosen</span>
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold leading-none text-[#333]">
      Dosen
    </h1>

    <p class="mt-3 text-gray-500">
      Data Dosen yang aktif
    </p>

    <!-- CARD -->
<div
  class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] px-8 py-5"
>

      <!-- HEADER -->
      <h2 class="mb-10 text-[32px] font-bold text-[#444]">
        Data Dosen
      </h2>

      <!-- TABLE -->
      <div class="overflow-x-auto">
        <table class="w-full">

          <!-- HEAD -->
          <thead>
            <tr class="text-left text-[15px] font-semibold text-[#555]">
              <th class="pb-4">No</th>
              <th class="pb-4">NIP</th>
              <th class="pb-4">Nama Dosen</th>
              <th class="pb-4">Jurusan</th>
            </tr>
          </thead>

          <!-- BODY -->
          <tbody>

            <!-- EMPTY -->
            <tr v-if="dosenList.length === 0">
              <td
                colspan="4"
                class="py-12 text-center text-gray-400"
              >
                Tidak ada data
              </td>
            </tr>

            <!-- DATA -->
            <tr
              v-for="(item, index) in dosenList"
              :key="item.id"
              class="text-[15px] text-[#444]"
            >
              <td class="py-4">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <td class="py-4 font-medium">
                {{ item.nip }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.nama_dosen }}
              </td>

              <td class="py-4 font-medium">
                {{ item.jurusan }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAGINATION -->
      <div class="mt-52 flex items-center justify-between">

        <!-- SELECT -->
        <!-- <select
          v-model.number="perPage"
          @change="() => { currentPage = 1; getDosen() }"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select> -->

        <!-- PAGINATION -->
        <div class="flex items-center gap-2">

          <!-- PREV -->
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="text-sm text-gray-400"
          >
            ← Previous
          </button>

          <!-- PAGE -->
          <template v-for="p in pages" :key="p">

            <span
              v-if="p === '...'"
              class="px-2 text-gray-400"
            >
              ...
            </span>

            <button
              v-else
              @click="goToPage(p as number)"
              class="flex h-9 w-9 items-center justify-center rounded-lg text-sm font-medium"
              :class="
                currentPage === p
                  ? 'bg-[#2447a8] text-white'
                  : 'text-gray-600 hover:bg-gray-100'
              "
            >
              {{ p }}
            </button>

          </template>

          <!-- NEXT -->
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="text-sm text-gray-600"
          >
            Next →
          </button>

        </div>
      </div>
    </div>
  </div>
</template>