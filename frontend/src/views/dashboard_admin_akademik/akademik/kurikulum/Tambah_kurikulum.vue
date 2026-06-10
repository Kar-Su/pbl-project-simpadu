<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const BASE_URL = "https://be.karlearn.site"

// ================= INTERFACE =================
interface TahunAkademik {
  id: number
  tipee_semester: string
  tahun_awal: string
  tahun_akhir: string
  status: string
}

interface MataKuliahOption {
  id: string
  kode: string
  name: string
  sks: number
}

interface MKEntry {
  mk_kode: string      // kode MK yang dipilih
  semester: number | null
  wajib: boolean       // true = Wajib, false = Pilihan
}

// ================= STATE =================
const namaKurikulum = ref("")
const kodeKurikulum = ref("")
const tahunAkademikId = ref("")
const jurusanId = ref("")
const prodiName = ref("")         // nama prodi (format kebab, untuk payload)

const jurusanList = ref<any[]>([])
const prodiList = ref<any[]>([])
const tahunAkademikList = ref<TahunAkademik[]>([])
const mataKuliahOptions = ref<MataKuliahOption[]>([])

const mkEntries = ref<MKEntry[]>([
  { mk_kode: "", semester: null, wajib: true },
])

const submitLoading = ref(false)
const submitError = ref("")
const submitSuccess = ref("")

// ================= HEADER =================
const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= COMPUTED =================
const filteredProdiList = computed(() => {
  if (!jurusanId.value) return []
  return prodiList.value.filter(
    (p) => String(p.jurusan?.id) === String(jurusanId.value)
  )
})

watch(jurusanId, () => {
  prodiName.value = ""
})

// ================= FETCH =================
const getTahunAkademik = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/tahun-akademik`, { headers: getHeaders() })
    const json = await res.json()
    tahunAkademikList.value = Array.isArray(json?.data) ? json.data : []
  } catch (err) {
    console.error("GET TAHUN AKADEMIK ERROR:", err)
  }
}

const getProdi = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/prodi`, { headers: getHeaders() })
    const json = await res.json()
    console.log("PRODI RAW:", json)
    prodiList.value = Array.isArray(json?.data)
      ? json.data
      : json?.data?.items ?? []

    // Build jurusanList dari prodi
    const map = new Map<number, any>()
    prodiList.value.forEach((p: any) => {
      const j = p.jurusan
      if (j?.id && !map.has(j.id)) map.set(j.id, { id: j.id, name: j.name })
    })
    jurusanList.value = Array.from(map.values())
  } catch (err) {
    console.error("GET PRODI ERROR:", err)
  }
}

const getMataKuliah = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/mata-kuliah`, { headers: getHeaders() })
    const json = await res.json()
    console.log("MK RAW:", json)
    const list = Array.isArray(json?.data)
      ? json.data
      : json?.data?.items ?? []
    mataKuliahOptions.value = list.map((m: any) => ({
      id: m.id,
      kode: m.kode ?? m.code ?? "",
      name: (m.name ?? m.nama ?? "").replace(/-/g, " "),
      sks: m.sks ?? 0,
    }))
  } catch (err) {
    console.error("GET MK ERROR:", err)
  }
}

onMounted(async () => {
  await Promise.all([getTahunAkademik(), getProdi(), getMataKuliah()])
})

// ================= MK ENTRIES =================
const tambahMK = () => {
  mkEntries.value.push({ mk_kode: "", semester: null, wajib: true })
}

const hapusMK = (index: number) => {
  mkEntries.value.splice(index, 1)
}

// ================= FORMAT =================
const formatTahun = (item: TahunAkademik) => {
  const awal = item.tahun_awal?.slice(0, 4) ?? "?"
  const akhir = item.tahun_akhir?.slice(0, 4) ?? "?"
  const tipe = item.tipee_semester
    ? item.tipee_semester.charAt(0).toUpperCase() + item.tipee_semester.slice(1)
    : ""
  return `${awal}/${akhir} – ${tipe}`
}

// ================= SIMPAN =================
const simpanKurikulum = async () => {
  submitError.value = ""
  submitSuccess.value = ""

  // Validasi
  if (!namaKurikulum.value.trim()) return (submitError.value = "Nama kurikulum wajib diisi.")
  if (!kodeKurikulum.value.trim()) return (submitError.value = "Kode kurikulum wajib diisi.")
  if (!tahunAkademikId.value) return (submitError.value = "Tahun akademik wajib dipilih.")
  if (!jurusanId.value) return (submitError.value = "Jurusan wajib dipilih.")
  if (!prodiName.value) return (submitError.value = "Prodi wajib dipilih.")

  const adaMKKosong = mkEntries.value.some((mk) => !mk.mk_kode || mk.semester === null)
  if (adaMKKosong) return (submitError.value = "Semua mata kuliah harus diisi lengkap (MK & Semester).")

  submitLoading.value = true

  try {
    // ── STEP 1: POST kurikulum ──
    const payloadKurikulum = {
      kode: kodeKurikulum.value.trim().replace(/\s+/g, "-").toLowerCase(),
      name: namaKurikulum.value.trim(),
      prodi_name: prodiName.value,
    }

    console.log("PAYLOAD KURIKULUM:", payloadKurikulum)

    const resKurikulum = await fetch(`${BASE_URL}/api/kurikulum`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payloadKurikulum),
    })

    if (!resKurikulum.ok) {
      const json = await resKurikulum.json()
      submitError.value = json?.message || "Gagal membuat kurikulum."
      return
    }

    // ── STEP 2: POST tiap mata kuliah ──
    for (const mk of mkEntries.value) {
      const payloadMK = {
        kurikulum_kode: payloadKurikulum.kode,
        mk_kode: mk.mk_kode,
        semester: Number(mk.semester),
        wajib: mk.wajib,
      }

      console.log("PAYLOAD MK:", payloadMK)

      const resMK = await fetch(`${BASE_URL}/api/kurikulum/mata-kuliah`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(payloadMK),
      })

      if (!resMK.ok) {
        const json = await resMK.json()
        submitError.value = `Gagal menambahkan MK ${mk.mk_kode}: ${json?.message ?? ""}`
        return
      }
    }

    submitSuccess.value = "Kurikulum berhasil disimpan!"
    setTimeout(() => router.push("/dashboard-admin/kurikulum"), 1000)
  } catch (err) {
    console.error("SIMPAN ERROR:", err)
    submitError.value = "Terjadi kesalahan jaringan."
  } finally {
    submitLoading.value = false
  }
}

const kembali = () => router.push("/dashboard-admin/kurikulum")
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <div class="mb-2 text-sm text-gray-400">
      Mahasiswa > Kurikulum > Tambah Kurikulum
    </div>

    <!-- TITLE -->
    <div class="flex items-center justify-between mb-1">
      <h1 class="text-4xl font-bold text-gray-800">Tambah Kurikulum</h1>
      <button @click="kembali"
        class="flex items-center gap-2 px-4 py-2 rounded-xl border border-gray-300 bg-white hover:bg-gray-100 text-gray-600 text-sm font-medium transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Kembali
      </button>
    </div>
    <p class="mb-6 mt-1 text-gray-500">Pengelolaan Data</p>

    <!-- ALERT -->
    <div v-if="submitError" class="mb-4 px-4 py-3 bg-red-50 border border-red-300 text-red-700 rounded-xl text-sm">
      ❌ {{ submitError }}
    </div>
    <div v-if="submitSuccess" class="mb-4 px-4 py-3 bg-green-50 border border-green-300 text-green-700 rounded-xl text-sm">
      ✅ {{ submitSuccess }}
    </div>

    <!-- ══════════════════════════════════ -->
    <!-- FORM KURIKULUM                     -->
    <!-- ══════════════════════════════════ -->
    <div class="mb-5 rounded-2xl border border-blue-100 bg-white shadow-sm">

      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">Form Kurikulum</h2>
      </div>

      <div class="grid grid-cols-1 gap-5 p-5 md:grid-cols-2">

        <!-- NAMA KURIKULUM -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Nama Kurikulum <span class="text-red-500">*</span>
          </label>
          <input v-model="namaKurikulum" type="text" placeholder="Isi Nama Kurikulum ..."
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
        </div>

        <!-- TAHUN AKADEMIK -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Tahun Akademik <span class="text-red-500">*</span>
          </label>
          <select v-model="tahunAkademikId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option value="">Pilih Tahun Akademik</option>
            <option v-for="item in tahunAkademikList" :key="item.id" :value="item.id">
              {{ formatTahun(item) }}
            </option>
          </select>
        </div>

        <!-- JURUSAN -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Jurusan <span class="text-red-500">*</span>
          </label>
          <select v-model="jurusanId"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option value="">Pilih Jurusan</option>
            <option v-for="j in jurusanList" :key="j.id" :value="j.id">
              {{ j.name.replace(/-/g, " ") }}
            </option>
          </select>
        </div>

        <!-- PRODI -->
        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Prodi <span class="text-red-500">*</span>
          </label>
          <select v-model="prodiName" :disabled="!jurusanId"
            :class="!jurusanId ? 'bg-gray-100 cursor-not-allowed text-gray-400' : ''"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option value="">{{ jurusanId ? "Pilih Prodi" : "Pilih Jurusan terlebih dahulu" }}</option>
            <option v-for="p in filteredProdiList" :key="p.id" :value="p.name">
              {{ p.name.replace(/-/g, " ") }}
            </option>
          </select>
        </div>

        <!-- KODE KURIKULUM -->
        <div class="md:col-span-2">
          <label class="mb-2 block text-sm font-medium text-gray-700">
            Kode Kurikulum <span class="text-red-500">*</span>
          </label>
          <input v-model="kodeKurikulum" type="text" placeholder="contoh: kur-2025-01"
            class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
          <p class="text-xs text-gray-400 mt-1">Spasi otomatis diubah jadi tanda hubung</p>
        </div>

      </div>
    </div>

    <!-- ══════════════════════════════════ -->
    <!-- FORM MATA KULIAH                   -->
    <!-- ══════════════════════════════════ -->
    <div class="rounded-2xl border border-blue-100 bg-white shadow-sm">

      <div class="border-b border-gray-100 px-5 py-4">
        <h2 class="text-2xl font-semibold text-gray-700">Form Matakuliah</h2>
      </div>

      <div class="p-5">

        <div v-for="(mk, index) in mkEntries" :key="index"
          class="mb-6 pb-6 border-b border-gray-100 last:border-0 last:mb-0 last:pb-0">

          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">

            <!-- DROPDOWN MATAKULIAH -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">
                Matakuliah <span class="text-red-500">*</span>
              </label>
              <select v-model="mk.mk_kode"
                class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
                <option value="">Pilih Matakuliah</option>
                <option v-for="m in mataKuliahOptions" :key="m.id" :value="m.kode">
                  {{ m.name }} ({{ m.kode }}) — {{ m.sks }} SKS
                </option>
              </select>
            </div>

            <!-- STATUS WAJIB / PILIHAN -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">
                Status <span class="text-red-500">*</span>
              </label>
              <div class="flex gap-3">
                <button type="button" @click="mk.wajib = true"
                  :class="mk.wajib
                    ? 'bg-[#1f3c93] text-white border-[#1f3c93]'
                    : 'bg-white text-gray-600 border-gray-300 hover:border-[#1f3c93] hover:text-[#1f3c93]'"
                  class="flex-1 rounded-xl border px-4 py-3 text-sm font-semibold transition">
                  Wajib
                </button>
                <button type="button" @click="mk.wajib = false"
                  :class="!mk.wajib
                    ? 'bg-[#1f3c93] text-white border-[#1f3c93]'
                    : 'bg-white text-gray-600 border-gray-300 hover:border-[#1f3c93] hover:text-[#1f3c93]'"
                  class="flex-1 rounded-xl border px-4 py-3 text-sm font-semibold transition">
                  Pilihan
                </button>
              </div>
            </div>

            <!-- SEMESTER -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">
                Semester <span class="text-red-500">*</span>
              </label>
              <input v-model="mk.semester" type="number" min="1" max="8" placeholder="Contoh: 1"
                class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
            </div>

            <!-- TOMBOL HAPUS -->
            <div class="flex items-end">
              <button v-if="mkEntries.length > 1" type="button" @click="hapusMK(index)"
                class="w-full rounded-xl bg-red-500 hover:bg-red-600 px-4 py-3 text-sm font-semibold text-white transition">
                🗑 Hapus
              </button>
            </div>

          </div>
        </div>

        <!-- TOMBOL TAMBAH MK -->
        <div class="mt-4 flex flex-col items-center gap-1">
          <button type="button" @click="tambahMK"
            class="flex h-12 w-12 items-center justify-center rounded-full border-2 border-gray-400 text-2xl text-gray-600 transition hover:bg-gray-100">
            +
          </button>
          <span class="text-sm text-gray-500">Tambah</span>
        </div>

      </div>
    </div>

    <!-- TOMBOL SIMPAN -->
    <div class="mt-6 flex gap-3">
      <button @click="simpanKurikulum" :disabled="submitLoading"
        class="rounded-xl bg-green-500 hover:bg-green-600 px-6 py-3 text-sm font-semibold text-white shadow disabled:opacity-50 transition">
        {{ submitLoading ? "Menyimpan..." : "💾 Simpan" }}
      </button>
      <button @click="kembali"
        class="rounded-xl border border-gray-300 bg-white hover:bg-gray-100 px-6 py-3 text-sm font-semibold text-gray-600 transition">
        Batal
      </button>
    </div>

  </div>
</template>