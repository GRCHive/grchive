<template>
    <span>
        <a
            v-if="!!handle.resourceUri"
            :href="handle.resourceUri"
        >
            {{ handle.displayText }}
        </a>

        <span v-else>
            {{ handle.displayText }}
        </span>
    </span>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ResourceHandle } from '../../ts/resourceUtils'

const Props = Vue.extend({
    props: {
        handle: {
            type: Object,
            default: () => Object() as ResourceHandle
        }
    }
})

@Component
export default class ResourceHandleRenderer extends Props {
}

</script>
