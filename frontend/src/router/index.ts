import { createRouter, createWebHistory } from "vue-router";

import Login from "@/views/Login.vue";

// super admin
import DashboardSuperAdmin from "@/views/dashboard_super_admin/Dashboardsuperadmin.vue";

// admin akademik
import DashboardAdmin from "@/views/dashboard_admin_akademik/Dashboard_admin_akademik.vue";

// akademik
import KHS from "@/views/dashboard_admin_akademik/akademik/KHS/KHS.vue";
import Nilai from "@/views/dashboard_admin_akademik/akademik/nilai/Nilai.vue";
import Presensi from "@/views/dashboard_admin_akademik/akademik/Presensi.vue";
import Tahun_akademik from "@/views/dashboard_admin_akademik/akademik/tahun-akademik/Tahun_akademik.vue";
import Tambah_tahunakademik from "../views/dashboard_admin_akademik/akademik/tahun-akademik/Tambah_tahunakademik.vue";

// kurikulum
import Kurikulum from "@/views/dashboard_admin_akademik/akademik/kurikulum/Kurikulum.vue";
import Tambah_kurikulum from "@/views/dashboard_admin_akademik/akademik/kurikulum/Tambah_kurikulum.vue";

// mahasiswa
import Tambah_kelas from "@/views/dashboard_admin_akademik/akademik/kelas/Tambah_kelas.vue";

// kelas
import Kelas from "@/views/dashboard_admin_akademik/akademik/kelas/Kelas.vue";

// peserta kelas
import Peserta_kelas from "@/views/dashboard_admin_akademik/akademik/Peserta-kelas/Peserta_kelas.vue";
import Detail_pesertakelas from "@/views/dashboard_admin_akademik/akademik/Peserta-kelas/Detail_pesertakelas.vue";

// pegawai
import Dosen from "@/views/dashboard_admin_akademik/akademik/pegawai/Dosen.vue";
import Dashboard_kurikulum from "../views/dashboard_admin_akademik/Dashboard_kurikulum.vue";
import Dashboard_tahunakademik from "../views/dashboard_admin_akademik/Dashboard_tahunakademik.vue";

// jurusan
import Jurusan from "@/views/dashboard_admin_akademik/akademik/Jurusan/Jurusan.vue";

// prodi
import Prodi from "@/views/dashboard_admin_akademik/akademik/Prodi/Prodi.vue";

import profile from "../views/dashboard_admin_akademik/profile.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: Login },

    {
      path: "/dashboard-superadmin",
      component: DashboardSuperAdmin,
      children: [
        // akun

        {
          path: "reset_password",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/password/Reset_password.vue"),
        },
        {
          path: "akun",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/akun/Akun.vue"),
        },
        {
          path: "edit_akun/:id",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/akun/Edit_akun.vue"),
        },
        {
          path: "tambah_akun",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/akun/Tambah_akun.vue"),
        },
        {
          path: "role",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/role/Role.vue"),
        },
        {
          path: "tambah_role",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/role/Tambah_role.vue"),
        },
        {
          path: "edit_role/:id",
          component: () =>
            import("@/views/dashboard_super_admin/akademik/role/Edit_role.vue"),
        },
      ],
    },

    {
      path: "/dashboard-admin",
      component: DashboardAdmin,
      children: [
        { path: "dashboard_kurikulum", component: Dashboard_kurikulum },
        { path: "dashboard_tahunakademik", component: Dashboard_tahunakademik },
        // akademik
        { path: "khs", component: KHS },
        { path: "nilai", component: Nilai },
        { path: "presensi", component: Presensi },
        { path: "tahun_akademik", component: Tahun_akademik },
        { path: "tambah_tahunakademik", component: Tambah_tahunakademik },

        // kurikulum
        { path: "kurikulum", component: Kurikulum },
        {
          path: "detail_kurikulum/:id",
          component: () =>
            import("@/views/dashboard_admin_akademik/akademik/kurikulum/Detail_kurikulum.vue"),
        },

        // mahasiswa
        { path: "tambah_kelas", component: Tambah_kelas },

        // kelas
        { path: "kelas", component: Kelas },

        // peserta kelas
        { path: "peserta_kelas", component: Peserta_kelas },
        { path: "detail_pesertakelas", component: Detail_pesertakelas },

        // pegawai
        { path: "dosen", component: Dosen },

        //kurikulum
        { path: "tambah_kurikulum", component: Tambah_kurikulum },
        { path: "profile", component: profile },

        // prodi
        { path: "prodi", component: Prodi },

        // jurusan
        { path: "jurusan", component: Jurusan },
      ],
    },
  ],
});

export default router;
