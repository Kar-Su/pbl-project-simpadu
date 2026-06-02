<script setup lang="ts">
import { ref } from "vue"

const name = ref("Aka Demi")
const email = ref("akademi@gmail.id.sabar")
const role = ref("Admin Akademik")

const avatar = ref("https://i.pravatar.cc/300")

const fileInput = ref<HTMLInputElement | null>(null)

const triggerUpload = () => {
  fileInput.value?.click()
}

const onFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    avatar.value = URL.createObjectURL(file)
  }
}
</script>

<template>
  <div class="p-6">

    <!-- BREADCRUMB -->
    <div class="text-sm text-gray-400 mb-4">
      Dashboard > Profile
    </div>

    <!-- CARD -->
    <div class="bg-white rounded-2xl shadow-md border p-6 flex gap-10">

      <!-- LEFT -->
      <div class="flex flex-col items-center w-1/3">

        <img
          :src="avatar"
          class="w-40 h-40 rounded-full object-cover border-4 border-blue-200"
        />

        <button
          @click="triggerUpload"
          class="mt-4 bg-orange-500 hover:bg-orange-600 text-white px-4 py-2 rounded-lg flex items-center gap-2"
        >
          ✏️ Edit
        </button>

        <input
          type="file"
          ref="fileInput"
          class="hidden"
          accept="image/*"
          @change="onFileChange"
        />

      </div>

      <!-- RIGHT -->
      <div class="w-2/3 space-y-4">

        <div>
          <label class="text-sm text-gray-500">Nama Kelas</label>
          <input
            v-model="name"
            class="w-full border rounded-lg px-3 py-2"
          />
        </div>

        <div>
          <label class="text-sm text-gray-500">Email</label>
          <input
            v-model="email"
            class="w-full border rounded-lg px-3 py-2"
          />
        </div>

        <div>
          <label class="text-sm text-gray-500">Role</label>
          <input
            v-model="role"
            class="w-full border rounded-lg px-3 py-2"
            disabled
          />
        </div>

      </div>

    </div>

  </div>
</template>