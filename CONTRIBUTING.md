# Guía de Contribución a Claw CLI 🦾

¡Gracias por tu interés en mejorar **Claw**! Las contribuciones de la comunidad son lo que hace que el desarrollo open-source sea tan increíble. Aquí te explicamos cómo puedes ayudarnos y participar en el proyecto.

## 🚀 Cómo Empezar

### Requisitos Previos
Asegúrate de tener instalado en tu máquina local:
- [Node.js](https://nodejs.org/) (versión v18 o superior recomendada)
- npm o pnpm / yarn

### Paso 1: Bifurcar (Fork) el Proyecto
1. Haz un fork de este repositorio a tu propia cuenta de GitHub.
2. Clona tu bifurcación localmente:
   ```bash
   git clone https://github.com/TU-USUARIO/claw.git
   cd claw
   ```

### Paso 2: Configurar tu Entorno
Instala las dependencias del proyecto:
```bash
npm install
```

### Paso 3: Probar los Cambios Localmente
Puedes ejecutar el CLI directamente desde el código fuente sin compilar o publicar, usando Node:
```bash
node index.js
```
O bien, puedes enlazarlo globalmente de forma temporal en tu máquina:
```bash
npm link
```
Una vez enlazado, podrás ejecutar el comando asignado en tu `package.json` (por ejemplo, `crear-proyecto` o `claw`) de forma global desde cualquier carpeta en tu sistema de archivos.

---

## 🛠️ Flujo de Trabajo para Contribuciones

1. **Crea una nueva rama** para tu funcionalidad o corrección de error:
   ```bash
   git checkout -b feature/nueva-tecnologia
   ```
2. **Realiza tus cambios** y asegúrate de seguir las buenas prácticas de desarrollo.
3. **Prueba a fondo**: Verifica que los archivos andamiados se creen de forma correcta y que no rompa las opciones existentes.
4. **Haz commit** de tus cambios con mensajes descriptivos:
   ```bash
   git commit -m "feat: añadir opción de Fastify para el backend"
   ```
5. **Sube (Push) tu rama** a tu repositorio fork:
   ```bash
   git push origin feature/nueva-tecnologia
   ```
6. Abre un **Pull Request (PR)** hacia la rama principal del repositorio original.

---

## 📋 Directrices de Código

- **Estética Limpia**: Mantenemos un diseño interactivo sumamente visual y limpio en la consola. Cualquier nuevo prompt o mensaje debe respetar la paleta de colores y la personalidad de nuestra mascota `Clawy` (`picocolors` y `@clack/prompts`).
- **Validaciones**: Siempre valida que las entradas del usuario no estén vacías y tengan formatos válidos antes de proceder a la creación física en disco.
- **Formateo**: Te animamos a mantener un código limpio y legible.

Si tienes dudas o necesitas discutir alguna propuesta técnica compleja antes de programarla, siéntete libre de abrir un Issue con la etiqueta `enhancement`. ¡Esperamos tu Pull Request!
