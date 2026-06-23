<script setup lang="ts">
import { ref, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

// ================= STATE =================
const idTahunAkademik = ref("")
const tipeSemester = ref("") // "ganjil" | "genap"
const status = ref("")

const awalTanggal = ref("")
const awalBulan = ref("")
const awalTahun = ref("")

const akhirTanggal = ref("")
const akhirBulan = ref("")
const akhirTahun = ref("")

const loading = ref(false)
const warning = ref("")

// ================= UTIL =================
const setWarning = (msg: string) => {
  warning.value = msg
  setTimeout(() => (warning.value = ""), 2500)
}

const isLeapYear = (year: number) => {
  return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0
}

const getMaxDay = (month: number, year: number) => {
  if (!month) return 31

  const days = [
    31,
    isLeapYear(year) ? 29 : 28,
    31,
    30,
    31,
    30,
    31,
    31,
    30,
    31,
    30,
    31,
  ]

  return days[month - 1]
}

// ================= INPUT VALIDATION =================

// BULAN
const onBulanInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  let value = target.value.replace(/\D/g, "").slice(0, 2)

  if (Number(value) > 12) {
    value = "12"
    setWarning("Bulan maksimal 12")
  }

  awalBulan.value = value

  const maxDay = getMaxDay(Number(value), Number(awalTahun.value))

  if (awalTanggal.value && Number(awalTanggal.value) > maxDay) {
    awalTanggal.value = String(maxDay)
    setWarning(`Tanggal disesuaikan ke ${maxDay}`)
  }
}

// TANGGAL
const onTanggalInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  let value = target.value.replace(/\D/g, "").slice(0, 2)

  const maxDay = getMaxDay(Number(awalBulan.value), Number(awalTahun.value))

  if (Number(value) > maxDay) {
    value = String(maxDay)
    setWarning(`Tanggal maksimal bulan ini ${maxDay}`)
  }

  awalTanggal.value = value
}

// TAHUN
const onTahunInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  awalTahun.value = target.value.replace(/\D/g, "").slice(0, 4)

  const maxDay = getMaxDay(Number(awalBulan.value), Number(awalTahun.value))

  if (awalTanggal.value && Number(awalTanggal.value) > maxDay) {
    awalTanggal.value = String(maxDay)
    setWarning("Tanggal disesuaikan karena perubahan tahun")
  }
}

// ================= AUTO TAHUN AKHIR =================
watch(
  () => [awalTanggal.value, awalBulan.value, awalTahun.value],
  () => {
    if (
      awalTanggal.value &&
      awalBulan.value &&
      awalTahun.value.length === 4
    ) {
      akhirTanggal.value = awalTanggal.value
      akhirBulan.value = awalBulan.value
      akhirTahun.value = String(Number(awalTahun.value) + 1)
    } else {
      akhirTanggal.value = ""
      akhirBulan.value = ""
      akhirTahun.value = ""
    }
  }
)

// ================= AUTO ID (tahun + kode semester) =================
// contoh: tahun 2026 + genap(2) => 20262
watch(
  () => [awalTahun.value, tipeSemester.value],
  () => {
    if (awalTahun.value.length === 4 && tipeSemester.value) {
      const kode = tipeSemester.value === "ganjil" ? "1" : "2"
      idTahunAkademik.value = `${awalTahun.value}${kode}`
    }
  }
)

// ================= HEADER =================
const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= VALIDASI TANGGAL =================
const isValidDate = (t: string) => {
  const [y, m, d] = t.split("-").map(Number)
  return y > 0 && m >= 1 && m <= 12 && d >= 1 && d <= 31
}

// ================= SIMPAN =================
const simpanData = async () => {
  try {
    if (
      !idTahunAkademik.value ||
      !tipeSemester.value ||
      !status.value ||
      !awalTanggal.value ||
      !awalBulan.value ||
      !awalTahun.value
    ) {
      setWarning("Semua field wajib diisi")
      return
    }

    const tAwal = `${awalTahun.value}-${awalBulan.value.padStart(2, "0")}-${awalTanggal.value.padStart(2, "0")}`
    const tAkhir = `${akhirTahun.value}-${akhirBulan.value.padStart(2, "0")}-${akhirTanggal.value.padStart(2, "0")}`

    if (!isValidDate(tAwal) || !isValidDate(tAkhir)) {
      setWarning("Format tanggal tidak valid")
      return
    }

    loading.value = true

    const payload = {
      id: Number(idTahunAkademik.value),
      tipe_semester: tipeSemester.value.toLowerCase(),
      status: status.value.toLowerCase(),
      tahun_awal: tAwal,
      tahun_akhir: tAkhir,
    }

    console.log("PAYLOAD FINAL:", payload)

    const res = await fetch(
      "https://be.karlearn.site/api/tahun-akademik",
      {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(payload),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      console.log("API ERROR:", json)
      setWarning(json.message || "Gagal tambah data")
      return
    }

    alert("Berhasil tambah tahun akademik")
    router.push("/dashboard-admin/tahun_akademik")
  } catch (err) {
    console.error(err)
    setWarning("Terjadi kesalahan server")
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 p-6">

    <!-- Breadcrumb -->
    <div class="text-sm text-gray-500 mb-2">
      Akademik > Tahun Akademik > Tambah
    </div>

    <h1 class="text-3xl font-bold text-slate-800">
      Tahun Akademik
    </h1>

    <p class="text-gray-500 mb-6">
      Tambah Data Tahun Akademik
    </p>

    <!-- CARD -->
    <div class="bg-white rounded-xl shadow border p-6">

      <h2 class="text-xl font-semibold mb-6">
        Form Tambah
      </h2>

      <div class="grid md:grid-cols-2 gap-6">

        <!-- ID (auto, bisa diedit manual) -->
        <div>
          <label class="block text-sm font-medium mb-2">
            ID Tahun Akademik *
          </label>
          <input
            v-model="idTahunAkademik"
            inputmode="numeric"
            placeholder="Otomatis dari Tahun + Semester"
            class="input"
          />
        </div>

        <!-- TIPE SEMESTER -->
        <div>
          <label class="block text-sm font-medium mb-2">
            Tipe Semester *
          </label>
          <select v-model="tipeSemester" class="input">
            <option value="">Pilih</option>
            <option value="ganjil">Ganjil</option>
            <option value="genap">Genap</option>
          </select>
        </div>

        <!-- STATUS -->
        <div>
          <label class="block text-sm font-medium mb-2">
            Status *
          </label>
          <select v-model="status" class="input">
            <option value="">Pilih</option>
            <option value="aktif">Aktif</option>
            <option value="nonaktif">Non Aktif</option>
          </select>
        </div>

        <!-- TAHUN AWAL -->
        <div>
          <label class="block text-sm font-medium mb-2">
            Tahun Awal *
          </label>

          <div class="flex gap-2">
            <input
              v-model="awalTanggal"
              @input="onTanggalInput"
              maxlength="2"
              placeholder="Tgl"
              class="input text-center"
            />

            <input
              v-model="awalBulan"
              @input="onBulanInput"
              maxlength="2"
              placeholder="Bln"
              class="input text-center"
            />

            <input
              v-model="awalTahun"
              @input="onTahunInput"
              maxlength="4"
              placeholder="Thn"
              class="input text-center"
            />
          </div>
        </div>

        <!-- TAHUN AKHIR -->
        <div>
          <label class="block text-sm font-medium mb-2">
            Tahun Akhir (Auto)
          </label>

          <div class="flex gap-2">
            <input v-model="akhirTanggal" readonly class="input text-center bg-gray-100" />
            <input v-model="akhirBulan" readonly class="input text-center bg-gray-100" />
            <input v-model="akhirTahun" readonly class="input text-center bg-gray-100" />
          </div>
        </div>

      </div>

      <!-- WARNING -->
      <p v-if="warning" class="text-red-500 text-sm mt-3">
        {{ warning }}
      </p>

      <!-- BUTTON -->
      <div class="flex gap-3 mt-6">
        <button
          @click="router.push('/dashboard-admin/tahun_akademik')"
          class="px-4 py-2 border rounded-lg"
        >
          Batal
        </button>

        <button
          @click="simpanData"
          :disabled="loading"
          class="px-4 py-2 bg-blue-700 text-white rounded-lg"
        >
          {{ loading ? "Menyimpan..." : "Simpan" }}
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
.input {
  width: 100%;
  border: 1px solid #ddd;
  padding: 10px;
  border-radius: 8px;
}
</style>