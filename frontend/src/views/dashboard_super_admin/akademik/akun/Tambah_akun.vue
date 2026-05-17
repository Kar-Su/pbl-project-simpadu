<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const email = ref("")
const nama = ref("")
const jabatan = ref("")
const password = ref("")
const errorMsg = ref("")
const passwordError = ref("")
const loading = ref(false)

const validatePassword = (pwd: string) => {
  const hasUpper = /[A-Z]/.test(pwd)
  const hasNumber = /[0-9]/.test(pwd)
  const hasSymbol = /[^A-Za-z0-9]/.test(pwd)
  if (!hasUpper || !hasNumber || !hasSymbol) return "Password harus ada huruf besar, angka, dan simbol"
  return ""
}

const handleSimpan = async () => {
  if (!email.value || !nama.value || !jabatan.value || !password.value) {
    errorMsg.value = "Semua field wajib diisi"
    return
  }
  if (jabatan.value === "super-admin") {
    errorMsg.value = "Tidak boleh membuat super-admin baru"
    return
  }
  const token = localStorage.getItem("token")
  if (!token) { errorMsg.value = "Token tidak ditemukan, login ulang"; return }

  passwordError.value = validatePassword(password.value)
  if (passwordError.value) return

  loading.value = true
  errorMsg.value = ""

  try {
    const response = await fetch(`/api/super/user`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({
        name: nama.value,
        email: email.value,
        password: password.value,
        role_name: jabatan.value,
      })
    })

    let data: any = {}
    try { data = await response.json() } catch { data = {} }

    if (!response.ok) throw new Error(data?.message || "Gagal menambah akun")

    router.push("/dashboard-superadmin/akun")
  } catch (err: any) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}
</script>
<template>
  <div class="min-h-screen bg-[#f5f7fb] p-6">

    <p class="text-sm text-gray-400 mb-2">
      Akademik > Akun > Tambah Akun
    </p>

    <h1 class="text-3xl font-bold text-gray-800 mb-1">
      Tambah Akun
    </h1>

    <p class="text-gray-500 mb-8">
      Tambahkan Data yang Diinginkan
    </p>

    <p v-if="errorMsg" class="text-red-500 text-sm mb-4">
      {{ errorMsg }}
    </p>

    <!-- EMAIL -->
    <div class="mb-5">
      <label>Email</label>
      <input v-model="email" type="email" class="w-96 border px-4 py-2 rounded" />
    </div>

    <!-- NAMA -->
    <div class="mb-5">
      <label>Nama</label>
      <input v-model="nama" type="text" class="w-96 border px-4 py-2 rounded" />
    </div>

    <!-- PASSWORD -->
    <div class="mb-5">
      <label>Password</label>
      <input v-model="password" type="text" class="w-96 border px-4 py-2 rounded" />

      <p v-if="passwordError" class="text-red-500 text-xs mt-1">
        {{ passwordError }}
      </p>
    </div>

    <!-- JABATAN -->
    <div class="mb-8">
      <label>Jabatan</label>
      <select v-model="jabatan" class="w-96 border px-4 py-2 rounded">
        <option disabled value="">Pilih Jabatan...</option>
        <option value="super-admin">Super admin</option>
        <option value="admin-akademik">Admin akademik</option>
        <option value="admin-pegawai">Admin pegawai</option>
        <option value="admin-mahasiswa">Admin mahasiswa</option>
        <option value="admin-keuangan">Admin keuangan</option>
        <option value="dosen">Dosen</option>
        <option value="mahasiswa">Mahasiswa</option>
      </select>
    </div>

    <!-- BUTTON -->
    <button @click="handleSimpan" class="bg-green-500 text-white px-5 py-2 rounded">
      Simpan
    </button>

  </div>
</template>