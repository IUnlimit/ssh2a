<template>
  <div class="admin-page">
    <header class="header">
      <h1>SSH2A Admin</h1>
      <router-link to="/" class="back-link">Back to Login</router-link>
    </header>

    <!-- Stats Overview -->
    <div v-if="stats" class="stats-grid">
      <div class="stat-card">
        <span class="stat-value">{{ stats.total_ssh_attempts }}</span>
        <span class="stat-label">SSH Attempts</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.total_rejected }}</span>
        <span class="stat-label">Rejected</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.total_honeypot }}</span>
        <span class="stat-label">Honeypot</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.total_forwarded }}</span>
        <span class="stat-label">Forwarded</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.total_auth_success }}</span>
        <span class="stat-label">Auth Success</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.unique_honeypot_creds }}</span>
        <span class="stat-label">Unique Creds</span>
      </div>
      <div class="stat-card">
        <span class="stat-value">{{ stats.unique_rejected_ips }}</span>
        <span class="stat-label">Unique Rejected IPs</span>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab', { active: activeTab === tab.key }]"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Honeypot Credentials -->
    <div v-if="activeTab === 'honeypot'" class="table-wrap">
      <table class="table">
        <thead>
          <tr>
            <th>Username</th>
            <th>Password</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, i) in honeypotData" :key="i">
            <td class="mono">{{ item.username }}</td>
            <td class="mono">{{ item.password }}</td>
            <td>{{ item.count }}</td>
          </tr>
          <tr v-if="!honeypotData.length">
            <td colspan="3" class="empty">No data</td>
          </tr>
        </tbody>
      </table>
      <Pagination :total="honeypotTotal" :page="honeypotPage" @change="loadHoneypot" />
    </div>

    <!-- Rejected IPs -->
    <div v-if="activeTab === 'rejected'" class="table-wrap">
      <table class="table">
        <thead>
          <tr>
            <th>IP</th>
            <th>Reject Count</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, i) in rejectedData" :key="i">
            <td class="mono">{{ item.ip }}</td>
            <td>{{ item.count }}</td>
          </tr>
          <tr v-if="!rejectedData.length">
            <td colspan="2" class="empty">No data</td>
          </tr>
        </tbody>
      </table>
      <Pagination :total="rejectedTotal" :page="rejectedPage" @change="loadRejected" />
    </div>

    <!-- Verified IPs -->
    <div v-if="activeTab === 'verified'" class="table-wrap">
      <table class="table">
        <thead>
          <tr>
            <th>IP</th>
            <th>Method</th>
            <th>Time</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, i) in verifiedData" :key="i">
            <td class="mono">{{ item.ip }}</td>
            <td>{{ item.method }}</td>
            <td>{{ item.created_at }}</td>
          </tr>
          <tr v-if="!verifiedData.length">
            <td colspan="3" class="empty">No data</td>
          </tr>
        </tbody>
      </table>
      <Pagination :total="verifiedTotal" :page="verifiedPage" @change="loadVerified" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import {
  getStats,
  getHoneypotCredentials,
  getRejectedIPs,
  getVerifiedIPs,
  type Stats,
  type CredentialStat,
  type IPStat,
  type VerifiedIP,
} from '../api'
import Pagination from '../components/Pagination.vue'

const tabs = [
  { key: 'honeypot', label: 'Honeypot Credentials' },
  { key: 'rejected', label: 'Rejected IPs' },
  { key: 'verified', label: 'Verified IPs' },
]

const activeTab = ref('honeypot')
const stats = ref<Stats | null>(null)

const honeypotData = ref<CredentialStat[]>([])
const honeypotTotal = ref(0)
const honeypotPage = ref(1)

const rejectedData = ref<IPStat[]>([])
const rejectedTotal = ref(0)
const rejectedPage = ref(1)

const verifiedData = ref<VerifiedIP[]>([])
const verifiedTotal = ref(0)
const verifiedPage = ref(1)

const loadStats = async () => {
  try {
    const res = await getStats()
    stats.value = res.data
  } catch { /* ignore */ }
}

const loadHoneypot = async (page = 1) => {
  honeypotPage.value = page
  try {
    const res = await getHoneypotCredentials(page)
    honeypotData.value = res.data.data || []
    honeypotTotal.value = res.data.total
  } catch { /* ignore */ }
}

const loadRejected = async (page = 1) => {
  rejectedPage.value = page
  try {
    const res = await getRejectedIPs(page)
    rejectedData.value = res.data.data || []
    rejectedTotal.value = res.data.total
  } catch { /* ignore */ }
}

const loadVerified = async (page = 1) => {
  verifiedPage.value = page
  try {
    const res = await getVerifiedIPs(page)
    verifiedData.value = res.data.data || []
    verifiedTotal.value = res.data.total
  } catch { /* ignore */ }
}

watch(activeTab, (tab) => {
  if (tab === 'honeypot') loadHoneypot()
  else if (tab === 'rejected') loadRejected()
  else if (tab === 'verified') loadVerified()
})

onMounted(() => {
  loadStats()
  loadHoneypot()
})
</script>

<style scoped>
.admin-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.header h1 {
  font-size: 22px;
  font-weight: 600;
  color: var(--accent);
}

.back-link {
  color: var(--text-dim);
  font-size: 13px;
  text-decoration: none;
}

.back-link:hover {
  color: var(--accent);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  margin-bottom: 28px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text);
}

.stat-label {
  display: block;
  font-size: 11px;
  color: var(--text-dim);
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 16px;
  border-bottom: 1px solid var(--border);
  padding-bottom: 0;
}

.tab {
  padding: 10px 18px;
  background: none;
  border: none;
  color: var(--text-dim);
  font-size: 13px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.tab:hover {
  color: var(--text);
}

.tab.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

.table-wrap {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th {
  text-align: left;
  padding: 12px 16px;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-dim);
  background: var(--bg);
  border-bottom: 1px solid var(--border);
}

.table td {
  padding: 10px 16px;
  font-size: 13px;
  border-bottom: 1px solid var(--border);
}

.table tr:last-child td {
  border-bottom: none;
}

.mono {
  font-family: 'JetBrains Mono', monospace;
}

.empty {
  text-align: center;
  color: var(--text-dim);
  padding: 24px 16px;
}
</style>
