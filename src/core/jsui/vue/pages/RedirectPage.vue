<template>
    <v-content>
        <v-row
            align="center"
            justify="center"
            class="max-height"
        >

            <v-card width="50%" class="ma-4">
                <v-card-title>
                    {{ title }}
                </v-card-title>
                <v-card-text class="long-text">
                    {{ subtext }}. Redirecting in {{ seconds }} seconds...
                    Click <a :href="redirectUrl">here</a> if you are not redirected automatically.
                </v-card-text>
            </v-card>
        </v-row>

    </v-content>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    props: {
        title: String,
        subtext: String,
        redirectUrl: String,
    },
    data : () => ({
        seconds: 5
    }),
    mounted() {
        window.setInterval(() => {
            this.$data.seconds = this.$data.seconds - 1
            if (this.$data.seconds <= 0) {
                window.location.replace(this.$props.redirectUrl);
            }
        }, 1000);
    }
})
</script>
