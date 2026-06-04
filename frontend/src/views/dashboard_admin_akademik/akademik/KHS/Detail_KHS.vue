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
// DATA KHS
// ─────────────────────────────────────────────
interface KhsItem {
  id: number
  nim: string
  nama_mahasiswa: string
  sks: number
  ips: number
  ipk: number
}

const khsList = ref<KhsItem[]>([])

// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const currentPage = ref<number>(1)
const perPage = ref<number>(5)
const totalItems = ref<number>(0)

  const BASE_URL = "https://be.karlearn.site"

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
// HIT API KHS
// Endpoint :
// GET /api/khs?page=1&per_page=5
// ─────────────────────────────────────────────
const getKhs = async (): Promise<void> => {
  try {
    const res = await fetch(
      `/api/khs?page=${currentPage.value}&per_page=${perPage.value}`,
      {
        headers: getHeaders(),
      }
    )

    const data = await res.json()

    khsList.value = data.data ?? []
    totalItems.value = data.meta?.total ?? data.total ?? 0
  } catch (err) {
    console.error("getKhs:", err)
  }
}

// ─────────────────────────────────────────────
// PAGINATION
// ─────────────────────────────────────────────
const goToPage = (page: number): void => {
  if (page < 1 || page > totalPages.value) return

  currentPage.value = page
  getKhs()
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
  getKhs()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
      <span>KHS</span>
      <span>›</span>
      <span class="text-gray-700">Detail KHS</span>
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold leading-none text-[#333]">
      Kartu Hasil Studi
    </h1>

    <p class="mt-3 text-gray-500">
      Data hasil studi mahasiswa
    </p>

    <!-- CARD -->
    <div
      class="mt-8 rounded-2xl border border-[#d8e1f0] bg-white p-5 shadow-sm"
    >

      <!-- HEADER -->
      <h2 class="mb-10 text-[32px] font-bold text-[#444]">
        Data KHS
      </h2>

      <!-- TABLE -->
      <div class="overflow-x-auto">
        <table class="w-full">

          <!-- HEAD -->
          <thead>
            <tr class="text-left text-[15px] font-semibold text-[#555]">
              <th class="pb-4">No</th>
              <th class="pb-4">NIM</th>
              <th class="pb-4">Nama Mahasiswa</th>
              <th class="pb-4">SKS</th>
              <th class="pb-4">IPS</th>
              <th class="pb-4">IPK</th>
            </tr>
          </thead>

          <!-- BODY -->
          <tbody>

            <!-- EMPTY -->
            <tr v-if="khsList.length === 0">
              <td
                colspan="6"
                class="py-12 text-center text-gray-400"
              >
                Tidak ada data
              </td>
            </tr>

            <!-- DATA -->
            <tr
              v-for="(item, index) in khsList"
              :key="item.id"
              class="text-[15px] text-[#444]"
            >
              <td class="py-4">
                {{ (currentPage - 1) * perPage + index + 1 }}
              </td>

              <td class="py-4 font-medium">
                {{ item.nim }}
              </td>

              <td class="py-4 font-medium">
                {{ item.nama_mahasiswa }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.sks }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.ips }}
              </td>

              <td class="py-4 font-semibold">
                {{ item.ipk }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAGINATION -->
      <div class="mt-72 flex items-center justify-between">

        <!-- SELECT -->
        <select
          v-model.number="perPage"
          @change="() => { currentPage = 1; getKhs() }"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select>

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