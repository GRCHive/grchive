<template>
    <v-navigation-drawer app clipped
        :mini-variant="mini"
        :expand-on-hover="mini"
        mini-variant-width=50
        :width="width"
        permanent
        ref="drawer"
    >
        <slot></slot>
        <v-list class="py-0" expand>
            <div v-for="(item, i) in finalNavLinks" 
                 :key="i"
                 :style="!!item.hidden ? `display: none;` : ``"
            >
                <v-list-group v-if="!!item.children && item.children.length > 0"
                              :prepend-icon="item.icon"
                              no-action
                              value="true"
                              color="rgba(0, 0, 0, 0.87) !important"
                >

                    <template v-slot:activator>
                        <v-list-item-title>{{ item.title }} </v-list-item-title>
                    </template>

                    <generic-nav-bar-item v-for="(child, ci) in item.children"
                                          :key="ci"
                                          :item="child"
                                          :primary-color="primaryColor">
                    </generic-nav-bar-item>
                </v-list-group>

                <generic-nav-bar-item :item="item"
                                      :primary-color="primaryColor"
                                      v-else>
                </generic-nav-bar-item>
            </div>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import GenericNavBarItem from './GenericNavBarItem.vue'

export default Vue.extend({
    props : {
        mini : Boolean,
        navLinks: Array,
        primaryColor: {
            type: String,
            default: "primary"
        },
        width: {
            type: Number,
            default: 256
        }
    },
    components: {
        GenericNavBarItem
    },
    data: () => ({
        dispMini: false,
    }),
    computed: {
        finalNavLinks(): Array<any> {
            if (!this.dispMini) {
                return this.navLinks
            } else {
                let fullArray = new Array<any>()

                let expandFn = (ele : any) => {
                    if (!!ele.children && ele.children.length > 0) {
                        let newParent = Object.assign({}, ele)
                        newParent.children = []
                        fullArray.push(newParent)
                        ele.children.forEach(expandFn)
                    } else {
                        fullArray.push(ele)
                    }
                }

                this.navLinks.forEach(expandFn)
                return fullArray
            }
        },
    },
    mounted() {
        this.$watch(
            () => {
                //@ts-ignore
                return this.$refs.drawer.isMiniVariant
            },
            (v : boolean) => {
                this.dispMini = v
            }
        )

        //@ts-ignore
        this.dispMini = this.$refs.drawer.isMiniVariant
    },
})

</script>

<style scoped>

>>>.v-list-group__header .v-list-group__header__prepend-icon .v-icon {
    color: rgba(0, 0, 0, 0.54) !important;
}

>>>.v-navigation-drawer__content {
    overflow-y: hidden;
}

>>>.v-navigation-drawer__content:hover {
    overflow-y: auto;
}

</style>
