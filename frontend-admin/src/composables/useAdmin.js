import axios from 'axios'

const ADMIN_TOKEN_KEY = 'admin_token'

function getStoredAdminToken() {
  return localStorage.getItem(ADMIN_TOKEN_KEY) || ''
}

function saveAdminToken(token) {
  const value = (token || '').trim()
  if (value) {
    localStorage.setItem(ADMIN_TOKEN_KEY, value)
    return
  }
  localStorage.removeItem(ADMIN_TOKEN_KEY)
}

function adminHeaders() {
  const token = getStoredAdminToken()
  if (!token) return {}
  return { 'X-Admin-Token': token }
}

export function useAdmin() {
  // --- Inspirations
  async function listAdminInspirations(params = {}) {
    const { data } = await axios.get('/api/admin/inspirations', {
      params,
      headers: adminHeaders()
    })
    return data
  }

  async function reviewInspiration(postID, payload = {}) {
    const { data } = await axios.post(`/api/admin/inspirations/${postID}/review`, payload, {
      headers: adminHeaders()
    })
    return data
  }

  // --- Users
  async function listUsers(params = {}) {
    const { data } = await axios.get('/api/admin/users', {
      params,
      headers: adminHeaders()
    })
    return data
  }

  async function updateUserCredits(userID, payload = {}) {
    const { data } = await axios.post(`/api/admin/users/${userID}/credits`, payload, {
      headers: adminHeaders()
    })
    return data
  }

  async function updateUserStatus(userID, payload = {}) {
    const { data } = await axios.put(`/api/admin/users/${userID}/status`, payload, {
      headers: adminHeaders()
    })
    return data
  }

  // --- Generations
  async function listGenerations(params = {}) {
    const { data } = await axios.get('/api/admin/generations', {
      params,
      headers: adminHeaders()
    })
    return data
  }

  // --- Settings
  async function getSettings() {
    const { data } = await axios.get('/api/admin/settings', {
      headers: adminHeaders()
    })
    return data
  }

  async function updateSetting(payload = {}) {
    const { data } = await axios.put('/api/admin/settings', payload, {
      headers: adminHeaders()
    })
    return data
  }

  return {
    getStoredAdminToken,
    saveAdminToken,
    
    listAdminInspirations,
    reviewInspiration,

    listUsers,
    updateUserCredits,
    updateUserStatus,

    listGenerations,

    getSettings,
    updateSetting
  }
}
