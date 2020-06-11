<template>
    <div id="table">
        <p class="subtitle-2 header header-text ma-0 pl-3">{{ table.TableName }}</p>
        <v-simple-table :height="height" fixed-header>
            <template v-slot:default>
                <thead class="header">
                    <tr>
                        <th class="header-text">Column</th>
                        <th class="header-text">Type</th>
                    </tr>
                </thead>
                <tbody>
                    <tr
                        v-for="(col, i) in table.Columns"
                        :key="i"
                    >
                        <td>{{ col.Name }}</td>
                        <td>{{ normalizeType(col.Type) }}</td>
                    </tr>
                </tbody>
            </template>
        </v-simple-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    DbTable,
} from '../../ts/sql'

const Props = Vue.extend({
    props: {
        table: {
            type: Object,
            default: () => null as DbTable | null
        },
        height: {
            type: Number,
            default: undefined
        }
    }
})

@Component
export default class DatabaseTableViewer extends Props {

    normalizeType(typ: string) : string{
        typ = typ.toUpperCase()
        if (typ == "TIMESTAMP WITH TIME ZONE") {
            typ = "TIMESTAMPTZ"
        } else if (typ == "CHARACTER VARYING") {
            typ = "VARCHAR"
        } else if (typ == "TIMESTAMP WITHOUT TIME ZONE") {
            typ = "TIMESTAMP"
        }
        return typ
    }
}

</script>

<style scoped>

#table {
    border: 1px solid black;
    overflow: auto;
}

.header {
    background-color: #1976d2;
}

.header-text {
    background-color: #1976d2 !important;
    color: white !important;
    height: 2.0em;
}

</style>
