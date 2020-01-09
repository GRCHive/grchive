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
}

</script>
