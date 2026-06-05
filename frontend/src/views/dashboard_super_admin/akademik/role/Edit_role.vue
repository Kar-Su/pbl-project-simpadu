<script setup lang="ts">
import { ref, onMounted, watch } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()

// role lama yang dipilih
const selectedRole = ref("")

// role baru hasil edit
const editedRole = ref("")

const errorMsg = ref("")
const loading = ref(false)

// list role dropdown
const roleOptions = ref<string[]>([])

const roleId = route.params.id

// ================= GET ROLE =================
const fetchRoles = async () => {
  try {
    const token = localStorage.getItem("token")
    const BASE_URL = "https://be.karlearn.site"

    const response = await fetch(`${BASE_URL}/api/roles`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    const data = await response.json()

    roleOptions.value =
      data?.data?.role?.map((item: any) => item.name) || []

    // cari role sesuai id dari URL
    const currentRole = data?.data?.role?.find(
      (item: any) => item.id == roleId
    )

    if (currentRole) {
      selectedRole.value = currentRole.name
      editedRole.value = currentRole.name
    }

  } catch (err) {
    console.log(err)
  }
}

onMounted(() => {
  fetchRoles()
})

// ================= AUTO ISI INPUT EDIT =================
watch(selectedRole, (newValue) => {
  editedRole.value = newValue
})

// ================= SIMPAN =================
const handleSimpan = async () => {
  if (!selectedRole.value || !editedRole.value) {
    errorMsg.value = "Role wajib diisi"
    return
  }

  const token = localStorage.getItem("token")

  if (!token) {
    errorMsg.value = "Token tidak ditemukan"
    return
  }

  loading.value = true
  errorMsg.value = ""

  try {
    const response = await fetch(
      `/api/roles/${selectedRole.value}`,
      {
        method: "PUT",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
          authorization: `Bearer ${token}`
        },
        body: JSON.stringify({
          role_name: editedRole.value
        })
      }
    )

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || "Gagal edit role")
    }

    alert("Role berhasil diubah")

    router.push("/dashboard-superadmin/role")

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
      Akademik > Role > Edit Role
    </p>

    <p class="text-gray-500 mb-8">
      Edit Role
    </p>

    <!-- ERROR -->
    <p v-if="errorMsg" class="text-red-500 text-sm mb-4">
      {{ errorMsg }}
    </p>

    <!-- DROPDOWN ROLE -->
    <div class="mb-8">
      <label>Pilih Role</label>

      <select
        v-model="selectedRole"
        class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm"
      >
        <option disabled value="">
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

    <!-- EDIT ROLE -->
    <div class="mb-8">
      <label>Edit Role</label>

      <input
        v-model="editedRole"
        type="text"
        placeholder="Masukkan role baru..."
        class="w-full rounded-xl border border-gray-400 px-4 py-3 text-sm focus:outline-none focus:border-blue-500"
      />
    </div>

    <!-- BUTTON -->
    <button
      @click="handleSimpan"
      :disabled="loading"
      class="bg-green-500 text-white px-5 py-2 rounded"
    >
      {{ loading ? "Loading..." : "Simpan" }}
    </button>

  </div>
</template>