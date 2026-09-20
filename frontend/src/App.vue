<script setup lang="ts">
const title = "Bill Splitter"
</script>

<template>
  <main>
    <h1>{{ title }}</h1>
    <p>The application is ready for bill-share editing.</p>
  </main>
</template>

<style>
:root {
  font-family: system-ui, sans-serif;
  color: #17202a;
  background: #f7f8fa;
}

body {
  margin: 0;
}

main {
  max-width: 42rem;
  margin: 5rem auto;
  padding: 2rem;
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 0.5rem 2rem rgb(23 32 42 / 8%);
}
</style>
