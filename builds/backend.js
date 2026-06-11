async function installBackend(type, targetDir, folderName) {
    if (type === 'nestjs') {
      // NestJS CLI crea su propia carpeta
      await exec(`npx @nestjs/cli new "${folderName}" --package-manager npm --strict`, { cwd: targetDir });
    } 
    
    else if (type === 'express') {
      const backendPath = path.join(targetDir, folderName);
      await fs.mkdir(backendPath, { recursive: true });
      
      // Inicializar un proyecto de Express minimalista con TypeScript
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
      // Crear archivo de requerimientos y main.py para Python
      await fs.writeFile(path.join(backendPath, 'requirements.txt'), 'fastapi\nuvicorn\n');
      await fs.writeFile(
        path.join(backendPath, 'main.py'),
        `from fastapi import FastAPI\n\napp = FastAPI()\n\n@app.get("/")\ndef read_root():\n    return {"message": "API Backend FastAPI en funcionamiento"}\n`
      );
    }
    else if (type === 'actix') {
      // Inicializar Rust con cargo
      await exec(`cargo new "${folderName}" --bin`, { cwd: targetDir });
      const backendPath = path.join(targetDir, folderName);
      
      // Configurar Cargo.toml básico
      await fs.writeFile(
        path.join(backendPath, 'Cargo.toml'),
        `[package]\nname = "${folderName}"\nversion = "0.1.0"\nedition = "2021"\n\n[dependencies]\nactix-web = "4"\n`
      );
      
      // Crear main.rs básico
      await fs.writeFile(
        path.join(backendPath, 'src', 'main.rs'),
        `use actix_web::{get, App, HttpResponse, HttpServer, Responder};\n\n#[get("/")]\nasync fn hello() -> impl Responder {\n    HttpResponse::Ok().body("API Backend Actix-web en funcionamiento")\n}\n\n#[actix_web::main]\nasync fn main() -> std::io::Result<()> {\n    HttpServer::new(|| {\n        App::new().service(hello)\n    })\n    .bind(("127.0.0.1", 8080))?\n    .run()\n    .await\n}\n`
      );
    }
  }