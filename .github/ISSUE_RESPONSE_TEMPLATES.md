# 📬 Koko CLI — Issue Response & Triage Templates

Guía oficial de respuestas estandarizadas para el equipo de desarrollo y producto al interactuar con issues en GitHub. Diseñadas para mantener una comunicación empática, técnica, clara y profesional.

---

## 📑 Índice de Plantillas

1. [✅ Resolución / Cierre Exitoso de Feature o Bugfix](#1--resolución--cierre-exitoso-de-feature-o-bugfix)
2. [🎯 Aceptación para Roadmap / Ready for Dev](#2--aceptación-para-roadmap--ready-for-dev)
3. [🔍 Solicitud de Información / Reproducción Mínima](#3--solicitud-de-información--reproducción-mínima)
4. [🔁 Issue Duplicada](#4--issue-duplicada)
5. [🛑 Fuera de Alcance / Decisión de Producto (Won't Fix)](#5--fuera-de-alcance--decisión-de-producto-wont-fix)

---

## 1. ✅ Resolución / Cierre Exitoso de Feature o Bugfix

> **Cuándo usar**: Al resolver una issue mediante un commit o Pull Request cerrado y validado.

```markdown
✨ **¡Completado con éxito!**

Esta funcionalidad / corrección ha sido implementada y validada en el motor de Koko CLI:
- **Resumen:** <Breve descripción de lo implementado>
- **Módulos afectados:** `<internal/...>` / `<cmd/...>`

🔗 **Commits / PRs de referencia:**
- [`<commit-hash>`](https://github.com/BlasVernazza06/koko-cli/commit/<commit-hash>) / #<pr-number>

¡Muchas gracias por la contribución a la comunidad de Koko! Cerramos esta issue. 🚀
```

---

## 2. 🎯 Aceptación para Roadmap / Ready for Dev

> **Cuándo usar**: Cuando una propuesta de la comunidad o feature request es aprobada por el equipo de producto para entrar al backlog activo.

```markdown
🎯 **Propuesta Aceptada para el Roadmap**

¡Hola @<username>! Hemos evaluado esta propuesta desde el equipo de producto y encaja perfectamente con la visión de Developer Experience (DX) de **Koko CLI**.

### 📋 Próximos pasos:
- **Milestone asignado:** `v2.x.x`
- **Etiquetas:** `status/accepted`, `area/scaffold`
- **Alcance inicial:** <Resumen de los requisitos acordados>

Si deseas contribuir con el Pull Request para esta funcionalidad, consulta nuestra guía en [CONTRIBUTING.md](CONTRIBUTING.md). ¡Seguimos avanzando! 🛠️
```

---

## 3. 🔍 Solicitud de Información / Reproducción Mínima

> **Cuándo usar**: Cuando un reporte de bug no tiene suficiente detalle para reproducir el fallo.

```markdown
👋 **¡Hola @<username>! Gracias por el reporte.**

Para poder reproducir y solucionar este comportamiento de forma precisa, ¿podrías compartirnos los siguientes detalles de tu entorno?

1. **Versión exacta de Koko CLI:** (ej. `koko --version` o `npx koko-app --version`)
2. **Sistema Operativo:** (Windows PowerShell / macOS / Linux WSL)
3. **Comando o Flags exactos utilizados:** (ej. `koko init my-app --recipe saas --pm pnpm`)
4. **Log de error completo:** (ejecuta agregando `--verbose` o `--debug`)

Quedamos atentos a tus comentarios para asistirte de inmediato. 🔍
```

---

## 4. 🔁 Issue Duplicada

> **Cuándo usar**: Cuando el reporte o feature request ya existe en otra issue abierta o cerrada.

```markdown
👋 **¡Hola @<username>! Gracias por abrir esta issue.**

Hemos identificado que este tema ya se encuentra en seguimiento en la siguiente issue activa:
👉 #<issue-number> (<Título de la issue principal>)

Para centralizar la discusión técnica y evitar duplicación de esfuerzos, procedemos a cerrar esta issue y continuar el tracking en el hilo principal. ¡Te invitamos a dejar tu feedback allí! 🔗
```

---

## 5. 🛑 Fuera de Alcance / Decisión de Producto (Won't Fix)

> **Cuándo usar**: Cuando una solicitud entra en conflicto con los principios de producto (ej. mantener el binario ligero, combinaciones incompatibles u obsoletas).

```markdown
👋 **¡Hola @<username>! Gracias por tu propuesta.**

Hemos analizado detalladamente este requerimiento con el equipo de producto. En este momento, hemos decidido no incorporar esta funcionalidad directamente en el core de **Koko CLI** por las siguientes razones:

- **Motivo principal:** <Explicar: peso del binario / compatibilidad de monorepos / stack no alineado con la filosofía de Koko>
- **Alternativa recomendada:** <Sugerir cómo puede lograrlo el usuario manualmente o mediante addons post-init>

Agradecemos enormemente tu tiempo y feedback para hacer crecer la herramienta. Cerramos esta issue por decisión de diseño. 🛡️
```
