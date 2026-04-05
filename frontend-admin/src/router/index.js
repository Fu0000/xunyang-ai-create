import { createRouter, createWebHistory } from 'vue-router'
import { useAdmin } from '../composables/useAdmin'

const Login = () => import('../views/Login.vue')
const AdminLayout = () => import('../views/AdminLayout.vue')
const InspirationReview = () => import('../views/InspirationReview.vue')
const UserList = () => import('../views/UserList.vue')
const GenerationList = () => import('../views/GenerationList.vue')

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/',
    component: AdminLayout,
    children: [
      {
        path: '',
        redirect: '/inspirations'
      },
      {
        path: 'inspirations',
        name: 'inspirations',
        component: InspirationReview,
        meta: { title: '灵感内容审核', subTitle: '发布内容审核、筛选与处理', moduleKey: 'inspiration_review' }
      },
      {
        path: 'users',
        name: 'users',
        component: UserList,
        meta: { title: '用户列表', subTitle: '账号、状态与资产管理模块', moduleKey: 'user_list' }
      },
      {
        path: 'generations',
        name: 'generations',
        component: GenerationList,
        meta: { title: '生成列表', subTitle: '任务、产物和调用拦截追踪', moduleKey: 'generation_list' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

router.beforeEach((to, from, next) => {
  const { getStoredAdminToken } = useAdmin()
  const hasToken = !!getStoredAdminToken()
  
  if (to.name !== 'login' && !hasToken) {
    next({ name: 'login' })
  } else if (to.name === 'login' && hasToken) {
    next({ path: '/' })
  } else {
    next()
  }
})

export default router
