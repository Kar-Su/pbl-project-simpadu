<script setup lang="ts">
defineProps<{
  email?: string
}>()

const emit = defineEmits(["close", "confirm"])
</script>

<template>
  <div
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
  >
    <div
      class="bg-white rounded-2xl shadow-lg px-8 py-8 w-72 flex flex-col items-center text-center"
    >

      <div class="mb-4 text-red-500 text-6xl">
        ⚠
      </div>

      <p class="text-gray-900 font-semibold text-base leading-snug mb-6">
        Yakin ingin reset password?<br />
        {{ email }}
      </p>

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