<template>
    <v-snackbar 
        v-model="showSnack"
        :color="snackIsError ? 'error' : 'success'"
        bottom
        :timeout="timeout"
    >
        {{ snackText }} 

        <v-btn
            color="primary"
            v-if="snackShowPrimaryButton"
            :href="snackPrimaryUrl"
        >
            {{ snackPrimaryButtonText }}
        </v-btn>

        <v-btn
            color="secondary"
            @click="showSnack = false"
        >
            Close
        </v-btn>
    </v-snackbar>
</template>

<script lang="ts">

import Vue from 'vue';

export default Vue.extend({
    data: () => ({
        // Snackbar options for giving visual feedback to the user.
        showSnack: false,
        snackText: "",
        snackShowPrimaryButton: false,
        snackPrimaryButtonText: "",
        snackPrimaryUrl: "",
        snackIsError: false,
        timeout: 10000,
    }),
    methods: {
        showSnackBar(
            text: string,
            showPrimaryButton: boolean,
            primaryButtonText: string,
            primaryUrl: string,
            isError: boolean) {

            this.$data.showSnack = true
            this.$data.snackText = text
            this.$data.snackShowPrimaryButton = showPrimaryButton
            this.$data.snackPrimaryButtonText = primaryButtonText
            this.$data.snackPrimaryUrl = primaryUrl
            this.$data.snackIsError = isError
        }
    }
});

</script>
