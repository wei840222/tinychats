<template lang="pug">
div 123
router-view
</template>

<script>
import { provide, onMounted } from "vue";
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
          Authorization: window.localStorage.getItem("accessToken") || null,
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

    onMounted(async () => {
      if (process.env.NODE_ENV !== "production") {
        window.localStorage.setItem(
          "accessToken",
          process.env.VUE_APP_ACCESS_TOKEN
        );
        return;
      }
      await liff.init({ liffId: "1656247924-eX5ZOvN0" });
      if (!liff.isLoggedIn()) {
        liff.login();
        return;
      }
      window.localStorage.setItem("accessToken", liff.getAccessToken());
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
