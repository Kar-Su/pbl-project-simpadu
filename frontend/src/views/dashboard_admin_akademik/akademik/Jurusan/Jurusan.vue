<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const BASE_URL = "https://be.karlearn.site"

// ========================
// TABLE DATA
// ========================
const jurusanData = ref<any[]>([])

const currentPage = ref(1)
const perPage = ref(10)
const totalPages = ref(1)

const getHeaders = () => ({
  "Content-Type": "application/json",
  accept: "application/json",
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

// ========================
// GET JURUSAN
// ========================
const getJurusan = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan?page=${currentPage.value}&per_page=${perPage.value}`,
      {
        method: "GET",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    const data =
      json?.data?.items ??
      json?.data ??
      []

    jurusanData.value = data.map((item: any) => ({
      id: item.id,
      name: item.name.replace(/-/g, " "),
    }))

    totalPages.value =
      json?.data?.pagination?.total_pages ?? 1
  } catch (err) {
    console.error(err)
  }
}

// ========================
// TAMBAH
// ========================
const tambahData = () => {
  router.push(
    "/dashboard-admin/Tambah_jurusan"
  )
}

// ========================
// EDIT
// ========================
const showEditModal = ref(false)

const editForm = ref({
  oldName: "",
  newName: "",
})

const editData = (item: any) => {
  editForm.value = {
    oldName: item.name,
    newName: item.name,
  }

  showEditModal.value = true
}

const submitEdit = async () => {
  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan/${editForm.value.oldName}`,
      {
        method: "PUT",
        headers: getHeaders(),
        body: JSON.stringify({
          new_name: editForm.value.newName,
        }),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json.message)
      return
    }

    showEditModal.value = false
    await getJurusan()
  } catch (err) {
    console.error(err)
  }
}

// ========================
// HAPUS
// ========================
const hapusJurusan = async (nama: string) => {
  const konfirmasi = confirm(
    `Hapus jurusan ${nama}?`
  )

  if (!konfirmasi) return

  try {
    const res = await fetch(
      `${BASE_URL}/api/jurusan/${nama}`,
      {
        method: "DELETE",
        headers: getHeaders(),
      }
    )

    const json = await res.json()

    if (!res.ok) {
      alert(json.message)
      return
    }

    await getJurusan()
  } catch (err) {
    console.error(err)
  }
}

// ========================
// PAGINATION
// ========================
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 4)
    return Array.from(
      { length: total },
      (_, i) => i + 1
    )

  const pages: (number | string)[] = []

  pages.push(1)

  const start = Math.max(
    2,
    current - 1
  )

  const end = Math.min(
    total - 1,
    current + 1
  )

  if (start > 2) pages.push("...")

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  if (end < total - 1) {
    pages.push("...")
  }

  pages.push(total)

  return pages
})

const goToPage = async (page: number) => {
  currentPage.value = page
  await getJurusan()
}

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await getJurusan()
  }
}

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--
    await getJurusan()
  }
}

onMounted(async () => {
  await getJurusan()
})
</script>

<template>
  <div class="min-h-screen bg-[#eef4fb] p-5">

    <div class="text-sm text-gray-500 mb-2">
      Akademik > Jurusan
    </div>

    <h1
      class="text-[42px] font-bold text-[#404040]"
    >
      Jurusan
    </h1>

    <p
      class="text-gray-500 text-sm mt-2 mb-6"
    >
      Pengelolaan Data Jurusan
    </p>

    <div
      class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]"
    >
      <!-- HEADER -->
      <div class="p-5">
        <div
          class="flex justify-between items-center"
        >
          <h2
            class="text-[32px] font-semibold text-[#505050]"
          >
            Data Jurusan
          </h2>

          <button
            @click="tambahData"
            class="bg-[#29479d] hover:bg-[#1d377f] text-white px-6 py-3 rounded-xl font-semibold"
          >
            + Tambah
          </button>
        </div>
      </div>

      <!-- TABLE -->
      <div class="px-5">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-600 border-b border-gray-300">
              <th class="py-4">No</th>
              <th>Nama Jurusan</th>
              <th class="text-center">
                Aksi
              </th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="(item,index) in jurusanData"
              :key="item.id"
              class=" hover:bg-gray-50"
            >
              <td class="py-4">
                {{
                  (currentPage - 1) *
                    perPage +
                  index +
                  1
                }}
              </td>

              <td>
                {{ item.name }}
              </td>

              <td
                class="flex justify-center gap-2 py-3"
              >
                <button
                  @click="editData(item)"
                  class="bg-[#f3a317] hover:bg-[#d78e0f] text-white px-4 py-2 rounded-lg"
                >
                  Edit
                </button>

                <button
                  @click="
                    hapusJurusan(item.name)
                  "
                  class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg"
                >
                  Hapus
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- PAGINATION -->
      <div
        class="flex justify-end items-center gap-2 p-5"
      >
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
        >
          Previous
        </button>

        <template
          v-for="page in visiblePages"
          :key="page"
        >
          <span
            v-if="page === '...'"
          >
            ...
          </span>

          <button
            v-else
            @click="
              goToPage(Number(page))
            "
            :class="
              currentPage === Number(page)
                ? 'bg-[#29479d] text-white'
                : 'bg-white'
            "
            class="w-8 h-8 rounded"
          >
            {{ page }}
          </button>
        </template>

        <button
          @click="nextPage"
          :disabled="
            currentPage === totalPages
          "
        >
          Next
        </button>
      </div>
    </div>

    <!-- EDIT MODAL -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-black/40 flex justify-center items-center"
    >
      <div
        class="bg-white rounded-xl p-6 w-[400px]"
      >
        <h2
          class="text-xl font-bold mb-4"
        >
          Edit Jurusan
        </h2>

        <input
          v-model="editForm.newName"
          type="text"
          class="w-full border rounded-lg px-3 py-2"
        />

        <div
          class="flex justify-end gap-2 mt-5"
        >
          <button
            @click="
              showEditModal = false
            "
            class="border px-4 py-2 rounded"
          >
            Batal
          </button>

          <button
            @click="submitEdit"
            class="bg-[#29479d] text-white px-4 py-2 rounded"
          >
            Simpan
          </button>
        </div>
      </div>
    </div>

  </div>
</template>