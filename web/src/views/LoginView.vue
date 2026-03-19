<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="title">SSH2A</h1>
      <p class="subtitle">SSH Access Authorization</p>

      <div v-if="status" class="status-info">
        <div class="status-row">
          <span class="label">Your IP</span>
          <span class="value">{{ status.ip }}</span>
        </div>
        <div class="status-row">
          <span class="label">SSH Record</span>
          <span :class="['badge', status.has_record ? 'badge-warn' : 'badge-ok']">
            {{ status.has_record ? 'Detected' : 'None' }}
          </span>
        </div>
        <div class="status-row">
          <span class="label">Status</span>
          <span :class="['badge', status.verified ? 'badge-ok' : 'badge-dim']">
            {{ status.verified ? 'Verified' : 'Not Verified' }}
          </span>
        </div>
      </div>

      <div v-if="status?.verified" class="success-msg">
        SSH access granted. You can now connect via SSH.
      </div>

      <div v-else-if="status && !status.has_record" class="hint-msg">
        No SSH access attempt detected from your IP. Please try connecting via SSH first.
      </div>

      <form v-else-if="status && status.has_record" @submit.prevent="handleAuth" class="auth-form">
        <input
          v-model="password"
          type="password"
          placeholder="Password or 2FA code"
          class="input"
          autofocus
          :disabled="loading"
        />
        <button type="submit" class="btn" :disabled="loading">
          {{ loading ? '...' : 'Verify' }}
        </button>
      </form>

      <p v-if="error" class="error-msg">{{ error }}</p>

      <router-link to="/admin" class="admin-link">Admin Panel</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getStatus, postAuth, type StatusResponse } from '../api'

const status = ref<StatusResponse | null>(null)
const password = ref('')
const error = ref('')
const loading = ref(false)

const fetchStatus = async () => {
  try {
    const res = await getStatus()
    status.value = res.data
  } catch {
    error.value = 'Failed to fetch status'
  }
}

const handleAuth = async () => {
  if (!password.value || !status.value) return
  loading.value = true
  error.value = ''
  try {
    await postAuth({ password: password.value, remoteIP: status.value.ip })
    await fetchStatus()
    password.value = ''
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Authentication failed'
  } finally {
    loading.value = false
  }
}

onMounted(fetchStatus)
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
  text-align: center;
}

.title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 2px;
  color: var(--accent);
}

.subtitle {
  color: var(--text-dim);
  margin-top: 4px;
  margin-bottom: 28px;
  font-size: 14px;
}

.status-info {
  background: var(--bg);
  border-radius: var(--radius);
  padding: 16px;
  margin-bottom: 24px;
  text-align: left;
}

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}

.status-row .label {
  color: var(--text-dim);
  font-size: 13px;
}

.status-row .value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
}

.badge {
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.badge-ok {
  background: rgba(81, 207, 102, 0.15);
  color: var(--success);
}

.badge-warn {
  background: rgba(252, 196, 25, 0.15);
  color: var(--warning);
}

.badge-dim {
  background: rgba(139, 143, 163, 0.15);
  color: var(--text-dim);
}

.auth-form {
  display: flex;
  gap: 8px;
}

.input {
  flex: 1;
  padding: 10px 14px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.input:focus {
  border-color: var(--accent);
}

.btn {
  padding: 10px 20px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: var(--radius);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-msg {
  color: var(--danger);
  font-size: 13px;
  margin-top: 12px;
}

.success-msg {
  background: rgba(81, 207, 102, 0.1);
  border: 1px solid rgba(81, 207, 102, 0.3);
  border-radius: var(--radius);
  padding: 14px;
  color: var(--success);
  font-size: 14px;
  margin-bottom: 16px;
}

.hint-msg {
  background: rgba(139, 143, 163, 0.1);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 14px;
  color: var(--text-dim);
  font-size: 14px;
  margin-bottom: 16px;
}

.admin-link {
  display: inline-block;
  margin-top: 24px;
  color: var(--text-dim);
  font-size: 13px;
  text-decoration: none;
  transition: color 0.2s;
}

.admin-link:hover {
  color: var(--accent);
}
</style>
