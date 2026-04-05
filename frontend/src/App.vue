<script setup>
import { onMounted, computed } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme, zhCN, dateZhCN, enUS, dateEnUS } from 'naive-ui'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from './stores/user'
import { useThemeStore } from './stores/theme'
import { useLocaleStore } from './stores/locale'
import AuthModal from './components/AuthModal.vue'
import InviteModal from './components/InviteModal.vue'
import RedeemModal from './components/RedeemModal.vue'
import PricingModal from './components/PricingModal.vue'
import BindEmailModal from './components/BindEmailModal.vue'

const userStore = useUserStore()
const themeStore = useThemeStore()
const localeStore = useLocaleStore()
const router = useRouter()
const route = useRoute()

const handleLoginSuccess = (user, token) => {
  userStore.loginSuccess(user, token)
  // 登录后如果在落地页，跳转到创作页
  if (route.name === 'landing') {
    const redirect = route.query.redirect
    router.push(redirect || '/inspiration')
  }
}

const handlePricingToRedeem = () => {
  userStore.closePricing()
  userStore.openRedeem()
}

const handleRedeemToPricing = () => {
  userStore.closeRedeem()
  userStore.openPricing()
}

// 在 setup 阶段立即初始化用户状态（早于子组件 onMounted），
// 避免子组件在 onMounted 中读到空 token。
userStore.init()

const naiveTheme = computed(() => themeStore.isDark ? darkTheme : null)
const naiveLocale = computed(() => localeStore.locale === 'en' ? enUS : zhCN)
const naiveDateLocale = computed(() => localeStore.locale === 'en' ? dateEnUS : dateZhCN)

// Naive UI 主题覆盖 - 根据当前主题返回不同配置
const themeOverrides = computed(() => {
  if (themeStore.isDark) {
    return {
      common: {
        primaryColor: '#6366f1',
        primaryColorHover: '#818cf8',
        primaryColorPressed: '#4f46e5',
        primaryColorSuppl: '#6366f1',
        successColor: '#10b981',
        successColorHover: '#34d399',
        errorColor: '#ef4444',
        errorColorHover: '#f87171',
        warningColor: '#f59e0b',
        warningColorHover: '#fbbf24',
        bodyColor: '#0a0e17',
        cardColor: 'rgba(21, 27, 43, 0.95)',
        modalColor: 'rgba(21, 27, 43, 0.98)',
        popoverColor: '#1e2538',
        inputColor: 'rgba(255, 255, 255, 0.05)',
        inputColorDisabled: 'rgba(255, 255, 255, 0.02)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        hoverColor: 'rgba(255, 255, 255, 0.08)',
        borderRadius: '16px',
        borderRadiusSmall: '8px',
        fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
        fontSize: '14px',
        textColorBase: '#f8fafc',
        textColor1: 'rgba(255, 255, 255, 0.92)',
        textColor2: 'rgba(255, 255, 255, 0.65)',
        textColor3: 'rgba(255, 255, 255, 0.45)',
        placeholderColor: 'rgba(255, 255, 255, 0.25)',
        closeIconColor: 'rgba(255, 255, 255, 0.4)',
        closeIconColorHover: 'rgba(255, 255, 255, 0.9)',
        closeColorHover: 'rgba(255, 255, 255, 0.1)',
        closeColorPressed: 'rgba(255, 255, 255, 0.15)'
      },
      Button: {
        borderRadiusMedium: '12px',
        borderRadiusSmall: '8px',
        fontWeight: '600',
        fontSizeMedium: '14px',
        heightMedium: '42px',
        paddingMedium: '0 20px',
        colorPrimary: '#6366f1',
        colorHoverPrimary: '#818cf8',
        colorPressedPrimary: '#4f46e5',
        textColorPrimary: '#fff',
        textColorGhostPrimary: '#818cf8',
        borderPrimary: 'none',
        borderHoverPrimary: 'none',
        borderPressedPrimary: 'none'
      },
      Input: {
        borderRadius: '12px',
        heightMedium: '42px',
        fontSizeMedium: '14px',
        color: 'rgba(255,255,255,0.05)',
        colorFocus: 'rgba(255,255,255,0.07)',
        border: '1px solid rgba(255,255,255,0.08)',
        borderHover: '1px solid rgba(255,255,255,0.18)',
        borderFocus: '1px solid rgba(99,102,241,0.5)',
        boxShadowFocus: '0 0 0 4px rgba(99,102,241,0.15)',
        placeholderColor: 'rgba(255,255,255,0.3)',
        colorDisabled: 'rgba(255,255,255,0.03)',
        textColorDisabled: 'rgba(255,255,255,0.35)',
        caretColor: '#818cf8'
      },
      Card: {
        borderRadius: '24px',
        color: 'rgba(21, 27, 43, 0.95)',
        borderColor: 'rgba(255, 255, 255, 0.08)',
        paddingMedium: '0',
        paddingSmall: '0',
        titleFontSizeMedium: '18px',
        titleFontWeight: '700',
        titleTextColor: '#f8fafc',
        closeIconColor: 'rgba(255, 255, 255, 0.4)',
        closeIconColorHover: 'rgba(255, 255, 255, 0.9)',
        closeColorHover: 'rgba(255, 255, 255, 0.1)'
      },
      Modal: {
        color: 'rgba(21, 27, 43, 0.98)',
        borderRadius: '24px',
        boxShadow: '0 40px 120px rgba(0,0,0,0.8), 0 0 80px rgba(99,102,241,0.1), inset 0 1px 0 0 rgba(255, 255, 255, 0.1)'
      },
      Tabs: {
        tabBorderRadius: '10px',
        tabFontSizeMedium: '14px',
        tabFontWeightActive: '600',
        tabGapMediumSegment: '4px',
        tabPaddingMediumSegment: '8px 16px',
        colorSegment: 'rgba(255,255,255,0.05)',
        tabColorSegment: 'transparent',
        tabTextColorSegment: 'rgba(255,255,255,0.45)',
        tabTextColorActiveSegment: '#fff',
        tabTextColorHoverSegment: 'rgba(255,255,255,0.75)'
      },
      Alert: {
        borderRadius: '12px',
        fontSize: '14px',
        iconSizeMedium: '20px',
        padding: '12px 16px'
      },
      Spin: { color: '#818cf8' },
      Tag: { borderRadius: '8px' },
      Empty: { textColor: 'rgba(255,255,255,0.4)', iconColor: 'rgba(255,255,255,0.2)', iconSizeMedium: '48px' }
    }
  } else {
    // Light theme overrides
    return {
      common: {
        primaryColor: '#4f46e5',
        primaryColorHover: '#818cf8',
        primaryColorPressed: '#3730a3',
        primaryColorSuppl: '#4f46e5',
        successColor: '#059669',
        successColorHover: '#10b981',
        errorColor: '#dc2626',
        errorColorHover: '#ef4444',
        warningColor: '#d97706',
        warningColorHover: '#f59e0b',
        bodyColor: '#f8fafc',
        cardColor: '#ffffff',
        modalColor: '#ffffff',
        popoverColor: '#ffffff',
        inputColor: '#ffffff',
        inputColorDisabled: 'rgba(0, 0, 0, 0.02)',
        borderColor: 'rgba(0, 0, 0, 0.08)',
        hoverColor: 'rgba(0, 0, 0, 0.05)',
        borderRadius: '16px',
        borderRadiusSmall: '8px',
        fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
        fontSize: '14px',
        textColorBase: '#0f172a',
        textColor1: 'rgba(0, 0, 0, 0.88)',
        textColor2: 'rgba(0, 0, 0, 0.65)',
        textColor3: 'rgba(0, 0, 0, 0.45)',
        placeholderColor: 'rgba(0, 0, 0, 0.3)',
        closeIconColor: 'rgba(0, 0, 0, 0.4)',
        closeIconColorHover: 'rgba(0, 0, 0, 0.8)',
        closeColorHover: 'rgba(0, 0, 0, 0.08)',
        closeColorPressed: 'rgba(0, 0, 0, 0.12)'
      },
      Button: {
        borderRadiusMedium: '12px',
        borderRadiusSmall: '8px',
        fontWeight: '600',
        fontSizeMedium: '14px',
        heightMedium: '42px',
        paddingMedium: '0 20px',
        colorPrimary: '#4f46e5',
        colorHoverPrimary: '#6366f1',
        colorPressedPrimary: '#3730a3',
        textColorPrimary: '#fff',
        textColorGhostPrimary: '#4f46e5',
        borderPrimary: 'none',
        borderHoverPrimary: 'none',
        borderPressedPrimary: 'none'
      },
      Input: {
        borderRadius: '12px',
        heightMedium: '42px',
        fontSizeMedium: '14px',
        color: '#ffffff',
        colorFocus: '#ffffff',
        border: '1px solid rgba(0,0,0,0.1)',
        borderHover: '1px solid rgba(0,0,0,0.2)',
        borderFocus: '1px solid rgba(79,70,229,0.5)',
        boxShadowFocus: '0 0 0 4px rgba(79,70,229,0.1)',
        placeholderColor: 'rgba(0,0,0,0.35)',
        colorDisabled: 'rgba(0,0,0,0.03)',
        textColorDisabled: 'rgba(0,0,0,0.3)',
        caretColor: '#4f46e5'
      },
      Card: {
        borderRadius: '24px',
        color: '#ffffff',
        borderColor: 'rgba(0, 0, 0, 0.08)',
        paddingMedium: '0',
        paddingSmall: '0',
        titleFontSizeMedium: '18px',
        titleFontWeight: '700',
        titleTextColor: '#0f172a',
        closeIconColor: 'rgba(0, 0, 0, 0.4)',
        closeIconColorHover: 'rgba(0, 0, 0, 0.8)',
        closeColorHover: 'rgba(0, 0, 0, 0.08)'
      },
      Modal: {
        color: '#ffffff',
        borderRadius: '24px',
        boxShadow: '0 40px 120px rgba(0,0,0,0.12), 0 0 80px rgba(79,70,229,0.06)'
      },
      Tabs: {
        tabBorderRadius: '10px',
        tabFontSizeMedium: '14px',
        tabFontWeightActive: '600',
        tabGapMediumSegment: '4px',
        tabPaddingMediumSegment: '8px 16px',
        colorSegment: 'rgba(0,0,0,0.05)',
        tabColorSegment: 'transparent',
        tabTextColorSegment: 'rgba(0,0,0,0.45)',
        tabTextColorActiveSegment: '#0f172a',
        tabTextColorHoverSegment: 'rgba(0,0,0,0.7)'
      },
      Alert: {
        borderRadius: '12px',
        fontSize: '14px',
        iconSizeMedium: '20px',
        padding: '12px 16px'
      },
      Spin: { color: '#4f46e5' },
      Tag: { borderRadius: '8px' },
      Empty: { textColor: 'rgba(0,0,0,0.4)', iconColor: 'rgba(0,0,0,0.2)', iconSizeMedium: '48px' }
    }
  }
})

onMounted(() => {
  // 检查 URL 邀请码
  const urlParams = new URLSearchParams(window.location.search)
  const inviteCode = urlParams.get('invite')
  if (inviteCode) {
    localStorage.setItem('pendingInviteCode', inviteCode)
    window.history.replaceState({}, document.title, window.location.pathname)
    if (!userStore.isLoggedIn) {
      userStore.openAuth()
    }
  }
})
</script>

<template>
  <NConfigProvider
    :theme="naiveTheme"
    :theme-overrides="themeOverrides"
    :locale="naiveLocale"
    :date-locale="naiveDateLocale"
  >
    <NMessageProvider>
      <NDialogProvider>
        <div class="app-root">
          <!-- 全局粒子背景 -->
          <div class="particle-field">
            <div
              v-for="i in 30"
              :key="i"
              class="particle"
              :style="{
                '--x': Math.random() * 100 + '%',
                '--y': Math.random() * 100 + '%',
                '--size': (Math.random() * 3 + 1) + 'px',
                '--duration': (Math.random() * 20 + 10) + 's',
                '--delay': (Math.random() * 10) + 's',
                '--opacity': Math.random() * 0.5 + 0.1
              }"
            />
          </div>
          <router-view />

          <!-- 全局弹窗 -->
          <AuthModal
            v-if="userStore.showAuthModal"
            @login-success="handleLoginSuccess"
            @close="userStore.closeAuth()"
          />
          <RedeemModal
            v-if="userStore.showRedeemModal && userStore.isLoggedIn"
            @close="userStore.closeRedeem()"
            @success="userStore.fetchUserInfo()"
            @open-pricing="handleRedeemToPricing"
          />
          <PricingModal
            v-if="userStore.showPricingModal"
            @close="userStore.closePricing()"
            @open-redeem="handlePricingToRedeem"
          />
          <InviteModal
            v-if="userStore.showInviteModal"
            @close="userStore.closeInvite()"
            @credits-updated="userStore.fetchUserInfo()"
          />
          <BindEmailModal
            v-if="userStore.showBindEmailModal"
          />
        </div>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style scoped>
.app-root {
  width: 100%;
  min-height: 100vh;
  position: relative;
}

/* 全局粒子背景 */
.particle-field {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.particle {
  position: absolute;
  left: var(--x);
  top: var(--y);
  width: var(--size);
  height: var(--size);
  background: rgba(var(--particle-color), var(--opacity));
  border-radius: 50%;
  animation: particleFloat var(--duration) ease-in-out var(--delay) infinite;
}

@keyframes particleFloat {
  0%, 100% { transform: translate(0, 0) scale(1); opacity: var(--opacity); }
  25% { transform: translate(30px, -40px) scale(1.2); }
  50% { transform: translate(-20px, -80px) scale(0.8); opacity: calc(var(--opacity) * 1.5); }
  75% { transform: translate(40px, -40px) scale(1.1); }
}
</style>
