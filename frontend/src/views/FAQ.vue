<script setup>
import { NCard, NCollapse, NCollapseItem } from 'naive-ui'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Generate generic FAQ data, easily translated if needed.
const faqs = [
  {
    title: '什么是寻氧AI？',
    content: '寻氧AI 是一个提供先进 AI 生成能力的内容平台。我们目前提供文字生成图片、逆向提示词、各类小工具等能力，帮助创作者释放无限灵感。'
  },
  {
    title: '每次生成消耗多少钻石(积分)？',
    content: '基础文生图可能每次消耗 1-5 个钻石不等，具体取决于分辨率、使用的大模型级别及是否开启高级调优。请在生成界面的参数旁查看具体耗费估算。'
  },
  {
    title: '钻石有有效期吗？如何购买？',
    content: '购买和赠送的钻石目前没有过期时间。您可以通过平台右上角的“购买钻石”按钮或前往个人主页进行充值。对于活跃用户，我们亦提供每日签到获取免费钻石的任务。'
  },
  {
    title: '如果生成失败会自动退款吗？',
    content: '是的。所有因为服务器超时、内容检测拦截或不可预知的错误导致的生成失败，系统会自动识别并全额退还该次请求锁定的钻石。'
  },
  {
    title: '生成的内容版权归谁所有？',
    content: '平台不主张对您生成内容的版权。您拥有通过合法提示词在平台上生成的所有内容的使用权、发布权及商业转化权利。'
  }
]
</script>

<template>
  <div class="faq-page">
    <div class="faq-container">
      <h1 class="page-title">
        常见问题 (FAQ)
      </h1>
      <p class="page-subtitle">
        在这里找到您需要的帮助和解答
      </p>

      <NCard
        class="faq-card"
        :bordered="false"
      >
        <NCollapse
          accordion
          :default-expanded-names="['0']"
        >
          <NCollapseItem 
            v-for="(faq, index) in faqs" 
            :key="index"
            :title="faq.title" 
            :name="index.toString()"
          >
            <div class="faq-content">
              {{ faq.content }}
            </div>
          </NCollapseItem>
        </NCollapse>
      </NCard>
    </div>
  </div>
</template>

<style scoped>
.faq-page {
  min-height: 100%;
  padding: 40px 24px;
  overflow-y: auto;
}

.faq-container {
  max-width: 800px;
  margin: 0 auto;
}

.page-title {
  font-size: 32px;
  font-weight: 800;
  text-align: center;
  color: var(--color-text-primary);
  margin-bottom: 12px;
  transform: skewX(-3deg);
}

.page-subtitle {
  text-align: center;
  color: var(--color-text-secondary);
  font-size: 16px;
  margin-bottom: 48px;
}

.faq-card {
  background: var(--color-card-bg, rgba(255, 255, 255, 0.02));
  border: 1px solid var(--color-tint-white-08);
  border-radius: 16px;
  padding: 12px;
}

:deep(.n-collapse-item) {
  margin-bottom: 16px;
}

:deep(.n-collapse-item__header-main) {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  padding: 8px 0;
}

.faq-content {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  padding-bottom: 12px;
}

@media (max-width: 768px) {
  .faq-page {
    padding: 24px 16px;
  }
  
  .page-title {
    font-size: 26px;
  }
}
</style>
