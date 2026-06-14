<script setup lang="ts">
import { ref, onMounted } from "vue"

const name = ref("")
const email = ref("")
const role = ref("")
const userId = ref("")
const detailId = ref("")

const BASE_URL = "https://be.karlearn.site"

const avatar = ref("https://i.pravatar.cc/300")

const fileInput = ref<HTMLInputElement | null>(null)

const isEditing = ref(false)

// ================= HELPER =================
const getHeaders = (): Record<string, string> => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ================= DECODE JWT =================
const getEmailFromToken = (): string => {
  try {
    const token = localStorage.getItem("token") ?? ""
    const payload = token.split(".")[1]
    const decoded = JSON.parse(atob(payload))
    return decoded.user_email ?? ""
  } catch {
    return ""
  }
}

const getDetailIdFromToken = (): string => {
  try {
    const token = localStorage.getItem("token") ?? ""
    const payload = token.split(".")[1]
    const decoded = JSON.parse(atob(payload))
    return decoded.detail_id ?? ""
  } catch {
    return ""
  }
}

// ================= FORMAT ROLE =================
// Hilangkan tanda "-" dari role, contoh: "admin-akademik" → "Admin Akademik"
const formatRole = (raw: string): string => {
  return raw
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ")
}

// ================= GET PROFILE =================
const getProfile = async (): Promise<void> => {
  try {
    const emailLogin = getEmailFromToken()

    if (!emailLogin) {
      console.warn("Email tidak ditemukan di token")
      return
    }

    const res = await fetch(
      `${BASE_URL}/api/users/search?email=${encodeURIComponent(emailLogin)}`,
      { headers: getHeaders() }
    )

    const data = await res.json()

    console.log("PROFILE:", data)

    name.value   = data.data?.name      ?? ""
    email.value  = data.data?.email     ?? ""
    role.value   = formatRole(data.data?.role_name ?? "")
    userId.value = data.data?.id        ?? ""

    if (data.data?.image_url) {
      avatar.value = data.data.image_url
    }

    // Simpan nama ke localStorage supaya dashboard bisa baca
    localStorage.setItem("user_name", name.value)
  } catch (err) {
    console.error("getProfile:", err)
  }
}

// ================= TOGGLE EDIT =================
const toggleEdit = () => {
  isEditing.value = !isEditing.value
}

const triggerUpload = () => {
  if (!isEditing.value) return
  fileInput.value?.click()
}

const onFileChange = (e: Event) => {
  if (!isEditing.value) return

  const target = e.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    avatar.value = URL.createObjectURL(file)
  }
}

// ================= SAVE PROFILE =================
const saveProfile = async (): Promise<void> => {
  try {
    const res = await fetch(`${BASE_URL}/api/users/${userId.value}`, {
      method: "PUT",
      headers: getHeaders(),
      body: JSON.stringify({
        detail_id: detailId.value,
        name: name.value,
      }),
    })

    const data = await res.json()

    console.log("SAVE PROFILE:", data)

    if (data.success) {
      name.value = data.data?.name ?? name.value

      // Update localStorage supaya dashboard ikut berubah
      localStorage.setItem("user_name", name.value)

      isEditing.value = false
    } else {
      console.error("Gagal simpan:", data.message)
      alert(data.message ?? "Gagal menyimpan perubahan")
    }
  } catch (err) {
    console.error("saveProfile:", err)
  }
}

// ================= ON MOUNTED =================
onMounted(() => {
  detailId.value = getDetailIdFromToken()
  getProfile()
})
</script>

<template>
  <div class="p-6">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-400 mb-4">
      Dashboard > Profile
    </div>

    <!-- CARD -->
    <div
      class="bg-white rounded-2xl shadow-md border p-10 flex items-center gap-16 min-h-[500px]"
    >

      <!-- FOTO -->
      <div class="flex flex-col items-center justify-center w-2/5">

        <img
          :src="avatar"
          @click="triggerUpload"
          class="w-64 h-64 rounded-full object-cover border-[6px] border-blue-200 transition shadow-lg"
          :class="isEditing
            ? 'cursor-pointer hover:opacity-80'
            : 'cursor-default'"
        />

        <p v-if="isEditing" class="text-xs text-gray-500 mt-3 text-center">
          Klik foto untuk mengganti
        </p>

        <input
          type="file"
          ref="fileInput"
          class="hidden"
          accept="image/*"
          @change="onFileChange"
        />

      </div>

      <!-- FORM -->
      <div class="w-2/3 space-y-4">

        <!-- NAMA -->
        <div>
          <label class="text-sm text-gray-500">Nama</label>

          <input
            v-model="name"
            :disabled="!isEditing"
            class="w-full border rounded-lg px-3 py-2 mt-1"
            :class="!isEditing
              ? 'bg-gray-100 text-gray-500 cursor-not-allowed'
              : 'bg-white'"
          />
        </div>

        <!-- EMAIL -->
        <div>
          <label class="text-sm text-gray-500">Email</label>

          <input
            v-model="email"
            disabled
            class="w-full border rounded-lg px-3 py-2 mt-1 bg-gray-100 text-gray-500 cursor-not-allowed"
          />
        </div>

        <!-- ROLE -->
        <div>
          <label class="text-sm text-gray-500">Role</label>

          <input
            v-model="role"
            disabled
            class="w-full border rounded-lg px-3 py-2 mt-1 bg-gray-100 text-gray-500 cursor-not-allowed"
          />
        </div>

        <!-- BUTTON -->
        <!-- <div class="flex gap-3 pt-2">

          <button
            @click="toggleEdit"
            class="bg-orange-500 hover:bg-orange-600 text-white px-5 py-2 rounded-lg font-medium"
          >
            {{ isEditing ? "Batal Edit" : "✏️ Edit Profil" }}
          </button>

          <button
            v-if="isEditing"
            @click="saveProfile"
            class="bg-[#243e90] hover:bg-[#1d377f] text-white px-5 py-2 rounded-lg font-medium"
          >
            Simpan Perubahan
          </button>

        </div> -->

      </div>

    </div>

  </div>
</template>