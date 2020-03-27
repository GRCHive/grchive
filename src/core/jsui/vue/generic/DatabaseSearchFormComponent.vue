<template>

<v-autocomplete
     filled
    :value="value"
    :items="items"
    :rules="rules"
    :label="label"
    @input="onInput"
>

</v-autocomplete>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { Database } from '../../ts/databases'
import { TAllDatabaseOutputs, allDatabase } from '../../ts/api/apiDatabases'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import { ComparisonOperators } from '../../ts/filters'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => Object() as Database | null
        },
        rules: {
            type: Array,
            default: () => []
        },
        typeId: {
            type: Number,
            default: -1
        },
        label: {
            type: String,
            default: "",
        }
    }
})

@Component
export default class DatabaseSearchFormComponent extends Props {
    allDbs : Database[] = []

    mounted() {
        this.refreshData()
    }

    @Watch('typeId')
    refreshData() {
        allDatabase({
            orgId: PageParamsStore.state.organization!.Id,
            filter: {
                Type: {
                    Op: this.typeId == -1 ? ComparisonOperators.Disabled : ComparisonOperators.Equal,
                    Target: this.typeId,
                }
            }
        }).then((resp : TAllDatabaseOutputs) => {
            this.allDbs = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    get items() : any[] {
        return this.allDbs.map((ele : Database) => ({
            text: ele.Name,
            value: ele,
        }))
    }

    onInput(val : Database) {
        this.$emit('input', val)
    }
}

</script>
