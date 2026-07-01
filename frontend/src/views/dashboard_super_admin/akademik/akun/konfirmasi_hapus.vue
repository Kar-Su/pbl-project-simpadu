<script setup lang="ts">
defineProps<{
  message?: string
}>()

const emit = defineEmits(["close", "confirm"])

const handleClose = () => {
  emit("close")
}

const handleConfirm = () => {
  emit("confirm")
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
    <div
      class="w-80 rounded-2xl bg-white px-8 py-8 shadow-lg flex flex-col items-center text-center"
    >
      <!-- ICON -->
      <div class="mb-4 text-6xl text-red-500">
        ⚠️
      </div>

      <!-- TEXT -->
      <p class="mb-6 text-base font-semibold leading-relaxed text-gray-900">
        {{ message }}
      </p>

      <!-- BUTTON -->
      <div class="flex gap-3">
        <button
          @click="handleConfirm"
          class="rounded bg-red-500 px-5 py-2 text-sm font-semibold text-white hover:bg-red-600 active:scale-95 transition"
        >
          Yakin
        </button>

        <button
          @click="handleClose"
          class="rounded bg-gray-500 px-5 py-2 text-sm font-semibold text-white hover:bg-gray-600 active:scale-95 transition"
        >
          Kembali
        </button>
      </div>
    </div>
  </div>
</template>