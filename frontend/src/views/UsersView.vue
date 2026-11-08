<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { deleteUserAccount, listAllUserAccounts } from '../services/storefront'
import { useSession } from '../composables/useSession'
import type { UserAccount } from '../models/storefront'
import '../styles/pages.css'

const { currentUser } = useSession()

const users = ref<UserAccount[]>([])
const isLoading = ref(true)
const errorMessage = ref('')
const actionError = ref('')
const removingId = ref<number | null>(null)

onMounted(async () => {
  try {
    // The backend can send null instead of [] when there are no accounts.
    users.value = (await listAllUserAccounts()) ?? []
  } catch (error) {
    console.error(error)
    errorMessage.value =
      'Couldn’t load the users. Check that the server is running, then reload the page.'
  } finally {
    isLoading.value = false
  }
})

async function removeUser(user: UserAccount) {
  if (removingId.value !== null) return
  if (!window.confirm(`Remove the account “${user.username}”? This can’t be undone.`)) return

  removingId.value = user.id
  actionError.value = ''

  try {
    await deleteUserAccount(user.id)
    users.value = users.value.filter((item) => item.id !== user.id)
  } catch (error) {
    console.error(error)
    actionError.value = `Couldn’t remove “${user.username}”. The account may still have listings or cart items.`
  } finally {
    removingId.value = null
  }
}
</script>

<template>
  <main class="ui-page">
    <h1 class="ui-title">Users</h1>

    <p v-if="actionError" class="ui-message ui-message--error" role="alert">{{ actionError }}</p>

    <p v-if="isLoading" class="ui-message" role="status">Loading users…</p>
    <p v-else-if="errorMessage" class="ui-message ui-message--error" role="alert">
      {{ errorMessage }}
    </p>
    <p v-else-if="users.length === 0" class="ui-message">There are no user accounts.</p>

    <ul v-else class="ui-list">
      <li v-for="user in users" :key="user.id" class="ui-item">
        <div class="ui-item__info">
          <h2 class="ui-item__name">
            {{ user.username }}<span v-if="user.id === currentUser?.id"> (you)</span>
          </h2>
          <p class="ui-item__detail">{{ user.passwordHash }}</p>
          <p class="ui-item__detail">{{ user.userType }}</p>
        </div>

        <button
          type="button"
          class="ui-btn ui-btn--red"
          :disabled="removingId !== null || user.id === currentUser?.id"
          :title="user.id === currentUser?.id ? 'You can’t remove your own account' : undefined"
          @click="removeUser(user)"
        >
          Remove
        </button>
      </li>
    </ul>
  </main>
</template>