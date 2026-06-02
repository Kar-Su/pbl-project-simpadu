<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= FILTER =================
const jurusan = ref<string>("")
const prodi = ref<string>("")
const tahunAkademik = ref<string>("")

// ================= DATA =================
const kelasData = ref<any[]>([])

// ================= FETCH =================
const getKelas = async () => {
  try {
    const res = await fetch("/api/kelas/", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })

    const json = await res.json()

    console.log("KELAS RAW:", json)

    kelasData.value =
      json?.data?.items ??
      json?.data?.data ??
      json?.data ??
      json?.result ??
      []

    if (kelasData.value.length > 0) {
      console.log(
        "KELAS ITEM SAMPLE:",
        JSON.stringify(kelasData.value[0], null, 2)
      )
    }
  } catch (err) {
    console.error("GET KELAS ERROR:", err)
  }
}

// ================= DROPDOWN OTOMATIS =================
const jurusanList = computed(() => {
  const map = new Map()

  kelasData.value.forEach((item) => {
    const jurusan = item?.prodi?.jurusan

    if (jurusan?.id && !map.has(jurusan.id)) {
      map.set(jurusan.id, {
        id: String(jurusan.id),
        name: jurusan.name,
      })
    }
  })

  return Array.from(map.values())
})

const prodiList = computed(() => {
  const map = new Map()

  kelasData.value.forEach((item) => {
    const prodi = item?.prodi

    if (prodi?.id && !map.has(prodi.id)) {
      map.set(prodi.id, {
        id: String(prodi.id),
        name: prodi.name,
      })
    }
  })

  return Array.from(map.values())
})

const tahunAkademikList = computed(() => {
  const map = new Map()

  kelasData.value.forEach((item) => {
    const tahun = item?.tahun_akademik

    if (tahun?.id && !map.has(tahun.id)) {
      map.set(tahun.id, {
        id: String(tahun.id),
        label: `${tahun.tipe_semester} ${tahun.tahun_awal} - ${tahun.tahun_akhir}`,
      })
    }
  })

  return Array.from(map.values())
})

// ================= HELPER =================
const getJurusanName = (item: any): string => {
  return item?.prodi?.jurusan?.name ?? "-"
}

const getProdiName = (item: any): string => {
  return item?.prodi?.name ?? "-"
}

const getTahunName = (item: any): string => {
  const tahun = item?.tahun_akademik

  if (!tahun) return "-"

  return `${tahun.tipe_semester} ${tahun.tahun_awal} - ${tahun.tahun_akhir}`
}

// ================= FILTER COMPUTED =================
const filteredData = computed(() => {
  return kelasData.value.filter((item) => {
    const jurusanId = String(item?.prodi?.jurusan?.id ?? "")
    const prodiId = String(item?.prodi?.id ?? "")
    const tahunId = String(item?.tahun_akademik?.id ?? "")

    return (
      (!jurusan.value || jurusanId === jurusan.value) &&
      (!prodi.value || prodiId === prodi.value) &&
      (!tahunAkademik.value || tahunId === tahunAkademik.value)
    )
  })
})

// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(5)

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredData.value.length / perPage.value))
)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return filteredData.value.slice(start, start + perPage.value)
})

const rowNumber = (index: number) =>
  (currentPage.value - 1) * perPage.value + index + 1

const pages = computed((): (number | "...")[] => {
  const total = totalPages.value
  const current = currentPage.value
  const result: (number | "...")[] = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      result.push(i)
    }
    return result
  }

  result.push(1)

  if (current > 3) {
    result.push("...")
  }

  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)

  for (let i = start; i <= end; i++) {
    result.push(i)
  }

  if (current < total - 2) {
    result.push("...")
  }

  result.push(total)

  return result
})

const goToPage = (page: number): void => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
}

const prevPage = (): void => goToPage(currentPage.value - 1)
const nextPage = (): void => goToPage(currentPage.value + 1)

const pilihData = () => {
  currentPage.value = 1
}

// ================= ACTION =================
const tambahData = () => {
  router.push("/dashboard-admin/tambah_kelas")
}

const editData = (item: any) => {
  console.log("EDIT:", item)
}

const hapusData = (item: any) => {
  console.log("HAPUS:", item)
}

// ================= MOUNT =================
onMounted(() => {
  getKelas()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-500 font-medium mb-2">
      Akademik > Kelas
    </div>

    <!-- TITLE -->
    <h1 class="text-[42px] font-bold text-[#404040] leading-none">
      Kelas
    </h1>

    <p class="text-gray-500 text-sm mt-2 mb-6">
      Pengelolaan Data
    </p>

    <!-- CARD -->
    <div class="bg-white border border-[#d5e1f0] rounded-2xl min-h-[760px] shadow-[0_4px_12px_rgba(0,0,0,0.08)] overflow-hidden">

      <!-- HEADER -->
      <div class="px-5 pt-4">
        <h2 class="text-[36px] font-semibold text-[#505050]">
          Data Kelas
        </h2>
      </div>

      <!-- FILTER -->
      <div class="px-5 pt-5 flex items-center gap-4 flex-wrap">

        <!-- JURUSAN -->
        <select
          v-model="jurusan"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="">Pilih Jurusan</option>

          <option
            v-for="j in jurusanList"
            :key="j.id"
            :value="j.id"
          >
            {{ j.name }}
          </option>
        </select>

        <!-- PRODI -->
        <select
          v-model="prodi"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="">Pilih Prodi</option>

          <option
            v-for="p in prodiList"
            :key="p.id"
            :value="p.id"
          >
            {{ p.name }}
          </option>
        </select>

        <!-- TAHUN AKADEMIK -->
        <select
          v-model="tahunAkademik"
          class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
        >
          <option value="">Pilih Tahun Akademik</option>

          <option
            v-for="t in tahunAkademikList"
            :key="t.id"
            :value="t.id"
          >
            {{ t.label }}
          </option>
        </select>

        <!-- BUTTON PILIH -->
        <button
          @click="pilihData"
          class="h-[54px] px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          Pilih
        </button>

        <!-- BUTTON TAMBAH -->
        <button
          @click="tambahData"
          class="h-[54px] px-6 bg-[#29479d] hover:bg-[#1d377f] rounded-xl text-white font-semibold text-[18px] shadow-md transition"
        >
          + Tambah
        </button>

      </div>

      <!-- TABLE -->
      <div class="px-5 pt-8">
        <table class="w-full border-separate border-spacing-y-5">

          <thead>
            <tr class="text-left text-gray-600">
              <th class="text-[18px] font-semibold">No</th>
              <th class="text-[18px] font-semibold">Nama Kelas</th>
              <th class="text-[18px] font-semibold">Jurusan</th>
              <th class="text-[18px] font-semibold">Prodi</th>
              <th class="text-[18px] font-semibold">Tahun Akademik</th>
              <th class="text-[18px] font-semibold text-center">Aksi</th>
            </tr>
          </thead>

          <tbody>

            <tr
              v-for="(item, index) in paginatedData"
              :key="item.id"
              class="text-[#505050]"
            >

              <td class="text-[18px]">
                {{ rowNumber(index) }}
              </td>

              <td class="text-[18px] font-medium">
                {{ item.name ?? "-" }}
              </td>

              <td class="text-[18px] font-medium">
                {{ getJurusanName(item) }}
              </td>

              <td class="text-[18px] font-medium">
                {{ getProdiName(item) }}
              </td>

              <td class="text-[18px] font-medium">
                {{ getTahunName(item) }}
              </td>

              <td class="flex items-center justify-center gap-3">

                <button
                  @click="editData(item)"
                  class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-5 py-2 rounded-xl text-[16px] font-semibold shadow-md transition"
                >
                  ✎ Edit
                </button>

                <button
                  @click="hapusData(item)"
                  class="bg-[#ef4d43] hover:bg-[#d93d34] text-white px-5 py-2 rounded-xl text-[16px] font-semibold shadow-md transition"
                >
                  🗑 Hapus
                </button>

              </td>

            </tr>

            <tr v-if="paginatedData.length === 0">
              <td
                colspan="6"
                class="text-center text-gray-400 py-10 text-[16px]"
              >
                Tidak ada data kelas
              </td>
            </tr>

          </tbody>

        </table>
      </div>

      <!-- FOOTER -->
      <div class="flex items-center justify-between px-5 pt-10 pb-5">

        <select
          v-model="perPage"
          @change="currentPage = 1"
          class="w-[90px] h-[42px] border border-gray-300 rounded-lg px-3 text-sm outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select>

        <div class="flex items-center gap-2">

          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="px-3 py-1 border rounded-lg disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Previous
          </button>

          <template v-for="p in pages" :key="p">

            <span
              v-if="p === '...'"
              class="px-1 text-gray-400"
            >
              ...
            </span>

            <button
              v-else
              @click="goToPage(p as number)"
              class="w-8 h-8 rounded-lg"
              :class="currentPage === p
                ? 'bg-blue-500 text-white'
                : 'bg-gray-100'"
            >
              {{ p }}
            </button>

          </template>

          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 border rounded-lg disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Next
          </button>

        </div>

      </div>

    </div>

  </div>
</template>

<style scoped>
table {
  border-collapse: separate;
}
</style>