<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { BillsApiError, getBillShares, replaceBillShares } from "./api/bills"
import { allocateDraft } from "./allocation"
import ShareRow from "./components/ShareRow.vue"
import { formatPercentage, parsePercentage } from "./decimal"
import type { BillShares, DraftShare, RowErrors } from "./types"

const billID = 1
let nextDraftKey = 1

const serverAllocation = ref<BillShares | null>(null)
const draft = ref<DraftShare[]>([])
const loading = ref(true)
const saving = ref(false)
const loadError = ref("")
const saveError = ref("")
const successMessage = ref("")

const rowErrors = computed<Map<number, RowErrors>>(() => {
  const errors = new Map<number, RowErrors>()
  const nameCounts = new Map<string, number>()

  for (const row of draft.value) {
    const normalizedName = row.personName.trim().toLowerCase()
    if (normalizedName) {
      nameCounts.set(normalizedName, (nameCounts.get(normalizedName) ?? 0) + 1)
    }
  }

  for (const row of draft.value) {
    const rowError: RowErrors = {}
    const normalizedName = row.personName.trim().toLowerCase()
    if (!normalizedName) {
      rowError.personName = "Enter a person name."
    } else if ((nameCounts.get(normalizedName) ?? 0) > 1) {
      rowError.personName = "Person names must be unique."
    }
    if (parsePercentage(row.percentage) === null) {
      rowError.percentage = "Enter a value from 0.01 to 100.00, using up to two decimals."
    }
    errors.set(row.key, rowError)
  }
  return errors
})

const percentageTotal = computed(() =>
  draft.value.reduce(
    (total, row) => total + (parsePercentage(row.percentage) ?? 0),
    0,
  ),
)
const formattedPercentageTotal = computed(() => formatPercentage(percentageTotal.value))
const allocatedAmounts = computed(() =>
  serverAllocation.value
    ? allocateDraft(serverAllocation.value.bill.totalAmount, draft.value)
    : draft.value.map(() => null),
)
const hasRowErrors = computed(() =>
  [...rowErrors.value.values()].some(
    (errors) => Boolean(errors.personName) || Boolean(errors.percentage),
  ),
)
const canSave = computed(
  () =>
    !loading.value &&
    !saving.value &&
    draft.value.length > 0 &&
    !hasRowErrors.value &&
    percentageTotal.value === 10_000,
)

function draftFor(allocation: BillShares): DraftShare[] {
  return allocation.shares.map((share) => ({
    key: nextDraftKey++,
    personName: share.personName,
    percentage: share.percentage,
  }))
}

function messageFor(error: unknown): string {
  return error instanceof BillsApiError
    ? error.message
    : "Something went wrong. Please try again."
}

function clearSaveStatus() {
  saveError.value = ""
  successMessage.value = ""
}

async function load() {
  loading.value = true
  loadError.value = ""
  successMessage.value = ""
  try {
    const allocation = await getBillShares(billID)
    serverAllocation.value = allocation
    draft.value = draftFor(allocation)
  } catch (error) {
    loadError.value = messageFor(error)
  } finally {
    loading.value = false
  }
}

function updateRow(updatedRow: DraftShare) {
  const index = draft.value.findIndex((row) => row.key === updatedRow.key)
  if (index !== -1) {
    draft.value[index] = updatedRow
  }
  clearSaveStatus()
}

function addShare() {
  draft.value.push({ key: nextDraftKey++, personName: "", percentage: "" })
  clearSaveStatus()
}

function removeShare(key: number) {
  draft.value = draft.value.filter((row) => row.key !== key)
  clearSaveStatus()
}

async function save() {
  if (!canSave.value) {
    return
  }

  saving.value = true
  clearSaveStatus()
  try {
    const allocation = await replaceBillShares(
      billID,
      draft.value.map((row) => ({
        personName: row.personName.trim(),
        percentage: row.percentage,
      })),
    )
    serverAllocation.value = allocation
    draft.value = draftFor(allocation)
    successMessage.value = "Shares saved."
  } catch (error) {
    saveError.value = messageFor(error)
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <main class="page-shell">
    <header class="page-header">
      <p class="eyebrow">Bill splitter</p>
      <h1>{{ serverAllocation?.bill.description ?? "Bill shares" }}</h1>
      <p v-if="serverAllocation" class="bill-total">
        Bill total <strong>€{{ serverAllocation.bill.totalAmount }}</strong>
      </p>
    </header>

    <section v-if="loading" class="state-card" aria-live="polite">
      <p>Loading bill shares…</p>
    </section>

    <section v-else-if="loadError" class="state-card error-state" role="alert">
      <p>{{ loadError }}</p>
      <button type="button" @click="load">Retry</button>
    </section>

    <form v-else class="editor" @submit.prevent="save">
      <div class="editor-heading">
        <div>
          <h2>Shares</h2>
          <p>Allocate exactly 100.00% of the bill.</p>
        </div>
        <button type="button" class="secondary-button" :disabled="saving" @click="addShare">
          Add person
        </button>
      </div>

      <ul v-if="draft.length" class="share-list">
        <ShareRow
          v-for="(row, index) in draft"
          :key="row.key"
          :row="row"
          :row-number="index + 1"
          :amount="allocatedAmounts[index]"
          :errors="rowErrors.get(row.key) ?? {}"
          :disabled="saving"
          @update="updateRow"
          @remove="removeShare(row.key)"
        />
      </ul>
      <p v-else class="empty-state">Add at least one person to split this bill.</p>

      <div class="summary">
        <span>Percentage total</span>
        <strong :class="{ 'invalid-total': percentageTotal !== 10_000 }">
          {{ formattedPercentageTotal }}%
        </strong>
      </div>
      <p v-if="percentageTotal !== 10_000" class="total-hint">
        Percentages must total exactly 100.00% before saving.
      </p>

      <div class="status-region" aria-live="polite" aria-atomic="true">
        <p v-if="saveError" class="status-message error-message">{{ saveError }}</p>
        <p v-else-if="successMessage" class="status-message success-message">
          {{ successMessage }}
        </p>
      </div>

      <div class="form-actions">
        <button type="submit" class="primary-button" :disabled="!canSave">
          {{ saving ? "Saving…" : "Save shares" }}
        </button>
      </div>
    </form>
  </main>
</template>
