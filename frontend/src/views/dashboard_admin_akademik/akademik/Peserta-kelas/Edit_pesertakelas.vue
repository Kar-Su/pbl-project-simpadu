<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()

const BASE_URL = "https://be.karlearn.site"
const kelasIdFromRoute = route.params.id as string

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

interface Kelas {
  id: string
  name: string
}

interface Mahasiswa {
  id: string
  nama: string
  nim?: string
  email?: string
}

interface Peserta {
  mahasiswa_id: string
  nama: string
  nim?: string
}

const kelasList = ref<Kelas[]>([])
const mahasiswaList = ref<Mahasiswa[]>([])
const pesertaKelas = ref<Peserta[]>([])

const selectedKelas = ref("")
const selectedMahasiswa = ref("")

const isLoading = ref(false)

const getKelasData = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/kelas`, {
      headers: getHeaders(),
    })

    const data = await res.json()

    console.log("DATA KELAS", data)

    const payload = data.data ?? []

    kelasList.value = Array.isArray(payload)
      ? payload
      : payload.items ?? []
  } catch (err) {
    console.error(err)
  }
}

const getMahasiswaData = async () => {
  try {
    const res = await fetch(`${BASE_URL}/api/mahasiswa`, {
      headers: getHeaders(),
    })

    const data = await res.json()

    console.log("DATA MAHASISWA", data)

    const payload = data.data ?? []

    mahasiswaList.value = Array.isArray(payload)
      ? payload
      : payload.items ?? []
  } catch (err) {
    console.error(err)
  }
}

const getDetailKelas = async () => {
  if (!kelasIdFromRoute) return

  try {
    isLoading.value = true

    const res = await fetch(
      `${BASE_URL}/api/kelas/${kelasIdFromRoute}`,
      {
        headers: getHeaders(),
      }
    )

    const data = await res.json()

    console.log("DETAIL KELAS", data)

    const kelas = data.data

    selectedKelas.value = kelas.id

    pesertaKelas.value =
      kelas.mahasiswa?.map((m: any) => ({
        mahasiswa_id: m.mahasiswa_id,
        nama: m.name,
        nim: m.email,
      })) ?? []
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const handleSelectMahasiswa = () => {
  if (!selectedMahasiswa.value) return

  const mahasiswa = mahasiswaList.value.find(
    (m) => String(m.id) === String(selectedMahasiswa.value)
  )

  if (!mahasiswa) return

  const sudahAda = pesertaKelas.value.some(
    (p) =>
      String(p.mahasiswa_id) === String(mahasiswa.id)
  )

  if (sudahAda) {
    alert("Mahasiswa sudah ada di kelas")
    selectedMahasiswa.value = ""
    return
  }

  pesertaKelas.value.push({
    mahasiswa_id: mahasiswa.id,
    nama: mahasiswa.nama,
    nim: mahasiswa.nim,
  })

  selectedMahasiswa.value = ""
}

const removeMahasiswaItem = (index: number) => {
  pesertaKelas.value.splice(index, 1)
}

const handleSimpan = async () => {
  try {
    const payload = {
      kelas_id: selectedKelas.value,
      mahasiswa_ids: pesertaKelas.value.map(
        (item) => item.mahasiswa_id
      ),
    }

    console.log(payload)

    const res = await fetch(
      `${BASE_URL}/api/peserta-kelas/update-bulk`,
      {
        method: "PUT",
        headers: getHeaders(),
        body: JSON.stringify(payload),
      }
    )

    const data = await res.json()

    if (res.ok) {
      alert("Berhasil update peserta kelas")
      router.push("/dashboard-admin/peserta_kelas")
    } else {
      alert(data.message || "Gagal update")
    }
  } catch (err) {
    console.error(err)
    alert("Terjadi kesalahan")
  }
}

onMounted(async () => {
  await getKelasData()
  await getMahasiswaData()
  await getDetailKelas()
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 p-6">
    <div class="max-w-4xl mx-auto">

      <h1 class="text-3xl font-bold mb-6">
        Edit Peserta Kelas
      </h1>

      <div class="bg-white rounded-xl shadow border p-6 space-y-5">

        <!-- KELAS -->
        <div>
          <label class="block mb-2 text-sm font-medium">
            Kelas
          </label>

          <select
            v-model="selectedKelas"
            class="w-full h-11 rounded-lg border px-4"
          >
            <option value="">
              Pilih Kelas
            </option>

            <option
              v-for="kelas in kelasList"
              :key="kelas.id"
              :value="kelas.id"
            >
              {{ kelas.name }}
            </option>
          </select>
        </div>

        <!-- MAHASISWA -->
        <div>
          <label class="block mb-2 text-sm font-medium">
            Tambah Mahasiswa
          </label>

          <select
            v-model="selectedMahasiswa"
            @change="handleSelectMahasiswa"
            class="w-full h-11 rounded-lg border px-4"
          >
            <option value="">
              Pilih Mahasiswa
            </option>

            <option
              v-for="mhs in mahasiswaList"
              :key="mhs.id"
              :value="mhs.id"
            >
              {{ mhs.nama }}
            </option>
          </select>
        </div>

        <!-- LIST PESERTA -->
        <div>
          <div
            v-if="isLoading"
            class="text-center py-5"
          >
            Loading...
          </div>

          <div
            v-else-if="pesertaKelas.length === 0"
            class="text-center py-5 text-gray-500"
          >
            Belum ada peserta kelas
          </div>

          <div
            v-for="(item,index) in pesertaKelas"
            :key="index"
            class="flex items-center justify-between bg-blue-50 border rounded-lg px-4 py-3 mb-2"
          >
            <div>
              <div class="font-medium">
                {{ item.nama }}
              </div>

              <div
                v-if="item.nim"
                class="text-sm text-gray-500"
              >
                {{ item.nim }}
              </div>
            </div>

            <button
              @click="removeMahasiswaItem(index)"
              class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
            >
              Hapus
            </button>
          </div>
        </div>

        <!-- ACTION -->
        <div class="flex gap-3 pt-2">
          <button
            @click="handleSimpan"
            class="bg-blue-900 text-white px-6 py-2 rounded-lg"
          >
            Simpan
          </button>

          <button
            @click="router.push('/dashboard-admin/peserta_kelas')"
            class="bg-gray-200 px-6 py-2 rounded-lg"
          >
            Batal
          </button>
        </div>

      </div>
    </div>
  </div>
</template>