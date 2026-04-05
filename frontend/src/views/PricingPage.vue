<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import PricingModal from '../components/PricingModal.vue'
import RedeemModal from '../components/RedeemModal.vue'

const router = useRouter()
const showRedeem = ref(false)

const handleClose = () => {
  // Try to go back, but fallback to home if directly landed here
  if (window.history.state && window.history.state.back) {
    router.back()
  } else {
    router.push('/inspiration') // Or landing based on auth, handled by router anyway
  }
}
</script>

<template>
  <div class="pricing-page-wrapper">
    <!-- Render the modal directly. It has fixed positioning so it will cover the screen. -->
    <PricingModal
      @close="handleClose"
      @open-redeem="showRedeem = true"
    />

    <RedeemModal 
      v-if="showRedeem" 
      @close="showRedeem = false" 
    />
  </div>
</template>

<style scoped>
.pricing-page-wrapper {
  /* This is just an anchor element, the modal positions itself fixed */
  display: none;
}
</style>
