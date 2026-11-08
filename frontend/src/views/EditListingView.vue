<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { getSingleProductDetails, updateProduct } from '../services/storefront'
import ListingForm from '../components/ListingForm.vue'
import { useCategories } from '../composables/useCategories'
import { useSession } from '../composables/useSession'
import type { UpdateProductRequest } from '../models/storefront'
import '../styles/pages.css'

const route = useRoute()
const router = useRouter()
const { currentUser } = useSession()
const { categories, categoriesError, loadCategories } = useCategories()

const productId = Number(route.params.id)

const initial = ref<{
  name: string
  price: number
  description: string
  categoryId: number | null
} | null>(null)

const isLoading = ref(true)
const loadError = ref('')
const isSubmitting = ref(false)
const errorMessage = ref('')

onMounted(async () => {
  const user = currentUser.value

  if (!user || !Number.isInteger(productId) || productId <= 0) {
    loadError.value = 'This listing doesn’t exist or has been removed.'
    isLoading.value = false
    return
  }

  try {
    const [product] = await Promise.all([getSingleProductDetails(productId), loadCategories()])

    // The URL can be typed by hand, so check the listing really belongs to this seller.
    if (product.sellerId !== user.id) {
      loadError.value = 'You can only edit your own listings.'
      return
    }

    // Match the category by id, falling back to its name in case the id field isn't
    // present in the product JSON (see the note about the Go `categoryID` tag).
    const category =
      categories.value.find((item) => item.id === product.categoryId) ??
      categories.value.find((item) => item.name === product.categoryName)

    initial.value = {
      name: product.name,
      price: product.price,
      description: product.description ?? '',
      categoryId: category?.id ?? null,
    }
  } catch (error) {
    console.error(error)

    const isNotFound = error instanceof Error && error.message.toLowerCase().includes('not found')
    loadError.value = isNotFound
      ? 'This listing doesn’t exist or has been removed.'
      : 'Couldn’t load the listing. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }
})

async function saveListing(values: UpdateProductRequest) {
  isSubmitting.value = true
  errorMessage.value = ''

  try {
    await updateProduct(productId, values)
    await router.push('/my-listings')
  } catch (error) {
    console.error(error)

    errorMessage.value =
      error instanceof Error && error.message.toLowerCase() === 'title is required'
        ? 'The product name can’t be blank.'
        : 'Couldn’t update the listing. Check the details and try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="ui-page ui-page--form">
    <h1 class="ui-title">Edit Listing</h1>

    <p v-if="isLoading" class="ui-message" role="status">Loading…</p>
    <p v-else-if="loadError" class="ui-message ui-message--error" role="alert">{{ loadError }}</p>

    <template v-else-if="initial">
      <p v-if="categoriesError" class="ui-message ui-message--error" role="alert">
        {{ categoriesError }}
      </p>

      <ListingForm
        :categories="categories"
        :initial="initial"
        submit-label="Update listing"
        submitting-label="Updating…"
        :is-submitting="isSubmitting"
        :error-message="errorMessage"
        @submit="saveListing"
      />
    </template>
  </main>
</template>