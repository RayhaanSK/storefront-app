<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import {
  addProductToCustomerCart,
  listAllCustomerCartItems,
  listProductCatalogue,
} from '../services/storefront'
import { useSession } from '../composables/useSession'
import type { Product } from '../models/storefront'

type SortOption = 'default' | 'price-asc' | 'price-desc' | 'name-asc'

const { currentUser, role } = useSession()
const isCustomer = computed(() => role.value === 'customer')

const products = ref<Product[]>([])
const isLoading = ref(true)
const errorMessage = ref('')

// The search runs when the user clicks search (or presses Enter), so the text box
// is kept separate from the query that is actually applied.
const searchDraft = ref('')
const searchQuery = ref('')

const sortOption = ref<SortOption>('default')
const selectedCategory = ref('') // '' means "All"

// Customer-only state
const cartProductIds = ref<Set<number>>(new Set())
const addingProductId = ref<number | null>(null)
const cartError = ref('')

async function loadProducts() {
  try {
    // The backend can send null instead of [] when there are no products.
    products.value = (await listProductCatalogue()) ?? []
  } catch (error) {
    console.error(error)
    errorMessage.value =
      'Couldn’t load the products. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }
}

// Used to show "In cart" on products the customer has already added.
async function loadCartProductIds() {
  if (!isCustomer.value || !currentUser.value) return

  try {
    const cartProducts = (await listAllCustomerCartItems(currentUser.value.id)) ?? []
    cartProductIds.value = new Set(cartProducts.map((product) => product.id))
  } catch (error) {
    // Not fatal: the catalogue still works without it.
    console.error(error)
  }
}

onMounted(() => {
  loadProducts()
  loadCartProductIds()
})

// There is no "list categories" endpoint, so the dropdown is built from the products.
const categories = computed(() =>
  [...new Set(products.value.map((product) => product.categoryName))].sort((a, b) =>
    a.localeCompare(b),
  ),
)

const visibleProducts = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()

  // filter() returns a new array, so sorting it does not touch products.value
  const matching = products.value.filter(
    (product) =>
      (selectedCategory.value === '' || product.categoryName === selectedCategory.value) &&
      product.name.toLowerCase().includes(query),
  )

  switch (sortOption.value) {
    case 'price-asc':
      return matching.sort((a, b) => a.price - b.price)
    case 'price-desc':
      return matching.sort((a, b) => b.price - a.price)
    case 'name-asc':
      return matching.sort((a, b) => a.name.localeCompare(b.name))
    default:
      return matching
  }
})

function submitSearch() {
  searchQuery.value = searchDraft.value
}

function formatPrice(price: number): string {
  return `$${price.toFixed(2)}`
}

async function addToCart(product: Product) {
  const user = currentUser.value
  if (!user || addingProductId.value !== null) return

  addingProductId.value = product.id
  cartError.value = ''

  try {
    await addProductToCustomerCart({ customerId: user.id, productId: product.id })
    cartProductIds.value = new Set(cartProductIds.value).add(product.id)
  } catch (error) {
    console.error(error)
    cartError.value = `Couldn’t add “${product.name}” to your cart. Please try again.`
  } finally {
    addingProductId.value = null
  }
}

function cartButtonLabel(product: Product): string {
  if (cartProductIds.value.has(product.id)) return 'In cart'
  if (addingProductId.value === product.id) return 'Adding…'
  return 'Add to cart'
}
</script>

<template>
  <div class="search-bar">
    <form class="search-bar__form" role="search" @submit.prevent="submitSearch">
      <input
        v-model="searchDraft"
        type="search"
        class="search-bar__input"
        placeholder="Search"
        aria-label="Search products by name"
      />
      <button type="submit" class="search-bar__button" aria-label="Search">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor"
          stroke-width="2" stroke-linecap="round" aria-hidden="true">
          <circle cx="11" cy="11" r="7" />
          <line x1="16.5" y1="16.5" x2="21" y2="21" />
        </svg>
      </button>
    </form>
  </div>

  <main class="catalogue">
    <h1 class="catalogue__title">Catalogue</h1>

    <div class="catalogue__filters">
      <div class="field">
        <label for="sort" class="field__label">Sort</label>
        <select id="sort" v-model="sortOption" class="field__select field__select--sort">
          <option value="default">Default</option>
          <option value="price-asc">Price - Low to High</option>
          <option value="price-desc">Price - High to Low</option>
          <option value="name-asc">Name - A to Z</option>
        </select>
      </div>

      <div class="field">
        <label for="category" class="field__label">Category</label>
        <select id="category" v-model="selectedCategory" class="field__select">
          <option value="">All</option>
          <option v-for="category in categories" :key="category" :value="category">
            {{ category }}
          </option>
        </select>
      </div>
    </div>

    <p v-if="cartError" class="catalogue__message catalogue__message--error" role="alert">
      {{ cartError }}
    </p>

    <p v-if="isLoading" class="catalogue__message" role="status">Loading products…</p>
    <p v-else-if="errorMessage" class="catalogue__message" role="alert">{{ errorMessage }}</p>
    <p v-else-if="products.length === 0" class="catalogue__message">
      No products are listed yet.
    </p>
    <p v-else-if="visibleProducts.length === 0" class="catalogue__message">
      No products match your search. Try a different name or category.
    </p>

    <ul v-else class="product-list">
      <li v-for="product in visibleProducts" :key="product.id" class="product">
        <div class="product__info">
          <h2 class="product__name">{{ product.name }}</h2>
          <p class="product__price">{{ formatPrice(product.price) }}</p>
        </div>

        <div class="product__actions">
          <button
            v-if="isCustomer"
            type="button"
            class="button button--cart"
            :disabled="cartProductIds.has(product.id) || addingProductId !== null"
            @click="addToCart(product)"
          >
            {{ cartButtonLabel(product) }}
          </button>

          <RouterLink :to="`/product/${product.id}`" class="button">Details</RouterLink>
        </div>
      </li>
    </ul>
  </main>
</template>

<style scoped>
/* Search strip */
.search-bar {
  height: 60px;
  background: #cccccc;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 1rem;
}

.search-bar__form {
  position: relative;
  width: 100%;
  max-width: 31rem;
}

.search-bar__input {
  width: 100%;
  height: 2.5rem;
  padding: 0 3rem 0 1.25rem;
  border: 1px solid #b5b5b5;
  border-radius: 999px;
  background: #ffffff;
  font: inherit;
  font-size: 1rem;
}

/* Remove the browser's built-in clear "x" so only our search icon shows */
.search-bar__input::-webkit-search-cancel-button {
  display: none;
}

.search-bar__button {
  position: absolute;
  top: 0;
  right: 0.25rem;
  height: 100%;
  width: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: transparent;
  color: #000000;
  cursor: pointer;
  border-radius: 999px;
}

.search-bar__input:focus-visible,
.search-bar__button:focus-visible {
  outline: 2px solid #2b2b2b;
  outline-offset: 2px;
}

/* Page body */
.catalogue {
  max-width: 46rem;
  margin: 0 auto;
  padding: 2.25rem 1.5rem 4rem;
}

.catalogue__title {
  margin: 0 0 1.5rem;
  font-size: 2rem;
  font-weight: 700;
  line-height: 1.2;
}

.catalogue__filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem 2.5rem;
  margin-bottom: 1.5rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field__label {
  font-size: 1rem;
}

.field__select {
  min-width: 9rem;
  padding: 0.6rem 0.8rem;
  border: 1px solid #cfcfcf;
  border-radius: 0.6rem;
  background: #ffffff;
  font: inherit;
  font-size: 1rem;
  cursor: pointer;
}

.field__select--sort {
  min-width: 13rem;
}

.field__select:focus-visible {
  outline: 2px solid #2b2b2b;
  outline-offset: 2px;
}

.catalogue__message {
  margin: 2rem 0 0;
  font-size: 1.1rem;
  color: #444444;
}

.catalogue__message--error {
  color: #b00020;
}

/* Product list */
.product-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.product {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.5rem 0;
}

.product__info {
  flex: 1 1 12rem;
  min-width: 0;
}

.product__name {
  margin: 0;
  font-size: 1.3rem;
  font-weight: 700;
  color: #222222;
  overflow-wrap: anywhere;
}

.product__price {
  margin: 0.15rem 0 0;
  font-size: 1.1rem;
  color: #555555;
}

.product__actions {
  display: flex;
  gap: 0.5rem;
}

/* Shared by the Details link and the Add to cart button */
.button {
  display: inline-block;
  flex-shrink: 0;
  padding: 0.65rem 1.1rem;
  border: 1.5px solid #2b2b2b;
  border-radius: 0.5rem;
  background: #2b2b2b;
  color: #ffffff;
  font: inherit;
  font-size: 1rem;
  line-height: 1.2;
  text-decoration: none;
  cursor: pointer;
}

.button:hover {
  background: #000000;
  border-color: #000000;
}

.button--cart {
  border-color: #1a9b4b;
  background: #ccf5d3;
  color: #12331d;
}

.button--cart:hover:not(:disabled) {
  border-color: #1a9b4b;
  background: #b4eebf;
}

.button--cart:disabled {
  border-color: #bdbdbd;
  background: #ececec;
  color: #555555;
  cursor: default;
}

.button:focus-visible {
  outline: 2px solid #2b2b2b;
  outline-offset: 2px;
}
</style>