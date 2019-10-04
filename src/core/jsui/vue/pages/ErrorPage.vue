<template>
    <v-container fluid class="max-height">

        <v-row align="center"
               justify="center"
               class="max-height"
        >
            <v-card width="60%">
                <v-card-title>
                    {{ errorTitle }}
                </v-card-title>

                <v-card-text>
                    {{ errorSubstring }} If you believe this to be in error, please
                    <a :href="contactUsUrl">contact us</a>.
                    Click <a onclick="window.history.back();">here</a> to go back to where you were before.
                </v-card-text>
            </v-card>

        </v-row>
    </v-container>
</template>

<script lang="ts">

import { contactUsUrl } from "../../ts/url"

export default {
    props: {
        errorTitle: String,
        errorSubstring: String
    },
    data : () => ({
        contactUsUrl
    }),
}

</script>
