<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= FILTER =================
const jurusan = ref<string>("")
const prodi = ref<string>("")
const tahunAkademik = ref<string>("")
const prodiMap = ref<Record<number, any>>({})
const tahunMap = ref<Record<number, any>>({})



// ================= DATA =================
const kelasData = ref<any[]>([])

// ================= PAGINATION =================
const currentPage = ref(1)
const perPage = ref(5)
const totalPages = ref(1)

// ================= FETCH =================
const getKelas = async () => {
  try {
    if (!prodi.value) {
      kelasData.value = []
      return
    }

    const BASE_URL = "https://be.karlearn.site"

    const prodiNama =
      prodiMap.value[Number(prodi.value)]?.name

    const url =
      `${BASE_URL}/api/kelas/prodi/${encodeURIComponent(prodiNama)}?page=${currentPage.value}`

    const res = await fetch(url, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
      },
    })

    const json = await res.json()

    console.log("KELAS RESPONSE:", json)

    kelasData.value = json?.data?.items ?? json?.data ?? []
    totalPages.value =
      json?.data?.pagination?.total_pages ?? 1

  } catch (err) {
    console.error("GET KELAS ERROR:", err)
    kelasData.value = []
  }
}

const getProdi = async () => {
  const res = await fetch("https://be.karlearn.site/api/prodi", {
    headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
  })
  const json = await res.json()
  const list = json?.data?.items ?? json?.data ?? []
  list.forEach((p: any) => {
    prodiMap.value[p.id] = p // p harus punya { id, name, jurusan: { id, name } }
  })
}

const getTahunAkademik = async () => {
  const res = await fetch("https://be.karlearn.site/api/tahun-akademik?per_page=10", {
    headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
  })
  const json = await res.json()
  console.log("TAHUN RESPONSE:", JSON.stringify(json, null, 2)) // ← lihat struktur lengkap
  
  const list = Array.isArray(json?.data) ? json.data : json?.data?.items ?? []
  console.log("LIST:", list) // ← apakah array-nya ada isinya?
  
  list.forEach((t: any) => {
    tahunMap.value[t.id] = t
  })
  console.log("TAHUN MAP:", tahunMap.value) // ← apakah map terisi?
}


// ================= DROPDOWN LIST =================
const jurusanList = computed(() => {
  const map = new Map()
  Object.values(prodiMap.value).forEach((p: any) => {
    const j = p?.jurusan
    if (j?.id && !map.has(j.id)) map.set(j.id, { id: String(j.id), name: j.name })
  })
  return Array.from(map.values())
})

const prodiList = computed(() => {
  return Object.values(prodiMap.value).map((p: any) => ({
    id: String(p.id),
    name: p.name,
    jurusanId: String(p.jurusan?.id ?? ""),
  }))
})

const tahunAkademikList = computed(() => {
  return Object.values(tahunMap.value).map((t: any) => ({
    id: String(t.id),
    label: `${new Date(t.tahun_awal).getFullYear()}/${new Date(t.tahun_akhir).getFullYear()}`,
  }))
})


// const tahunAkademikList = computed(() => {
//   const map = new Map()

//   kelasData.value.forEach((item) => {
//     const t = item?.tahun_akademik

//     if (t?.id && !map.has(t.id)) {
//       map.set(t.id, {
//         id: String(t.id),
//         label: `${t.tipe_semester} ${t.tahun_awal} - ${t.tahun_akhir}`,
//       })
//     }
//   })

//   return Array.from(map.values())
// })

// ================= WATCH =================
import { watch } from "vue"

watch(currentPage, () => {
  getKelas()
})
watch(jurusan, () => {
  prodi.value = ""
})

watch(prodi, () => {
  currentPage.value = 1
  getKelas()
})

// ================= HELPER (INI YANG KURANG) =================
const getJurusanName = (item: any) => {
  return (item.prodi?.jurusan?.name ?? "-").replace(/-/g, " ")
}

const getProdiName = (item: any) => {
  return (item.prodi?.name ?? "-").replace(/-/g, " ")
}

const getTahunName = (item: any) => {
  const t = item.tahun_akademik
  if (!t) return "-"
  const awal = t.tahun_awal?.slice(0, 4) ?? "?"
  const akhir = t.tahun_akhir?.slice(0, 4) ?? "?"
  return `${awal}/${akhir}`
}



// ================= FILTER CLIENT =================
const filteredData = computed(() => {
  return kelasData.value.filter((item) => {
    const p = prodiMap.value[item.prodi.id]
    const jurusanId = String(p?.jurusan?.id ?? "")
    const prodiId = String(item.prodi.id ?? "")
    const tahunId = String(item.tahun_akademik?.id ?? "")

    return (
      (!jurusan.value || jurusanId === jurusan.value) &&
      (!prodi.value || prodiId === prodi.value) &&
      (!tahunAkademik.value || tahunId === tahunAkademik.value)
    )
  })
})

const filteredProdiList = computed(() => {
  if (!jurusan.value) return prodiList.value

  return prodiList.value.filter(
    (p) => p.jurusanId === jurusan.value
  )
})

// ================= PAGINATION =================
const paginatedData = filteredData

const rowNumber = (index: number) =>
  (currentPage.value - 1) * perPage.value + index + 1
// ================= PAGINATION BUTTON (INI YANG KURANG) =================
const pages = computed<(number | "...")[]>(() => {
  const total = totalPages.value
  const current = currentPage.value
  const result: (number | "...")[] = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) result.push(i)
    return result
  }

  result.push(1)

  if (current > 3) result.push("...")

  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)

  for (let i = start; i <= end; i++) result.push(i)

  if (current < total - 2) result.push("...")

  result.push(total)

  return result
})

// ================= PAGINATION CONTROL =================
const goToPage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
}

const prevPage = () => goToPage(currentPage.value - 1)
const nextPage = () => goToPage(currentPage.value + 1)

// ================= ACTION =================
const pilihData = () => {
  currentPage.value = 1
  getKelas()
}

const tambahData = () => {
  router.push("/dashboard-admin/tambah_kelas")
}

const editData = (item: any) => {
  console.log("EDIT:", item)
}

const hapusData = (item: any) => {
  console.log("HAPUS:", item)
}

// ================= INIT =================
onMounted(async () => {
  await Promise.all([
    getProdi(),
    getTahunAkademik(),
  ])

  if (prodi.value) {
    getKelas()
  }
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
<div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">

  <!-- HEADER BIRU -->
  <div class="bg-[#243e90] px-5 py-4">
    <h2 class="text-white text-2xl font-bold">
      Data Kelas
    </h2>

    <p class="text-white text-sm mt-1">
      Kumpulan data keals yang tersimpan
    </p>
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
<!-- PRODI -->
<select
  v-model="prodi"
  class="w-[240px] h-[54px] border border-gray-300 rounded-xl px-4"
>
  <option value="">Pilih Prodi</option>
<option
  v-for="p in filteredProdiList"
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
            <tr class="text-left text-gray-600 border-b border-gray-300">
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
      <div class="flex items-center justify-end px-5 pt-10 pb-5">

        <!-- <select
          v-model="perPage"
          @change="currentPage = 1"
          class="w-[90px] h-[42px] border border-gray-300 rounded-lg px-3 text-sm outline-none"
        >
          <option :value="5">5 Baris</option>
          <option :value="10">10 Baris</option>
          <option :value="25">25 Baris</option>
        </select> -->

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