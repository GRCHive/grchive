<template>
    <v-form class="ma-4">
        <v-text-field
            label="Email"
            filled
        ></v-text-field>

        <v-btn
            color="success"
        >
            Next
        </v-btn>

        <p class="body-1 ma-2">
            Don't have an account?
            <a :href="getStartedUrl">Get started.</a>
        </p>
    </v-form>
</template>

<script lang="ts">

import { createGetStartedUrl } from '../../ts/url'

export default {
    computed: {
        getStartedUrl() : string {
            return createGetStartedUrl();
        },
    }
}

</script>
