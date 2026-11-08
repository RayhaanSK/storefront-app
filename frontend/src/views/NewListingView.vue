<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { createProduct } from '../services/storefront'
import ListingForm from '../components/ListingForm.vue'
import { useCategories } from '../composables/useCategories'
import { useSession } from '../composables/useSession'
import type { UpdateProductRequest } from '../models/storefront'
import '../styles/pages.css'

const router = useRouter()
const { currentUser } = useSession()
const { categories, categoriesError, loadCategories } = useCategories()

const isLoading = ref(true)
const isSubmitting = ref(false)
const errorMessage = ref('')

onMounted(async () => {
  await loadCategories()
  isLoading.value = false
})

async function addListing(values: UpdateProductRequest) {
  const user = currentUser.value
  if (!user) return

  isSubmitting.value = true
  errorMessage.value = ''

  try {
    await createProduct({ ...values, sellerId: user.id })
    await router.push('/my-listings')
  } catch (error) {
    console.error(error)

    errorMessage.value =
      error instanceof Error && error.message.toLowerCase() === 'title is required'
        ? 'The product name can’t be blank.'
        : 'Couldn’t add the listing. Check the details and try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="ui-page ui-page--form">
    <h1 class="ui-title">New Listing</h1>

    <p v-if="isLoading" class="ui-message" role="status">Loading…</p>

    <template v-else>
      <p v-if="categoriesError" class="ui-message ui-message--error" role="alert">
        {{ categoriesError }}
      </p>

      <ListingForm
        :categories="categories"
        submit-label="Add listing"
        submitting-label="Adding…"
        :is-submitting="isSubmitting"
        :error-message="errorMessage"
        @submit="addListing"
      />
    </template>
  </main>
</template>