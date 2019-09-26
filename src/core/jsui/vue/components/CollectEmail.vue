<template>
    <v-form>
        <v-text-field
            label="Name"
            filled
        ></v-text-field>

        <v-text-field
            label="Email"
            filled
        ></v-text-field>

        <v-switch
            :label="`I agree to ${companyName} collecting and storing my personal information to send me updates about future services and products.`"
        >
        </v-switch>

        <v-btn
            color="success"
        >
            Submit
        </v-btn>

    </v-form>
</template>

<script lang="ts">

export default {
    props: [
        'companyName'
    ]
}

</script>
