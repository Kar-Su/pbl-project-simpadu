```vue
<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"

const router = useRouter()
const route = useRoute()

const kelasId = route.params.id as string

const BASE_URL = "https://be.karlearn.site"

const getHeaders = () => ({
    "Content-Type": "application/json",
    accept: "application/json",
    Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
})

const loading = ref(false)
const saving = ref(false)

const tahunAkademik = ref("")
const jurusan = ref("")
const prodi = ref("")
const namaKelas = ref("")
const kurikulum = ref("")
const semester = ref("")

const tahunAkademikList = ref<any[]>([])
const jurusanList = ref<any[]>([])
const prodiList = ref<any[]>([])
const kurikulumList = ref<any[]>([])

const filteredProdiList = computed(() => {
    if (!jurusan.value) return prodiList.value

    return prodiList.value.filter(
        (item) =>
            String(item.jurusan?.id) === String(jurusan.value)
    )
})

const getTahunAkademik = async () => {
    try {
        const res = await fetch(
            `${BASE_URL}/api/tahun-akademik?per_page=100`,
            {
                headers: getHeaders(),
            }
        )

        const data = await res.json()

        tahunAkademikList.value =
            data?.data?.items ??
            data?.data ??
            []
    } catch (err) {
        console.error(err)
    }
}

const getJurusan = async () => {
    try {
        const res = await fetch(
            `${BASE_URL}/api/jurusan`,
            {
                headers: getHeaders(),
            }
        )

        const data = await res.json()

        jurusanList.value =
            data?.data?.items ??
            data?.data ??
            []
    } catch (err) {
        console.error(err)
    }
}

const getProdi = async () => {
    try {
        const res = await fetch(
            `${BASE_URL}/api/prodi`,
            {
                headers: getHeaders(),
            }
        )

        const data = await res.json()

        prodiList.value =
            data?.data?.items ??
            data?.data ??
            []
    } catch (err) {
        console.error(err)
    }
}

const getKurikulum = async () => {
    try {
        const res = await fetch(
            `${BASE_URL}/api/kurikulum?per_page=100`,
            {
                headers: getHeaders(),
            }
        )

        const data = await res.json()

        kurikulumList.value =
            data?.data?.items ??
            data?.data ??
            []
    } catch (err) {
        console.error(err)
    }
}

const getDetailKelas = async () => {
    if (!kelasId) return

    loading.value = true

    try {
        const res = await fetch(
            `${BASE_URL}/api/kelas/${kelasId}`,
            {
                headers: getHeaders(),
            }
        )

        const data = await res.json()

        const item = data?.data

        if (!item) return

        namaKelas.value = item.name ?? ""

        jurusan.value =
            String(item.prodi?.jurusan?.id ?? "")

        prodi.value =
            String(item.prodi?.id ?? "")

        tahunAkademik.value =
            String(item.tahun_akademik?.id ?? "")

        kurikulum.value =
            String(item.kurikulum?.kode ?? "")

        semester.value =
            String(item.semester ?? "")
    } catch (err) {
        console.error(err)
    } finally {
        loading.value = false
    }
}

const simpanData = async () => {
    try {
        saving.value = true

        const payload = {
            name: namaKelas.value,
            prodi_id: prodi.value,
            tahun_akademik_id: Number(tahunAkademik.value),
            kurikulum_kode: kurikulum.value,
            semester: Number(semester.value),
        }

        const res = await fetch(
            `${BASE_URL}/api/kelas/${kelasId}`,
            {
                method: "PUT",
                headers: getHeaders(),
                body: JSON.stringify(payload),
            }
        )

        const data = await res.json()

        if (!res.ok) {
            alert(data?.message ?? "Gagal update kelas")
            return
        }

        alert("Data kelas berhasil diperbarui")

        router.push("/dashboard-admin/kelas")
    } catch (err) {
        console.error(err)
        alert("Terjadi kesalahan")
    } finally {
        saving.value = false
    }
}

watch(jurusan, () => {
    const masihAda = filteredProdiList.value.some(
        (item) => String(item.id) === String(prodi.value)
    )

    if (!masihAda) {
        prodi.value = ""
    }
})

onMounted(async () => {
    await Promise.all([
        getTahunAkademik(),
        getJurusan(),
        getProdi(),
        getKurikulum(),
    ])

    await getDetailKelas()
})
</script>

<template>
    <div class="min-h-screen bg-[#eef4fb] p-5">

        <div class="text-sm text-gray-500 mb-2">
            Akademik > Kelas > Edit Kelas
        </div>

        <h1 class="text-[40px] font-bold text-[#404040]">
            Edit Kelas
        </h1>

        <p class="text-gray-500 mb-6">
            Pembaharuan data kelas
        </p>

        <!-- FORM AKADEMIK -->
        <div
            class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden mb-5">
            <div class="bg-[#243e90] px-5 py-3">
                <h2 class="text-white text-2xl font-bold">
                    Form Akademik
                </h2>

                <p class="text-white text-sm">
                    Perubahan data akademik
                </p>
            </div>

            <div class="p-5 grid grid-cols-3 gap-4">

                <div>
                    <label class="block mb-2 text-sm">
                        Tahun Akademik *
                    </label>

                    <select v-model="tahunAkademik" class="w-full h-[45px] border rounded-lg px-3">
                        <option value="">
                            Pilih Tahun Akademik
                        </option>

                        <option v-for="item in tahunAkademikList" :key="item.id" :value="item.id">
                            {{ item.tahun_awal }} -
                            {{ item.tahun_akhir }}
                        </option>
                    </select>
                </div>

                <div>
                    <label class="block mb-2 text-sm">
                        Jurusan *
                    </label>

                    <select v-model="jurusan" class="w-full h-[45px] border rounded-lg px-3">
                        <option value="">
                            Pilih Jurusan
                        </option>

                        <option v-for="item in jurusanList" :key="item.id" :value="item.id">
                            {{ item.name }}
                        </option>
                    </select>
                </div>

                <div>
                    <label class="block mb-2 text-sm">
                        Prodi *
                    </label>

                    <select v-model="prodi" class="w-full h-[45px] border rounded-lg px-3">
                        <option value="">
                            Pilih Prodi
                        </option>

                        <option v-for="item in filteredProdiList" :key="item.id" :value="item.id">
                            {{ item.name }}
                        </option>
                    </select>
                </div>

            </div>
        </div>

        <!-- FORM KELAS -->
        <div class="bg-[#ececec] rounded-xl shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc] overflow-hidden">
            <div class="bg-[#243e90] px-5 py-3">
                <h2 class="text-white text-2xl font-bold">
                    Form Kelas
                </h2>

                <p class="text-white text-sm">
                    Perubahan data kelas
                </p>
            </div>

            <div class="p-5">

                <div class="mb-4">
                    <label class="block mb-2 text-sm">
                        Nama Kelas *
                    </label>

                    <input v-model="namaKelas" type="text" placeholder="Isi Nama Kelas"
                        class="w-full h-[45px] border rounded-lg px-3" />
                </div>

                <div class="grid grid-cols-2 gap-4">

                    <div>
                        <label class="block mb-2 text-sm">
                            Kurikulum *
                        </label>

                        <select v-model="kurikulum" class="w-full h-[45px] border rounded-lg px-3">
                            <option value="">
                                Pilih Kurikulum
                            </option>

                            <option v-for="item in kurikulumList" :key="item.kode" :value="item.kode">
                                {{ item.name }}
                            </option>
                        </select>
                    </div>

                    <div>
                        <label class="block mb-2 text-sm">
                            Semester *
                        </label>

                        <select v-model="semester" class="w-full h-[45px] border rounded-lg px-3">
                            <option value="">
                                Pilih Semester
                            </option>

                            <option v-for="n in 14" :key="n" :value="n">
                                Semester {{ n }}
                            </option>
                        </select>
                    </div>

                </div>

                <div class="mt-6">
                    <button @click="simpanData" :disabled="saving"
                        class="bg-[#243e90] text-white px-6 py-2 rounded-lg font-semibold hover:bg-[#1d3275]">
                        {{ saving ? "Menyimpan..." : "💾 Simpan" }}
                    </button>
                </div>

            </div>
        </div>

    </div>
</template>
```
