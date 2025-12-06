# PWA Setup - Install UMess di Android

Aplikasi UMess sudah dikonfigurasi sebagai Progressive Web App (PWA) dan bisa diinstall di Android.

## Langkah-langkah Setup

### 1. Generate Icon (PENTING!)

Sebelum aplikasi bisa diinstall, Anda perlu membuat icon:

**Option A: Menggunakan generate-icons.html (Mudah)**
1. Buka file `frontend/public/generate-icons.html` di browser
2. Klik tombol "Download 192x192" dan "Download 512x512"
3. Simpan file dengan nama `icon-192.png` dan `icon-512.png`
4. Pindahkan ke folder `frontend/public/`

**Option B: Online Tools**
- Kunjungi https://www.pwabuilder.com/imageGenerator
- Upload logo/icon Anda
- Download icon 192x192 dan 512x512
- Simpan di `frontend/public/` dengan nama `icon-192.png` dan `icon-512.png`

**Option C: Manual**
- Buat logo dengan ukuran minimal 512x512px
- Export sebagai PNG
- Resize ke 192x192 untuk icon-192.png
- Copy ke 512x512 untuk icon-512.png
- Simpan di `frontend/public/`

### 2. Build Aplikasi

```bash
cd frontend
npm run build
```

### 3. Deploy ke Server

Deploy folder `frontend/dist` ke server dengan HTTPS (PWA memerlukan HTTPS).

**Development (Local):**
- PWA bisa di-test di localhost dengan Chrome
- Untuk test di Android, gunakan IP local network (192.168.x.x) dengan HTTPS atau tunnel

**Production:**
- Deploy ke server dengan HTTPS (misalnya Vercel, Netlify, atau server sendiri)
- Pastikan server mengirim header yang benar untuk PWA

### 4. Install di Android

1. Buka aplikasi di **Chrome Android** (bukan browser lain)
2. Tap menu (3 titik) di kanan atas
3. Pilih **"Add to Home screen"** atau **"Install app"**
4. Konfirmasi install
5. Aplikasi akan muncul di home screen seperti aplikasi native!

## Fitur PWA yang Tersedia

✅ **Installable** - Bisa diinstall seperti aplikasi native
✅ **Offline Support** - Service Worker untuk caching
✅ **Fast Loading** - Assets di-cache untuk loading cepat
✅ **App-like Experience** - Fullscreen, no browser bar
✅ **Auto Update** - Aplikasi otomatis update saat ada versi baru

## Troubleshooting

**Icon tidak muncul?**
- Pastikan file `icon-192.png` dan `icon-512.png` ada di `frontend/public/`
- Rebuild aplikasi: `npm run build`

**Tidak bisa install?**
- Pastikan menggunakan HTTPS (atau localhost untuk development)
- Pastikan menggunakan Chrome browser
- Cek console browser untuk error

**Service Worker tidak bekerja?**
- Pastikan server mengirim header yang benar
- Cek di Chrome DevTools > Application > Service Workers

## Testing PWA

1. Buka Chrome DevTools (F12)
2. Tab "Application" > "Manifest" - cek manifest sudah benar
3. Tab "Application" > "Service Workers" - cek service worker terdaftar
4. Tab "Lighthouse" > Run audit untuk PWA - harus score 90+

## Catatan

- PWA memerlukan HTTPS untuk production (kecuali localhost)
- Icon harus berupa PNG, bukan SVG
- Service Worker akan otomatis di-generate oleh vite-plugin-pwa saat build

