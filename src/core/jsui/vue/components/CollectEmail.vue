<template>
    <v-form v-model="formValid" ref="form">
        <v-text-field
            v-model="name"
            label="Name"
            filled
            counter="100"
            :rules="[rules.required, rules.createMaxLength(100)]"
            required
        ></v-text-field>

        <v-text-field
            v-model="email"
            label="Email"
            type="email"
            filled
            counter="100"
            :rules="[rules.required, rules.createMaxLength(100), rules.email]"
            required
        ></v-text-field>

        <v-switch
            v-model="agree"
            :label="`I agree to ${companyName} collecting and storing my personal information to send me updates about future services and products.`"
            required
            :rules="[rules.required]"
        >
        </v-switch>

        <v-btn
            color="success"
            class="my-2"
            :disabled="!formValid || !name || !email || !agree"
            @click="submit"
        >
            Submit
        </v-btn>

        <v-snackbar 
            v-model="showSnack"
            :color="snackIsError ? 'error' : 'success'"
            bottom
            :style="'transform: translateY('+snackTransformY+');'"
            timeout=10000
        >
            {{ snackText }} 

            <v-btn
                color="primary"
                v-if="snackShowContact"
                :href="createContactUsUrl()"
            >
                Contact Us
            </v-btn>

            <v-btn
                color="secondary"
                @click="showSnack = false"
            >
                Close
            </v-btn>
        </v-snackbar>
    </v-form>
</template>

<script lang="ts">

import * as rules from "../../ts/formRules"
import { createContactUsUrl } from "../../ts/url"
import axios from 'axios';
import Vue from 'vue';

export default Vue.extend({
    props: {
        'companyName' : String,
        'snackTransformY' : {
            type: String,
            default: "0px"
        }
    },
    data: () => ({
        name: undefined,
        email: undefined,
        agree: false,
        rules: rules,
        formValid: false,

        // Snackbar options for giving visual feedback to the user.
        showSnack: false,
        snackText: "",
        snackShowContact: false,
        snackIsError: false
    }),
    methods: {
        submit() {
            // ts-ignore is needed since as of Vuetify 2.0.19 we can't type assert the form
            // to a Vuetify VForm to have TypeScript properly fin the validate function.
            // @ts-ignore
            if (!this.$refs.form.validate()) {
                this.$data.snackText = "Oops! Something went wrong. Please refresh the page and try again.";
                this.$data.showSnack = true;
                return;
            }

            axios.post('#', {
                name: this.$data.name,
                email: this.$data.email
            }).then(() => {
                this.$data.snackShowContact = false;
                this.$data.snackText = "Success! We will reach out soon.";
                this.$data.snackIsError = false;
            }).catch((err) => {
                this.$data.snackShowContact = true;
                this.$data.snackIsError = true;
                if (err.response.data.IsDuplicate) {
                    this.$data.snackText = "Oops! It looks like you already gave us your email address. Contact us if you still need something.";
                } else {
                    this.$data.snackText = "Oops! It looks like something went wrong on our end. Try again later or get in touch directly.";
                }
            }).finally(() => {
                this.$data.showSnack = true;
            })
        },
        createContactUsUrl
    }
});

</script>
