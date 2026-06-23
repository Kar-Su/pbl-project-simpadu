<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const email = ref("")
const nama = ref("")
const jabatan = ref("")
const password = ref("")


const errorMsg = ref("")
const errors = ref({
  email: "",
  nama: "",
  password: "",
  jabatan: ""
})
// const passwordError = ref("")
const loading = ref(false)
const roleOptions = ref<string[]>([])

// const validatePassword = (pwd: string) => {
//   // const hasUpper = /[A-Z]/.test(pwd)
//   // const hasNumber = /[0-9]/.test(pwd)
//   // const hasSymbol = /[^A-Za-z0-9]/.test(pwd)

//   if (pwd.length < 8) {
//     return "Password minimal 8 karakter"
//   } 
//   // else if (!hasUpper || !hasNumber || !hasSymbol) {
//   //   return "Password harus ada huruf besar, angka, dan simbol"
//   // }

//   return ""
// }

const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

const detailId = ref(generateUUID())

const handleSimpan = async () => {
  // if (!email.value || !nama.value || !jabatan.value || !password.value) {
  //   errorMsg.value = "Semua field wajib diisi"
  //   return
  // }
errors.value = {
  email: "",
  nama: "",
  password: "",
  jabatan: ""
}

let isValid = true

if (!email.value.trim()) {
  errors.value.email = "Email wajib diisi"
  isValid = false
}

if (!nama.value.trim()) {
  errors.value.nama = "Nama wajib diisi"
  isValid = false
}

if (!password.value.trim()) {
  errors.value.password = "Password wajib diisi"
  isValid = false
}

if (!jabatan.value.trim()) {
  errors.value.jabatan = "Role wajib dipilih"
  isValid = false
}

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

if (email.value && !emailRegex.test(email.value)) {
  errors.value.email = "Format email tidak valid"
  isValid = false
}

if (password.value && password.value.length < 8) {
  errors.value.password = "Password minimal 8 karakter"
  isValid = false
}

if (!isValid) return
  if (jabatan.value === "super-admin") {
    errorMsg.value = "Tidak boleh membuat super-admin baru"
    return
  }

  const token = localStorage.getItem("token")
  if (!token) {
    errorMsg.value = "Token tidak ditemukan, login ulang"
    return
  }


  // passwordError.value = validatePassword(password.value)
  // if (passwordError.value) return

  loading.value = true
  errorMsg.value = ""

  try {
    const BASE_URL = 'https://be.karlearn.site'
    const response = await fetch(`${BASE_URL}/api/users/admins`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        detail_id: detailId.value,
        email: email.value,
        name: nama.value,
        password: password.value,
        role_name: jabatan.value
      })
    })

    
    let data: any = {}
    try {
      data = await response.json()
    } catch {
      data = {}
    }
    console.log("STATUS:", response.status)
    console.log("RESPONSE:", data)
    console.log("PAYLOAD:", {
      detail_id: detailId.value,
      email: email.value,
      name: nama.value,
      password: password.value,
      role_name: jabatan.value
    })

if (!response.ok) {
  console.error("ERROR RESPONSE:", data)

  const message =
    data?.error ||
    data?.message ||
    ""

  if (
    message.toLowerCase().includes("email") ||
    message.toLowerCase().includes("already") ||
    message.toLowerCase().includes("exists")
  ) {
    errors.value.email = "Email sudah digunakan"
    return
  }

  errorMsg.value = message || "Gagal menambahkan akun"
  return
}

    router.push("/dashboard-superadmin/akun")
  } catch (err: any) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

const getRoles = async () => {
  try {
    const token = localStorage.getItem("token")
    

    const res = await fetch("https://be.karlearn.site/api/roles", {
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: "application/json"
      }
    })

    const data = await res.json()

    console.log("ROLE RESPONSE:", data)

    if (!res.ok) {
      throw new Error(data.message || "Gagal mengambil role")
    }

    // lihat struktur sebenarnya
    console.log("DATA:", data.data)

    roleOptions.value = (data?.data?.role || []).map(
      (item: any) => item.name)
.filter((name: string) => name !== "super-admin")
    console.log("ROLE OPTIONS:", roleOptions.value)

  } catch (err) {
    console.error("GET ROLE ERROR:", err)
  }
}
onMounted(() => {
  getRoles()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- HEADER -->
    <p class="text-sm text-gray-500 mb-2">
      Akademik > Akun > Tambah Akun
    </p>

    <h1 class="text-4xl font-bold text-gray-900">
      Tambah Akun
    </h1>

    <p class="text-gray-600 mt-2 mb-8">
      Tambahkan Data yang Diinginkan
    </p>

    <!-- ERROR GLOBAL -->
    <p v-if="errorMsg" class="mb-4 text-sm text-red-500">
      {{ errorMsg }}
    </p>

    <!-- FORM -->
<div class="w-full max-w-xl">

  <!-- EMAIL -->
  <div class="mb-6">
    <label class="block text-sm text-gray-700 mb-2">
      Email
    </label>

    <input
      v-model="email"
      type="email"
      placeholder="Isi Email..."
      :class="[
        'w-full rounded-xl px-4 py-3 text-sm focus:outline-none',
        errors.email
          ? 'border border-red-500'
          : 'border border-gray-400 focus:border-blue-500'
      ]"
    />

    <p
      v-if="errors.email"
      class="text-red-500 text-sm mt-1"
    >
      {{ errors.email }}
    </p>
  </div>

  <!-- NAMA -->
  <div class="mb-6">
    <label class="block text-sm text-gray-700 mb-2">
      Nama
    </label>

    <input
      v-model="nama"
      type="text"
      placeholder="Isi Nama..."
      :class="[
        'w-full rounded-xl px-4 py-3 text-sm focus:outline-none',
        errors.nama
          ? 'border border-red-500'
          : 'border border-gray-400 focus:border-blue-500'
      ]"
    />

    <p
      v-if="errors.nama"
      class="text-red-500 text-sm mt-1"
    >
      {{ errors.nama }}
    </p>
  </div>

  <!-- PASSWORD -->
  <div class="mb-6">
    <label class="block text-sm text-gray-700 mb-2">
      Password
    </label>

    <input
      v-model="password"
      type="text"
      placeholder="Isi Password..."
      :class="[
        'w-full rounded-xl px-4 py-3 text-sm focus:outline-none',
        errors.password
          ? 'border border-red-500'
          : 'border border-gray-400 focus:border-blue-500'
      ]"
    />

    <p
      v-if="errors.password"
      class="text-red-500 text-sm mt-1"
    >
      {{ errors.password }}
    </p>
  </div>

  <!-- JABATAN -->
  <div class="mb-8">
    <label class="block text-sm text-gray-700 mb-2">
      Jabatan
    </label>

    <select
      v-model="jabatan"
      :class="[
        'w-full rounded-xl px-4 py-3 text-sm focus:outline-none',
        errors.jabatan
          ? 'border border-red-500'
          : 'border border-gray-400 focus:border-blue-500'
      ]"
    >
      <option
        disabled
        value=""
      >
        Pilih Role
      </option>

      <option
        v-for="role in roleOptions"
        :key="role"
        :value="role"
      >
        {{ role }}
      </option>
    </select>

    <p
      v-if="errors.jabatan"
      class="text-red-500 text-sm mt-1"
    >
      {{ errors.jabatan }}
    </p>
  </div>

  <!-- BUTTON -->
<!-- BUTTON -->
<div class="flex items-center gap-3">

  <button
    @click="handleSimpan"
    :disabled="loading"
    class="flex items-center gap-2 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 text-white px-5 py-3 rounded-xl text-sm font-semibold"
  >
    {{ loading ? "Menyimpan..." : "Simpan" }}
  </button>

  <button
    type="button"
    @click="router.push('/dashboard-superadmin/akun')"
    class="flex items-center gap-2 bg-gray-300 hover:bg-gray-400 text-gray-800 px-5 py-3 rounded-xl text-sm font-semibold"
  >
    Kembali
  </button>

</div>

</div>

  </div>
</template>