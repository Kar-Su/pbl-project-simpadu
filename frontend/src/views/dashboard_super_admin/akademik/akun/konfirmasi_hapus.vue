<script setup lang="ts">
defineProps<{
  email?: string
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
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
  >

    <!-- MODAL -->
    <div
      class="w-72 rounded-2xl bg-white px-8 py-8 shadow-lg flex flex-col items-center text-center"
    >

      <!-- ICON -->
      <div class="mb-4 text-6xl text-red-500">
        ⚠
      </div>

      <!-- TEXT -->
      <p class="mb-6 text-base font-semibold leading-snug text-gray-900">
        Hapus bayar 100K
        DANA 0828282828<br />
        {{ email }}
      </p>

      <!-- BUTTON -->
      <div class="flex gap-3">

        <!-- YAKIN -->
        <button
          @click="handleConfirm"
          class="rounded bg-red-500 px-5 py-1.5 text-sm font-semibold text-white hover:bg-red-600 active:scale-95 transition"
        >
          Yakin
        </button>

        <!-- KEMBALI -->
        <button
          type="button"
          @click="handleClose"
          class="rounded bg-gray-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-gray-600 active:scale-95 transition"
        >
          Kembali
        </button>

      </div>

    </div>

  </div>
</template>