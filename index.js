#!/usr/bin/env node
import { intro, select, text, multiselect, outro, spinner, cancel, isCancel } from '@clack/prompts';
import color from 'picocolors';
import { scaffoldProject } from './scaffolder.js';

// Mascot design and dialog bubble helper
const MASCOT = `
     ▄███████▄
    ██  • ◡ •  ██   🦾  { CLAWY }
    █████████████
    ██  ▀███▀  ██
     ▀█████████▀
     ▄█▀     ▀█▄
`;

function say(text) {
  return `${color.cyan('🦾 Clawy:')} "${color.italic(text)}"`;
}

async function main() {
  console.clear(); 

  // Welcome screen
  console.log(color.cyan(MASCOT));
  intro(color.bgCyan(color.black(' 🛠️  CREADOR DE PROYECTOS - CLAW CLI  ')));

  console.log(say('¡Hola! Soy tu asistente de andamiaje. Vamos a estructurar tu próximo proyecto.'));
  console.log();

  // 1. Project name
  const projectName = await text({
    message: color.bold('¿Cómo se llamará el proyecto?'),
    placeholder: 'mi-super-app',
    validate(value) {
      if (value.trim().length === 0) return 'Por favor, ingresa un nombre para el proyecto.';
      if (/[^a-zA-Z0-9-_]/.test(value)) return 'El nombre solo puede contener letras, números, guiones y guiones bajos.';
    }
  });

  if (isCancel(projectName)) {
    cancel('Operación cancelada.');
    return;
  }

  console.log(say(`¡Excelente nombre! Crearemos la estructura en la carpeta: ${color.green(projectName)}`));
  console.log();

  // 2. Frontend Selection
  const frontend = await select({
    message: color.bold('¿Qué Frontend deseas usar?'),
    options: [
      { value: 'nextjs', label: '⚡ Next.js', hint: 'React, App Router, SSR' },
      { value: 'vite-react', label: '⚛️ React (Vite)', hint: 'SPA ligero y rápido' },
      { value: 'angular', label: '🅰️ Angular', hint: 'Framework robusto de Google' },
      { value: 'none', label: '❌ Sin Frontend', hint: 'Solo API / Servicio de backend o CLI' }
    ]
  });

  if (isCancel(frontend)) {
    cancel('Operación cancelada.');
    return;
  }

  // 3. Backend Selection
  const backend = await select({
    message: color.bold('¿Qué tecnología para el Backend?'),
    options: [
      { value: 'nestjs', label: '🦁 NestJS', hint: 'TypeScript, arquitectura modular' },
      { value: 'express', label: '🚂 Express.js', hint: 'Mínimo y flexible en JS/TS' },
      { value: 'fastapi', label: '⚡ FastAPI', hint: 'Python, alto rendimiento, autodocumentado' },
      { value: 'actix', label: '🦀 Actix-web', hint: 'Rust, extremadamente veloz y seguro' },
      { value: 'none', label: '❌ Sin Backend', hint: 'Desarrollo puramente frontend o estático' }
    ]
  });

  if (isCancel(backend)) {
    cancel('Operación cancelada.');
    return;
  }

  // 4. Additional tools (now includes zod, better-auth, motion, svgl, lucide)
  const tools = await multiselect({
    message: color.bold('Selecciona herramientas de entorno y librerías adicionales:'),
    options: [
      { value: 'docker', label: '🐳 Docker Compose', hint: 'Contenedores listos para bases de datos' },
      { value: 'turborepo', label: '🏎️ Turborepo / Monorrepo', hint: 'Gestión inteligente de múltiples subproyectos' },
      { value: 'eslint-prettier', label: '✨ ESLint + Prettier', hint: 'Linter y formateador de código automático' },
      { value: 'tailwind', label: '🎨 TailwindCSS', hint: 'Framework CSS de utilidad instantánea' },
      { value: 'zod', label: '🛡️ Zod', hint: 'Validación de esquemas de datos TypeScript' },
      { value: 'better-auth', label: '🔑 Better Auth', hint: 'Framework de autenticación moderno' },
      { value: 'motion', label: '🎬 Motion', hint: 'Animaciones fluidas y ligeras (Framer Motion)' },
      { value: 'svgl', label: '📐 SVGL', hint: 'Librería e integración de logos SVG' },
      { value: 'lucide', label: '🎨 Lucide Icons', hint: 'Iconos vectoriales modernos y limpios' }
    ],
    required: false
  });

  if (isCancel(tools)) {
    cancel('Operación cancelada.');
    return;
  }

  console.log();
  const s = spinner();
  s.start(color.yellow('Generando la estructura real del proyecto en disco (esto puede tardar unos segundos)...'));
  
  try {
    // Ejecutar la lógica de andamiaje real
    await scaffoldProject({ projectName, frontend, backend, tools });
    s.stop(color.green('¡Estructura de archivos y configuración creada con éxito!'));
  } catch (error) {
    s.stop(color.red('Ocurrió un error al configurar el proyecto en el disco.'));
    console.error(color.red(`\nDetalle del error:\n${error.message}`));
    return;
  }

  console.log();
  console.log(color.cyan(`
  🦾 ${color.bold('Clawy dice:')} 
  "¡Todo listo! He diseñado e instalado la siguiente configuración para ti:"
  `));

  console.log(color.gray(`  ┌──────────────────────────────────────────┐`));
  console.log(`  │  📁 Carpeta:    ${color.cyan(projectName.padEnd(25))} │`);
  console.log(`  │  💻 Frontend:   ${color.cyan(frontend.padEnd(25))} │`);
  console.log(`  │  🔌 Backend:    ${color.cyan(backend.padEnd(25))} │`);
  console.log(`  │  🛠️  Herramientas: ${color.cyan((tools.length ? tools.join(', ') : 'Ninguna').slice(0, 23).padEnd(23))} │`);
  console.log(color.gray(`  └──────────────────────────────────────────┘`));
  console.log();

  outro(color.bgCyan(color.black('  ¡Proyecto configurado y listo! Abre la carpeta para comenzar.  ')));
}

main().catch(console.error);