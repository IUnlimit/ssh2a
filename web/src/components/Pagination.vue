<template>
  <div v-if="totalPages > 1" class="pagination">
    <button
      class="page-btn"
      :disabled="props.page <= 1"
      @click="emit('change', props.page - 1)"
    >
      Prev
    </button>
    <span class="page-info">{{ props.page }} / {{ totalPages }}</span>
    <button
      class="page-btn"
      :disabled="props.page >= totalPages"
      @click="emit('change', props.page + 1)"
    >
      Next
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  total: number
  page: number
  pageSize?: number
}>()

const emit = defineEmits<{
  change: [page: number]
}>()

const size = computed(() => props.pageSize || 20)
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / size.value)))
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px;
  border-top: 1px solid var(--border);
}

.page-btn {
  padding: 6px 14px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 12px;
  color: var(--text-dim);
}
</style>
