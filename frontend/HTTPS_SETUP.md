# Setup HTTPS untuk Development

## Masalah: "Sambungan untuk situs ini tidak aman"

Ini terjadi karena Vite menggunakan self-signed certificate yang tidak terpercaya oleh browser.

## Solusi 1: Terima Certificate Warning (Cepat)

1. Buka `https://192.168.1.6:3000` di browser
2. Akan muncul warning "Your connection is not private"
3. Klik "Advanced" atau "Lanjutkan ke situs (tidak aman)"
4. Klik "Proceed to 192.168.1.6 (unsafe)" atau "Lanjutkan"
5. Situs akan terbuka dengan HTTPS

**Catatan:** Warning ini normal untuk self-signed certificate di development.

## Solusi 2: Setup Trusted Certificate dengan mkcert (Recommended)

### Windows:

1. **Install mkcert:**
   ```powershell
   # Via Chocolatey
   choco install mkcert
   
   # Atau download dari: https://github.com/FiloSottile/mkcert/releases
   ```

2. **Install CA (Certificate Authority):**
   ```powershell
   mkcert -install
   ```

3. **Generate Certificate:**
   ```powershell
   cd frontend
   mkcert localhost 192.168.1.6 ::1
   ```
   Ini akan membuat file:
   - `localhost+2.pem` (certificate)
   - `localhost+2-key.pem` (private key)

4. **Update vite.config.ts:**
   ```typescript
   server: {
     https: {
       key: './localhost+2-key.pem',
       cert: './localhost+2.pem',
     },
     host: '0.0.0.0',
     port: 3000,
     // ... rest of config
   }
   ```

5. **Restart dev server:**
   ```bash
   npm run dev
   ```

6. **Test:** Buka `https://192.168.1.6:3000` - tidak ada warning!

## Solusi 3: Tetap Pakai HTTP untuk Development

Jika HTTPS terlalu kompleks untuk development, bisa tetap pakai HTTP:

1. **Comment HTTPS di vite.config.ts:**
   ```typescript
   server: {
     // https: true, // Disable untuk development
     host: '0.0.0.0',
     port: 3000,
   }
   ```

2. **PWA akan tetap bekerja di localhost** (Chrome mengizinkan PWA di localhost tanpa HTTPS)

3. **Untuk production, gunakan HTTPS** (deploy ke server dengan SSL certificate)

## Rekomendasi

**Untuk Development:**
- Gunakan **Solusi 2 (mkcert)** untuk trusted certificate tanpa warning
- Atau **Solusi 3 (HTTP)** jika tidak perlu HTTPS

**Untuk Production:**
- Gunakan HTTPS dengan certificate dari Let's Encrypt atau provider SSL lainnya
- Deploy ke server dengan SSL (Vercel, Netlify, atau server sendiri dengan nginx)

## Troubleshooting

**Certificate tidak bekerja?**
- Pastikan mkcert sudah di-install: `mkcert -version`
- Pastikan CA sudah di-install: `mkcert -install`
- Pastikan file certificate ada di folder `frontend/`

**Masih ada warning?**
- Clear browser cache
- Restart browser
- Pastikan menggunakan IP yang sama dengan certificate (192.168.1.6)

**Port 3000 tidak bisa diakses?**
- Pastikan firewall mengizinkan port 3000
- Cek apakah port sudah digunakan: `netstat -ano | findstr :3000`

