<script setup>
import zhCN from 'ant-design-vue/es/locale/zh_CN'
</script>

<template>
  <a-config-provider
    :locale="zhCN"
    :theme="{
      token: {
        colorPrimary: '#4f46e5',
        borderRadius: 12,
        fontFamily: '\'Inter\', \'Manrope\', \'PingFang SC\', \'Microsoft YaHei\', sans-serif',
        colorTextBase: '#0f172a',
        colorBgBase: '#ffffff'
      },
      components: {
        Layout: {
          headerBg: 'rgba(255, 255, 255, 0.8)',
          siderBg: '#ffffff',
          bodyBg: '#f8fafc'
        },
        Card: {
          boxShadowTertiary: '0 4px 20px rgba(0, 0, 0, 0.03)',
          colorBorderSecondary: 'rgba(0,0,0,0.03)'
        },
        Menu: {
          itemBorderRadius: 8
        }
      }
    }"
  >
    <a-app>
      <router-view />
    </a-app>
  </a-config-provider>
</template>
