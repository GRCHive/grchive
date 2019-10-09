<template>
    <v-toolbar flat height="30px">
        <v-toolbar-items>
            <v-divider vertical></v-divider>
            <v-btn text color="accent">
                File
                <v-icon small color="accent">mdi-chevron-down</v-icon>
            </v-btn>
            <v-divider vertical></v-divider>
            <v-btn text color="accent">
                View
                <v-icon small color="accent">mdi-chevron-down</v-icon>
            </v-btn>
            <v-divider vertical></v-divider>
        </v-toolbar-items>
        <div class="flex-grow-1"></div>
        <v-toolbar-items>
            <v-divider vertical></v-divider>
            <v-btn text color="accent">
                Add Node
                <v-icon small color="accent">mdi-plus</v-icon>
            </v-btn>
            <v-divider vertical></v-divider>
        </v-toolbar-items>
    </v-toolbar>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
})

</script>
