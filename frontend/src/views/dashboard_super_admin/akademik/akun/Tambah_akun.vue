<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const email = ref("")
const nama = ref("")
const jabatan = ref("")
// const password = ref("")

const errorMsg = ref("")
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

const handleSimpan = async () => {
  // if (!email.value || !nama.value || !jabatan.value || !password.value) {
  //   errorMsg.value = "Semua field wajib diisi"
  //   return
  // }

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
    const BASE_URL = import.meta.env.VITE_API_BASE_URL
    const response = await fetch(`${BASE_URL}/api/users/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        name: nama.value,
        email: email.value,
        // password: password.value,
        role_name: jabatan.value
      })
    })

    let data: any = {}
    try {
      data = await response.json()
    } catch {
      data = {}
    }

    if (!response.ok) {
      throw new Error(data?.message || "Gagal menambah akun")
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
    const BASE_URL = 'https://be.karlearn.site'
    const res = await fetch(`${BASE_URL}/api/roles`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    const data = await res.json()

    console.log("ROLES:", data)

    if (!res.ok) {
      alert(data.message || "Gagal mengambil role")
      return
    }

    roleOptions.value = data.data.map(
      (item: any) => item.name
    )

  } catch (err) {
    console.error(err)
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
        <label class="block text-sm text-gray-700 mb-2">Email</label>
        <input
          v-model="email"
          type="email"
          placeholder="Isi Email..."
          class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
      </div>

      <!-- NAMA -->
      <div class="mb-6">
        <label class="block text-sm text-gray-700 mb-2">Nama</label>
        <input
          v-model="nama"
          type="text"
          placeholder="Isi Nama..."
          class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />
      </div>

      <!-- PASSWORD -->
      <!-- <div class="mb-6">
        <label class="block text-sm text-gray-700 mb-2">Password</label>

        <input
          v-model="password"
          type="text"
          placeholder="Isi Password..."
          class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        />

        <p
          v-if="passwordError"
          class="text-red-500 text-xs mt-2"
        >
          {{ passwordError }}
        </p>
      </div> -->

      <!-- JABATAN -->
      <div class="mb-8">
        <label class="block text-sm text-gray-700 mb-2">Jabatan</label>

        <select
          v-model="jabatan"
          class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
        >
        <option
  v-for="role in roleOptions"
  :key="role"
  :value="role"
>
  {{ role }}
</option>
        </select>
      </div>

      <!-- BUTTON -->
      <button
        @click="handleSimpan"
        class="flex items-center gap-2 bg-green-500 hover:bg-green-600 text-white px-5 py-3 rounded-xl text-sm font-semibold"
      >
        Simpan
      </button>

    </div>

  </div>
</template>