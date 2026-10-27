<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { logIntoUserAccount } from '../services/storefront'
import { useSession } from '../composables/useSession'
import { hashPassword } from '../utils/passwordHash'
import '../styles/auth.css'

const route = useRoute()
const router = useRouter()
const { login } = useSession()

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const isSubmitting = ref(false)

// The signup page redirects here with ?registered=1 after an account is created.
const justRegistered = route.query.registered === '1'

// Messages the backend sends when the username or password is wrong.
const CREDENTIAL_ERRORS = ['user account not found', 'incorrect username or password']

function loginErrorMessage(error: unknown): string {
  const message = error instanceof Error ? error.message.toLowerCase().replace(/\.$/, '') : ''

  if (CREDENTIAL_ERRORS.includes(message)) {
    return 'Incorrect username or password.'
  }

  return 'Couldn’t log in. Check your details and try again.'
}

async function submit() {
  errorMessage.value = ''

  if (username.value.trim() === '' || password.value === '') {
    errorMessage.value = 'Enter your username and password.'
    return
  }

  isSubmitting.value = true

  try {
    const response = await logIntoUserAccount({
      username: username.value.trim(),
      passwordHash: await hashPassword(password.value),
    })

    if (!response.successStatus) {
      throw new Error('Login was not successful')
    }

    const account = response.userAccount

    login({
      id: account.id,
      username: account.username,
      userType: account.userType,
      userTypeId: account.userTypeId,
    })

    await router.push('/')
  } catch (error) {
    console.error(error)
    errorMessage.value = loginErrorMessage(error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <h1 class="auth-title">Login</h1>

    <form class="auth-form" @submit.prevent="submit">
      <div class="auth-field">
        <label for="username" class="auth-label">Username</label>
        <input
          id="username"
          v-model="username"
          type="text"
          class="auth-input"
          autocomplete="username"
          maxlength="30"
        />
      </div>

      <div class="auth-field">
        <label for="password" class="auth-label">Password</label>
        <input
          id="password"
          v-model="password"
          type="password"
          class="auth-input"
          autocomplete="current-password"
        />
      </div>

      <p v-if="justRegistered && !errorMessage" class="auth-message auth-message--success" role="status">
        Account created. Log in to continue.
      </p>
      <p v-if="errorMessage" class="auth-message auth-message--error" role="alert">
        {{ errorMessage }}
      </p>

      <button type="submit" class="auth-button" :disabled="isSubmitting">
        {{ isSubmitting ? 'Logging in…' : 'Login' }}
      </button>
    </form>

    <RouterLink to="/signup" class="auth-alt">Sign up new account</RouterLink>
  </main>
</template>