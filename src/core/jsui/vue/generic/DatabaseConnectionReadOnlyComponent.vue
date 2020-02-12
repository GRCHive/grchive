<template>

<div class="pa-4">
    <p>
        <span class="font-weight-bold">
            Host:
        </span>

        <span>
            {{ connection.Host }}
        </span>
    </p>
    
    <p>
        <span class="font-weight-bold">
            Port:
        </span>

        <span>
            {{ connection.Port }}
        </span>
    </p>

    <p>
        <span class="font-weight-bold">
            Database Name:
        </span>

        <span>
            {{ connection.DbName }}
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

    <span class="font-weight-bold">Extra Parameters</span>
    <v-divider></v-divider>
    <string-dict-form-component
        :value="connection.Parameters"
        readonly
    >
    </string-dict-form-component>
</div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { DatabaseConnection } from '../../ts/databases'
import StringDictFormComponent from './StringDictFormComponent.vue'

const VueComponent = Vue.extend({
    props: {
        connection: Object as () => DatabaseConnection
    }
})

@Component({
    components: {
        StringDictFormComponent
    }
})
export default class DatabaseConnectionReadOnlyComponent extends VueComponent {
}

</script>
