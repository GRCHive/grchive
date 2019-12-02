<template>

<div class="pa-4">
    <p>
        <span class="font-weight-bold">
            Connection String:
        </span>

        <span>
            {{ connection.ConnString }}
        </span>
    </p>

    <p>
        <span class="font-weight-bold">
            Username:
        </span>

        <span>
            {{ connection.Username }}
        </span>
    </p>

    <p>
        <span class="font-weight-bold">
            Password:
        </span>

        <input :type="showPassword ? `text` : `password`" :value="connection.Password" readonly>

        <v-btn icon @click="showPassword = !showPassword">
            <v-icon>
                {{ showPassword ? 'mdi-eye-off' : 'mdi-eye' }}
            </v-icon>
        </v-btn>
    </p>

</div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { DatabaseConnection } from '../../ts/databases'

const VueComponent = Vue.extend({
    props: {
        connection: Object as () => DatabaseConnection
    }
})

@Component
export default class DatabaseConnectionReadOnlyComponent extends VueComponent {
    showPassword: boolean = false
}

</script>
