<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { deleteProduct, listAllSellerProducts } from '../services/storefront'
import { useSession } from '../composables/useSession'
import type { Product } from '../models/storefront'
import { formatPrice } from '../utils/format'
import '../styles/pages.css'

const { currentUser } = useSession()

const listings = ref<Product[]>([])
const isLoading = ref(true)
const errorMessage = ref('')
const actionError = ref('')
const removingId = ref<number | null>(null)

onMounted(async () => {
  const user = currentUser.value
  if (!user) return

  try {
    // The backend can send null instead of [] when the seller has no listings.
    listings.value = (await listAllSellerProducts(user.id)) ?? []
  } catch (error) {
    console.error(error)
    errorMessage.value =
      'Couldn’t load your listings. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }
})

async function removeListing(listing: Product) {
  if (removingId.value !== null) return
  if (!window.confirm(`Remove “${listing.name}”? This can’t be undone.`)) return

  removingId.value = listing.id
  actionError.value = ''

  try {
    await deleteProduct(listing.id)
    listings.value = listings.value.filter((item) => item.id !== listing.id)
  } catch (error) {
    console.error(error)
    actionError.value = `Couldn’t remove “${listing.name}”. It may be in a customer’s cart.`
  } finally {
    removingId.value = null
  }
}
</script>

<template>
  <main class="ui-page">
    <div class="ui-heading-row">
      <h1 class="ui-title">My Listings</h1>
      <RouterLink to="/my-listings/new" class="ui-btn ui-btn--green">New listing</RouterLink>
    </div>

    <p v-if="actionError" class="ui-message ui-message--error" role="alert">{{ actionError }}</p>

    <p v-if="isLoading" class="ui-message" role="status">Loading your listings…</p>
    <p v-else-if="errorMessage" class="ui-message ui-message--error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-else-if="listings.length === 0" class="ui-message">
      You haven’t listed anything yet. Select “New listing” to add your first product.
    </p>

    <ul v-else class="ui-list">
      <li v-for="listing in listings" :key="listing.id" class="ui-item">
        <div class="ui-item__info">
          <h2 class="ui-item__name">{{ listing.name }}</h2>
          <p class="ui-item__detail">{{ formatPrice(listing.price) }}</p>
          <p class="ui-item__detail">{{ listing.sold ? 'Sold' : 'Not sold' }}</p>
        </div>

        <div class="ui-actions">
          <button
            type="button"
            class="ui-btn ui-btn--red"
            :disabled="removingId !== null"
            @click="removeListing(listing)"
          >
            Remove
          </button>
          <RouterLink :to="`/my-listings/${listing.id}/edit`" class="ui-btn">Edit</RouterLink>
        </div>
      </li>
    </ul>
  </main>
</template>