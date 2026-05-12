<script setup lang="ts">
const emit = defineEmits(["close", "confirm"])
</script>

<template>
  <div
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
  >
    <div
      class="bg-white rounded-2xl shadow-lg px-8 py-8 w-72 flex flex-col items-center text-center"
    >

      <!-- ICON -->
      <div class="mb-4 text-red-500 text-6xl">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-12">
  <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" />
</svg>

      </div>

      <!-- MESSAGE -->
      <p class="text-gray-900 font-semibold text-base leading-snug mb-6">
        Yakin ingin keluar?
      </p>

      <!-- BUTTON -->
      <div class="flex gap-3">

        <button
          @click="emit('confirm')"
          class="bg-red-500 hover:bg-red-600 text-white text-sm font-semibold px-5 py-1.5 rounded"
        >
          Yakin
        </button>

        <button
          @click="emit('close')"
          class="bg-gray-500 hover:bg-gray-600 text-white text-sm font-semibold px-4 py-1.5 rounded"
        >
          Kembali
        </button>

      </div>

    </div>
  </div>
</template>