<script setup lang="ts">
import type { DraftShare, RowErrors } from "../types"

const props = defineProps<{
  row: DraftShare
  rowNumber: number
  amount: string | null
  errors: RowErrors
  disabled: boolean
}>()

const emit = defineEmits<{
  update: [row: DraftShare]
  remove: []
}>()

function update(field: "personName" | "percentage", value: string) {
  emit("update", { ...props.row, [field]: value })
}
</script>

<template>
  <li class="share-list-item">
    <div class="share-row" role="group" :aria-label="`Share ${rowNumber}`">
    <div class="field name-field">
      <label :for="`person-name-${row.key}`">Person name</label>
      <input
        :id="`person-name-${row.key}`"
        :value="row.personName"
        type="text"
        autocomplete="off"
        :disabled="disabled"
        :aria-invalid="Boolean(errors.personName)"
        :aria-describedby="errors.personName ? `person-name-error-${row.key}` : undefined"
        @input="update('personName', ($event.target as HTMLInputElement).value)"
      />
      <p
        v-if="errors.personName"
        :id="`person-name-error-${row.key}`"
        class="field-error"
      >
        {{ errors.personName }}
      </p>
    </div>

    <div class="field percentage-field">
      <label :for="`percentage-${row.key}`">Percentage</label>
      <div class="percentage-input">
        <input
          :id="`percentage-${row.key}`"
          :value="row.percentage"
          type="text"
          inputmode="decimal"
          :disabled="disabled"
          :aria-invalid="Boolean(errors.percentage)"
          :aria-describedby="errors.percentage ? `percentage-error-${row.key}` : undefined"
          @input="update('percentage', ($event.target as HTMLInputElement).value)"
        />
        <span aria-hidden="true">%</span>
      </div>
      <p
        v-if="errors.percentage"
        :id="`percentage-error-${row.key}`"
        class="field-error"
      >
        {{ errors.percentage }}
      </p>
    </div>

    <div class="amount-field">
      <span class="field-label">Allocated amount</span>
      <output :aria-label="`Allocated amount for ${row.personName || `share ${rowNumber}`}`">
        {{ amount === null ? "—" : `€${amount}` }}
      </output>
    </div>

    <button
      class="remove-button"
      type="button"
      :disabled="disabled"
      :aria-label="`Remove ${row.personName || `share ${rowNumber}`}`"
      @click="emit('remove')"
    >
      Remove
    </button>
    </div>
  </li>
</template>
