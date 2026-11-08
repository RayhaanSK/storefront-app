<script setup lang="ts">
import { computed, ref } from 'vue'

import type { Category, UpdateProductRequest } from '../models/storefront'
import '../styles/pages.css'

// Shared by the New Listing and Edit Listing pages.
// The form is created after its data has loaded, so `initial` is already final here.
const props = defineProps<{
  categories: Category[]
  submitLabel: string
  submittingLabel: string
  isSubmitting: boolean
  errorMessage: string
  initial?: {
    name: string
    price: number
    description: string
    categoryId: number | null
  }
}>()

const emit = defineEmits<{
  submit: [values: UpdateProductRequest]
}>()

const name = ref(props.initial?.name ?? '')
// Vue turns the value of a type="number" input into a number (or '' when empty),
// so this holds either, and is converted with String() before checking it.
const priceText = ref<string | number>(props.initial ? props.initial.price : '')
const description = ref(props.initial?.description ?? '')
const categoryId = ref<number | ''>(props.initial?.categoryId ?? '')

const validationError = ref('')
const shownError = computed(() => validationError.value || props.errorMessage)

function submit() {
  validationError.value = ''

  const trimmedName = name.value.trim()
  const priceInput = String(priceText.value).trim()
  const price = Number(priceInput)

  if (trimmedName === '') {
    validationError.value = 'Enter a product name.'
    return
  }

  if (priceInput === '' || !Number.isFinite(price)) {
    validationError.value = 'Enter a price, for example 19.99.'
    return
  }

  if (price < 0) {
    validationError.value = 'The price can’t be negative.'
    return
  }

  if (categoryId.value === '') {
    validationError.value = 'Select a category.'
    return
  }

  emit('submit', {
    name: trimmedName,
    price: Math.round(price * 100) / 100, // whole cents
    description: description.value,
    categoryId: categoryId.value,
  })
}
</script>

<template>
  <form class="ui-form" @submit.prevent="submit">
    <div class="ui-field">
      <label for="listing-name" class="ui-label">Product Name</label>
      <input id="listing-name" v-model="name" type="text" class="ui-input" maxlength="50" />
    </div>

    <div class="ui-field">
      <label for="listing-price" class="ui-label">Price ($)</label>
      <input
        id="listing-price"
        v-model="priceText"
        type="number"
        class="ui-input"
        min="0"
        step="0.01"
        inputmode="decimal"
      />
    </div>

    <div class="ui-field">
      <label for="listing-category" class="ui-label">Category</label>
      <select id="listing-category" v-model="categoryId" class="ui-input">
        <option value="" disabled>Select a category</option>
        <option v-for="category in categories" :key="category.id" :value="category.id">
          {{ category.name }}
        </option>
      </select>
    </div>

    <div class="ui-field">
      <label for="listing-description" class="ui-label">Description</label>
      <textarea
        id="listing-description"
        v-model="description"
        class="ui-input ui-textarea"
      ></textarea>
    </div>

    <p v-if="shownError" class="ui-message ui-message--error form-message" role="alert">
      {{ shownError }}
    </p>

    <div>
      <button type="submit" class="ui-btn" :disabled="isSubmitting">
        {{ isSubmitting ? submittingLabel : submitLabel }}
      </button>
    </div>
  </form>
</template>

<style scoped>
/* The shared message style has a top margin meant for standalone messages */
.form-message {
  margin: 0;
  font-size: 0.95rem;
}
</style>