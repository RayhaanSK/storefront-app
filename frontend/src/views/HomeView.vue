<script setup lang="ts">
import { computed } from 'vue'

import { useSession, type Role } from '../composables/useSession'

const { currentUser, isLoggedIn, role } = useSession()

const roleMessages: Record<Role, string> = {
  customer: 'You are logged in as a customer, so try adding items to your cart from the catalogue!',
  seller:
    'You are logged in as a seller, so try managing your listings from the my listings page!',
  admin: 'You are logged in as an admin, so try managing users from the users page!',
}

// Falls back to a plain message if the account has a user type this page doesn't know about.
const roleMessage = computed(() => {
  if (role.value) return roleMessages[role.value]
  return `You are logged in as ${currentUser.value?.username ?? 'a user'}.`
})
</script>

<template>
  <main class="home">
    <h1 class="home__title">Online Storefront Application</h1>

    <div class="home__text">
      <p>
        This web application allows you to view listings, manage your shopping cart, manage your
        own listings, and manage users.
      </p>

      <p v-if="!isLoggedIn">
        You can begin by browsing the catalogue, but login or sign up to access more specific
        features.
      </p>

      <template v-else>
        <p>{{ roleMessage }}</p>
        <p>Login to a different user type to access other features.</p>
      </template>
    </div>
  </main>
</template>

<style scoped>
.home {
  padding: clamp(3rem, 18vh, 12rem) 1.5rem 3rem;
  text-align: center;
}

.home__title {
  margin: 0;
  font-size: clamp(1.75rem, 4vw, 2.75rem);
  font-weight: 700;
  line-height: 1.2;
  text-transform: uppercase;
}

.home__text {
  max-width: 68rem;
  margin: clamp(2.5rem, 12vh, 5.5rem) auto 0;
  font-size: clamp(1.1rem, 2.2vw, 1.75rem);
  line-height: 1.25;
}

.home__text p {
  margin: 0 0 1.25em;
}

.home__text p:last-child {
  margin-bottom: 0;
}
</style>