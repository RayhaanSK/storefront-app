<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'

import { useSession, type Role } from './composables/useSession'

const router = useRouter()
const { isLoggedIn, role, logout } = useSession()

interface NavLink {
  label: string
  to: string
}

// The one extra link each type of user gets, between Catalogue and Logout.
const roleLinks: Record<Role, NavLink> = {
  customer: { label: 'Cart', to: '/cart' },
  seller: { label: 'My Listings', to: '/my-listings' },
  admin: { label: 'Users', to: '/users' },
}

const links = computed<NavLink[]>(() => {
  const shared: NavLink[] = [
    { label: 'Home', to: '/' },
    { label: 'Catalogue', to: '/catalogue' },
  ]

  if (!isLoggedIn.value) {
    return [...shared, { label: 'Login', to: '/login' }]
  }

  return role.value ? [...shared, roleLinks[role.value]] : shared
})

async function handleLogout() {
  logout()
  await router.push('/')
}
</script>

<template>
  <header class="navbar">
    <nav class="navbar__links" aria-label="Main">
      <RouterLink v-for="link in links" :key="link.to" :to="link.to" class="navbar__link">
        {{ link.label }}
      </RouterLink>

      <button v-if="isLoggedIn" type="button" class="navbar__link" @click="handleLogout">
        Logout
      </button>
    </nav>
  </header>

  <!-- The page for the current URL is rendered here -->
  <RouterView />
</template>

<style>
/* Global reset (unscoped so it can reach <body>) */
*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  background: #ffffff;
  color: #000000;
  font-family: 'Inter', system-ui, -apple-system, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}
</style>

<style scoped>
.navbar {
  height: 60px;
  background: #5a5a5a;
  display: flex;
  align-items: center;
  justify-content: center;
}

.navbar__links {
  display: flex;
  gap: clamp(1.25rem, 4.4vw, 3.5rem);
}

/* Applies to both the <a> links and the Logout <button>, so they look identical. */
.navbar__link {
  padding: 0.25rem 0;
  border: 0;
  border-bottom: 2px solid transparent;
  background: none;
  color: #ffffff;
  font: inherit;
  font-size: clamp(1.05rem, 2vw, 1.625rem);
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
}

.navbar__link:hover,
.navbar__link.router-link-exact-active {
  border-bottom-color: #ffffff;
}

.navbar__link:focus-visible {
  outline: 2px solid #ffffff;
  outline-offset: 4px;
}
</style>