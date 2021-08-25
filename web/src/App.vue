<template>
  <div id="app" class="leading-normal tracking-normal text-white gradient">
    <LandingPage />
  </div>
</template>

<script>
import LandingPage from "./components/Navbar.vue";

export default {
  name: "App",
  components: {
    LandingPage,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
