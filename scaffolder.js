import { exec as execCb } from 'node:child_process';
import { promisify } from 'node:util';
import fs from 'node:fs/promises';
import path from 'node:path';

const exec = promisify(execCb);

/**
 * Orquesta la creación de la estructura del proyecto y la instalación de tecnologías.
 */
export async function scaffoldProject({ projectName, frontend, backend, tools }) {
  const isMonorepo = tools.includes('turborepo');
  const hasTailwind = tools.includes('tailwind');
  const hasDocker = tools.includes('docker');

  const rootDir = path.join(process.cwd(), projectName);

  if (isMonorepo) {
    // --- FLUJO MONORREPO (TURBOREPO) ---
    
    // 1. Inicializar Turborepo en la raíz
    await exec(`npx create-turbo@latest "${projectName}" --package-manager=npm --skip-install`);

    // 2. Limpiar las aplicaciones de ejemplo por defecto (web y docs)
    const appsDir = path.join(rootDir, 'apps');
    await fs.rm(path.join(appsDir, 'web'), { recursive: true, force: true }).catch(() => {});
    await fs.rm(path.join(appsDir, 'docs'), { recursive: true, force: true }).catch(() => {});

    // 3. Instalar Frontend dentro de la carpeta 'apps'
    if (frontend !== 'none') {
      await installFrontend(frontend, appsDir, 'frontend', hasTailwind);
    }

    // 4. Instalar Backend dentro de la carpeta 'apps'
    if (backend !== 'none') {
      await installBackend(backend, appsDir, 'backend');
    }

  } else {
    // --- FLUJO ESTÁNDAR (SIN MONORREPO) ---
    
    // 1. Crear carpeta raíz
    await fs.mkdir(rootDir, { recursive: true });

    // 2. Instalar Frontend en la raíz
    if (frontend !== 'none') {
      await installFrontend(frontend, rootDir, 'frontend', hasTailwind);
    }

    // 3. Instalar Backend en la raíz
    if (backend !== 'none') {
      await installBackend(backend, rootDir, 'backend');
    }
  }

  // --- INSTALACIÓN DE HERRAMIENTAS Y LIBRERÍAS ADICIONALES ---
  await installTools(tools, rootDir, frontend, backend, isMonorepo);

  // --- CONFIGURACIONES ADICIONALES (POST-INSTALACIÓN) ---
  
  // Docker Compose
  if (hasDocker) {
    await setupDockerCompose(rootDir, backend);
  }
}

/**
 * Instala el Frontend configurando dinámicamente TailwindCSS si se solicita.
 */
async function installFrontend(type, targetDir, folderName, hasTailwind) {
  if (type === 'nextjs') {
    const tailwindFlag = hasTailwind ? '--tailwind' : '--no-tailwind';
    const command = `npx create-next-app@latest "${folderName}" --ts --eslint --app --src-dir --import-alias "@/*" --use-npm ${tailwindFlag}`;
    
    await exec(command, { cwd: targetDir });
  } 
  
  else if (type === 'vite-react') {
    await exec(`npm create vite@latest "${folderName}" -- --template react-ts`, { cwd: targetDir });
    
    if (hasTailwind) {
      const viteProjectDir = path.join(targetDir, folderName);
      await exec('npm install -D tailwindcss postcss autoprefixer', { cwd: viteProjectDir });
      await exec('npx tailwindcss init -p', { cwd: viteProjectDir });
      await setupViteTailwindConfig(viteProjectDir);
    }
  }

  else if (type === 'angular') {
    await exec(`npx @angular/cli@latest new "${folderName}" --routing --style=css --skip-git --skip-install`, { cwd: targetDir });
  }
}

/**
 * Instala la tecnología de Backend seleccionada.
 */
async function installBackend(type, targetDir, folderName) {
  if (type === 'nestjs') {
    await exec(`npx @nestjs/cli new "${folderName}" --package-manager npm --strict`, { cwd: targetDir });
  } 
  
  else if (type === 'express') {
    const backendPath = path.join(targetDir, folderName);
    await fs.mkdir(backendPath, { recursive: true });
    
    await exec('npm init -y', { cwd: backendPath });
    await exec('npm install express', { cwd: backendPath });
    await exec('npm install -D typescript @types/express @types/node ts-node-dev', { cwd: backendPath });
    await exec('npx tsc --init', { cwd: backendPath });
    
    // Crear entrypoint básico
    const srcDir = path.join(backendPath, 'src');
    await fs.mkdir(srcDir, { recursive: true });
    await fs.writeFile(
      path.join(srcDir, 'index.ts'),
      `import express from 'express';\nconst app = express();\nconst PORT = process.env.PORT || 3001;\n\napp.get('/', (req, res) => {\n  res.send('API Backend Express + TS en funcionamiento');\n});\n\napp.listen(PORT, () => {\n  console.log(\`Servidor corriendo en http://localhost:\${PORT}\`);\n});`
    );
  }

  else if (type === 'fastapi') {
    const backendPath = path.join(targetDir, folderName);
    await fs.mkdir(backendPath, { recursive: true });
    await fs.writeFile(path.join(backendPath, 'requirements.txt'), 'fastapi\nuvicorn\n');
    await fs.writeFile(
      path.join(backendPath, 'main.py'),
      `from fastapi import FastAPI\n\napp = FastAPI()\n\n@app.get("/")\ndef read_root():\n    return {"message": "API Backend FastAPI en funcionamiento"}\n`
    );
  }

  else if (type === 'actix') {
    await exec(`cargo new "${folderName}" --bin`, { cwd: targetDir });
    const backendPath = path.join(targetDir, folderName);
    
    await fs.writeFile(
      path.join(backendPath, 'Cargo.toml'),
      `[package]\nname = "${folderName}"\nversion = "0.1.0"\nedition = "2021"\n\n[dependencies]\nactix-web = "4"\n`
    );
    
    await fs.writeFile(
      path.join(backendPath, 'src', 'main.rs'),
      `use actix_web::{get, App, HttpResponse, HttpServer, Responder};\n\n#[get("/")]\nasync fn hello() -> impl Responder {\n    HttpResponse::Ok().body("API Backend Actix-web en funcionamiento")\n}\n\n#[actix_web::main]\nasync fn main() -> std::io::Result<()> {\n    HttpServer::new(|| {\n        App::new().service(hello)\n    })\n    .bind(("127.0.0.1", 8080))?\n    .run()\n    .await\n}\n`
    );
  }
}

/**
 * Instala librerías seleccionadas (zod, better-auth, motion, svgl, lucide) en el entorno correcto.
 */
async function installTools(tools, rootDir, frontend, backend, isMonorepo) {
  const frontDir = isMonorepo ? path.join(rootDir, 'apps', 'frontend') : path.join(rootDir, 'frontend');
  const backDir = isMonorepo ? path.join(rootDir, 'apps', 'backend') : path.join(rootDir, 'backend');

  const hasFront = frontend !== 'none';
  const hasBack = backend !== 'none';

  for (const tool of tools) {
    if (tool === 'zod') {
      if (hasFront) await exec('npm install zod', { cwd: frontDir }).catch(() => {});
      if (hasBack && (backend === 'nestjs' || backend === 'express')) {
        await exec('npm install zod', { cwd: backDir }).catch(() => {});
      }
    }

    if (tool === 'better-auth') {
      if (hasFront) await exec('npm install better-auth', { cwd: frontDir }).catch(() => {});
      if (hasBack && (backend === 'nestjs' || backend === 'express')) {
        await exec('npm install better-auth', { cwd: backDir }).catch(() => {});
      }
    }

    if (tool === 'motion' && hasFront) {
      await exec('npm install motion', { cwd: frontDir }).catch(() => {});
    }

    if (tool === 'svgl' && hasFront) {
      // Usar svgl-react si es React (NextJS o Vite React), sino la biblioteca vanilla svgl
      const pkgName = (frontend === 'nextjs' || frontend === 'vite-react') ? 'svgl-react' : 'svgl';
      await exec(`npm install ${pkgName}`, { cwd: frontDir }).catch(() => {});
    }

    if (tool === 'lucide' && hasFront) {
      const pkgName = (frontend === 'nextjs' || frontend === 'vite-react') ? 'lucide-react' : 'lucide';
      await exec(`npm install ${pkgName}`, { cwd: frontDir }).catch(() => {});
    }
  }
}

/**
 * Helpers para inyectar configuraciones
 */
async function setupViteTailwindConfig(projectDir) {
  const configContent = `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`;
  await fs.writeFile(path.join(projectDir, 'tailwind.config.js'), configContent);

  const cssDirectives = `@tailwind base;\n@tailwind components;\n@tailwind utilities;\n`;
  const srcCssPath = path.join(projectDir, 'src', 'index.css');
  await fs.writeFile(srcCssPath, cssDirectives).catch(() => {});
}

async function setupDockerCompose(rootDir, backend) {
  let dbServices = '';
  
  if (backend !== 'none') {
    dbServices = `
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: database
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
`;
  }

  const dockerComposeContent = `version: '3.8'

services:${dbServices || '  # Agrega tus contenedores aquí'}

${dbServices ? 'volumes:\n  pgdata:\n' : ''}`;

  await fs.writeFile(path.join(rootDir, 'docker-compose.yml'), dockerComposeContent);
}