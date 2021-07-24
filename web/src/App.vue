<template lang="pug">
#nav
  router-link(to="/")
    van-button Home
  router-link(to="/about") 
    van-button About
router-view
</template>

<script>
import { provide } from "vue";
import {
  ApolloClient,
  createHttpLink,
  InMemoryCache,
} from "@apollo/client/core";
import { createPersistedQueryLink } from "@apollo/client/link/persisted-queries";
import { sha256 } from "crypto-hash";
import { DefaultApolloClient } from "@vue/apollo-composable";

export default {
  setup() {
    const apolloClient = new ApolloClient({
      link: createPersistedQueryLink({
        useGETForHashedQueries: true,
        sha256,
      }).concat(
        createHttpLink({
          uri: "https://wei840222-todo.herokuapp.com/graphql",
        })
      ),
      cache: new InMemoryCache(),
    });
    provide(DefaultApolloClient, apolloClient);
  },
};
</script>

<style lang="sass">
#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  text-align: center
  color: #2c3e50

#nav
  padding: 30px
  a
    font-weight: bold
    color: #2c3e50
    &.router-link-exact-active
      color: #42b983
</style>
