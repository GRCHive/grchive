<template>
    <v-toolbar
        dense
        dark
        id="toolbar"
    >
        <v-toolbar-items id="leftToolbar">
            <v-btn
                icon
                small
                @click="changePage(page - 1)"
            >
                <v-icon>mdi-arrow-up</v-icon>
            </v-btn>

            <v-text-field
                id="pageText"
                :value="pageText"
                hide-details
                solo
                flat
                dense
                @input="changePage(Number(arguments[0]))"
            >
            </v-text-field>

            <span id="pageLabel">
                of {{ totalPages }}
            </span>

            <v-btn
                icon
                small
                @click="changePage(page + 1)"
            >
                <v-icon>mdi-arrow-down</v-icon>
            </v-btn>
        </v-toolbar-items>

        <v-toolbar-items id="middleToolbar">
            <v-btn
                icon
                small
                @click="changeScale(scale - 0.1)"
            >
                <v-icon>mdi-minus</v-icon>
            </v-btn>

            <v-select
                :value="scale"
                :items="scaleOptions"
                hide-details
                solo
                flat
                @input="changeScale"
                dense
            >
            </v-select>

            <v-btn
                icon
                small
                @click="changeScale(scale + 0.1)"
            >
                <v-icon>mdi-plus</v-icon>
            </v-btn>
        </v-toolbar-items>

    </v-toolbar>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        scale: Number,
        page: Number,
        totalPages: Number,
    }
})

@Component
export default class PdfToolbar extends Props {
    get pageText() : string {
        return (this.page + 1).toString()
    }

    get scaleText() : string {
        return Math.round(this.scale * 100.0).toString() + "%"
    }

    get scaleOptions() : any[] {
        let items = [
            {
                text: "50%",
                value: 0.5,
            },
            {
                text: "75%",
                value: 0.75,
            },
            {
                text: "100%",
                value: 1.0,
            },
            {
                text: "125%",
                value: 1.25,
            },
            {
                text: "150%",
                value: 1.5,
            },
            {
                text: "200%",
                value: 2.0,
            },
            {
                text: "300%",
                value: 3.0,
            },
            {
                text: "400%",
                value: 4.0,
            },
            {
                text: this.scaleText,
                value: this.scale
            },
        ]

        items.sort((a, b) => {
            if (a.value < b.value) {
                return -1
            } else if (a.value > b.value) {
                return 1
            } else {
                return 0
            }
        })

        return items
    }

    changeScale(v : number) {
        v = Math.min(Math.max(v, 0.1), 10.0)
        this.$emit("update:scale", v)
    }

    changePage(p : number) {
        if (p < 0) {
            p = 0
        } else if (p >= this.totalPages) {
            p = this.totalPages - 1
        }
        this.$emit("update:page", p)
    }
}

</script>

<style scoped>

.toolbarDiv {
    display: flex;
}

>>>.v-input {
    align-items: center !important;
}

#pageLabel {
    display: flex;
    align-items: center;
}

>>>#pageText {
    border-width: 1px;
    border-style: solid;
    border-color: gray;
    width: 30px !important;
    text-align: right;
}

>>>.v-select {
    width: 150px !important;
}

#toolbar {
    position: relative;
}

#middleToolbar {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
}

</style>
