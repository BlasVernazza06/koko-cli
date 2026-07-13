<script setup lang="ts">
import { ref } from 'vue'
const count = ref(0)
</script>

<template>
  <div class="container">
    <header>
      <h1>[[.ProjectName]]</h1>
      <p>Welcome to your Vue 3 + Vite + TypeScript frontend, bootstrapped with Koko CLI.</p>
    </header>
    <main>
      <div class="card">
        <button type="button" @click="count++">count is {{ count }}</button>
      </div>
    </main>
  </div>
</template>

<style scoped>
.container {
  font-family: sans-serif;
  text-align: center;
  padding: 2rem;
}
.card {
  margin: 2rem 0;
}
button {
  padding: 0.6em 1.2em;
  font-size: 1em;
  font-weight: 500;
  cursor: pointer;
}
</style>
