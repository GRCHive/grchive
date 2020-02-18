<template>
    <div>
        <div v-if="result.data.Success">
            <v-text-field outlined
                          v-model="filterString"
                          prepend-inner-icon="mdi-magnify"
                          label="Search"
                          hide-details
                          clearable
            ></v-text-field>

            <v-data-table
                :headers="headers"
                :items="items"
                :search="filterString"
            >
            </v-data-table>
        </div>

        <v-row v-else align="center" justify="center">
            <div class="max-width">
                <p class="subtitle-1 text-center">Oops! We were unable to run your SQL query.</p>
                <v-divider class="mb-3"></v-divider>
                <pre>{{ result.data.Data }}</pre>
            </div>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { TRunSqlQueryOutput } from '../../ts/api/apiSqlQueries'
import { SqlResult } from '../../ts/sql'
import Papa from 'papaparse'

const Props = Vue.extend({
    props: {
        result: {
            type: Object,
            default: () => Object() as TRunSqlQueryOutput
        }
    }
})

@Component
export default class SqlResultDisplay extends Props {
    filterString : string = ""

    get headers() : any[] {
        if (!this.result.data.Success) {
            return []
        }
        return (<SqlResult>this.result.data.Data).Columns.map((ele : string) => ({
            text: ele,
            value: ele,
        }))
    }

    get items() : any[] {
        if (!this.result.data.Success) {
            return []
        }

        let allColumns = (<SqlResult>this.result.data.Data).Columns
        let parsedData : any = Papa.parse((<SqlResult>this.result.data.Data).CsvText)
        return parsedData.data.filter((ele : Array<any>) => ele.length == allColumns.length).map(
            (ele : Array<any>) => {
                let obj = Object()
                for (let i = 0; i < ele.length; ++i) {
                    obj[allColumns[i]] = ele[i]
                }
                return obj
            })
    }
}

</script>
