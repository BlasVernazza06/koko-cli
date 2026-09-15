# 🦾 Koko CLI - Guía de Comandos para Landing y Distribución

Comandos oficiales para inicializar e instalar **Koko** a través de diferentes ecosistemas de desarrollo:

---

## ⚡ 1. Ejecución Instantánea (Sin instalación previa)
La forma más rápida para que un desarrollador inicialice un proyecto sin tener que instalar nada globalmente en su sistema:

### Node.js (NPX)
```bash
npx koko-app init
```

### PNPM
```bash
pnpm dlx koko-app init
```

### Bun
```bash
bunx koko-app init
```

---

## 🚀 2. Instalación Nativa en Go
Para desarrolladores en el ecosistema de Go que prefieren instalar el binario nativo directamente desde el código fuente:

```bash
go install github.com/BlasVernazza06/koko-cli@latest
```

---

## 🌍 3. Instalación Global (Para usar siempre el comando `koko`)
Permite registrar el comando de terminal `koko` globalmente en el sistema:

### Con NPM
```bash
npm install -g koko-app
```

### Con PNPM
```bash
pnpm add -g koko-app
```

### Con Bun
```bash
bun add -g koko-app
```

*Una vez instalado globalmente, se ejecuta en cualquier carpeta con:*
```bash
koko init
```

---

## 📦 4. Descarga Directa de Binarios Precompilados
Para usuarios que no utilizan Node ni Go y prefieren el ejecutable directo:

* **Descarga desde GitHub Releases:**  
  [https://github.com/BlasVernazza06/koko-cli/releases](https://github.com/BlasVernazza06/koko-cli/releases)
