<template>
    <v-data-table
        v-model="selected"
        :headers="userTableHeaders"
        :items="userTableItems"
        :show-select="selectable"
        :single-select="!multi"
        @input="changeInput">
    </v-data-table>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
    props: {
        users: Array,
        value: {
            type: Array,
            default: () => []
        },
        selectable: {
            type: Boolean,
            default: false
        },
        multi: {
            type: Boolean,
            default: false
        },
    },
    data : () => ({
        userTableHeaders: [
            {
                text: 'Name',
                value: 'fullName'
            },
            {
                text: 'Email',
                value: 'email'
            },
        ],
        selected: []
    }),
    computed: {
        userTableItems() : any[] {
            return this.users.map((ele : any) => ({
                fullName: `${ele.FirstName} ${ele.LastName}`,
                email: ele.Email,
                user: ele,
                id: ele.Id
            }))
        }
    },
    methods: {
        changeInput(items : any[]) {
            this.$emit('input', items.map((ele : any) => ele.user))
        }
    }
})

</script>
