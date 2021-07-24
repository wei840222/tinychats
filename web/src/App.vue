<template lang="pug">
router-view
</template>

<script>
import { provide, onBeforeMount } from "vue";
import {
  ApolloClient,
  ApolloLink,
  createHttpLink,
  InMemoryCache,
  from,
} from "@apollo/client/core";
import { createPersistedQueryLink } from "@apollo/client/link/persisted-queries";
import { sha256 } from "crypto-hash";
import { DefaultApolloClient } from "@vue/apollo-composable";
import liff from "@line/liff";

export default {
  setup() {
    const authMiddleware = new ApolloLink((operation, forward) => {
      operation.setContext(({ headers = {} }) => ({
        headers: {
          ...headers,
          Authorization: (() => {
            if (process.env.NODE_ENV !== "production") {
              return process.env.VUE_APP_ACCESS_TOKEN;
            }
            return liff.isLoggedIn() ? liff.getAccessToken() : null;
          })(),
        },
      }));
      return forward(operation);
    });

    const apolloClient = new ApolloClient({
      link: from([
        authMiddleware,
        createPersistedQueryLink({
          useGETForHashedQueries: true,
          sha256,
        }),
        createHttpLink({
          uri: "https://tinychats.herokuapp.com/graphql",
        }),
      ]),
      cache: new InMemoryCache(),
    });
    provide(DefaultApolloClient, apolloClient);

    onBeforeMount(async () => {
      if (process.env.NODE_ENV !== "production") {
        return;
      }
      await liff.init({ liffId: "1656247924-eX5ZOvN0" });
      if (!liff.isLoggedIn()) {
        liff.login();
      }
    });
  },
};
</script>

<style lang="sass">
#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  color: #2c3e50
</style>
