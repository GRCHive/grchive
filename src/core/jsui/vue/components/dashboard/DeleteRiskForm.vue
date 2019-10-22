<template>

<v-card>
    <v-card-title>
        Are you sure you wish to delete these ({{ risksToDelete.length }}) risks?
    </v-card-title>
    <v-divider></v-divider>

    <section class="ma-2">
        <v-list>
            <v-list-item v-for="(item, index) in risksToDelete"
                         :key="index">
                <v-list-item-content>
                    <v-list-item-title>{{index+1}}.&nbsp;{{item.Name}}</v-list-item-title>
                </v-list-item-content>
            </v-list-item>
        </v-list>

        <v-checkbox v-model="globalDelete"
                    color="error"
                    class="subtitle-1"
                    label="Global Deletion?"
        >
        </v-checkbox>
        <p class="subtitle-2 red--text" v-if="globalDelete">
            <v-icon>mdi-alert</v-icon>Global deletion will remove these risks from all process flow nodes.
        </p>
    </section>
    <v-divider></v-divider>
    <v-card-actions>
        <v-btn color="primary" @click="onCancel">
            Cancel
        </v-btn>
        <div class="flex-grow-1"></div>
        <v-btn color="error" @click="onDelete">
            Delete
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    props: {
        risksToDelete : Array
    },
    data : () => ({
        globalDelete: false,
    }),
    methods: {
        onCancel() {
            this.$emit('do-cancel')
        },
        onDelete() {
            this.$emit('do-delete', this.globalDelete)
        }
    }
})

</script>
