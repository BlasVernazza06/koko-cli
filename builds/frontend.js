export async function installFrontend(type, targetDir, dirName, hasTailwind) {
    if (type === 'nextjs') {
        const tailwindFlag = hasTailwind ? '--tailwind' : '--no-tailwind';
        const command = `npx create-next-app@latest "${folderName}" --ts --eslint --app --src-dir --import-alias "@/*" --use-npm ${tailwindFlag}`;
        
        await exec(command, { cwd: targetDir });
    } 

    else if (type === 'vite-react') {
        // Nota: create-vite no configura tailwind por defecto con una bandera simple.
        // Inicializamos React + TypeScript con Vite.
        await exec(`npm create vite@latest "${folderName}" -- --template react-ts`, { cwd: targetDir });
        
        // Si eligió Tailwind, podemos instalarlo y configurar los archivos básicos
        if (hasTailwind) {
            const viteProjectDir = path.join(targetDir, folderName);
            // Instalamos tailwindcss y sus dependencias
            await exec('npm install -D tailwindcss postcss autoprefixer', { cwd: viteProjectDir });
            await exec('npx tailwindcss init -p', { cwd: viteProjectDir });
            
            // Aquí se podrían inyectar plantillas de configuración básicas de Tailwind para Vite
            await setupViteTailwindConfig(viteProjectDir);
        }
    }
    else if (type === 'angular') {
        // Instalación de Angular
        await exec(`npx @angular/cli@latest new "${folderName}" --routing --style=css --skip-git --skip-install`, { cwd: targetDir });
    }
}