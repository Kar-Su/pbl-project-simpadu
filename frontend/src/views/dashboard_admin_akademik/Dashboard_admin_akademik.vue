<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import KonfirmasiKeluar from "@/views/dashboard_super_admin/akademik/konfirmasi_keluar.vue";

const router = useRouter();
const route = useRoute();

const goTo = (path: string) => router.push(path);
const isActive = (path: string) => route.path === path;

// ─────────────────────────────────────────────
// SIDEBAR STATE
// ─────────────────────────────────────────────
const isSidebarOpen = ref<boolean>(true);
const openAkademik = ref<boolean>(true);
const openMahasiswa = ref<boolean>(false);
const showLogoutPopup = ref<boolean>(false);

// ─────────────────────────────────────────────
// NAMA USER (untuk header)
// ─────────────────────────────────────────────
const userName = ref<string>(localStorage.getItem("user_name") ?? "")
const userAvatar = ref<string>(localStorage.getItem("user_avatar") ?? "https://i.pravatar.cc/100")

const refreshUserInfo = (): void => {
    userName.value = localStorage.getItem("user_name") ?? ""
    userAvatar.value = localStorage.getItem("user_avatar") ?? "https://i.pravatar.cc/100"
}

// ─────────────────────────────────────────────
// HELPER: ambil token dari localStorage
// ─────────────────────────────────────────────
const getHeaders = (): Record<string, string> => ({
    "Content-Type": "application/json",
    accept: "application/json",
    Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
});

// ─────────────────────────────────────────────
// DECODE JWT → ambil email untuk GET profile
// ─────────────────────────────────────────────
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

// ─────────────────────────────────────────────
// GET PROFILE (untuk isi nama & avatar di header)
// ─────────────────────────────────────────────
const getProfile = async (): Promise<void> => {
    try {
        const emailLogin = getEmailFromToken()
        if (!emailLogin) return

        const BASE_URL = "https://be.karlearn.site"
        const res = await fetch(
            `${BASE_URL}/api/users/search?email=${encodeURIComponent(emailLogin)}`,
            { headers: getHeaders() }
        )

        const data = await res.json()

        if (data.data?.name) {
            userName.value = data.data.name
            localStorage.setItem("user_name", data.data.name)
        }

        if (data.data?.image_url) {
            userAvatar.value = data.data.image_url
            localStorage.setItem("user_avatar", data.data.image_url)
        }
    } catch (err) {
        console.error("getProfile:", err)
    }
}

// ─────────────────────────────────────────────
// LOGOUT
// ─────────────────────────────────────────────
const handleLogout = (): void => {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("role");
    localStorage.removeItem("user_name");
    localStorage.removeItem("user_avatar");
    router.push("/");
};

// ─────────────────────────────────────────────
// STAT CARDS
// ─────────────────────────────────────────────
const totalPegawai = ref<number>(0);
const totalMahasiswa = ref<number>(0);
const totalDosen = ref<number>(0);

const getAllUsers = async (): Promise<any[]> => {
    try {
        const BASE_URL = "https://be.karlearn.site";
        let page = 1;
        let lastPage = 1;
        let allUsers: any[] = [];

        do {
            const res = await fetch(`${BASE_URL}/api/users?page=${page}`, {
                headers: getHeaders(),
            });
            const data = await res.json();
            const items = data.data.items ?? [];
            allUsers = [...allUsers, ...items];
            lastPage = data.data.pagination?.total_pages ?? 1;
            page++;
        } while (page <= lastPage);

        return allUsers;
    } catch (err) {
        console.error("getAllUsers:", err);
        return [];
    }
};

const getTotalPegawaiExternal = async (): Promise<number> => {
    try {
        const res = await fetch(
            "https://api-pegawai-4a.akufarish.my.id:1234/api/employees/info/count",
            { headers: getHeaders() }
        );
        const data = await res.json();
        return data.data?.total_employee ?? 0;
    } catch (err) {
        console.error("getTotalPegawaiExternal:", err);
        return 0;
    }
};

const getMahasiswaExternal = async (): Promise<any[]> => {
    try {
        const res = await fetch(
            "https://api-mahasiswa-4a.akufarish.my.id:8874/api/mahasiswa",
            { headers: getHeaders() }
        );
        const data = await res.json();
        return Array.isArray(data.data) ? data.data : [];
    } catch (err) {
        console.error("getMahasiswaExternal:", err);
        return [];
    }
};

// ─────────────────────────────────────────────
// TAHUN AKADEMIK
// ─────────────────────────────────────────────
interface TahunAkademik {
    id: number;
    tahun_awal: string;
    tahun_akhir: string;
    tipee_semester: string;
}
const tahunAkademik = ref<TahunAkademik[]>([]);

const getTahunAkademik = async (): Promise<void> => {
    try {
        const res = await fetch("https://be.karlearn.site/api/tahun-akademik", {
            headers: getHeaders(),
        });
        const data = await res.json();
        tahunAkademik.value = Array.isArray(data.data) ? data.data : [];
    } catch (err) {
        console.error("getTahunAkademik:", err);
    }
};

// ─────────────────────────────────────────────
// KURIKULUM
// ─────────────────────────────────────────────
interface Kurikulum {
    id: string;
    name: string;
}
const kurikulum = ref<Kurikulum[]>([]);

const formatYear = (date: string): string => {
    return new Date(date).getFullYear().toString();
};

const getKurikulum = async (): Promise<void> => {
    try {
        const res = await fetch("https://be.karlearn.site/api/kurikulum?page=1", {
            headers: getHeaders(),
        });
        const data = await res.json();
        kurikulum.value = data.data.items ?? [];
    } catch (err) {
        console.error("getKurikulum:", err);
    }
};

// ─────────────────────────────────────────────
// MAHASISWA
// ─────────────────────────────────────────────
interface MahasiswaItem {
    id: string;
    nim: string;
    name: string;
    kelas: string;
    prodi: string;
    angkatan: string;
}

const akunList = ref<MahasiswaItem[]>([]);
const currentPage = ref<number>(1);
const perPage = ref<number>(10);
const totalItems = ref<number>(0);
const allMahasiswa = ref<MahasiswaItem[]>([]);

const totalPages = computed<number>(() =>
    Math.max(1, Math.ceil(totalItems.value / perPage.value))
);

const mapMahasiswa = (raw: any[]): MahasiswaItem[] =>
    raw.map((item: any) => ({
        id: item.id_mahasiswa ?? "-",
        nim: item.nim ?? "-",
        name: item.nama_mahasiswa ?? "-",
        kelas: "-",
        prodi: item.prodi_id ? String(item.prodi_id) : "-",
        angkatan: item.tahunakademik_id ? String(item.tahunakademik_id) : "-",
    }));

const pages = computed<(number | string)[]>(() => {
    const total = totalPages.value;
    const cur = currentPage.value;
    if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

    const result: (number | string)[] = [1, 2];
    if (cur > 4) result.push("...");
    for (let i = Math.max(3, cur - 1); i <= Math.min(total - 2, cur + 1); i++) {
        result.push(i);
    }
    if (cur < total - 3) result.push("...");
    result.push(total - 1, total);
    return [...new Set(result)];
});

const goToPage = (page: number): void => {
    if (page < 1 || page > totalPages.value) return;
    currentPage.value = page;
    const start = (page - 1) * perPage.value;
    const end = start + perPage.value;
    akunList.value = allMahasiswa.value.slice(start, end);
};

const prevPage = (): void => goToPage(currentPage.value - 1);
const nextPage = (): void => goToPage(currentPage.value + 1);

// ─────────────────────────────────────────────
// ON MOUNTED
// ─────────────────────────────────────────────
onMounted(async (): Promise<void> => {
    // Refresh nama dari localStorage dulu (hasil simpan di profile)
    refreshUserInfo()

    // Lalu fetch ulang dari API untuk data terbaru
    getProfile()

    const [users, externalPegawaiTotal, mahasiswaRaw] = await Promise.all([
        getAllUsers(),
        getTotalPegawaiExternal(),
        getMahasiswaExternal(),
    ]);

    totalMahasiswa.value = mahasiswaRaw.length;

    const karlearnPegawai = users.filter((item: any) => {
        const role = item.role_name?.toLowerCase()?.trim();
        return role !== "mahasiswa";
    }).length;
    totalPegawai.value = karlearnPegawai + externalPegawaiTotal;

    const dosenKarlearn = users.filter((item: any) => {
        const role = item.role_name?.toLowerCase()?.trim();
        return role === "dosen";
    }).length;

    totalDosen.value = dosenKarlearn;

    allMahasiswa.value = mapMahasiswa(mahasiswaRaw);
    totalItems.value = allMahasiswa.value.length;
    akunList.value = allMahasiswa.value.slice(0, perPage.value);

    getTahunAkademik();
    getKurikulum();
});
</script>

<template>
    <div class="bg-[#eef3fb] min-h-screen overflow-y-auto overflow-x-hidden">

        <!-- ═══════════════════════════════════════
             HEADER
        ════════════════════════════════════════ -->
        <header
            class="fixed top-0 left-0 right-0 h-18 bg-[#1f3c93] flex items-center justify-between px-6 z-50 shadow-md">

            <!-- KIRI -->
            <div class="flex items-center gap-5">
                <div class="flex items-center gap-3">
                    <img src="@/assets/images/logo.png" alt="logo" class="w-10 h-10 object-contain" />
                    <h1 class="text-2xl font-bold text-white tracking-wide">
                        SABAR
                    </h1>
                </div>

                <button @click="isSidebarOpen = !isSidebarOpen"
                    class="text-white hover:bg-white/10 p-2 rounded-lg transition">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                        stroke="currentColor" class="size-7">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
                    </svg>
                </button>
            </div>

            <!-- KANAN: nama & avatar dari akun yang login -->
            <div class="flex items-center gap-4 cursor-pointer hover:opacity-80 transition mr-4"
                @click="goTo('/dashboard-admin/profile')">
                <div class="w-10 h-10 rounded-full bg-white overflow-hidden">
                    <img :src="userAvatar" class="w-full h-full object-cover" />
                </div>
                <div class="text-white font-medium">
                    <!-- {{ userName || "..." }} -->

                    <h1>
                        Admin Akademik</h1>
                </div>
            </div>
        </header>

        <!-- ═══════════════════════════════════════
             SIDEBAR
        ════════════════════════════════════════ -->
        <aside :class="[
            'fixed left-0 top-18 bottom-0 bg-[#c8d8ee] flex flex-col justify-between transition-all duration-300 overflow-hidden',
            isSidebarOpen ? 'w-62.5' : 'w-20',
        ]">
            <!-- MENU -->
            <div class="flex-1 overflow-y-auto p-3 space-y-1">
                <!-- DASHBOARD -->
                <div @click="goTo('/dashboard-admin')"
                    :class="isActive('/dashboard-admin') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                    class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 transition">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                        class="icon icon-tabler icons-tabler-outline icon-tabler-layout-dashboard">
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path d="M5 4h4a1 1 0 0 1 1 1v6a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-6a1 1 0 0 1 1 -1" />
                        <path d="M5 16h4a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-2a1 1 0 0 1 1 -1" />
                        <path d="M15 12h4a1 1 0 0 1 1 1v6a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-6a1 1 0 0 1 1 -1" />
                        <path d="M15 4h4a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-2a1 1 0 0 1 1 -1" />
                    </svg>

                    <span v-if="isSidebarOpen">Dashboard</span>
                </div>

                <!-- AKADEMIK -->
                <div>
                    <div @click="openAkademik = !openAkademik"
                        class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 text-gray-700 hover:bg-[#b8c9e2] transition">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 16 16">
                            <path d="M0 0h16v16H0z" fill="none" />
                            <path fill="currentColor"
                                d="M15.09 12.79a1 1 0 0 0-.086-1.638L15 5.33L14 6v5.15a1 1 0 0 0-.092 1.629l-.378.502a2.48 2.48 0 0 0-.53 1.498v1.222h.815a.88.88 0 0 0 .853-.664l.331-1.336v2h1v-1.21a2.5 2.5 0 0 0-.534-1.505zM8 0L0 4l8 5l8-5z" />
                            <path fill="currentColor" d="M8 10L3 6.67v1.71C3 9.29 5.94 12 8 12s5-2.71 5-3.62V6.67z" />
                        </svg>

                        <span v-if="isSidebarOpen" class="flex-1">Akademik</span>
                        <svg v-if="isSidebarOpen" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            stroke-width="1.5" stroke="currentColor"
                            :class="['size-4 transition-transform', openAkademik ? 'rotate-180' : '']">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25 12 15.75 4.5 8.25" />
                        </svg>
                    </div>

                    <div v-if="openAkademik && isSidebarOpen" class="ml-3 mt-1 space-y-1">

                        <!-- Tahun Akademik -->
                        <div @click="goTo('/dashboard-admin/tahun_akademik')"
                            :class="isActive('/dashboard-admin/tahun_akademik') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" class="size-6" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M8 4v4M16 4v4M4 11h16M5 6h14a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1" />
                            </svg>
                            <span>Tahun Akademik</span>
                        </div>

                        <!-- Presensi -->
                        <div @click="goTo('/dashboard-admin/presensi')"
                            :class="isActive('/dashboard-admin/presensi') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                                    stroke-width="2"
                                    d="m7 12l3-3l3 3l4-4M8 21l4-4l4 4M3 4h18M4 4h16v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1z" />
                            </svg>

                            <span>Presensi</span>
                        </div>

                        <!-- Kurikulum -->
                        <div @click="goTo('/dashboard-admin/kurikulum')"
                            :class="isActive('/dashboard-admin/kurikulum') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="currentColor" fill-rule="evenodd"
                                    d="M20.75 16.714a1 1 0 0 1-.014.143a.75.75 0 0 1-.736.893H6a1.25 1.25 0 1 0 0 2.5h14a.75.75 0 0 1 0 1.5H6A2.75 2.75 0 0 1 3.25 19V5A2.75 2.75 0 0 1 6 2.25h13.4c.746 0 1.35.604 1.35 1.35zM9 6.25a.75.75 0 0 0 0 1.5h6a.75.75 0 0 0 0-1.5z"
                                    clip-rule="evenodd" />
                            </svg>

                            <span>Kurikulum</span>
                        </div>

                        <!-- Jurusan -->
                        <div @click="goTo('/dashboard-admin/jurusan')"
                            :class="isActive('/dashboard-admin/jurusan') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 640 640">
                                <path d="M0 0h640v640H0z" fill="none" />
                                <path fill="currentColor"
                                    d="m80 259.8l209.2 86.1c9.8 4 20.2 6.1 30.8 6.1s21-2.1 30.8-6.1l242.4-99.8c9-3.7 14.8-12.4 14.8-22.1s-5.8-18.4-14.8-22.1l-242.4-99.8c-9.8-4-20.2-6.1-30.8-6.1s-21 2.1-30.8 6.1L46.8 201.9c-9 3.7-14.8 12.4-14.8 22.1v296c0 13.3 10.7 24 24 24s24-10.7 24-24zm48 71.7V448c0 53 86 96 192 96s192-43 192-96V331.4l-142.9 58.9c-15.6 6.4-32.2 9.7-49.1 9.7s-33.5-3.3-49.1-9.7L128 331.4z" />
                            </svg>

                            <span>Jurusan</span>
                        </div>

                        <!-- Prodi -->
                        <div @click="goTo('/dashboard-admin/prodi')"
                            :class="isActive('/dashboard-admin/prodi') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 1024 1024">
                                <path d="M0 0h1024v1024H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M1021.06 839.968L833.798 493.664c19.104-43.36 29.792-91.28 29.792-141.696C863.59 157.664 706.07.16 511.782.16c-194.336 0-351.84 157.52-351.84 351.808c0 51.632 11.216 100.624 31.184 144.784L3.03 839.808c-6.065 11.024-5.057 24.624 2.527 34.688c7.6 10.033 20.432 14.752 32.687 11.873l160.624-36.848l54.976 153.12c4.288 11.904 15.152 20.16 27.744 21.088c.817.064 1.6.096 2.368.096a32 32 0 0 0 28.192-16.88L475.844 701.97a355 355 0 0 0 35.92 1.808c11.12 0 22.095-.576 32.943-1.6l167.248 305.008a31.98 31.98 0 0 0 30.56 16.527c12.56-1.008 23.376-9.248 27.631-21.088l54.976-153.12l160.624 36.848c12.32 2.975 25.024-1.809 32.624-11.809c7.632-9.984 8.656-23.52 2.688-34.576zm-731.282 73.376L249.52 801.183c-5.504-15.248-21.471-24.128-37.28-20.368l-118.8 27.248l135.41-246.976c44.592 60.24 107.952 105.68 181.44 127.793zm-65.553-561.377c0-158.544 129.009-287.536 287.568-287.536c158.544 0 287.536 128.992 287.536 287.536S670.337 639.535 511.793 639.535c-158.576 0-287.568-129.024-287.568-287.568m587.52 428.847c-15.872-3.744-31.776 5.12-37.28 20.367l-40.529 112.976l-123.152-224.56c75.44-22.096 140.337-68.735 185.505-130.735L931.137 808.19z" />
                            </svg>

                            <span>Prodi</span>
                        </div>

                        <!-- Mata Kuliah -->
                        <div @click="goTo('/dashboard-admin/matakuliah')"
                            :class="isActive('/dashboard-admin/matakuliah') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 48 48">
                                <path d="M0 0h48v48H0z" fill="none" />
                                <g fill="none" stroke="currentColor" stroke-linejoin="round" stroke-width="4">
                                    <path d="M10 6a2 2 0 0 1 2-2h28a2 2 0 0 1 2 2v36a2 2 0 0 1-2 2H12a2 2 0 0 1-2-2z" />
                                    <path stroke-linecap="round" d="M34 6v36M6 14h8M6 24h8M6 34h8M27 4h12M27 44h12" />
                                </g>
                            </svg>

                            <span>Mata Kuliah</span>
                        </div>

                        <!-- Kelas -->
                        <div @click="goTo('/dashboard-admin/kelas')"
                            :class="isActive('/dashboard-admin/kelas') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" class="size-6" viewBox="0 0 24 24">
                                <path fill="none" stroke="currentColor" stroke-width="2" d="M3 5h18v14H3z" />
                                <path fill="none" stroke="currentColor" stroke-width="2" d="M7 9h10M7 13h6" />
                            </svg>
                            <span>Kelas</span>
                        </div>

                        <!-- Nilai -->
                        <!-- <div @click="goTo('/dashboard-admin/nilai')"
                            <!-- :class="isActive('/dashboard-admin/nilai') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" class="size-5" viewBox="0 0 24 24">
                                <path d="M3 3v18h18" fill="none" stroke="currentColor" stroke-width="2" />
                                <path d="M7 14l3-3l3 2l4-5" fill="none" stroke="currentColor" stroke-width="2" />
                            </svg>
                            <span>Nilai</span>
                        </div> -->

                        <!-- KHS -->
                        <div @click="goTo('/dashboard-admin/khs')"
                            :class="isActive('/dashboard-admin/khs') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M19.903 8.586a1 1 0 0 0-.196-.293l-6-6a1 1 0 0 0-.293-.196c-.03-.014-.062-.022-.094-.033a1 1 0 0 0-.259-.051C13.04 2.011 13.021 2 13 2H6c-1.103 0-2 .897-2 2v16c0 1.103.897 2 2 2h12c1.103 0 2-.897 2-2V9c0-.021-.011-.04-.013-.062a1 1 0 0 0-.051-.259q-.014-.048-.033-.093M16.586 8H14V5.414zM6 20V4h6v5a1 1 0 0 0 1 1h5l.002 10z" />
                                <path fill="currentColor" d="M8 12h8v2H8zm0 4h8v2H8zm0-8h2v2H8z" />
                            </svg>

                            <span>KHS</span>
                        </div>

                        <!-- Dosen -->
                        <div @click="goTo('/dashboard-admin/dosen')"
                            :class="isActive('/dashboard-admin/dosen') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M8 4a2 2 0 1 1-4 0a2 2 0 0 1 4 0M5 16v6H3V10a3 3 0 0 1 5.106-2.137l2.374 2.243l2.313-2.313l1.414 1.414l-3.687 3.687L9 11.46V22H7v-6zm5-11h9v9h-9v2h4.365l2.824 6h2.21l-2.823-6H20a1 1 0 0 0 1-1V4a1 1 0 0 0-1-1H10z" />
                            </svg>

                            <span>Dosen</span>
                        </div>

                    </div>

                    <!-- MAHASISWA -->
                    <div>
                        <div @click="openMahasiswa = !openMahasiswa"
                            class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-3 text-gray-700 hover:bg-[#b8c9e2] transition">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="27" viewBox="0 0 640 640">
                                <path d="M0 0h640v640H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M320 80c57.4 0 104 46.6 104 104s-46.6 104-104 104s-104-46.6-104-104S262.6 80 320 80M96 152c39.8 0 72 32.2 72 72s-32.2 72-72 72s-72-32.2-72-72s32.2-72 72-72M0 480c0-70.7 57.3-128 128-128c12.8 0 25.2 1.9 36.9 5.4C132 394.2 112 442.8 112 496v16c0 11.4 2.4 22.2 6.7 32H32c-17.7 0-32-14.3-32-32zm521.3 64c4.3-9.8 6.7-20.6 6.7-32v-16c0-53.2-20-101.8-52.9-138.6c11.7-3.5 24.1-5.4 36.9-5.4c70.7 0 128 57.3 128 128v32c0 17.7-14.3 32-32 32zM472 224c0-39.8 32.2-72 72-72s72 32.2 72 72s-32.2 72-72 72s-72-32.2-72-72M160 496c0-88.4 71.6-160 160-160s160 71.6 160 160v16c0 17.7-14.3 32-32 32H192c-17.7 0-32-14.3-32-32z" />
                            </svg>

                            <span v-if="isSidebarOpen" class="flex-1">Mahasiswa</span>
                            <svg v-if="isSidebarOpen" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                stroke-width="1.5" stroke="currentColor"
                                :class="['size-4 transition-transform', openMahasiswa ? 'rotate-180' : '']">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25 12 15.75 4.5 8.25" />
                            </svg>
                        </div>

                        <div v-if="openMahasiswa && isSidebarOpen" class="ml-3 mt-1 space-y-1">
                            <div @click="goTo('/dashboard-admin/peserta_kelas')"
                                :class="isActive('/dashboard-admin/peserta_kelas') ? 'bg-[#1f3c93] text-white' : 'text-gray-700 hover:bg-[#b8c9e2]'"
                                class="flex cursor-pointer items-center gap-3 rounded-lg px-4 py-2 transition">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="40" viewBox="0 0 16 16">
                                    <path d="M0 0h16v16H0z" fill="none" />
                                    <path fill="currentColor"
                                        d="M8.5 4.5a2.5 2.5 0 1 1-5 0a2.5 2.5 0 0 1 5 0m2.4 7.506c.11.542-.348.994-.9.994H2c-.553 0-1.01-.452-.902-.994a5.002 5.002 0 0 1 9.803 0M14.002 12h-1.59a3 3 0 0 0-.04-.29a6.5 6.5 0 0 0-1.167-2.603a3 3 0 0 1 3.633 1.911c.18.522-.283.982-.836.982M12 8a2 2 0 1 0 0-4a2 2 0 0 0 0 4" />
                                </svg>

                                <span>Peserta Kelas</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- LOGOUT -->
            <div @click="showLogoutPopup = true"
                class="border-t border-[#b0c4de] flex cursor-pointer items-center gap-3 p-5 text-gray-700 hover:bg-[#9fb5d6] hover:text-red-500 transition">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="size-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6A2.25 2.25 0 0 0 5.25 5.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
                </svg>
                <span v-if="isSidebarOpen">Keluar</span>
            </div>

        </aside>

        <!-- ═══════════════════════════════════════
             CONTENT
        ════════════════════════════════════════ -->
        <main :class="[
            'pt-22.5 p-6 transition-all duration-300',
            isSidebarOpen ? 'ml-62.5' : 'ml-20',
        ]">

            <!-- DASHBOARD HOME -->
            <div v-if="route.path === '/dashboard-admin'">
                <h1 class="text-3xl font-bold mb-1">Dashboard</h1>
                <p class="text-gray-600 mb-6">
                    Selamat Datang, {{ userName || "Admin Akademik" }}
                </p>

                <!-- STAT CARDS -->
                <div class="grid grid-cols-3 gap-4 mb-6">

                    <!-- Total Pegawai -->
                    <div
                        class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                        <div class="bg-[#a9c3e3] p-4 rounded-lg">
                            <svg xmlns="http://www.w3.org/2000/svg" width="35" height="35" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M12 12.25a3.75 3.75 0 1 1 3.75-3.75A3.75 3.75 0 0 1 12 12.25m0-6a2.25 2.25 0 1 0 2.25 2.25A2.25 2.25 0 0 0 12 6.25m7 13a.76.76 0 0 1-.75-.75c0-1.95-1.06-3.25-6.25-3.25s-6.25 1.3-6.25 3.25a.75.75 0 0 1-1.5 0c0-4.75 5.43-4.75 7.75-4.75s7.75 0 7.75 4.75a.76.76 0 0 1-.75.75" />
                            </svg>

                        </div>
                        <div>
                            <p class="text-sm font-semibold text-gray-500">TOTAL PEGAWAI</p>
                            <h2 class="text-3xl font-bold">{{ totalPegawai }}</h2>
                        </div>
                    </div>

                    <!-- Total Mahasiswa -->
                    <div
                        class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                        <div class="bg-[#a9c3e3] p-4 rounded-lg">
                            <i class="fi fi-rr-student text-3xl text-[#0f0f0f]"></i>
                        </div>
                        <div>
                            <p class="text-sm font-semibold text-gray-500">TOTAL MAHASISWA</p>
                            <h2 class="text-3xl font-bold">{{ totalMahasiswa }}</h2>
                        </div>
                    </div>

                    <!-- Total Dosen -->
                    <div
                        class="bg-[#ececec] rounded-xl p-4 flex items-center gap-4 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                        <div class="bg-[#a9c3e3] p-4 rounded-lg">
                            <svg xmlns="http://www.w3.org/2000/svg" width="35" height="35" viewBox="0 0 24 24">
                                <path d="M0 0h24v24H0z" fill="none" />
                                <path fill="currentColor"
                                    d="M8 4a2 2 0 1 1-4 0a2 2 0 0 1 4 0M5 16v6H3V10a3 3 0 0 1 5.106-2.137l2.374 2.243l2.313-2.313l1.414 1.414l-3.687 3.687L9 11.46V22H7v-6zm5-11h9v9h-9v2h4.365l2.824 6h2.21l-2.823-6H20a1 1 0 0 0 1-1V4a1 1 0 0 0-1-1H10z" />
                            </svg>
                        </div>
                        <div>
                            <p class="text-sm font-semibold text-gray-500">TOTAL DOSEN</p>
                            <h2 class="text-3xl font-bold">{{ totalDosen }}</h2>
                        </div>
                    </div>
                </div>

                <!-- CONTENT GRID -->
                <div class="grid grid-cols-4 gap-4">

                    <!-- TABLE MAHASISWA -->
                    <div
                        class="col-span-3 bg-[#ececec] rounded-xl p-5 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                        <h2 class="text-xl font-bold mb-4">Data Mahasiswa</h2>

                        <table class="w-full text-sm border-collapse">
                            <thead>
                                <tr class="bg-gray-50">
                                    <th class="py-3 px-2 text-left">No</th>
                                    <th class="py-3 px-2 text-left">NIM</th>
                                    <th class="py-3 px-2 text-left">Nama</th>
                                    <th class="py-3 px-2 text-left">Kelas</th>
                                    <th class="py-3 px-2 text-left">Prodi</th>
                                    <th class="py-3 px-2 text-left">Angkatan</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-if="akunList.length === 0">
                                    <td colspan="6" class="py-6 text-center text-gray-400">
                                        Tidak ada data
                                    </td>
                                </tr>
                                <tr v-for="(item, index) in akunList" :key="item.id"
                                    class="hover:bg-gray-50 border-b border-gray-100">
                                    <td class="py-3 px-2">
                                        {{ (currentPage - 1) * perPage + index + 1 }}
                                    </td>
                                    <td class="py-3 px-2">{{ item.nim }}</td>
                                    <td class="py-3 px-2">{{ item.name }}</td>
                                    <td class="py-3 px-2">{{ item.kelas }}</td>
                                    <td class="py-3 px-2">{{ item.prodi }}</td>
                                    <td class="py-3 px-2">{{ item.angkatan }}</td>
                                </tr>
                            </tbody>
                        </table>

                        <!-- PAGINATION -->
                        <div class="flex justify-end mt-5 pt-4">
                            <div class="flex items-center gap-2">
                                <button @click="prevPage" :disabled="currentPage === 1"
                                    class="px-3 py-1 border rounded-lg bg-white hover:bg-gray-100 disabled:opacity-50 text-sm">
                                    Previous
                                </button>

                                <template v-for="p in pages" :key="p">
                                    <span v-if="p === '...'" class="px-1 text-gray-400">...</span>
                                    <button v-else @click="goToPage(p as number)" class="w-8 h-8 rounded-lg text-sm"
                                        :class="currentPage === p ? 'bg-blue-500 text-white' : 'bg-gray-100 hover:bg-gray-200'">
                                        {{ p }}
                                    </button>
                                </template>

                                <button @click="nextPage" :disabled="currentPage === totalPages"
                                    class="px-3 py-1 border rounded-lg bg-white hover:bg-gray-100 disabled:opacity-50 text-sm">
                                    Next
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- SIDE PANEL -->
                    <div class="space-y-4">

                        <!-- Tahun Akademik -->
                        <div
                            class="bg-[#ececec] rounded-xl p-5 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                            <div class="flex items-center justify-between mb-3">
                                <h2 class="font-bold">Tahun Akademik</h2>
                                <button @click="goTo('/dashboard-admin/dashboard_tahunakademik')"
                                    class="text-blue-600 hover:text-blue-800">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                        stroke-width="1.5" stroke="currentColor" class="size-5">
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M2.25 12s3.75-7.5 9.75-7.5 9.75 7.5 9.75 7.5-3.75 7.5-9.75 7.5S2.25 12 2.25 12Z" />
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M12 15.75A3.75 3.75 0 1 0 12 8.25a3.75 3.75 0 0 0 0 7.5Z" />
                                    </svg>
                                </button>
                            </div>

                            <div v-for="item in tahunAkademik.slice(0, 4)" :key="item.id"
                                class="bg-gray-50 rounded-lg p-3 mb-2 text-sm">
                                {{ formatYear(item.tahun_awal) }} /
                                {{ formatYear(item.tahun_akhir) }}
                                {{ item.tipee_semester === 'ganjil' ? 'Ganjil' : 'Genap' }}
                            </div>

                            <div v-if="tahunAkademik.length === 0" class="text-gray-400 text-sm">
                                Tidak ada data
                            </div>
                        </div>

                        <!-- Kurikulum -->
                        <div
                            class="bg-[#ececec] rounded-xl p-5 shadow-sm border-l-[4px] border-b-[3px] border-[#9db9dc]">
                            <div class="flex items-center justify-between mb-3">
                                <h2 class="font-bold">Kurikulum</h2>
                                <button @click="goTo('/dashboard-admin/dashboard_kurikulum')"
                                    class="text-blue-600 hover:text-blue-800">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                        stroke-width="1.5" stroke="currentColor" class="size-5">
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M2.25 12s3.75-7.5 9.75-7.5 9.75 7.5 9.75 7.5-3.75 7.5-9.75 7.5S2.25 12 2.25 12Z" />
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M12 15.75A3.75 3.75 0 1 0 12 8.25a3.75 3.75 0 0 0 0 7.5Z" />
                                    </svg>
                                </button>
                            </div>

                            <div v-for="item in kurikulum.slice(0, 4)" :key="item.id"
                                class="bg-gray-50 rounded-lg p-3 mb-2 text-sm">
                                {{ item.name }}
                            </div>

                            <div v-if="kurikulum.length === 0" class="text-gray-400 text-sm">
                                Tidak ada data
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- ROUTER VIEW (halaman lain) -->
            <div v-else>
                <router-view />
            </div>
        </main>

        <!-- POPUP LOGOUT -->
        <KonfirmasiKeluar v-if="showLogoutPopup" @close="showLogoutPopup = false" @confirm="handleLogout" />
    </div>
</template>