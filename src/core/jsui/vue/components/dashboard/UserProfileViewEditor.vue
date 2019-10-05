<template>
    <v-form>
        <v-text-field v-model="firstName" label="First Name">
        </v-text-field>
        <v-text-field v-model="lastName" label="Last Name">
        </v-text-field>
        <v-text-field v-model="email" label="Email">
        </v-text-field>
    </v-form>

</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    data: function() {
        return {
            //@ts-ignore
            firstName: this.$root.userFirstName,
            //@ts-ignore
            lastName: this.$root.userLastName,
            //@ts-ignore
            email: this.$root.userEmail,
        }
    }
})

</script>
