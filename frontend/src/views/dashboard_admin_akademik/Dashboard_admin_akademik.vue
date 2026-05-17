# `Dashboard_admin_akademik.vue`

```vue
<template>
  <div class="dashboard-layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="logo-section">
        <img
          src="https://cdn-icons-png.flaticon.com/512/3135/3135715.png"
          alt="logo"
          class="logo"
        />
        <h1>SABAR</h1>
      </div>

      <nav class="menu">
        <ul>
          <li class="active">
            <i class="fas fa-th-large"></i>
            Dashboard
          </li>

          <li>
            <i class="fas fa-graduation-cap"></i>
            Akademik
            <i class="fas fa-chevron-down right"></i>
          </li>

          <li>
            <i class="fas fa-calendar"></i>
            Tahun Akademik
          </li>

          <li>
            <i class="fas fa-clipboard-check"></i>
            Presensi
          </li>

          <li>
            <i class="fas fa-book"></i>
            Kurikulum
          </li>

          <li>
            <i class="fas fa-pen-square"></i>
            Nilai
          </li>

          <li>
            <i class="fas fa-file-alt"></i>
            KHS
          </li>

          <li>
            <i class="fas fa-users"></i>
            Pegawai
            <i class="fas fa-chevron-down right"></i>
          </li>

          <li class="submenu">Dosen</li>

          <li>
            <i class="fas fa-user-graduate"></i>
            Mahasiswa
            <i class="fas fa-chevron-down right"></i>
          </li>

          <li class="submenu">Kelas</li>
          <li class="submenu">Peserta Kelas</li>
        </ul>
      </nav>

      <div class="logout">
        <i class="fas fa-sign-out-alt"></i>
        Keluar
      </div>
    </aside>

    <!-- Main -->
    <main class="main-content">
      <!-- Topbar -->
      <header class="topbar">
        <div class="menu-toggle">
          <i class="fas fa-bars"></i>
        </div>

        <div class="profile">
          <img
            src="https://i.pravatar.cc/100"
            alt="admin"
            class="profile-img"
          />
          <span>Admin Akademik</span>
        </div>
      </header>

      <!-- Content -->
      <section class="content">
        <div class="page-title">
          <h2>Dashboard</h2>
          <p>Selamat Datang Super Admin, Hafiz</p>
        </div>

        <!-- Cards -->
        <div class="cards">
          <div class="card-box">
            <div class="card-icon">
              <i class="fas fa-user"></i>
            </div>
            <div>
              <h4>TOTAL PEGAWAI</h4>
              <h2>120</h2>
            </div>
          </div>

          <div class="card-box">
            <div class="card-icon">
              <i class="fas fa-user-graduate"></i>
            </div>
            <div>
              <h4>TOTAL MAHASISWA</h4>
              <h2>120</h2>
            </div>
          </div>

          <div class="card-box">
            <div class="card-icon">
              <i class="fas fa-chalkboard-teacher"></i>
            </div>
            <div>
              <h4>TOTAL DOSEN</h4>
              <h2>120</h2>
            </div>
          </div>
        </div>

        <!-- Bottom Layout -->
        <div class="bottom-layout">
          <!-- Table -->
          <div class="table-section">
            <div class="table-header">
              <select>
                <option>Pilih Tabel Mana yang ingin di lihat</option>
              </select>
            </div>

            <h3>Data Akun</h3>

            <table>
              <thead>
                <tr>
                  <th>No</th>
                  <th>NIM</th>
                  <th>Nama</th>
                  <th>Kelas</th>
                  <th>Prodi</th>
                  <th>Angkatan</th>
                </tr>
              </thead>

              <tbody>
                <tr>
                  <td>1</td>
                  <td>C030324077</td>
                  <td>Budi Siregar</td>
                  <td>TI-4A</td>
                  <td>Teknik Informatika</td>
                  <td>2024</td>
                </tr>

                <tr>
                  <td>2</td>
                  <td>C030324812</td>
                  <td>Budi Lasegar</td>
                  <td>TI-4A</td>
                  <td>Teknik Informatika</td>
                  <td>2024</td>
                </tr>
              </tbody>
            </table>

            <div class="pagination">
              <select>
                <option>10 Baris</option>
              </select>

              <div class="pages">
                <span>← Previous</span>
                <button class="active-page">1</button>
                <button>2</button>
                <button>3</button>
                <span>...</span>
                <button>67</button>
                <button>68</button>
                <span>Next →</span>
              </div>
            </div>
          </div>

          <!-- Right Panel -->
          <div class="right-panel">
            <div class="info-card">
              <h4>Tahun Akademik</h4>

              <ul>
                <li>
                  Tahun 2023-2024
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  Tahun 2024-2025
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  Tahun 2025-2026
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  Tahun 2026-2027
                  <i class="fas fa-eye"></i>
                </li>
              </ul>
            </div>

            <div class="info-card">
              <h4>Kurikulum</h4>

              <ul>
                <li>
                  Merdeka
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  Teknik Informatika
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  SIKC
                  <i class="fas fa-eye"></i>
                </li>
                <li>
                  Elektronika
                  <i class="fas fa-eye"></i>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
@import url('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css');
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Poppins', sans-serif;
}

.dashboard-layout {
  display: flex;
  min-height: 100vh;
  background: #e9eef5;
}

/* SIDEBAR */
.sidebar {
  width: 250px;
  background: #c8d8ee;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.logo-section {
  background: #1f3c93;
  color: white;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 20px;
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.logo-section h1 {
  font-size: 28px;
  font-weight: 700;
}

.menu {
  padding: 12px;
}

.menu ul {
  list-style: none;
}

.menu li {
  padding: 12px 14px;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #4b5563;
  font-size: 14px;
}

.menu li.active {
  background: #1f3c93;
  color: white;
}

.menu li:hover {
  background: #b8c9e2;
}

.submenu {
  padding-left: 45px !important;
}

.right {
  margin-left: auto;
  font-size: 12px;
}

.logout {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #4b5563;
  cursor: pointer;
}

/* MAIN */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.topbar {
  height: 72px;
  background: #1f3c93;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  color: white;
}

.menu-toggle {
  font-size: 26px;
  cursor: pointer;
}

.profile {
  display: flex;
  align-items: center;
  gap: 12px;
}

.profile-img {
  width: 42px;
  height: 42px;
  border-radius: 50%;
}

.content {
  padding: 28px;
}

.page-title h2 {
  font-size: 36px;
  font-weight: 700;
  color: #111827;
}

.page-title p {
  margin-top: 6px;
  color: #374151;
}

/* CARDS */
.cards {
  display: flex;
  gap: 20px;
  margin-top: 24px;
}

.card-box {
  flex: 1;
  background: white;
  border-radius: 12px;
  padding: 18px;
  display: flex;
  align-items: center;
  gap: 18px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.card-icon {
  width: 90px;
  height: 90px;
  border-radius: 12px;
  background: #a8c4e9;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 38px;
  color: #3d3d3d;
}

.card-box h4 {
  font-size: 15px;
  color: #4b5563;
  margin-bottom: 8px;
}

.card-box h2 {
  font-size: 34px;
}

/* BOTTOM */
.bottom-layout {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 20px;
  margin-top: 24px;
}

.table-section {
  background: white;
  border-radius: 14px;
  padding: 18px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.table-header {
  margin-bottom: 16px;
}

.table-header select,
.pagination select {
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #ccc;
}

.table-section h3 {
  margin-bottom: 18px;
  color: #374151;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 14px 10px;
  text-align: left;
  font-size: 14px;
}

th {
  color: #4b5563;
}

tbody tr {
  border-top: 1px solid #f0f0f0;
}

.pagination {
  margin-top: 60px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pages {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pages button {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
}

.active-page {
  background: #1f3c93 !important;
  color: white;
}

/* RIGHT PANEL */
.right-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card {
  background: white;
  border-radius: 14px;
  padding: 18px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.info-card h4 {
  margin-bottom: 18px;
  color: #374151;
}

.info-card ul {
  list-style: none;
}

.info-card li {
  padding: 10px 0;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #4b5563;
  font-size: 14px;
}

/* RESPONSIVE */
@media (max-width: 1200px) {
  .bottom-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .sidebar {
    display: none;
  }

  .cards {
    flex-direction: column;
  }
}
</style>
```

# Struktur Folder

Simpan file di:

```bash
src/views/pegawai/Dashboard_admin_akademik.vue
```

# Install Font Awesome

```bash
npm install @fortawesome/fontawesome-free
```

Tambahkan di `main.js`

```js
import '@fortawesome/fontawesome-free/css/all.min.css'
```
