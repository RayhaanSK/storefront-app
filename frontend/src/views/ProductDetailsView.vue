<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import {
  addProductToCustomerCart,
  getSingleProductDetails,
  listAllCustomerCartItems,
} from "../services/storefront"
import { useSession } from '../composables/useSession'
import type { Product } from '../models/storefront'
import { formatDate, formatPrice } from '../utils/format'
import '../styles/pages.css'

const NOT_FOUND_MESSAGE = 'This product doesn’t exist or has been removed.'

const route = useRoute()
const { currentUser, role } = useSession()
const isCustomer = computed(() => role.value === 'customer')

const product = ref<Product | null>(null)
const isLoading = ref(true)
const errorMessage = ref('')

// Customer-only state
const isInCart = ref(false)
const isAdding = ref(false)
const cartError = ref('')

async function loadProduct() {
  // This watcher also fires while navigating away (the :id param disappears), so ignore that.
  if (route.name !== 'product-details') return

  isLoading.value = true
  errorMessage.value = ''
  cartError.value = ''
  product.value = null
  isInCart.value = false

  const id = Number(route.params.id)

  if (!Number.isInteger(id) || id <= 0) {
    errorMessage.value = NOT_FOUND_MESSAGE
    isLoading.value = false
    return
  }

  try {
    product.value = await getSingleProductDetails(id)
  } catch (error) {
    console.error(error)

    const isNotFound = error instanceof Error && error.message.toLowerCase().includes('not found')
    errorMessage.value = isNotFound
      ? NOT_FOUND_MESSAGE
      : 'Couldn’t load the product. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }

  if (product.value && isCustomer.value && currentUser.value) {
    try {
      const cartProducts = (await listAllCustomerCartItems(currentUser.value.id)) ?? []
      isInCart.value = cartProducts.some((cartProduct) => cartProduct.id === id)
    } catch (error) {
      // Not fatal: the page still works without knowing what's in the cart.
      console.error(error)
    }
  }
}

// Re-runs when the id in the URL changes, e.g. /product/1 -> /product/2
watch(() => route.params.id, loadProduct, { immediate: true })

async function addToCart() {
  const user = currentUser.value
  if (!user || !product.value || isAdding.value) return

  isAdding.value = true
  cartError.value = ''

  try {
    await addProductToCustomerCart({ customerId: user.id, productId: product.value.id })
    isInCart.value = true
  } catch (error) {
    console.error(error)
    cartError.value = 'Couldn’t add this product to your cart. Please try again.'
  } finally {
    isAdding.value = false
  }
}

const cartButtonLabel = computed(() => {
  if (isInCart.value) return 'In cart'
  if (isAdding.value) return 'Adding…'
  return 'Add to cart'
})
</script>

<template>
  <main class="ui-page">
    <p v-if="isLoading" class="ui-message" role="status">Loading product…</p>
    <p v-else-if="errorMessage" class="ui-message ui-message--error" role="alert">
      {{ errorMessage }}
    </p>

    <article v-else-if="product" class="product">
      <div class="product__header">
        <div class="product__heading">
          <h1 class="product__name">{{ product.name }}</h1>
          <p class="product__seller">Listed by: {{ product.sellerUsername }}</p>
        </div>

        <button
          v-if="isCustomer"
          type="button"
          class="ui-btn ui-btn--green"
          :disabled="isInCart || isAdding"
          @click="addToCart"
        >
          {{ cartButtonLabel }}
        </button>
      </div>

      <p v-if="cartError" class="ui-message ui-message--error" role="alert">{{ cartError }}</p>

      <p class="product__price">{{ formatPrice(product.price) }}</p>

      <h2 class="product__subheading">Description</h2>
      <p v-if="product.description" class="product__description">{{ product.description }}</p>
      <p v-else class="product__description product__description--empty">
        The seller didn’t add a description.
      </p>

      <p class="product__date">Date listed: {{ formatDate(product.listedAt) }}</p>
    </article>
  </main>
</template>

<style scoped>
.product__header {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.product__name {
  margin: 0;
  font-size: 2.25rem;
  font-weight: 700;
  line-height: 1.2;
  overflow-wrap: anywhere;
}

.product__seller {
  margin: 0.4rem 0 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #6b6b6b;
}

.product__price {
  margin: 2rem 0 0;
  font-size: 2rem;
  font-weight: 700;
}

.product__subheading {
  margin: 2.5rem 0 0.75rem;
  font-size: 1.6rem;
  font-weight: 700;
}

.product__description {
  margin: 0;
  max-width: 65ch;
  line-height: 1.5;
  /* Keeps the seller's own line breaks and paragraphs */
  white-space: pre-line;
  overflow-wrap: anywhere;
}

.product__description--empty {
  color: #6b6b6b;
}

.product__date {
  margin: 4rem 0 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #767676;
}
</style>