<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute } from "vue-router"

const route = useRoute()

const BASE_URL = "https://be.karlearn.site"

// Ambil NIM dari params route: /dashboard-admin/detail_khs_mahasiswa/:nim
const nim = route.params.nim as string

interface MahasiswaInfo {
  nama: string
  nim: string
  kelas: string
  prodi: string
  status: string
  kurikulum: string
  semester: number
  sksLulus: number
  ipk: number
  angkatan: number
}

interface MataKuliah {
  kode: string
  nama: string
  sks: number
  mutu: number
  bobot: number
  nilai: string
}

interface KhsDetail {
  periode: string
  matakuliah: MataKuliah[]
  totalSks: number
  totalBobot: number
  ips: number
}

const mahasiswa = ref<MahasiswaInfo | null>(null)
const khsDetail = ref<KhsDetail | null>(null)
const loading = ref(false)
const error = ref("")

const fetchDetail = async () => {
  loading.value = true
  error.value = ""
  try {
    const token = localStorage.getItem("token")
    const res = await fetch(`${BASE_URL}/api/khs/${nim}`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const json = await res.json()

    // Sesuaikan field mapping dengan response API
    const d = json.data ?? json
    mahasiswa.value = {
      nama: d.nama ?? d.mahasiswa?.nama,
      nim: d.nim ?? nim,
      kelas: d.kelas ?? d.mahasiswa?.kelas,
      prodi: d.prodi ?? d.mahasiswa?.prodi,
      status: d.status ?? "Aktif",
      kurikulum: d.kurikulum ?? "-",
      semester: d.semester ?? 0,
      sksLulus: d.sks_lulus ?? d.sksLulus ?? 0,
      ipk: d.ipk ?? 0,
      angkatan: d.angkatan ?? 0,
    }

    const mk: MataKuliah[] = (d.matakuliah ?? d.mata_kuliah ?? []).map((m: any) => ({
      kode: m.kode,
      nama: m.nama,
      sks: m.sks,
      mutu: m.mutu ?? m.nilai_mutu,
      bobot: m.bobot,
      nilai: m.nilai,
    }))

    const totalSks = mk.reduce((sum, m) => sum + m.sks, 0)
    const totalBobot = mk.reduce((sum, m) => sum + m.bobot, 0)
    const ips = totalSks > 0 ? totalBobot / totalSks : 0

    khsDetail.value = {
      periode: d.periode ?? d.tahun_akademik ?? "-",
      matakuliah: mk,
      totalSks,
      totalBobot,
      ips: parseFloat(ips.toFixed(2)),
    }
  } catch (e: any) {
    error.value = "Gagal memuat data KHS mahasiswa."
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchDetail)
</script>

<template>
  <div class="p-4">

    <!-- BREADCRUMB -->
    <p class="text-gray-500 text-sm">
      Akademik › KHS › {{ mahasiswa?.kelas ?? "..." }} › {{ nim }}
    </p>
    <h1 class="text-4xl font-bold mb-5 text-gray-700">Kartu Hasil Studi</h1>

    <!-- LOADING / ERROR -->
    <div v-if="loading" class="text-center py-10 text-gray-500">
      Memuat data...
    </div>

    <div v-else-if="error" class="text-center py-10 text-red-500">
      {{ error }}
      <button @click="fetchDetail" class="ml-2 underline text-blue-600">Coba lagi</button>
    </div>

    <template v-else-if="mahasiswa">

      <!-- INFO MAHASISWA -->
      <div class="bg-white rounded-xl shadow p-6 mb-5">
        <div class="grid grid-cols-2 gap-10">
          <div class="space-y-4">
            <div><b>Nama :</b> {{ mahasiswa.nama }}</div>
            <div><b>NIM :</b> {{ mahasiswa.nim }}</div>
            <div><b>Kelas :</b> {{ mahasiswa.kelas }}</div>
            <div><b>Prodi :</b> {{ mahasiswa.prodi }}</div>
            <div><b>Status :</b> {{ mahasiswa.status }}</div>
          </div>
          <div class="space-y-4">
            <div><b>Kurikulum :</b> {{ mahasiswa.kurikulum }}</div>
            <div><b>Semester :</b> {{ mahasiswa.semester }}</div>
            <div><b>SKS Lulus :</b> {{ mahasiswa.sksLulus }}</div>
            <div><b>IPK :</b> {{ mahasiswa.ipk }}</div>
            <div><b>Angkatan :</b> {{ mahasiswa.angkatan }}</div>
          </div>
        </div>
      </div>

      <!-- TABEL MATAKULIAH -->
      <div v-if="khsDetail" class="bg-white rounded-xl shadow overflow-hidden">
        <div class="bg-[#1f3c93] text-white p-5">
          <h2 class="text-3xl font-bold">Periode {{ khsDetail.periode }}</h2>
        </div>

        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-700">
              <th class="p-4">No</th>
              <th>Kode</th>
              <th>Matakuliah</th>
              <th>SKS</th>
              <th>Nilai Mutu</th>
              <th>Bobot</th>
              <th>Nilai</th>
              <th>Keterangan</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(item, index) in khsDetail.matakuliah"
              :key="item.kode"
            >
              <td class="p-4">{{ index + 1 }}</td>
              <td>{{ item.kode }}</td>
              <td>{{ item.nama }}</td>
              <td>{{ item.sks }}</td>
              <td>{{ item.mutu }}</td>
              <td>{{ item.bobot }}</td>
              <td>{{ item.nilai }}</td>
              <td>-</td>
            </tr>

            <tr v-if="khsDetail.matakuliah.length === 0">
              <td colspan="8" class="text-center py-6 text-gray-500">
                Tidak ada data matakuliah
              </td>
            </tr>

            <tr class="border-t">
              <td colspan="3" class="p-4 font-bold">Total SKS</td>
              <td>{{ khsDetail.totalSks }}</td>
              <td></td>
              <td>{{ khsDetail.totalBobot }}</td>
              <td colspan="2"></td>
            </tr>

            <tr>
              <td colspan="3" class="p-4 font-bold">Indeks Prestasi Semester</td>
              <td>{{ khsDetail.ips.toFixed(2) }}</td>
              <td colspan="4"></td>
            </tr>
          </tbody>
        </table>
      </div>

    </template>

  </div>
</template>