<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import {
  checkoutCustomerCart,
  listAllCustomerCartItems,
  removeProductFromCustomerCart,
} from '../services/storefront'
import { useSession } from '../composables/useSession'
import type { Product } from '../models/storefront'
import { formatPrice } from '../utils/format'
import '../styles/pages.css'

const { currentUser } = useSession()

const items = ref<Product[]>([])
const isLoading = ref(true)
const errorMessage = ref('')
const actionError = ref('')
const checkoutComplete = ref(false)

// While a request is running, the cart buttons are disabled so requests can't overlap.
const isBusy = ref(false)

onMounted(async () => {
  const user = currentUser.value
  if (!user) return

  try {
    // The backend can send null instead of [] when the cart is empty.
    items.value = (await listAllCustomerCartItems(user.id)) ?? []
  } catch (error) {
    console.error(error)
    errorMessage.value =
      'Couldn’t load your cart. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }
})

async function removeItem(product: Product) {
  const user = currentUser.value
  if (!user || isBusy.value) return

  isBusy.value = true
  actionError.value = ''

  try {
    await removeProductFromCustomerCart(user.id, product.id)
    items.value = items.value.filter((item) => item.id !== product.id)
  } catch (error) {
    console.error(error)
    actionError.value = `Couldn’t remove “${product.name}” from your cart. Please try again.`
  } finally {
    isBusy.value = false
  }
}

async function checkout() {
  const user = currentUser.value
  if (!user || isBusy.value) return

  isBusy.value = true
  actionError.value = ''

  try {
    await checkoutCustomerCart(user.id)
    items.value = []
    checkoutComplete.value = true
  } catch (error) {
    console.error(error)
    actionError.value = 'Couldn’t check out. Please try again.'
  } finally {
    isBusy.value = false
  }
}
</script>

<template>
  <main class="ui-page">
    <h1 class="ui-title">Shopping Cart</h1>

    <p v-if="actionError" class="ui-message ui-message--error" role="alert">{{ actionError }}</p>

    <p v-if="isLoading" class="ui-message" role="status">Loading your cart…</p>
    <p v-else-if="errorMessage" class="ui-message ui-message--error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-else-if="checkoutComplete" class="ui-message ui-message--success" role="status">
      Checkout complete. Your cart is now empty.
    </p>
    <p v-else-if="items.length === 0" class="ui-message">
      Your cart is empty. Browse the <RouterLink to="/catalogue">catalogue</RouterLink> to add
      items.
    </p>

    <template v-else>
      <ul class="ui-list">
        <li v-for="item in items" :key="item.id" class="ui-item">
          <div class="ui-item__info">
            <h2 class="ui-item__name">{{ item.name }}</h2>
            <p class="ui-item__detail">{{ formatPrice(item.price) }}</p>
          </div>

          <div class="ui-actions">
            <button
              type="button"
              class="ui-btn ui-btn--red"
              :disabled="isBusy"
              @click="removeItem(item)"
            >
              Remove
            </button>
            <RouterLink :to="`/product/${item.id}`" class="ui-btn">Details</RouterLink>
          </div>
        </li>
      </ul>

      <div class="ui-actions ui-actions--end">
        <button type="button" class="ui-btn ui-btn--green" :disabled="isBusy" @click="checkout">
          Checkout
        </button>
      </div>
    </template>
  </main>
</template>