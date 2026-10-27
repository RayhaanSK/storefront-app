<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { listAllUserTypes, signUpUserAccount } from '../services/storefront'
import type { UserType } from '../models/storefront'
import { hashPassword } from '../utils/passwordHash'
import '../styles/auth.css'

const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const selectedUserTypeId = ref<number | ''>('')

const userTypes = ref<UserType[]>([])
const errorMessage = ref('')
const isSubmitting = ref(false)

onMounted(async () => {
  try {
    // The backend can send null instead of [] when there are no rows.
    userTypes.value = (await listAllUserTypes()) ?? []
  } catch (error) {
    console.error(error)
    errorMessage.value =
      'Couldn’t load the account types. Check that the server is running, then reload the page.'
  }
})

function capitalise(text: string): string {
  return text.charAt(0).toUpperCase() + text.slice(1)
}

function signupErrorMessage(error: unknown): string {
  const message = error instanceof Error ? error.message.toLowerCase() : ''

  if (message === 'title is required') {
    return 'The username and password can’t be blank.'
  }

  return 'Couldn’t create the account. The username may already be taken, so try a different one.'
}

async function submit() {
  errorMessage.value = ''

  const trimmedUsername = username.value.trim()

  if (trimmedUsername === '' || password.value === '') {
    errorMessage.value = 'Enter a username and a password.'
    return
  }

  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'The passwords don’t match.'
    return
  }

  if (selectedUserTypeId.value === '') {
    errorMessage.value = 'Select an account type.'
    return
  }

  isSubmitting.value = true

  try {
    await signUpUserAccount({
      username: trimmedUsername,
      passwordHash: await hashPassword(password.value),
      userTypeId: selectedUserTypeId.value,
    })

    await router.push({ name: 'login', query: { registered: '1' } })
  } catch (error) {
    console.error(error)
    errorMessage.value = signupErrorMessage(error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <h1 class="auth-title">Sign up</h1>

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
          autocomplete="new-password"
        />
      </div>

      <div class="auth-field">
        <label for="confirm-password" class="auth-label">Confirm Password</label>
        <input
          id="confirm-password"
          v-model="confirmPassword"
          type="password"
          class="auth-input"
          autocomplete="new-password"
        />
      </div>

      <div class="auth-field">
        <label for="user-type" class="auth-label">Account type</label>
        <select id="user-type" v-model="selectedUserTypeId" class="auth-input">
          <option value="" disabled>Select an account type</option>
          <option v-for="userType in userTypes" :key="userType.id" :value="userType.id">
            {{ capitalise(userType.type) }}
          </option>
        </select>
      </div>

      <p v-if="errorMessage" class="auth-message auth-message--error" role="alert">
        {{ errorMessage }}
      </p>

      <button type="submit" class="auth-button" :disabled="isSubmitting">
        {{ isSubmitting ? 'Signing up…' : 'Sign up' }}
      </button>
    </form>
  </main>
</template>