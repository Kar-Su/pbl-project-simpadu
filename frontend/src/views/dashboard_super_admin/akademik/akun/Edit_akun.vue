<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()

// ================= STATE =================
const email = ref("")
const nama = ref("")
const jabatan = ref("")
console.log("ROUTE:", route.fullPath)
console.log("PARAMS:", route.params)
console.log("ID:", route.params.id)

const loading = ref(false)

// ================= ROLE OPTION =================
const roleOptions = ref<string[]>([])

// ================= GET USER =================
const getUser = async () => {
  try {

    loading.value = true

    const token = localStorage.getItem("token")
    const id = route.params.id

    const res = await fetch(`/api/user/super/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    const data = await res.json()

    console.log("DETAIL USER:", data)

    if (!res.ok) {
      alert(data.message || "Gagal mengambil data user")
      return
    }

    email.value = data.data.email || ""
    nama.value = data.data.name || ""
    jabatan.value = data.data.role || ""

  } catch (err) {
    console.error(err)
    alert("Terjadi error")
  } finally {
    loading.value = false
  }
}
const getRoles = async () => {
  try {
    const token = localStorage.getItem("token")

    console.log("TOKEN ROLE:", token)

    const res = await fetch("/api/role", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${token}`
      }
    })

    const data = await res.json()

    console.log("ROLES RESPONSE:", data)

    if (!res.ok) {
      alert(data.error || "Gagal mengambil role")
      return
    }

    roleOptions.value = data.data.map(
      (item: any) => item.name
    )

    console.log("ROLE OPTIONS:", roleOptions.value)

  } catch (err) {
    console.error(err)
  }
}

// ================= UPDATE USER =================
const handleSimpan = async () => {
  try {

    const token = localStorage.getItem("token")
    const id = route.params.id

    const res = await fetch(`/api/super/user/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        email: email.value,
        name: nama.value,
        role: jabatan.value
      })
    })

    const data = await res.json()

    console.log("UPDATE:", data)

    if (!res.ok) {
      alert(data.message || "Gagal update akun")
      return
    }

    alert("Akun berhasil diupdate")

    router.push("/dashboard-superadmin/akun")

  } catch (err) {
    console.error(err)
    alert("Terjadi error")
  }
}

onMounted(() => {
  getUser()
  getRoles()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef3fb] p-6">

    <!-- BREADCRUMB -->
    <p class="text-sm text-gray-500 mb-2">
      Dashboard > Akun > Edit Akun
    </p>

    <!-- TITLE -->
    <h1 class="text-4xl font-bold text-[#1d1d1d]">
      Edit Akun
    </h1>

    <p class="text-gray-500 mt-2 mb-10">
      Edit data akun user
    </p>

    <!-- FORM -->
    <div
      class="bg-white rounded-2xl shadow-sm p-8 w-full max-w-2xl"
    >

      <!-- EMAIL -->
      <div class="mb-6">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Email
        </label>

        <input
          v-model="email"
          type="text"
          class="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:border-[#243e90]"
        />

      </div>

      <!-- NAMA -->
      <div class="mb-6">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Nama
        </label>

        <input
          v-model="nama"
          type="text"
          class="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:border-[#243e90]"
        />

      </div>

      <!-- ROLE -->
      <div class="mb-8">

        <label class="block text-sm font-medium text-gray-700 mb-2">
          Role
        </label>

        <select
          v-model="jabatan"
          class="w-full border border-gray-300 rounded-xl px-4 py-3 focus:outline-none focus:border-[#243e90]"
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

      </div>

      <!-- BUTTON -->
      <div class="flex items-center gap-3">

        <button
          @click="handleSimpan"
          :disabled="loading"
          class="bg-[#243e90] hover:bg-[#1a2f6d] text-white px-6 py-3 rounded-xl font-medium transition-all"
        >
          Simpan
        </button>

        <button
          @click="router.back()"
          class="bg-gray-200 hover:bg-gray-300 px-6 py-3 rounded-xl font-medium transition-all"
        >
          Kembali
        </button>

      </div>

    </div>

  </div>
</template>