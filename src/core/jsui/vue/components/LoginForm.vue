<template>
    <v-form class="ma-4" v-model="formValid" ref="form">
        <v-text-field
            v-model="email"
            label="Email"
            type="email"
            filled
            :rules="[rules.required, rules.createMaxLength(320), rules.email]"
            required
        ></v-text-field>

        <v-btn
            color="success"
            class="my-2"
            :disabled="!formValid || !email"
            @click="submit"
        >
            Next
        </v-btn>

        <p class="body-1 my-2">
            Don't have an account?
            <a :href="getStartedUrl">Get started.</a>
        </p>
    </v-form>
</template>

<script lang="ts">

import { getStartedUrl } from '../../ts/url'
import * as rules from "../../ts/formRules"

export default {
    props: {
        csrf: String
    },
    data: () => ({
        getStartedUrl,
        rules,
        email: undefined,
        formValid: false,
    }),
    methods: {
        submit() {

        }
    }
}

</script>
