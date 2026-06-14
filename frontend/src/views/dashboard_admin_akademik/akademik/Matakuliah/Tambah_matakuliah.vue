<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

interface MataKuliahForm {
  kode: string
  name: string
  sks: number | null
}

const loading = ref(false)
const errorMessage = ref("")
const successMessage = ref("")

const form = ref<MataKuliahForm>({
  kode: "",
  name: "",
  sks: null,
})

const getHeaders = () => ({
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

const handleSksInput = (event: Event) => {
  const target = event.target as HTMLInputElement

  if (target.value === "") {
    form.value.sks = null
    return
  }

  let value = Number(target.value)

  if (isNaN(value)) {
    form.value.sks = null
    return
  }

  if (value > 3) value = 3
  if (value < 1) value = 1

  form.value.sks = value
}

const simpanMataKuliah = async () => {
  errorMessage.value = ""
  successMessage.value = ""

  if (!form.value.kode.trim()) {
    errorMessage.value = "Kode mata kuliah wajib diisi"
    return
  }

  if (!form.value.name.trim()) {
    errorMessage.value = "Nama mata kuliah wajib diisi"
    return
  }

  if (
    form.value.sks === null ||
    form.value.sks < 1 ||
    form.value.sks > 3
  ) {
    errorMessage.value = "SKS hanya boleh 1 sampai 3"
    return
  }

  loading.value = true

  try {
    const response = await fetch(
      `${BASE_URL}/api/mata-kuliah`,
      {
        method: "POST",
        headers: {
          ...getHeaders(),
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          kode: form.value.kode.trim(),
          name: form.value.name.trim(),
          sks: Number(form.value.sks),
        }),
      }
    )

    const result = await response.json()

    console.log("RESP MATA KULIAH:", result)

    if (!response.ok) {
      errorMessage.value =
        result?.message || "Gagal menambahkan mata kuliah"
      return
    }

    successMessage.value = "Berhasil menambahkan mata kuliah"

    form.value = {
      kode: "",
      name: "",
      sks: null,
    }

    setTimeout(() => {
      router.push("/dashboard-admin/matakuliah")
    }, 800)

  } catch (error) {
    console.error("SIMPAN ERROR:", error)
    errorMessage.value = "Terjadi kesalahan saat menyimpan"
  } finally {
    loading.value = false
  }
}

const kembali = () => {
  router.back()
}
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

```
<!-- BREADCRUMB -->
<div class="mb-2 text-sm text-gray-400">
  Akademik > Mata Kuliah > Tambah Mata Kuliah
</div>

<!-- TITLE -->
<h1 class="text-4xl font-bold text-gray-800">
  Tambah Mata Kuliah
</h1>

<p class="mb-6 mt-1 text-gray-500">
  Pengelolaan Data Mata Kuliah
</p>

<!-- FORM -->
<div class="rounded-2xl border border-blue-100 bg-white shadow-sm">

  <!-- HEADER -->
  <div class="border-b border-gray-100 px-5 py-4">
    <h2 class="text-2xl font-semibold text-gray-700">
      Form Mata Kuliah
    </h2>
  </div>

  <!-- CONTENT -->
  <div class="p-5">

    <div
      v-if="errorMessage"
      class="mb-4 rounded-lg bg-red-100 px-4 py-3 text-red-600"
    >
      {{ errorMessage }}
    </div>

    <div
      v-if="successMessage"
      class="mb-4 rounded-lg bg-green-100 px-4 py-3 text-green-600"
    >
      {{ successMessage }}
    </div>

    <div class="grid grid-cols-1 gap-5 md:grid-cols-2">

      <!-- KODE -->
      <div>
        <label class="mb-2 block text-sm font-medium text-gray-700">
          Kode Mata Kuliah
        </label>

        <input
          v-model="form.kode"
          type="text"
          placeholder="Contoh: MK001"
          class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
        />
      </div>

      <!-- NAMA -->
      <div>
        <label class="mb-2 block text-sm font-medium text-gray-700">
          Nama Mata Kuliah
        </label>

        <input
          v-model="form.name"
          type="text"
          placeholder="Masukkan nama mata kuliah"
          class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
        />
      </div>

      <!-- SKS -->
      <div>
        <label class="mb-2 block text-sm font-medium text-gray-700">
          SKS
        </label>

        <input
          v-model="form.sks"
          type="number"
          min="1"
          max="9"
          placeholder="Masukkan SKS"
          class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm outline-none focus:border-blue-500"
          @input="handleSksInput"
        />
      </div>

    </div>

  </div>

</div>

<!-- BUTTON -->
<div class="mt-6 flex gap-3">

  <button
    @click="kembali"
    class="rounded-xl border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 hover:bg-gray-100"
  >
    Kembali
  </button>

  <button
    @click="simpanMataKuliah"
    :disabled="loading"
    class="rounded-xl bg-green-500 px-6 py-3 text-sm font-semibold text-white shadow hover:bg-green-600 disabled:opacity-50"
  >
    {{ loading ? "Menyimpan..." : "💾 Simpan" }}
  </button>

</div>
```

  </div>
</template>