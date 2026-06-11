<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL =
  "https://be.karlearn.site"

const jurusanList = ref<any[]>([])

const form = ref({
  name: "",
  jenjang: "D3",
  jurusan_name: ""
})

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

const getJurusan = async () => {
  const res = await fetch(
    `${BASE_URL}/api/jurusan`,
    {
      headers: getHeaders()
    }
  )

  const json = await res.json()

  jurusanList.value =
    json.data?.items ??
    json.data ??
    []
}

const simpan = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/prodi`,
      {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify({
          name: form.value.name,
          jenjang: form.value.jenjang,
          jurusan_name:
            form.value.jurusan_name
        })
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json.message)
      return
    }

    alert("Berhasil tambah prodi")

    router.push(
      "/dashboard-admin/prodi"
    )
  } catch (err) {
    console.error(err)
  }
}

onMounted(() => {
  getJurusan()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <h1
      class="text-[42px] font-bold mb-6"
    >
      Tambah Prodi
    </h1>

    <div
      class="bg-white rounded-xl p-6 max-w-xl"
    >

      <div class="mb-4">
        <label>Nama Prodi</label>
        <input
          v-model="form.name"
          class="w-full border rounded-lg p-3 mt-1"
        />
      </div>

      <div class="mb-4">
        <label>Jenjang</label>

        <select
          v-model="form.jenjang"
          class="w-full border rounded-lg p-3 mt-1"
        >
          <option>D3</option>
          <option>D4</option>
          <option>S1</option>
          <option>S2</option>
        </select>
      </div>

      <div class="mb-6">
        <label>Jurusan</label>

        <select
          v-model="form.jurusan_name"
          class="w-full border rounded-lg p-3 mt-1"
        >
          <option value="">
            Pilih Jurusan
          </option>

          <option
            v-for="j in jurusanList"
            :key="j.id"
            :value="j.name"
          >
            {{ j.name }}
          </option>
        </select>
      </div>

      <div class="flex gap-3">
        <button
          @click="router.back()"
          class="px-5 py-3 border rounded-lg"
        >
          Kembali
        </button>

        <button
          @click="simpan"
          class="px-5 py-3 bg-[#29479d] text-white rounded-lg"
        >
          Simpan
        </button>
      </div>

    </div>
  </div>
</template>