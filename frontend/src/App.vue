<template>
  <div class="ambient-glow-bottom"></div>
  <router-view />
  <BackToTop />
  <MagisteriumFloatingModal />
  <AccessGateModal />
</template>

<script setup lang="ts">
import BackToTop from './components/common/BackToTop.vue';
import MagisteriumFloatingModal from './components/magisterium/MagisteriumFloatingModal.vue';
import AccessGateModal from './components/common/AccessGateModal.vue';
</script>