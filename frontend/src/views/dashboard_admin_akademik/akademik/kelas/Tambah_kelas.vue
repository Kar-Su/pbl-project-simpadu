<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= STATE =================
const tahunAkademikList = ref<any[]>([])
const prodiMap = ref<Record<number, any>>({})
const kurikulumList = ref<any[]>([])
const loading = ref(false)

// ================= FORM =================
const tahunAkademikId = ref<number | "">("")
const jurusanId = ref("")
const prodiName = ref("")
const namaKelas = ref("")
const kurikulumKode = ref("")
const semester = ref("")

// ================= GET TAHUN AKADEMIK =================
const getTahunAkademik = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik`, {
      headers: getHeaders(),
    })
    const json = await res.json()
    // Response: data adalah array langsung (bukan data.items)
    tahunAkademikList.value = Array.isArray(json?.data) ? json.data : []
    console.log("TAHUN AKADEMIK:", tahunAkademikList.value)
  } catch (err) {
    console.error("GET TAHUN AKADEMIK ERROR:", err)
  }
}

const tahunAkademikOptions = computed(() =>
  tahunAkademikList.value.map((t: any) => ({
    id: t.id,
    // Field asli dari backend adalah tipee_semester (typo tapi itulah yang ada)
    label: `${t.tahun_awal?.slice(0, 4) ?? "?"}/${t.tahun_akhir?.slice(0, 4) ?? "?"} - ${
      t.tipee_semester === "ganjil" ? "Ganjil"
      : t.tipee_semester === "genap" ? "Genap"
      : t.tipee_semester ?? "-"
    } (${t.status ?? "-"})`,
  }))
)

// ================= GET PRODI =================
const getProdi = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, {
      headers: getHeaders(),
    })
    const json = await res.json()
    const list: any[] = json?.data?.items ?? json?.data ?? []
    const map: Record<number, any> = {}
    list.forEach((p: any) => { map[p.id] = p })
    prodiMap.value = map
    console.log("PRODI:", list)
  } catch (err) {
    console.error("GET PRODI ERROR:", err)
  }
}

// Jurusan diturunkan dari prodi
const jurusanList = computed(() => {
  const map = new Map()
  Object.values(prodiMap.value).forEach((p: any) => {
    const j = p?.jurusan
    if (j?.id && !map.has(j.id)) {
      map.set(j.id, {
        id: String(j.id),
        name: (j.name ?? "-").replace(/-/g, " "),
      })
    }
  })
  return Array.from(map.values())
})

const prodiList = computed(() =>
  Object.values(prodiMap.value).map((p: any) => ({
    id: String(p.id),
    rawName: p.name ?? "",
    displayName: (p.name ?? "-").replace(/-/g, " "),
    jurusanId: String(p.jurusan?.id ?? ""),
  }))
)

const filteredProdiList = computed(() => {
  if (!jurusanId.value) return prodiList.value
  return prodiList.value.filter((p) => p.jurusanId === jurusanId.value)
})

watch(jurusanId, () => {
  const stillValid = filteredProdiList.value.some((p) => p.rawName === prodiName.value)
  if (!stillValid) prodiName.value = ""
})

// ================= GET KURIKULUM =================
const getKurikulum = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/kurikulum?per_page=100`, {
      headers: getHeaders(),
    })
    const json = await res.json()
    // Response: data.items
    kurikulumList.value = json?.data?.items ?? []
    console.log("KURIKULUM:", kurikulumList.value)
  } catch (err) {
    console.error("GET KURIKULUM ERROR:", err)
  }
}

const fmt = (str: string) => (str ?? "-").replace(/-/g, " ")

// ================= SIMPAN =================
const simpanKelas = async () => {
  if (!namaKelas.value.trim())   { alert("Nama kelas tidak boleh kosong."); return }
  if (!prodiName.value)          { alert("Prodi harus dipilih."); return }
  if (!tahunAkademikId.value)    { alert("Tahun akademik harus dipilih."); return }
  if (!kurikulumKode.value)      { alert("Kurikulum harus dipilih."); return }
  if (!semester.value)           { alert("Semester harus dipilih."); return }

  loading.value = true

  try {
    const payload = {
      name: namaKelas.value,
      prodi_name: prodiName.value,
      tahun_akademik_id: Number(tahunAkademikId.value),
      kurikulum_kode: kurikulumKode.value,
      semester: Number(semester.value),
    }

    console.log("PAYLOAD:", payload)

    const res = await fetch(`${BASE_URL}/api/kelas`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })

    const result = await res.json()
    console.log("CREATE KELAS RESULT:", result)

    if (!res.ok) {
      alert(result?.message || "Gagal menyimpan kelas.")
      return
    }

    alert("Kelas berhasil disimpan!")
    router.push("/dashboard-admin/kelas")
  } catch (err) {
    console.error("SIMPAN KELAS ERROR:", err)
    alert("Terjadi kesalahan jaringan saat menyimpan kelas.")
  } finally {
    loading.value = false
  }
}

// ================= MOUNTED =================
onMounted(() => {
  getTahunAkademik()
  getProdi()
  getKurikulum()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Akademik > Kelas > Tambah Kelas
    </div>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-gray-800">Tambah Kelas</h1>
    <p class="mb-6 mt-1 text-gray-500">Pengelolaan Data</p>

    <!-- DEBUG: hapus setelah konfirmasi data tampil -->
    <div class="mb-4 rounded-xl bg-yellow-50 border border-yellow-200 p-3 text-xs text-yellow-800 font-mono">
      <div>Tahun Akademik: {{ tahunAkademikList.length }} item</div>
      <div>Prodi: {{ Object.keys(prodiMap).length }} item</div>
      <div>Kurikulum: {{ kurikulumList.length }} item</div>
    </div>

    <!-- FORM AKADEMIK -->
    <div class="mb-5 rounded-2xl border border-blue-100 bg-white shadow-sm">
      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">Form Akademik</h2>
      </div>

      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-3">

        <!-- TAHUN AKADEMIK -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Tahun Akademik</label>
          <select
            v-model="tahunAkademikId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >
            <option value="" disabled>Pilih Tahun Akademik</option>
            <option v-for="item in tahunAkademikOptions" :key="item.id" :value="item.id">
              {{ item.label }}
            </option>
          </select>
        </div>

        <!-- JURUSAN -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Jurusan</label>
          <select
            v-model="jurusanId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >
            <option value="" disabled>Pilih Jurusan</option>
            <option v-for="item in jurusanList" :key="item.id" :value="item.id" class="capitalize">
              {{ item.name }}
            </option>
          </select>
        </div>

        <!-- PRODI -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Prodi</label>
          <select
            v-model="prodiName"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >
            <option value="" disabled>Pilih Prodi</option>
            <option
              v-for="item in filteredProdiList"
              :key="item.id"
              :value="item.rawName"
              class="capitalize"
            >
              {{ item.displayName }}
            </option>
          </select>
        </div>

      </div>
    </div>

    <!-- FORM KELAS -->
    <div class="rounded-2xl border border-blue-100 bg-white shadow-sm">
      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">Form Kelas</h2>
      </div>

      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-2">

        <!-- NAMA KELAS -->
        <div class="md:col-span-2">
          <label class="mb-2 block text-sm font-medium text-gray-700">Nama Kelas</label>
          <input
            v-model="namaKelas"
            type="text"
            placeholder="Isi Nama Kelas..."
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          />
        </div>

        <!-- KURIKULUM -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Kurikulum</label>
          <select
            v-model="kurikulumKode"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >
            <option value="" disabled>Pilih Kurikulum</option>
            <option
              v-for="item in kurikulumList"
              :key="item.kode"
              :value="item.kode"
              class="capitalize"
            >
              {{ fmt(item.name) }} ({{ item.kode?.toUpperCase() }})
            </option>
          </select>
        </div>

        <!-- SEMESTER -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">Semester</label>
          <select
            v-model="semester"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          >
            <option value="" disabled>Pilih Semester</option>
            <option v-for="n in 8" :key="n" :value="String(n)">Semester {{ n }}</option>
          </select>
        </div>

      </div>
    </div>

    <!-- BUTTON -->
    <div class="mt-5">
      <button
        @click="simpanKelas"
        :disabled="loading"
        class="rounded-xl bg-green-500 px-6 py-3 text-sm font-semibold text-white shadow hover:bg-green-600 disabled:opacity-50"
      >
        {{ loading ? "Menyimpan..." : "💾 Simpan" }}
      </button>
    </div>

  </div>
</template>